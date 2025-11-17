#!/bin/bash

# Quick verification script to test deployments
# Run this after deploying to verify everything works

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "🔍 apimock.codes Deployment Verification"
echo "========================================"
echo ""

# Get URLs
read -p "Enter your API URL (e.g., https://apimock-api.vercel.app): " API_URL
read -p "Enter your Website URL (e.g., https://apimock-web.vercel.app): " WEB_URL

echo ""
echo "Testing deployments..."
echo ""

# Test API
echo -e "${BLUE}Testing API...${NC}"
API_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/")
if [ "$API_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ API root endpoint: OK${NC}"
else
    echo -e "${RED}✗ API root endpoint: FAILED (Status: $API_STATUS)${NC}"
fi

API_USERS_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/api/users?count=1")
if [ "$API_USERS_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ API /api/users endpoint: OK${NC}"
else
    echo -e "${RED}✗ API /api/users endpoint: FAILED (Status: $API_USERS_STATUS)${NC}"
fi

# Test Website
echo ""
echo -e "${BLUE}Testing Website...${NC}"
WEB_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$WEB_URL/")
if [ "$WEB_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ Website home page: OK${NC}"
else
    echo -e "${RED}✗ Website home page: FAILED (Status: $WEB_STATUS)${NC}"
fi

DOCS_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$WEB_URL/docs")
if [ "$DOCS_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ Website /docs page: OK${NC}"
else
    echo -e "${RED}✗ Website /docs page: FAILED (Status: $DOCS_STATUS)${NC}"
fi

PLAYGROUND_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$WEB_URL/playground")
if [ "$PLAYGROUND_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ Website /playground page: OK${NC}"
else
    echo -e "${RED}✗ Website /playground page: FAILED (Status: $PLAYGROUND_STATUS)${NC}"
fi

echo ""
echo -e "${GREEN}======================================${NC}"
echo "Verification complete!"
echo ""
echo "Open in browser:"
echo "  API:        $API_URL"
echo "  Website:    $WEB_URL"
echo "  Docs:       $WEB_URL/docs"
echo "  Playground: $WEB_URL/playground"
echo ""
