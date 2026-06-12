#!/bin/bash
# Quick test script for prom2sar

set -e

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║         prom2sar Quick Test Script                          ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Binary exists
echo -n "Test 1: Binary exists... "
if [ -f "./bin/prom2sar" ]; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC} - Run 'make build-cli' first"
    exit 1
fi

# Test 2: Version
echo -n "Test 2: Version check... "
VERSION=$(./bin/prom2sar --version 2>&1)
if echo "$VERSION" | grep -q "prom2sar version"; then
    echo -e "${GREEN}✓ PASS${NC} ($VERSION)"
else
    echo -e "${RED}✗ FAIL${NC}"
    exit 1
fi

# Test 3: Help
echo -n "Test 3: Help output... "
if ./bin/prom2sar --help 2>&1 | grep -q "Usage: prom2sar"; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
    exit 1
fi

# Test 4: Error handling - no tsdb
echo -n "Test 4: Error handling (no tsdb)... "
if ./bin/prom2sar 2>&1 | grep -q "tsdb is required"; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
    exit 1
fi

# Test 5: Error handling - invalid path
echo -n "Test 5: Invalid path handling... "
if ./bin/prom2sar -tsdb /nonexistent/path/12345 2>&1 | grep -q "does not exist"; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
    exit 1
fi

# Test 6: Invalid time format
echo -n "Test 6: Invalid time format... "
if ./bin/prom2sar -tsdb /tmp -start "not-a-date" 2>&1 | grep -q "Invalid start time"; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
    exit 1
fi

# Test 7: Invalid profile
echo -n "Test 7: Invalid profile... "
if ./bin/prom2sar -tsdb /tmp -profile invalid 2>&1 | grep -q "Invalid profile"; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
    exit 1
fi

echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║         ${GREEN}All Tests Passed!${NC}                                   ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "Next Steps:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "To test with real Prometheus data:"
echo "  ${YELLOW}./bin/prom2sar -tsdb /path/to/prometheus -output ./results -verbose${NC}"
echo ""
echo "To test with Docker (creates sample data):"
echo "  ${YELLOW}See TESTING.md for Docker test setup${NC}"
echo ""
echo "For detailed testing guide:"
echo "  ${YELLOW}cat TESTING.md${NC}"
echo ""
