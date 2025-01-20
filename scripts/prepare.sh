#!/bin/bash

# Set up env variables
echo "Setting up env variables"
export $(grep -v '^#' .env | xargs)

# Install go dependencies
echo "Installing go dependencies"
go mod tidy

# Create prices
echo "Creating prices table"
PGPASSWORD=$POSTGRES_PASSWORD psql -U $POSTGRES_USER -h $POSTGRES_HOST -p $POSTGRES_PORT -d $POSTGRES_DB -c "
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at DATE NOT NULL
);"

echo "Done!"
