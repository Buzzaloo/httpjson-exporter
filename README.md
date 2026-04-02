# Custom OpenTelemetry Collector for Exaforce

This package contains a **ready-to-build** custom OpenTelemetry Collector with an HTTP JSON exporter specifically designed to send Windows Event Logs to Exaforce's HTTP endpoint.

## What This Solves

The standard OpenTelemetry Collector doesn't have a simple HTTP exporter that:
- Sends raw JSON (not OTLP format)
- Supports Bearer token authentication
- Works like the S3 exporter's `marshaler: body` option

This custom build includes the `httpjson` exporter that does exactly that!

## Quick Start (3 Steps)

### Step 1: Build the Custom Collector

```bash
# Make sure Docker is running
docker --version

# Build the custom collector image
./build.sh
```

This will create a Docker image called `otelcol-exaforce:latest` with your custom exporter included.

### Step 2: Configure Your Exaforce Endpoint

1. Copy the environment template:
```bash
cp .env.template .env
```

2. Edit `.env` and add your Exaforce API token:
```bash
EXAFORCE_API_TOKEN=your-actual-token-here
```

3. Edit `config.yaml` and update the endpoint:
```yaml
exporters:
  httpjson:
    endpoint: "https://YOUR-TENANT.us.app.exaforce.io/api/v1/logs/ingest"
```

### Step 3: Run It!

```bash
# Start the collector
docker-compose up -d

# Check logs
docker-compose logs -f

# Check health
curl http://localhost:13133
```

## For Windows Deployment

If you're running this on Windows servers (which you probably are for Windows Event Logs):

### Option A: Use This as a Gateway

Keep your existing OTEL collector on Windows, but point it at this container:

**On your Windows server**, use your current config but change the exporter:

```yaml
exporters:
  otlphttp:
    endpoint: "http://your-linux-host:4318"
```

Then this container forwards to Exaforce.

### Option B: Build Windows Binary

If you want the collector running directly on Windows:

1. On a machine with Go installed:
```bash
cd httpjsonexporter
go mod download
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o otelcol-windows.exe
```

2. Or extract the binary from the Docker image:
```bash
docker create --name temp otelcol-exaforce:latest
docker cp temp:/otelcol ./otelcol-windows.exe
docker rm temp
```

3. Copy to Windows and use your existing OTEL setup, just swap the exporter config.

## Configuration Reference

### HTTP JSON Exporter Options

```yaml
exporters:
  httpjson:
    # Required: Your Exaforce endpoint
    endpoint: "https://tenant.app.exaforce.io/api/v1/logs/ingest"
    
    # Required: Bearer token for authentication
    bearer_token: "${EXAFORCE_API_TOKEN}"
    
    # Optional: Compression (none or gzip)
    compression: gzip
    
    # Optional: Request timeout
    timeout: 30s
    
    # Optional: Logs per HTTP request
    batch_size: 100
    
    # Optional: Additional headers
    headers:
      X-Custom: "value"
    
    # Optional: Retry settings
    retry_on_failure:
      enabled: true
      initial_interval: 5s
      max_interval: 30s
      max_elapsed_time: 300s
    
    # Optional: Queue settings
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 1000
```

## Output Format

Logs are sent as newline-delimited JSON (NDJSON):

```json
{"channel":"Security","computer":"DC1","event_id":{"id":5152},"message":"The Windows Filtering Platform has blocked a packet.","timestamp":"2026-02-05T04:29:18.802684800Z"}
{"channel":"Security","computer":"DC1","event_id":{"id":5156},"message":"The Windows Filtering Platform has permitted a connection.","timestamp":"2026-02-05T04:29:19.123456700Z"}
```

## Troubleshooting

### Check if the collector is running:
```bash
docker ps | grep otelcol
```

### View logs:
```bash
docker-compose logs -f otelcol
```

### Test the health endpoint:
```bash
curl http://localhost:13133
```

### Test sending to Exaforce manually:
```bash
curl -X POST https://your-tenant.app.exaforce.io/api/v1/logs/ingest \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/x-ndjson" \
  -d '{"test":"log","timestamp":"2026-02-05T12:00:00Z"}'
```

### Common Issues

**"endpoint must be specified"**
- Make sure you updated the endpoint in `config.yaml`

**"bearer_token must be specified"**
- Make sure you created `.env` with your token
- Check: `docker-compose config` to verify env vars are loaded

**"connection refused"**
- Check if Exaforce endpoint is reachable
- Verify firewall rules

**"401 Unauthorized"**
- Check your bearer token is correct
- Verify the token has permissions to ingest logs

## What's Included

```
otelcol-exaforce/
├── Dockerfile                 # Builds the custom collector
├── builder-config.yaml        # Tells the builder what to include
├── docker-compose.yaml        # Easy deployment
├── config.yaml                # Collector configuration
├── .env.template              # Environment variables template
├── build.sh                   # One-command build script
├── httpjsonexporter/          # The custom exporter code
│   ├── config.go
│   ├── factory.go
│   ├── exporter.go
│   ├── go.mod
│   └── README.md
└── README.md                  # This file
```

## Architecture

```
Windows Server
    ↓
OpenTelemetry Collector (Windows Event Logs)
    ↓ (OTLP or direct)
Custom OTel Collector (with httpjson exporter)
    ↓ (HTTP POST with Bearer token)
Exaforce HTTP Endpoint
```

## Next Steps

1. **Build** the collector: `./build.sh`
2. **Configure** your endpoint and token
3. **Deploy** with `docker-compose up -d`
4. **Verify** logs are reaching Exaforce
5. **Scale** by adjusting batch_size and queue settings

## Getting Help

- Check the logs: `docker-compose logs -f`
- Verify config: `docker-compose config`
- Test manually with curl (see Troubleshooting section)

## License

Apache 2.0 (same as OpenTelemetry)
