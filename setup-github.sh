#!/bin/bash
set -e

echo "=========================================="
echo "Setting up GitHub Repository"
echo "=========================================="
echo ""

cd /mnt/user-data/outputs/otelcol-exaforce

# Initialize git repo if not already initialized
if [ ! -d .git ]; then
    echo "Initializing Git repository..."
    git init
    git config user.name "ExaForce"
    git config user.email "dev@exaforce.com"
fi

# Add all files
echo "Adding files..."
git add .

# Create initial commit
echo "Creating initial commit..."
git commit -m "Initial commit: Custom OTEL Collector with HTTP JSON Exporter

- Custom HTTP JSON exporter for sending logs to Exaforce endpoint
- Bearer token authentication support
- Windows build support with build-windows.ps1
- Docker build for Linux deployment
- GitHub Actions workflow for automated Windows builds
- Complete documentation and examples"

echo ""
echo "=========================================="
echo "Repository initialized!"
echo "=========================================="
echo ""
echo "Next steps:"
echo ""
echo "1. Go to: https://github.com/orgs/ExaForce/repositories/new"
echo ""
echo "2. Name it: otelcol-http-exporter"
echo "   (or whatever you prefer)"
echo ""
echo "3. Make it Public or Private"
echo ""
echo "4. DO NOT initialize with README (we already have one)"
echo ""
echo "5. Click 'Create repository'"
echo ""
echo "6. Then run these commands:"
echo ""
echo "   cd /mnt/user-data/outputs/otelcol-exaforce"
echo "   git remote add origin https://github.com/ExaForce/YOUR-REPO-NAME.git"
echo "   git branch -M main"
echo "   git push -u origin main"
echo ""
echo "Or if you want me to do it, just tell me the repo name!"
echo ""
