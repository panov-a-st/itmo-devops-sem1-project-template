#!/bin/bash

# Exit on error
set -e

# Build application
go build -o app ./main.go

# Run application in bg
setsid ./app > app.log 2>&1 &
#./app

echo "Running. PID: $!"
