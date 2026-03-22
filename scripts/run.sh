#!/bin/bash

echo "Running CLI application..."

# Navigate to code directory
cd "$(dirname "$0")/../code" || exit 1

# Clear previous results
echo "" > ../scripts/result.txt

# Run the CLI application with commands
{
    echo "normal"
    echo "vip"
    echo "normal"
    echo "addbot"
    echo "addbot"
    sleep 2
    echo "status"
    sleep 12
    echo "status"
    echo "exit"
} | go run . > ../scripts/result.txt 2>&1

echo "CLI application execution completed"