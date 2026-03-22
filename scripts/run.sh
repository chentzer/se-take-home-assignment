#!/bin/bash

echo "Running McDonald's Order Management System..."

cd "$(dirname "$0")/../code" || exit 1

# Clear previous output
echo "" > ../scripts/result.txt

# Run the application
./mcdonalds-bot