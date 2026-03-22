#!/bin/bash

echo "Building McDonald's Order Management System..."

cd "$(dirname "$0")/../code" || exit 1

# Build the application
go build -o mcdonalds-bot .

if [ $? -eq 0 ]; then
    echo "Build successful! Binary created: code/mcdonalds-bot"
else
    echo "Build failed!"
    exit 1
fi