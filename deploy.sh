#!/bin/bash

# Vercel Deployment Script for apimock.codes
# This script helps you deploy both API and Website to Vercel

set -e

echo "🚀 apimock.codes Vercel Deployment"
echo "===================================="
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if logged into Vercel
echo -e "${BLUE}Checking Vercel authentication...${NC}"
if ! vercel whoami > /dev/null 2>&1; then
    echo -e "${YELLOW}Not logged into Vercel. Please login:${NC}"
    vercel login
fi

echo -e "${GREEN}✓ Authenticated with Vercel${NC}"
echo ""

# Get the root directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
API_DIR="$SCRIPT_DIR/api"
WEB_DIR="$SCRIPT_DIR/web"

echo "Project directories:"
echo "  API: $API_DIR"
echo "  Web: $WEB_DIR"
echo ""

# Deploy API first
echo -e "${BLUE}Step 1: Deploying API${NC}"
echo "================================"
cd "$API_DIR"

echo "Creating Vercel project for API..."
vercel link --yes || true

echo ""
echo -e "${YELLOW}Deploying API to production...${NC}"
API_URL=$(vercel --prod 2>&1 | grep -Eo 'https://[^ ]+' | tail -1)

if [ -z "$API_URL" ]; then
    echo -e "${YELLOW}⚠️  Could not extract API URL automatically${NC}"
    echo "Please check the Vercel dashboard for your API URL"
    read -p "Enter your API production URL: " API_URL
fi

echo -e "${GREEN}✓ API deployed successfully!${NC}"
echo "  URL: $API_URL"
echo ""

# Deploy Website
echo -e "${BLUE}Step 2: Deploying Website${NC}"
echo "================================"
cd "$WEB_DIR"

echo "Creating Vercel project for Website..."
vercel link --yes || true

# Set environment variable for website
echo ""
echo "Setting environment variable for website..."
echo "  NEXT_PUBLIC_API_URL=$API_URL"

# Add environment variable
vercel env add NEXT_PUBLIC_API_URL production <<EOF
$API_URL
EOF

echo ""
echo -e "${YELLOW}Deploying Website to production...${NC}"
WEB_URL=$(vercel --prod 2>&1 | grep -Eo 'https://[^ ]+' | tail -1)

if [ -z "$WEB_URL" ]; then
    echo -e "${YELLOW}⚠️  Could not extract Web URL automatically${NC}"
    echo "Please check the Vercel dashboard for your website URL"
    read -p "Enter your Website production URL: " WEB_URL
fi

echo -e "${GREEN}✓ Website deployed successfully!${NC}"
echo "  URL: $WEB_URL"
echo ""

# Summary
echo ""
echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}🎉 Deployment Complete!${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""
echo "Your services are now live:"
echo ""
echo "  📡 API:        $API_URL"
echo "  🌐 Website:    $WEB_URL"
echo "  📚 Docs:       $WEB_URL/docs"
echo "  🎮 Playground: $WEB_URL/playground"
echo ""
echo "Next steps:"
echo "  1. Test your API:  curl $API_URL/"
echo "  2. Visit website:  open $WEB_URL"
echo "  3. Add custom domains in Vercel dashboard"
echo "     - api.apimock.codes → API project"
echo "     - apimock.codes → Website project"
echo ""
echo "For custom domain setup, see: VERCEL_SETUP.md"
echo ""
