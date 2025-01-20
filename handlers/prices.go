package handlers

import (
    "log"
    "bytes"
    "io"
    "strings"
    "fmt"
    "strconv"
    "time"
    "database/sql"
    "net/http"
    "archive/zip"
    "encoding/csv"
    "encoding/json"
)

// Define type for JSON response
type TypeResponse struct {
    TotalItems      int     `json:"total_items"`
    TotalCategories int     `json:"total_categories"`
    TotalPrice      float64 `json:"total_price"`
}

func PricesGET(db *sql.DB) http.HandlerFunc {
     return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        rows, err := db.Query(`SELECT id, product_name, category, price, created_at FROM prices`)
        if err != nil {
            log.Printf("Error reading db: %v", err)
            http.Error(w, "Error reading db", http.StatusInternalServerError)
            return
        }
    	defer rows.Close()

        // Read data for the csv using buffer
        var csvData bytes.Buffer
        csvWriter := csv.NewWriter(&csvData)
        for rows.Next() {
            var (
                productId    int
                productName  string
                category     string
                price        float64
                createdAt    time.Time
            )
            if err := rows.Scan(&productId, &productName, &category, &price, &createdAt); err != nil {
                log.Printf("Failed to get record: %v", err)
                http.Error(w, "Failed to get records", http.StatusInternalServerError)
                continue
            }

            record := []string{
                strconv.Itoa(productId),
                productName,
                category,
                fmt.Sprintf("%.2f", price),
                createdAt.Format("2006-01-02")}
            if err := csvWriter.Write(record); err != nil {
                log.Printf("Error appending CSV: %v", err)
                http.Error(w, "Error appending CSV", http.StatusInternalServerError)
                return
            }
        }
        csvWriter.Flush()
        if err := csvWriter.Error(); err != nil {
            log.Printf("Error flushing CSV writer: %v", err)
            http.Error(w, "Error writing CSV data", http.StatusInternalServerError)
            return
        }

        // Create zip file and send it
        zipBuffer := &bytes.Buffer{}
        zipWriter := zip.NewWriter(zipBuffer)
        if csvFile, err := zipWriter.Create("data.csv"); err != nil {
            log.Printf("Error creating file in ZIP: %v", err)
            http.Error(w, "Failed to create ZIP", http.StatusInternalServerError)
            return
        } else if _, err := csvFile.Write(csvData.Bytes()); err != nil {
            log.Printf("Error writing CSV to ZIP: %v", err)
            http.Error(w, "Failed to write ZIP", http.StatusInternalServerError)
            return
        }

        if err := zipWriter.Close(); err != nil {
            log.Printf("Error processing ZIP: %v", err)
            http.Error(w, "Failed to process ZIP", http.StatusInternalServerError)
            return
        }

    	// Send zip
    	w.Header().Set("Content-Type", "application/zip")
    	w.Header().Set("Content-Disposition", "attachment; filename=\"data.zip\"")
    	w.WriteHeader(http.StatusOK)
    	if _, err := w.Write(zipBuffer.Bytes()); err != nil {
    		log.Printf("Error generating file: %v", err)
            http.Error(w, "Error generating file", http.StatusInternalServerError)
    	}
    }
}

func PricesPOST(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var totalItems int
        var totalCats int
        var totalPrice float64

        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        file, _, err := r.FormFile("file")
        if err != nil {
            log.Printf("Error retrieving file: %v", err)
            http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
            return
        }
        defer file.Close()

        // Read archive
        zipBuffer := new(bytes.Buffer)
        if _, err := io.Copy(zipBuffer, file); err != nil {
            log.Printf("Error reading file: %v", err)
            http.Error(w, "Failed to read file", http.StatusInternalServerError)
            return
        }

        // Open zip file
        zipReader, err := zip.NewReader(bytes.NewReader(zipBuffer.Bytes()), int64(zipBuffer.Len()))
        if err != nil {
            log.Printf("Error opening zip: %v", err)
            http.Error(w, "Invalid zip file", http.StatusBadRequest)
            return
        }

        // Get records from zip archive
        var dataRecords [][]string

        for _, zipFile := range zipReader.File {
            if strings.HasSuffix(zipFile.Name, "data.csv") {
                log.Printf("Processing file: %s", zipFile.Name)
                f, err := zipFile.Open()
                if err != nil {
                    http.Error(w, "Can't open data.csv", http.StatusInternalServerError)
                    return
                }
                defer f.Close()

                csvReader := csv.NewReader(f)
                // Skip the first line
                _, err = csvReader.Read()
                if err == io.EOF {
                    http.Error(w, "CSV file is empty", http.StatusBadRequest)
                    return
                }
                if err != nil {
                    http.Error(w, "Error reading CSV", http.StatusBadRequest)
                    return
                }

                // Process lines
                skipped := 0
                for {
                    row, err := csvReader.Read()
                    if err != nil {
                        if err == io.EOF {
                            break
                        }
                        http.Error(w, "Error reading CSV", http.StatusInternalServerError)
                        return
                    }

                    if len(row) < 5 {
                        skipped++
                        continue
                    }

                    dataRecords = append(dataRecords, row)
                }
                log.Printf("CSV has been processed: %d ok, %d rows skipped", len(dataRecords), skipped)
            }
        }

        // Store data
        tx, err := db.Begin()
        if err != nil {
            log.Printf("Error opening transaction: %v", err)
            http.Error(w, "Error opening transaction", http.StatusInternalServerError)
            return
        }
        defer func() {
            if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
                log.Printf("Error rollbacking transaction: %v", rollbackErr)
            }
        }()

        // Process records
        for _, dataRecord := range dataRecords {
            _, err := tx.Exec(`
                INSERT INTO prices (id, product_name, category, price, created_at)
                VALUES ($1, $2, $3, $4, $5)`,
                dataRecord[0], dataRecord[1], dataRecord[2], dataRecord[3], dataRecord[4])
            if err != nil {
                log.Printf("Error inserting record: %v", err)
                continue
            }
        }

        // Get totals
        row := tx.QueryRow(`
            SELECT COUNT(DISTINCT product_name), COUNT(DISTINCT category), COALESCE(SUM(price), 0) FROM prices
        `)
        if err := row.Scan(&totalItems, &totalCats, &totalPrice); err != nil {
            log.Printf("Failed to get totals: %v", err)
            http.Error(w, "Failed to get totals", http.StatusInternalServerError)
            return
        }

        // Commit
        if err := tx.Commit(); err != nil {
            log.Printf("Error commiting transaction: %v", err)
            http.Error(w, "Error commiting transaction", http.StatusInternalServerError)
            return
        }

        // Send response
        resp := TypeResponse {
            TotalItems:      totalItems,
            TotalCategories: totalCats,
            TotalPrice:      totalPrice,
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error encoding JSON: %v", err)
        }
    }
}


