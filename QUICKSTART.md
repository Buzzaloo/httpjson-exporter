# 🚀 QUICKSTART - Get Running in 5 Minutes

I've built you a complete, ready-to-run custom OpenTelemetry Collector that sends Windows Event Logs directly to Exaforce's HTTP endpoint with Bearer token auth.

## What You Got

A complete package with:
- ✅ Custom OTEL Collector with HTTP JSON exporter
- ✅ Docker build setup (works on any platform)
- ✅ Config files ready to go
- ✅ Scripts to make it easy

## Run This (Literally 3 Commands)

### 1. Extract the package
```bash
tar -xzf otelcol-exaforce.tar.gz
cd otelcol-exaforce
```

### 2. Set your Exaforce token
```bash
cp .env.template .env
# Edit .env and put your real Exaforce token
nano .env  # or vim, or code, whatever you use
```

### 3. Build and run
```bash
./build.sh
docker-compose up -d
```

## That's It!

Now your collector is running and ready to forward logs to Exaforce.

## Configure Your Windows Collector

On your Windows server, update your OTEL collector config to use the new exporter:

**Open**: `C:\Program Files\OpenTelemetry Collector\config.yaml`

**Find the exporters section** and add:

```yaml
exporters:
  httpjson:
    endpoint: "https://YOUR-TENANT.us.app.exaforce.io/api/v1/logs/ingest"
    bearer_token: "${EXAFORCE_API_TOKEN}"
    compression: gzip
    timeout: 30s
```

**Update the service pipeline:**

```yaml
service:
  pipelines:
    logs:
      receivers: [windowseventlog/system, windowseventlog/application, windowseventlog/security]
      processors: [memory_limiter, resourcedetection, batch]
      exporters: [httpjson]  # Changed from awss3
```

**Set the environment variable:**

```powershell
[System.Environment]::SetEnvironmentVariable("EXAFORCE_API_TOKEN", "your-token-here", "Machine")
```

**Restart the service:**

```powershell
Restart-Service OpenTelemetryCollector
```

## 🪟 Running on Windows (Your Actual Question!)

**YES**, this will run on Windows just like otel-contrib! You just need to build a Windows `.exe` instead of a Linux binary.

### Three Ways to Get the Windows Binary:

### ✅ Option 1: Build on Windows (Easiest if you have Go)

**Requirements:** Go 1.21+ installed  
**Download:** https://go.dev/dl/

```powershell
cd otelcol-exaforce
.\build-windows.ps1
```

This creates `otelcol-custom\otelcol-custom.exe` - ready to use!

Then just **replace** your existing collector exe:
```powershell
Stop-Service OpenTelemetryCollector
Copy-Item ".\otelcol-custom\otelcol-custom.exe" "C:\Program Files\OpenTelemetry Collector\otelcol-contrib.exe"
# Update config.yaml (see WINDOWS-DEPLOYMENT.md)
Start-Service OpenTelemetryCollector
```

---

### ✅ Option 2: Use GitHub Actions (No Go Required!)

1. Push this code to GitHub
2. GitHub Actions automatically builds Windows .exe
3. Download from Actions artifacts
4. Deploy to Windows

**Workflow file included:** `.github/workflows/build-windows.yaml`

---

### ✅ Option 3: Docker as Gateway (Keep Current Setup)

Don't change Windows at all! Instead:
1. Run this Docker container on ANY Linux machine
2. Keep your existing Windows OTEL collector
3. Point Windows → Docker container → Exaforce

Your Windows config stays almost the same, just change endpoint to point at the Docker container.

---

### 📖 Full Windows Instructions

See **WINDOWS-DEPLOYMENT.md** for complete step-by-step instructions for all three options.

## Check If It's Working

```bash
# View logs
docker-compose logs -f

# Check health
curl http://localhost:13133

# Check if logs are going to Exaforce
# (Look in your Exaforce console)
```

## Common Questions

**Q: Do I need to rebuild my entire OTEL collector?**  
A: Not necessarily! You can run this as a separate container and forward to it.

**Q: Will this work with my existing config?**  
A: Yes! Just swap the exporter from `awss3` to `httpjson`.

**Q: What if I don't have Docker?**  
A: You can build directly with Go, or I can help you build a Windows .exe

**Q: Can I test without Docker?**  
A: Yes, but you'll need Go installed. Let me know and I'll help.

## Need Help?

1. Check `README.md` for full documentation
2. Look at logs: `docker-compose logs -f`
3. Test manually with curl (see README troubleshooting section)

## Files Included

```
otelcol-exaforce/
├── README.md              ← Full documentation
├── QUICKSTART.md          ← This file
├── build.sh               ← Run this to build
├── Dockerfile             ← Docker build instructions
├── docker-compose.yaml    ← Easy deployment
├── config.yaml            ← Collector config (edit the endpoint!)
├── .env.template          ← Put your token here
└── httpjsonexporter/      ← The custom exporter code
```

## You're All Set! 🎉

The hard part is done. Now you just need to:
1. Build it (`./build.sh`)
2. Configure it (edit `.env` and `config.yaml`)
3. Run it (`docker-compose up -d`)

That's it!
