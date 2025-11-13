#!/bin/bash

# Pre-commit hook for running CI checks locally
# This ensures code quality before committing

set -e

echo "🔍 Running pre-commit checks..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi

# Format check
echo -e "${YELLOW}📝 Checking formatting...${NC}"
if [ -n "$(gofmt -l .)" ]; then
    echo -e "${RED}❌ Code is not formatted. Run 'make fmt' to fix.${NC}"
    gofmt -l .
    exit 1
fi
echo -e "${GREEN}✅ Formatting check passed${NC}"

# Go vet
echo -e "${YELLOW}🔎 Running go vet...${NC}"
if ! go vet ./...; then
    echo -e "${RED}❌ go vet failed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ go vet passed${NC}"

# Linting
echo -e "${YELLOW}🔍 Running golangci-lint...${NC}"
if command -v golangci-lint &> /dev/null; then
    if ! golangci-lint run ./...; then
        echo -e "${RED}❌ Linting failed${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Linting passed${NC}"
else
    echo -e "${YELLOW}⚠️  golangci-lint not installed, skipping. Run 'make install-tools' to install.${NC}"
fi

# Run tests
echo -e "${YELLOW}🧪 Running tests...${NC}"
if ! go test -race ./...; then
    echo -e "${RED}❌ Tests failed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Tests passed${NC}"

echo -e "${GREEN}🎉 All pre-commit checks passed!${NC}"
exit 0
