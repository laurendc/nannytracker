#!/bin/bash

# Development build script for NannyTracker
# Quick build for local development and testing

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Parse command line arguments
MODE="web"
SKIP_DEPS=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --tui)
            MODE="tui"
            shift
            ;;
        --web)
            MODE="web"
            shift
            ;;
        --skip-deps)
            SKIP_DEPS=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  --web        Build and run web application (default)"
            echo "  --tui        Build and run TUI application"
            echo "  --skip-deps  Skip npm install for web mode"
            echo "  -h, --help   Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option $1"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🚀 NannyTracker Development Build${NC}"
echo "=================================="

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Error: go.mod not found. Please run this script from the project root.${NC}"
    exit 1
fi

if [ "$MODE" = "tui" ]; then
    # Build TUI application
    echo -e "${BLUE}📦 Building TUI application...${NC}"
    make build-dev

    # Verify the build
    if [ -f "./nannytracker" ]; then
        echo -e "${GREEN}✅ Build successful!${NC}"
        echo -e "${BLUE}📊 Binary info:${NC}"
        ls -lh ./nannytracker
        echo ""
        echo -e "${GREEN}🎯 Ready for development!${NC}"
        echo "   Run: ./nannytracker"
    else
        echo -e "${RED}❌ Build failed!${NC}"
        exit 1
    fi
else
    # Web application mode
    echo -e "${BLUE}🌐 Setting up web application development environment...${NC}"
    
    # Build backend server
    echo -e "${BLUE}📦 Building backend server...${NC}"
    go build -ldflags="-X github.com/laurendc/nannytracker/pkg/version.Version=dev" -o nannytracker-web ./cmd/web
    
    if [ ! -f "./nannytracker-web" ]; then
        echo -e "${RED}❌ Backend build failed!${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Backend build successful!${NC}"
    
    # Check if web directory exists
    if [ ! -d "web" ]; then
        echo -e "${RED}❌ Error: web directory not found.${NC}"
        exit 1
    fi
    
    # Install frontend dependencies if needed
    if [ "$SKIP_DEPS" = false ]; then
        echo -e "${BLUE}📦 Installing frontend dependencies...${NC}"
        cd web
        if [ ! -d "node_modules" ] || [ ! -f "package-lock.json" ]; then
            npm install
        else
            echo -e "${YELLOW}⚡ Dependencies already installed, skipping...${NC}"
        fi
        cd ..
    fi
    
    echo -e "${GREEN}✅ Web application setup complete!${NC}"
    echo ""
    echo -e "${GREEN}🎯 Ready for development!${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}🚀 To start the development servers:${NC}"
    echo ""
    echo -e "${BLUE}1. Start Backend Server (Terminal 1):${NC}"
    echo "   ./nannytracker-web"
    echo "   Server will run on: http://localhost:8080"
    echo ""
    echo -e "${BLUE}2. Start Frontend Server (Terminal 2):${NC}"
    echo "   cd web && npm run dev"
    echo "   Frontend will run on: http://localhost:3000"
    echo ""
    echo -e "${YELLOW}📱 To test mobile responsiveness:${NC}"
    echo "   • Open Chrome DevTools (F12)"
    echo "   • Click the device toolbar icon (📱) or press Ctrl+Shift+M"
    echo "   • Test different device sizes (iPhone, iPad, etc.)"
    echo "   • Test bottom navigation on mobile breakpoints"
    echo ""
    echo -e "${YELLOW}🔧 Quick Start (runs both servers):${NC}"
    echo "   ./scripts/dev-run.sh"
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
fi 