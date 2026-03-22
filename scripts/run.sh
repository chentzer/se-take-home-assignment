#!/bin/bash

echo "Running CLI application..."

# Navigate to code directory
cd "$(dirname "$0")/../code" || exit 1

# Clear previous results
> ../scripts/result.txt

# Run the CLI application in demo mode
# The -demo flag runs a predefined sequence that demonstrates all functionality
# The -output flag specifies where to write the log output
./order-controller -demo -output ../scripts/result.txt

echo "CLI application execution completed"
echo "Results written to scripts/result.txt"