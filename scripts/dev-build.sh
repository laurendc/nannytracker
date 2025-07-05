#!/bin/bash

# Development build script for NannyTracker
# Quick build for local development and testing

set -e

echo "🚀 NannyTracker Development Build"
echo "=================================="

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: go.mod not found. Please run this script from the project root."
    exit 1
fi

# Build for current platform (fast)
echo "📦 Building TUI application..."
make build-dev

# Verify the build
if [ -f "./nannytracker" ]; then
    echo "✅ Build successful!"
    echo "📊 Binary info:"
    ls -lh ./nannytracker
    echo ""
    echo "🎯 Ready for development!"
    echo "   Run: ./nannytracker"
else
    echo "❌ Build failed!"
    exit 1
fi 