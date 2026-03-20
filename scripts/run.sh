#!/bin/bash
# Run Script
# Runs the compiled CLI

echo "Running the CLI application..."

# Go to code folder
cd ../code || { echo "Code folder not found"; exit 1; }

# Run the CLI
./app.exe

if [ $? -eq 0 ]; then
    echo "CLI exited successfully"
else
    echo "CLI exited with errors"
    exit 1
fi