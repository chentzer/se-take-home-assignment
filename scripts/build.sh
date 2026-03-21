#!/bin/bash
set -e

echo "Building Go CLI..."

# Go to code folder
cd ../code || { echo "Code folder not found"; exit 1; }

# Compile CLI
go build -o app

echo "Build completed. Executable created as code/app"