package models

import (
    "fmt"
    "database/sql"
)

var db *sql.DB

func InitDB(config map[string]string) (*sql.DB, error) {
    connStr := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        config["POSTGRES_USER"],
        config["POSTGRES_PASSWORD"],
        config["POSTGRES_HOST"],
        config["POSTGRES_PORT"],
        config["POSTGRES_DB"],
    )

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, err
}
