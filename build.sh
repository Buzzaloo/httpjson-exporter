#!/bin/bash
set -e

echo "=========================================="
echo "Building Custom OpenTelemetry Collector"
echo "with HTTP JSON Exporter for Exaforce"
echo "=========================================="
echo ""

# Build the Docker image
echo "Building Docker image..."
docker build -t otelcol-exaforce:latest .

echo ""
echo "=========================================="
echo "Build Complete!"
echo "=========================================="
echo ""
echo "Docker image: otelcol-exaforce:latest"
echo ""
echo "Next steps:"
echo "1. Update config.yaml with your Exaforce endpoint and token"
echo "2. Run: docker-compose up -d"
echo ""
