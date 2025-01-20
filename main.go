package main

import (
//     "fmt"
    "log"
    "net/http"
    "os"
//     "database/sql"
//     "github.com/gorilla/mux"
    "github.com/joho/godotenv"
//     "project_sem/handlers"
    "project_sem/models"
    "project_sem/routes"
)

// import "fmt"

func main() {
    // Load config
    config := loadEnv()

    // Init DB
    models.InitDB(config)
    defer db.Close()

    // Register routes
    router := routes.RegisterRoutes(db)

    // Start server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}

func loadEnv() map[string]string {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, loading environment variables from the environment")
    }

    // Variables we need to get
    envVars := []string{"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"}
    config := make(map[string]string)

    // Load environment variables into a map
    for _, key := range envVars {
        value := os.Getenv(key)
        if value == "" {
            log.Printf("Error: %s is not set\n", key)
        }
        config[key] = value
    }

    return config
}
