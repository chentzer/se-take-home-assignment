# Test Script for PowerShell
Write-Host "Running unit tests..."

# Remember current folder and go to code folder
Push-Location ..\code

# Run Go tests
go test ./... -v

# Check result
if ($LASTEXITCODE -eq 0) {
    Write-Host "All tests passed!"
} else {
    Write-Host "Some tests failed!"
    Pop-Location  # return before exit
    exit 1
}

# Return to original folder
Pop-Location

Write-Host "Unit tests completed"