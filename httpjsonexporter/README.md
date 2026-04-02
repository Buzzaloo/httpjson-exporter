# HTTP JSON Exporter for OpenTelemetry Collector

A simple OpenTelemetry Collector exporter that sends log records as raw JSON to a generic HTTP endpoint with Bearer token authentication.

## Purpose

This exporter is designed for backends that:
- Accept JSON logs via HTTP POST
- Use Bearer token authentication
- Don't support OTLP protocol
- Need logs in a simple, flat JSON format (similar to what the `awss3` exporter with `marshaler: body` produces)

## Use Case

Perfect for sending Windows Event Logs (or any logs) from OpenTelemetry Collector to services like Exaforce that have custom HTTP ingest endpoints.

## Configuration

```yaml
exporters:
  httpjson:
    endpoint: "https://api.exaforce.io/ingest/logs"
    bearer_token: "${EXAFORCE_API_TOKEN}"
    timeout: 30s
    compression: gzip  # optional, default: none
    batch_size: 100    # optional, default: 100
    headers:           # optional additional headers
      X-Custom-Header: "value"
```

## Features

- **Simple JSON format**: Sends log records as newline-delimited JSON
- **Bearer token auth**: Built-in support for Bearer token authentication
- **Compression**: Optional gzip compression
- **Batching**: Configurable batch sizes
- **Retries**: Built-in retry logic with exponential backoff
- **TLS**: Full TLS support including custom CA certs

## Building

```bash
cd httpjsonexporter
go mod init github.com/ExaForce/httpjsonexporter
go mod tidy
go build
```

## Integration with OpenTelemetry Collector Contrib

To add this to your OpenTelemetry Collector Contrib build:

1. Add to `otelcontribcol/components.go`:
```go
import (
    httpjsonexporter "github.com/ExaForce/httpjsonexporter"
)

// In components() function:
exporters = append(exporters, httpjsonexporter.NewFactory())
```

2. Rebuild the collector:
```bash
make otelcontribcol
```

## Example Pipeline

```yaml
receivers:
  windowseventlog/system:
    channel: System
  windowseventlog/application:
    channel: Application  
  windowseventlog/security:
    channel: Security

processors:
  batch:
    send_batch_size: 100
    timeout: 10s

exporters:
  httpjson:
    endpoint: "https://your-tenant.app.exaforce.io/api/ingest/logs"
    bearer_token: "${EXAFORCE_API_TOKEN}"
    compression: gzip
    timeout: 30s

service:
  pipelines:
    logs:
      receivers: [windowseventlog/system, windowseventlog/application, windowseventlog/security]
      processors: [batch]
      exporters: [httpjson]
```

## Output Format

Each log record is sent as a JSON object on a single line (NDJSON format):

```json
{"channel":"Security","computer":"DC1.ad-buzzaloo.com","event_id":{"id":5152},"message":"The Windows Filtering Platform has blocked a packet.","system_time":"2026-02-05T04:29:18.802684800Z"}
{"channel":"Security","computer":"DC1.ad-buzzaloo.com","event_id":{"id":5156},"message":"The Windows Filtering Platform has permitted a connection.","system_time":"2026-02-05T04:29:19.123456700Z"}
```

Multiple records are sent in a single HTTP POST, separated by newlines.

## Comparison to Other Exporters

| Feature | httpjson | splunkhecexporter | otlphttpexporter | awss3exporter |
|---------|----------|-------------------|------------------|---------------|
| Simple JSON output | ✅ | ❌ (HEC format) | ❌ (OTLP format) | ✅ (with marshaler:body) |
| HTTP POST | ✅ | ✅ | ✅ | ❌ (S3 API) |
| Bearer token auth | ✅ | ✅ (via Token) | ❌ | ❌ (IAM) |
| No vendor lock-in | ✅ | ❌ (Splunk-specific) | ❌ (OTLP-specific) | ❌ (S3-specific) |

## License

Apache 2.0
