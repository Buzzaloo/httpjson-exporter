# Build script for Windows binary
# Run this on a machine with Go installed, or use GitHub Actions

Write-Host "=========================================="
Write-Host "Building Custom OTEL Collector for Windows"
Write-Host "=========================================="
Write-Host ""

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Go is not installed!" -ForegroundColor Red
    Write-Host "Please install Go from: https://go.dev/dl/"
    exit 1
}

# Install the OpenTelemetry Collector Builder
Write-Host "Installing OpenTelemetry Collector Builder..."
go install go.opentelemetry.io/collector/cmd/builder@v0.95.0

# Build the custom collector
Write-Host "Building custom collector for Windows..."
$env:GOOS = "windows"
$env:GOARCH = "amd64"
& "$env:USERPROFILE\go\bin\builder.exe" --config builder-config.yaml

Write-Host ""
Write-Host "=========================================="
Write-Host "Build Complete!"
Write-Host "=========================================="
Write-Host ""
Write-Host "Binary location: .\otelcol-custom\otelcol-custom.exe"
Write-Host ""
Write-Host "Next steps:"
Write-Host "1. Copy otelcol-custom.exe to C:\Program Files\OpenTelemetry Collector\"
Write-Host "2. Update your config.yaml to use the httpjson exporter"
Write-Host "3. Restart the OpenTelemetry Collector service"
Write-Host ""
