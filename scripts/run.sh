#!/bin/bash

# Exit on error
set -e

# Set up env variables
echo "Setting up env variables"
export $(grep -v '^#' .env | xargs)

# Create prices
echo "Create prices table if not exists"
PGPASSWORD=$POSTGRES_PASSWORD psql -U $POSTGRES_USER -h $POSTGRES_HOST -p $POSTGRES_PORT -d $POSTGRES_DB -c "
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at DATE NOT NULL
);"

# Build application
go build -o app ./main.go

# Run application in bg
setsid ./app > app.log 2>&1 &
#./app

echo "Running. PID: $!"
