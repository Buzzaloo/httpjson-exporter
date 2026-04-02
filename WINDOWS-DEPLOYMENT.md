# Deploying to Windows - Complete Guide

There are **3 ways** to run this on your Windows servers. Pick the one that works best for you.

---

## Option 1: Build Windows .exe Directly (Recommended if you have Go)

### Requirements
- Go 1.21+ installed on Windows (or any machine)
- Download from: https://go.dev/dl/

### Steps

1. **On a Windows machine with Go installed:**
```powershell
cd otelcol-exaforce
.\build-windows.ps1
```

2. **This creates:** `otelcol-custom\otelcol-custom.exe`

3. **Replace your existing collector:**
```powershell
# Stop the service
Stop-Service OpenTelemetryCollector

# Backup the old one
Copy-Item "C:\Program Files\OpenTelemetry Collector\otelcol-contrib.exe" `
          "C:\Program Files\OpenTelemetry Collector\otelcol-contrib.exe.backup"

# Copy the new one
Copy-Item ".\otelcol-custom\otelcol-custom.exe" `
          "C:\Program Files\OpenTelemetry Collector\otelcol-contrib.exe"

# Update config (see below)
# Set environment variable (see below)

# Start the service
Start-Service OpenTelemetryCollector
```

---

## Option 2: Use GitHub Actions to Build (No Go Required!)

I can give you a GitHub Actions workflow that builds the Windows binary automatically.

### Steps

1. **Create a GitHub repo** and push this code
2. **Add this workflow** (I'll create it below)
3. **Download the built .exe** from GitHub Actions artifacts
4. **Deploy to your Windows servers**

This is great if you don't want to install Go locally.

---

## Option 3: Extract from Docker (Quick & Dirty)

If you have Docker on any machine (Windows, Mac, Linux):

```bash
# Build the Linux version first
./build.sh

# Create a temporary container
docker create --name temp otelcol-exaforce:latest

# Try to extract (this gets you the Linux binary)
docker cp temp:/otelcol ./otelcol-linux

# Clean up
docker rm temp
```

**Problem:** This gives you a Linux binary, not Windows.

**Solution:** Cross-compile:
```bash
# On Linux/Mac/WSL with the source code
cd httpjsonexporter
GOOS=windows GOARCH=amd64 go build -o ../otelcol-windows.exe
```

---

## Updating Your Config

Once you have `otelcol-custom.exe` on Windows, update your config:

**File:** `C:\Program Files\OpenTelemetry Collector\config.yaml`

### Replace your awss3 exporter with:

```yaml
exporters:
  httpjson:
    endpoint: "https://your-tenant.us.app.exaforce.io/api/v1/logs/ingest"
    bearer_token: "${EXAFORCE_API_TOKEN}"
    compression: gzip
    timeout: 30s
    batch_size: 100
    retry_on_failure:
      enabled: true
      initial_interval: 5s
      max_interval: 30s
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 1000
```

### Update the pipeline:

```yaml
service:
  pipelines:
    logs:
      receivers: [windowseventlog/system, windowseventlog/application, windowseventlog/security]
      processors: [memory_limiter, resourcedetection, batch]
      exporters: [httpjson]  # Changed from awss3
```

### Set the environment variable:

```powershell
[System.Environment]::SetEnvironmentVariable(
    "EXAFORCE_API_TOKEN", 
    "your-actual-token-here", 
    "Machine"
)
```

### Restart the service:

```powershell
Restart-Service OpenTelemetryCollector
```

---

## Verify It's Working

### Check the service:
```powershell
Get-Service OpenTelemetryCollector
```

### Check logs:
```powershell
Get-Content "C:\ProgramData\OpenTelemetry Collector\otelcol.log" -Tail 50
```

### Test manually:
```powershell
$headers = @{
    "Authorization" = "Bearer your-token-here"
    "Content-Type" = "application/x-ndjson"
}

$body = '{"test":"log","timestamp":"2026-02-05T12:00:00Z","message":"test from windows"}'

Invoke-WebRequest -Uri "https://your-tenant.app.exaforce.io/api/v1/logs/ingest" `
    -Method Post `
    -Headers $headers `
    -Body $body
```

---

## Which Option Should I Use?

| Option | Best For | Difficulty |
|--------|----------|-----------|
| **Option 1** (build-windows.ps1) | You have Go installed | Easy |
| **Option 2** (GitHub Actions) | You want automation | Medium |
| **Option 3** (Docker extract) | Quick one-off | Easy but limited |

### My Recommendation:

**If you have multiple Windows servers:** Use GitHub Actions (Option 2) - build once, deploy everywhere

**If you have one Windows server and Go:** Use build-windows.ps1 (Option 1) - quickest

**If you're just testing:** Use Docker as gateway - keep your existing collector, just point it at the Docker container

---

## Need the GitHub Actions Workflow?

Let me know and I'll create a complete GitHub Actions workflow that:
1. Builds Windows .exe on every commit
2. Creates releases with downloadable binaries
3. Supports multiple architectures (amd64, arm64)

---

## Troubleshooting

### "otelcol-custom.exe is not recognized"
- Make sure it's in `C:\Program Files\OpenTelemetry Collector\`
- Or update the service to point to the new path

### "Service failed to start"
- Check logs: `C:\ProgramData\OpenTelemetry Collector\otelcol.log`
- Verify config.yaml syntax: Run `otelcol-custom.exe validate --config config.yaml`

### "Unknown exporter: httpjson"
- You're using the old binary
- Make sure you replaced the .exe with the custom one

### "Cannot find otelcol-contrib.exe"
- The new binary might have a different name
- Update the service to use `otelcol-custom.exe`
