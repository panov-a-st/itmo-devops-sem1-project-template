package models

import (
    "fmt"
    "database/sql"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(config map[string]string) (*sql.DB, error) {
    connStr := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        config["POSTGRES_USER"],
        config["POSTGRES_PASSWORD"],
        config["POSTGRES_HOST"],
        config["POSTGRES_PORT"],
        config["POSTGRES_DB"],
    )

    DB, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    if err := DB.Ping(); err != nil {
        return nil, err
    }

    return DB, err
}
