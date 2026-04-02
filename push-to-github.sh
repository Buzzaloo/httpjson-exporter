#!/bin/bash

echo "=========================================="
echo "Pushing to GitHub"
echo "=========================================="
echo ""
echo "Repository: https://github.com/Buzzaloo/httpjson-exporter"
echo ""

cd /mnt/user-data/outputs/otelcol-exaforce

echo "Files ready to push:"
git log --oneline -1
echo ""
git diff --stat HEAD~1 HEAD
echo ""

echo "Run this command to push:"
echo ""
echo "  cd /mnt/user-data/outputs/otelcol-exaforce"
echo "  git push -u origin main"
echo ""
echo "You'll be prompted for your GitHub credentials."
echo "Use a Personal Access Token as the password if prompted."
echo ""
