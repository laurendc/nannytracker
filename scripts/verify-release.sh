#!/bin/bash

# Release Artifact Verification Script
# This script downloads and verifies release artifacts for NannyTracker

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
VERSION=${1:-"v1.0.0"}
REPO="laurendc/nannytracker"
TEMP_DIR="./release-verify-${VERSION}"
PLATFORMS=(
    "linux-amd64"
    "linux-arm64" 
    "darwin-amd64"
    "darwin-arm64"
    "windows-amd64"
)

# Expected binaries
TUI_BINARIES=(
    "nannytracker-linux-amd64"
    "nannytracker-linux-arm64"
    "nannytracker-darwin-amd64"
    "nannytracker-darwin-arm64"
    "nannytracker-windows-amd64.exe"
)

WEB_BINARIES=(
    "nannytracker-web-linux-amd64"
    "nannytracker-web-linux-arm64"
    "nannytracker-web-darwin-amd64"
    "nannytracker-web-darwin-arm64"
    "nannytracker-web-windows-amd64.exe"
)

echo -e "${BLUE}🔍 NannyTracker Release Verification${NC}"
echo -e "${BLUE}Version: ${VERSION}${NC}"
echo -e "${BLUE}Repository: ${REPO}${NC}"
echo ""

# Create temporary directory
mkdir -p "$TEMP_DIR"
cd "$TEMP_DIR"

echo -e "${YELLOW}📥 Downloading release artifacts...${NC}"

# Function to download a file with retries
download_file() {
    local file=$1
    local url="https://github.com/${REPO}/releases/download/${VERSION}/${file}"
    local max_retries=3
    local retry_count=0
    
    while [ $retry_count -lt $max_retries ]; do
        echo -n "  Downloading ${file} (attempt $((retry_count + 1))/${max_retries})... "
        if curl -L -s -o "$file" "$url"; then
            echo -e "${GREEN}✓${NC}"
            return 0
        else
            echo -e "${YELLOW}✗ (retrying in 5s)${NC}"
            retry_count=$((retry_count + 1))
            if [ $retry_count -lt $max_retries ]; then
                sleep 5
            fi
        fi
    done
    
    echo -e "${RED}✗ (failed after ${max_retries} attempts)${NC}"
    return 1
}

# Function to verify binary
verify_binary() {
    local binary=$1
    local platform=$2
    
    echo -n "  Verifying ${binary}... "
    
    # Check if file exists and is not empty
    if [[ ! -f "$binary" ]] || [[ ! -s "$binary" ]]; then
        echo -e "${RED}✗ (file missing or empty)${NC}"
        return 1
    fi
    
    # Make executable (for Unix-like systems)
    if [[ "$binary" != *.exe ]]; then
        chmod +x "$binary"
    fi
    
    # Test version output
    if [[ "$binary" == *"nannytracker"* ]]; then
        if ./"$binary" --version 2>/dev/null | grep -q "$VERSION"; then
            echo -e "${GREEN}✓${NC}"
            return 0
        else
            echo -e "${RED}✗ (version mismatch)${NC}"
            return 1
        fi
    else
        # For web binaries, just check they're executable
        echo -e "${GREEN}✓${NC}"
        return 0
    fi
}

# Download and verify TUI binaries
echo -e "${YELLOW}📦 TUI Application Binaries:${NC}"
tui_success=0
tui_total=${#TUI_BINARIES[@]}

for binary in "${TUI_BINARIES[@]}"; do
    if download_file "$binary"; then
        if verify_binary "$binary"; then
            ((tui_success++))
        fi
    fi
done

echo ""

# Download and verify Web binaries
echo -e "${YELLOW}🌐 Web Server Binaries:${NC}"
web_success=0
web_total=${#WEB_BINARIES[@]}

for binary in "${WEB_BINARIES[@]}"; do
    if download_file "$binary"; then
        if verify_binary "$binary"; then
            ((web_success++))
        fi
    fi
done

echo ""

# Summary
echo -e "${BLUE}📊 Verification Summary:${NC}"
echo -e "  TUI Binaries: ${tui_success}/${tui_total} ✓"
echo -e "  Web Binaries: ${web_success}/${web_total} ✓"
echo -e "  Total: $((tui_success + web_success))/$((tui_total + web_total)) ✓"

# Overall result
if [[ $((tui_success + web_success)) -eq $((tui_total + web_total)) ]]; then
    echo -e "${GREEN}🎉 All artifacts verified successfully!${NC}"
    exit_code=0
else
    echo -e "${RED}❌ Some artifacts failed verification.${NC}"
    exit_code=1
fi

echo ""
echo -e "${YELLOW}🧹 Cleaning up...${NC}"
cd ..
rm -rf "$TEMP_DIR"

echo -e "${GREEN}✅ Verification complete!${NC}"
exit $exit_code 