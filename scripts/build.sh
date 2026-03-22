#!/bin/bash

echo "Building CLI application..."

cd "$(dirname "$0")/../code" || exit 1

# Build the application
go build -o order-controller .

if [ $? -eq 0 ]; then
    echo "Build successful! Binary created: code/order-controller"
else
    echo "Build failed!"
    exit 1
fi