Write-Host "Building Go application..."

Push-Location ..\code

# Compile CLI
go build -o app.exe

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build succeeded!"
} else {
    Write-Host "Build failed!"
    Pop-Location
    exit 1
}

Pop-Location