#!/bin/bash

echo "Running McDonald's Order Management System Tests..."

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../code" || { 
    echo "Code folder not found"
    exit 1
}

echo "Running from: $(pwd)"
echo ""

# Run tests without race detection first
echo "Running standard tests..."
go test -v -short

if [ $? -ne 0 ]; then
    echo "Standard tests failed!"
    exit 1
fi

echo ""
echo "Running concurrent operation tests..."
go test -v -run TestConcurrentOperations -timeout 2m

if [ $? -ne 0 ]; then
    echo "Concurrent operation tests failed!"
    exit 1
fi

echo ""
echo "All tests passed!"