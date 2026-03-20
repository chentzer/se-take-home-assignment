Write-Host "Running CLI application..."

Push-Location ..\code

Start-Process -NoNewWindow -Wait -FilePath .\app.exe

if ($LASTEXITCODE -eq 0) {
    Write-Host "CLI exited successfully"
} else {
    Write-Host "CLI exited with errors"
    Pop-Location
    exit 1
}

Pop-Location