package routes

import (
    "database/sql"
    "github.com/gorilla/mux"
    "net/http"

    "project_sem/handlers"
)

// register routes
func RegisterRoutes(db *sql.DB) *mux.Router {
    router := mux.NewRouter()

    // POST /api/v0/prices
    router.HandleFunc("/api/v0/prices", handlers.PricesPOST(db)).Methods(http.MethodPost)

    // GET /api/v0/prices
    router.HandleFunc("/api/v0/prices", handlers.PricesGET(db)).Methods(http.MethodGet)

    return router
}