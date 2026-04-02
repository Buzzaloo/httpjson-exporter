# Integration Guide: HTTP JSON Exporter with OpenTelemetry Collector Contrib

This guide shows how to build a custom OpenTelemetry Collector with the HTTP JSON exporter included.

## Option 1: Use OpenTelemetry Collector Builder (Recommended)

The easiest way to create a custom collector with your exporter is using the OpenTelemetry Collector Builder.

### Step 1: Install the Builder

```bash
go install go.opentelemetry.io/collector/cmd/builder@latest
```

### Step 2: Create a Builder Config

Create a file called `builder-config.yaml`:

```yaml
dist:
  name: otelcol-custom
  description: Custom OpenTelemetry Collector with HTTP JSON Exporter
  output_path: ./otelcol-custom
  otelcol_version: 0.95.0

exporters:
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awss3exporter v0.95.0
  - gomod: github.com/ExaForce/httpjsonexporter v0.1.0
    path: ./httpjsonexporter

receivers:
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowseventlogreceiver v0.95.0

processors:
  - gomod: go.opentelemetry.io/collector/processor/batchprocessor v0.95.0
  - gomod: go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.95.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.95.0

extensions:
  - gomod: go.opentelemetry.io/collector/extension/zpagesextension v0.95.0
```

### Step 3: Build the Collector

```bash
builder --config builder-config.yaml
```

This will create a custom collector binary at `./otelcol-custom/otelcol-custom`.

### Step 4: Test the Collector

```bash
cd otelcol-custom
./otelcol-custom --config ../example-config.yaml
```

## Option 2: Fork opentelemetry-collector-contrib

If you want to contribute this back to the community or need more control:

### Step 1: Clone the Repo

```bash
git clone https://github.com/open-telemetry/opentelemetry-collector-contrib.git
cd opentelemetry-collector-contrib
```

### Step 2: Add the Exporter

```bash
mkdir -p exporter/httpjsonexporter
cp -r /path/to/httpjsonexporter/* exporter/httpjsonexporter/
```

### Step 3: Register the Exporter

Edit `cmd/otelcontribcol/components.go` and add:

```go
import (
    // ... other imports ...
    httpjsonexporter "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/httpjsonexporter"
)

func components() (otelcol.Factories, error) {
    // ... existing code ...
    
    exporters, err := exporter.MakeFactoryMap(
        // ... existing exporters ...
        httpjsonexporter.NewFactory(),
    )
    
    // ... rest of the function ...
}
```

### Step 4: Update go.mod

```bash
go mod tidy
```

### Step 5: Build

```bash
make otelcontribcol
```

The binary will be at `./bin/otelcontribcol_linux_amd64` (or your platform).

## Option 3: Simple Standalone Build (Quickest for Testing)

If you just want to test quickly:

### Step 1: Create a Simple Main

Create `main.go`:

```go
package main

import (
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/otelcol"
	
	httpjsonexporter "github.com/ExaForce/httpjsonexporter"
	windowseventlogreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowseventlogreceiver"
)

func main() {
	factories, err := components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
	}

	info := component.BuildInfo{
		Command:     "otelcol",
		Description: "Custom OpenTelemetry Collector",
		Version:     "0.1.0",
	}

	app := otelcol.NewCommand(otelcol.CollectorSettings{
		BuildInfo: info,
		Factories: factories,
	})

	if err := app.Execute(); err != nil {
		log.Fatalf("collector server run finished with error: %v", err)
	}
}

func components() (otelcol.Factories, error) {
	var err error
	factories := otelcol.Factories{}

	factories.Receivers, err = receiver.MakeFactoryMap(
		windowseventlogreceiver.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Exporters, err = exporter.MakeFactoryMap(
		httpjsonexporter.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Processors, err = processor.MakeFactoryMap(
		batchprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	return factories, nil
}
```

### Step 2: Build

```bash
go build -o otelcol-custom
```

### Step 3: Run

```bash
export EXAFORCE_API_TOKEN="your-token-here"
./otelcol-custom --config config.yaml
```

## Deployment to Windows

Once you've built the custom collector:

1. Copy the binary to your Windows server
2. Update your collector config to use `httpjson` exporter instead of `awss3`
3. Set the `EXAFORCE_API_TOKEN` environment variable
4. Restart the OpenTelemetry Collector service

## Environment Variables

```powershell
# Set the Exaforce API token
[System.Environment]::SetEnvironmentVariable("EXAFORCE_API_TOKEN", "your-token-here", "Machine")

# Restart the service
Restart-Service OpenTelemetryCollector
```

## Troubleshooting

### Check if logs are being sent

Look at the collector logs:
```
C:\ProgramData\OpenTelemetry Collector\otelcol.log
```

### Enable debug logging

In your config.yaml:
```yaml
service:
  telemetry:
    logs:
      level: debug
```

### Test the HTTP endpoint manually

```powershell
$headers = @{
    "Authorization" = "Bearer your-token-here"
    "Content-Type" = "application/x-ndjson"
}

$body = '{"test":"log","timestamp":"2026-02-05T12:00:00Z"}'

Invoke-WebRequest -Uri "https://your-tenant.app.exaforce.io/api/ingest/logs" `
    -Method Post `
    -Headers $headers `
    -Body $body
```

## Next Steps

1. Test with a small volume of logs first
2. Monitor Exaforce console to verify logs are appearing
3. Adjust batch_size and timeout based on your log volume
4. Consider enabling compression for high-volume environments
