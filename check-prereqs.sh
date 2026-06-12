#!/bin/bash
# Check prerequisites for building prom2sar

echo "════════════════════════════════════════════════════════════════"
echo "  Checking Prerequisites for prom2sar"
echo "════════════════════════════════════════════════════════════════"
echo ""

ALL_GOOD=true

# Check Go
echo -n "Go (1.21+):        "
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    echo "✓ Installed ($GO_VERSION)"

    # Check if version is 1.21+
    MAJOR=$(echo $GO_VERSION | cut -d. -f1)
    MINOR=$(echo $GO_VERSION | cut -d. -f2)
    if [ "$MAJOR" -ge "1" ] && [ "$MINOR" -ge "21" ]; then
        echo "               ✓ Version OK (>= 1.21)"
    else
        echo "               ✗ Version too old (need >= 1.21)"
        ALL_GOOD=false
    fi
else
    echo "✗ Not found"
    echo "               Install from: https://go.dev/dl/"
    ALL_GOOD=false
fi

# Check Make
echo -n "Make:              "
if command -v make &> /dev/null; then
    MAKE_VERSION=$(make --version | head -1)
    echo "✓ Installed"
else
    echo "✗ Not found"
    echo "               Install: sudo apt-get install build-essential (Ubuntu)"
    echo "               or: sudo dnf install make (RHEL/Fedora)"
    ALL_GOOD=false
fi

# Check Git
echo -n "Git:               "
if command -v git &> /dev/null; then
    GIT_VERSION=$(git --version | awk '{print $3}')
    echo "✓ Installed ($GIT_VERSION)"
else
    echo "✗ Not found"
    echo "               Install: sudo apt-get install git"
    ALL_GOOD=false
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  Optional Tools"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Check Docker (optional)
echo -n "Docker:            "
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    echo "✓ Installed ($DOCKER_VERSION)"
else
    echo "○ Not found (optional - only for building container images)"
fi

# Check kubectl (optional)
echo -n "kubectl:           "
if command -v kubectl &> /dev/null; then
    echo "✓ Installed"
else
    echo "○ Not found (optional - only for Kubernetes operator)"
fi

# Check oc (optional)
echo -n "oc (OpenShift):    "
if command -v oc &> /dev/null; then
    echo "✓ Installed"
else
    echo "○ Not found (optional - only for OpenShift operator)"
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  Go Dependencies"
echo "════════════════════════════════════════════════════════════════"
echo ""

if [ -f "go.mod" ]; then
    echo "✓ go.mod found"

    # Check if go.sum exists
    if [ -f "go.sum" ]; then
        echo "✓ go.sum found (dependencies already downloaded)"
    else
        echo "○ go.sum not found (dependencies not yet downloaded)"
        echo "  Run 'make deps' to download dependencies"
    fi
else
    echo "✗ go.mod not found - are you in the project directory?"
    ALL_GOOD=false
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  Summary"
echo "════════════════════════════════════════════════════════════════"
echo ""

if [ "$ALL_GOOD" = true ]; then
    echo "✓ All required prerequisites are installed!"
    echo ""
    echo "Next steps:"
    echo "  1. Download dependencies: make deps"
    echo "  2. Build CLI:            make build-cli"
    echo "  3. Install:              sudo make install-cli"
    echo "  4. Verify:               prom2sar --version"
    echo ""
else
    echo "✗ Some required prerequisites are missing."
    echo ""
    echo "Please install missing items and run this script again."
    echo ""
    echo "See PREREQUISITES.md for detailed installation instructions:"
    echo "  cat PREREQUISITES.md"
    echo ""
fi

echo "════════════════════════════════════════════════════════════════"
