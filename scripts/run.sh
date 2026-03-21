#!/bin/bash
set -e

echo "Running CLI simulation..."

# Go to code folder
cd code || { echo "Code folder not found"; exit 1; }

# Feed commands to the CLI automatically
# Example: 2 normal orders, 2 VIP orders, add 2 bots, show status, exit
echo -e "normal\nvip\nnormal\nvip\naddbot\naddbot\nstatus\nexit" | ./app

echo "CLI finished. Logs written to ../scripts/result.txt"