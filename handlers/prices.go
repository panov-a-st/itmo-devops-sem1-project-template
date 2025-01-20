package handlers

import (
    "database/sql"
    "net/http"
)

func PricesGET(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func PricesPOST(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

