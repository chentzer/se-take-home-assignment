#!/bin/bash

# Run Script
# This script should execute your CLI application and output results to result.txt

echo "Running CLI application..."

cd "$(dirname "$0")/../code" || exit 1

# Clear previous results
echo "" > ../scripts/result.txt

# Send commands to the interactive CLI
{
    echo "normal"
    echo "normal"
    echo "vip"
    echo "vip"
    echo "status"
    echo "addbot"
    echo "addbot"
    echo "status"
    echo "removebot"
    echo "normal"
    echo "vip"
    echo "addbot"
    echo "status"
    echo "exit"
} | ./order-controller > ../scripts/result.txt 2>&1

echo "CLI application execution completed"