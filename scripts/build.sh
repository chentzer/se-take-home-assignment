#!/bin/bash
# Build Script
# Compiles Go CLI application

echo "Building Go application..."

# Go to code folder
cd ../code || { echo "Code folder not found"; exit 1; }

# Build the CLI executable
go build -o app.exe

if [ $? -eq 0 ]; then
    echo "Build succeeded!"
else
    echo "Build failed!"
    exit 1
fi