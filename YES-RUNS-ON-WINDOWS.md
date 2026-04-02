# YES, This Runs on Windows! 🪟

## Quick Answer

**YES!** This will run on Windows exactly like the otel-contrib collector does. You just need to build a Windows `.exe` file instead of the Linux binary.

## The Easiest Way (If You Have Go Installed)

1. Install Go: https://go.dev/dl/
2. Run this:
```powershell
cd otelcol-exaforce
.\build-windows.ps1
```
3. You get: `otelcol-custom.exe`
4. Replace your existing collector with it
5. Done!

## If You Don't Want to Install Go

### Use GitHub Actions (Automated Builds)

1. Push this code to a GitHub repo
2. GitHub automatically builds the Windows .exe for you
3. Download it from the Actions tab
4. Deploy to Windows

**Workflow is already included** in `.github/workflows/build-windows.yaml`

## Or Just Use Docker as a Gateway

The **simplest option** if you don't want to change anything on Windows:

1. Run the Docker container on ANY Linux machine (could even be WSL)
2. Keep your existing Windows OTEL collector unchanged  
3. Just point it at the Docker container instead of directly to Exaforce

Your Windows config changes from:
```yaml
exporters:
  awss3:
    # S3 config
```

To:
```yaml
exporters:
  otlphttp:
    endpoint: "http://your-linux-host:4318"  # Point at Docker container
```

The Docker container handles the HTTP JSON export to Exaforce.

## Which Should You Choose?

| Method | Best For | Setup Time |
|--------|----------|-----------|
| **build-windows.ps1** | You have Go on Windows | 5 minutes |
| **GitHub Actions** | Multiple servers, want automation | 10 minutes |
| **Docker Gateway** | Don't want to touch Windows | 3 minutes |

## Full Documentation

- **WINDOWS-DEPLOYMENT.md** - Complete step-by-step for all options
- **QUICKSTART.md** - General getting started
- **README.md** - Full documentation

## Bottom Line

This is **not** just a Docker thing. You can absolutely run it directly on Windows, just like your current OTEL collector. The code I wrote works on any platform - you just need to compile it for Windows.

The `build-windows.ps1` script does that for you automatically.
