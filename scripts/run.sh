#!/bin/bash

# Run Script
echo "Running CLI application..."

cd "$(dirname "$0")/.." || exit 1

echo "" > ./scripts/result.txt

# Send commands to the binary
{
    echo "normal"
    echo "vip"
    echo "normal"
    echo "addbot"
    echo "addbot"
    echo "status"
    echo "exit"
} | ./order-controller > ./scripts/result.txt 2>&1

echo "CLI application execution completed"