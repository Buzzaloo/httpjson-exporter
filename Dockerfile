FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy the exporter code
COPY httpjsonexporter/ ./httpjsonexporter/

# Create builder config
COPY builder-config.yaml .

# Install the OpenTelemetry Collector Builder
RUN go install go.opentelemetry.io/collector/cmd/builder@v0.95.0

# Build the custom collector
RUN /go/bin/builder --config builder-config.yaml

# Final stage - create minimal runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

# Copy the built collector
COPY --from=builder /build/otelcol-custom/otelcol-custom /otelcol

# Create directories
RUN mkdir -p /etc/otelcol /var/log/otelcol

# Set the entrypoint
ENTRYPOINT ["/otelcol"]
CMD ["--config", "/etc/otelcol/config.yaml"]
