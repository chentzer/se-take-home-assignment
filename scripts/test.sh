#!/bin/bash
# Test Script
# Runs Go unit tests

echo "Running unit tests..."

# Go to code folder
cd ../code || { echo "Code folder not found"; exit 1; }

# Run tests with verbose output
go test ./... -v

if [ $? -eq 0 ]; then
    echo "All tests passed!"
else
    echo "Some tests failed!"
    exit 1
fi

echo "Unit tests completed"