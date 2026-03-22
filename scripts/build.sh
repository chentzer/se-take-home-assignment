#!/bin/bash

echo "Building CLI application..."

cd "$(dirname "$0")/.." || exit 1

# Build the application
go build -o cmd/order-controller ./cmd

if [ $? -eq 0 ]; then
    echo "Build successful! Binary created: cmd/order-controller"
else
    echo "Build failed!"
    exit 1
fi