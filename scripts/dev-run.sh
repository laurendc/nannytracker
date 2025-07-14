#!/bin/bash

# Development run script for NannyTracker Web Application
# Runs both backend and frontend servers simultaneously

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}🚀 NannyTracker Web Development Server${NC}"
echo "========================================"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Error: go.mod not found. Please run this script from the project root.${NC}"
    exit 1
fi

# Build the application first
echo -e "${BLUE}📦 Building application...${NC}"
./scripts/dev-build.sh --web

echo ""
echo -e "${GREEN}🎯 Starting development servers...${NC}"
echo -e "${YELLOW}📝 Note: Use Ctrl+C to stop both servers${NC}"
echo ""

# Function to cleanup background processes
cleanup() {
    echo -e "\n${YELLOW}🛑 Stopping servers...${NC}"
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    exit 0
}

# Set up trap to cleanup on exit
trap cleanup SIGINT SIGTERM

# Start backend server in background
echo -e "${BLUE}🔧 Starting backend server...${NC}"
./nannytracker-web &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 2

# Check if backend is running
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo -e "${RED}❌ Backend server failed to start!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Backend server started (PID: $BACKEND_PID)${NC}"
echo -e "${BLUE}   📡 API: http://localhost:8080${NC}"

# Start frontend server in background
echo -e "${BLUE}🎨 Starting frontend server...${NC}"
cd web
npm run dev &
FRONTEND_PID=$!
cd ..

# Wait a moment for frontend to start
sleep 3

# Check if frontend is running
if ! kill -0 $FRONTEND_PID 2>/dev/null; then
    echo -e "${RED}❌ Frontend server failed to start!${NC}"
    kill $BACKEND_PID 2>/dev/null || true
    exit 1
fi

echo -e "${GREEN}✅ Frontend server started (PID: $FRONTEND_PID)${NC}"
echo -e "${BLUE}   🌐 Web App: http://localhost:3000${NC}"

echo ""
echo -e "${GREEN}🎉 Both servers are running!${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}📱 Mobile Testing Instructions:${NC}"
echo -e "${YELLOW}1. Open your browser to:${NC} http://localhost:3000"
echo -e "${YELLOW}2. Open Chrome DevTools:${NC} F12 or right-click → Inspect"
echo -e "${YELLOW}3. Enable Device Toolbar:${NC} Click 📱 icon or press Ctrl+Shift+M"
echo -e "${YELLOW}4. Test different devices:${NC}"
echo "   • iPhone SE (375px) - Mobile bottom navigation"
echo "   • iPhone 12 Pro (390px) - Mobile optimized"
echo "   • iPad (768px) - Tablet breakpoint"
echo "   • Desktop (1024px+) - Desktop sidebar"
echo -e "${YELLOW}5. Test features:${NC}"
echo "   • Add/edit/delete trips and expenses"
echo "   • Test mobile forms and touch targets"
echo "   • Verify responsive navigation"
echo "   • Test table-to-card transformation"
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}🔧 Development Tips:${NC}"
echo -e "${YELLOW}• Backend logs will show API requests${NC}"
echo -e "${YELLOW}• Frontend has hot reload - changes appear instantly${NC}"
echo -e "${YELLOW}• Use browser DevTools Console for debugging${NC}"
echo -e "${YELLOW}• Test touch interactions on mobile breakpoints${NC}"
echo ""
echo -e "${RED}Press Ctrl+C to stop both servers${NC}"

# Wait for user to stop the servers
while true; do
    sleep 1
    # Check if either server has stopped
    if ! kill -0 $BACKEND_PID 2>/dev/null; then
        echo -e "${RED}❌ Backend server stopped unexpectedly!${NC}"
        break
    fi
    if ! kill -0 $FRONTEND_PID 2>/dev/null; then
        echo -e "${RED}❌ Frontend server stopped unexpectedly!${NC}"
        break
    fi
done

cleanup 