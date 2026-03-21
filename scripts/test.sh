#!/bin/bash
set -e

echo "Running Go unit tests..."

# Go to code folder
cd code || { echo "Code folder not found"; exit 1; }

# Run all Go tests with verbose output
go test ./... -v

echo "Unit tests completed successfully"