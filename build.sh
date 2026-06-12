#!/bin/bash
# Build script for prom2sar CLI

set -e

VERSION="1.0.0"
BUILD_DIR="./bin"
DIST_DIR="./dist"

echo "=== Building prom2sar v${VERSION} ==="

# Create directories
mkdir -p "$BUILD_DIR"
mkdir -p "$DIST_DIR"

# Build for current platform
echo "Building for current platform..."
go build -o "${BUILD_DIR}/prom2sar" cmd/prom2sar/main.go
echo "✓ Built: ${BUILD_DIR}/prom2sar"

# Build for common platforms
echo ""
echo "Building release binaries..."

# Linux AMD64
echo "- Linux AMD64..."
GOOS=linux GOARCH=amd64 go build \
  -ldflags "-X main.version=${VERSION}" \
  -o "${DIST_DIR}/prom2sar-linux-amd64" \
  cmd/prom2sar/main.go

# Linux AMD64 (static)
echo "- Linux AMD64 (static)..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags "-X main.version=${VERSION} -extldflags '-static'" \
  -o "${DIST_DIR}/prom2sar-linux-amd64-static" \
  cmd/prom2sar/main.go

# Linux ARM64
echo "- Linux ARM64..."
GOOS=linux GOARCH=arm64 go build \
  -ldflags "-X main.version=${VERSION}" \
  -o "${DIST_DIR}/prom2sar-linux-arm64" \
  cmd/prom2sar/main.go

# Create checksums
echo ""
echo "Creating checksums..."
cd "$DIST_DIR"
sha256sum prom2sar-* > checksums.txt
cd - > /dev/null

# Display results
echo ""
echo "=== Build Complete ==="
echo ""
echo "Binaries:"
ls -lh "${DIST_DIR}"/prom2sar-*
echo ""
echo "Checksums:"
cat "${DIST_DIR}/checksums.txt"
echo ""
echo "Test the binary:"
echo "  ${BUILD_DIR}/prom2sar --version"
echo ""
echo "Install system-wide:"
echo "  sudo cp ${BUILD_DIR}/prom2sar /usr/local/bin/"
echo ""
