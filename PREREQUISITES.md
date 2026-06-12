# Prerequisites

This document lists all prerequisites needed to build and run the Prometheus to SAR Converter.

## System Requirements

### Operating System
- **Linux** (recommended)
- **macOS** (for development)
- **Windows WSL2** (for development)

### Required Software

#### 1. Go Programming Language

**Version:** 1.21 or later

**Check if installed:**
```bash
go version
```

**Install Go:**

**On Ubuntu/Debian:**
```bash
# Remove old Go (if any)
sudo rm -rf /usr/local/go

# Download and install Go 1.21
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

**On RHEL/Fedora:**
```bash
sudo dnf install golang
# or
sudo yum install golang

# Verify
go version
```

**On macOS:**
```bash
brew install go@1.21
```

#### 2. Make

**Check if installed:**
```bash
make --version
```

**Install Make:**

**On Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install build-essential
```

**On RHEL/Fedora:**
```bash
sudo dnf install make
```

**On macOS:**
```bash
xcode-select --install
```

#### 3. Git

**Check if installed:**
```bash
git --version
```

**Install Git:**

**On Ubuntu/Debian:**
```bash
sudo apt-get install git
```

**On RHEL/Fedora:**
```bash
sudo dnf install git
```

## Go Dependencies

The project requires several Go modules. These are **automatically installed** when you run `make build-cli` or `make build`.

### Main Dependencies

```
github.com/prometheus/prometheus      v0.48.0
k8s.io/api                           v0.28.3
k8s.io/apimachinery                  v0.28.3
k8s.io/client-go                     v0.28.3
sigs.k8s.io/controller-runtime       v0.16.3
github.com/go-logr/logr              v1.2.4
github.com/prometheus/common          v0.45.0
```

### Automatic Installation

Dependencies are automatically downloaded when you build:

```bash
# Build CLI (automatically downloads dependencies)
make build-cli

# Or manually download dependencies first
make deps
```

### Manual Dependency Installation

If you need to manually install dependencies:

```bash
# Download all dependencies
go mod download

# Tidy up go.mod and go.sum
go mod tidy

# Verify dependencies
go mod verify
```

## Kubernetes/OpenShift (For Operator Only)

If you want to deploy the **Kubernetes Operator** (not needed for CLI):

### Kubernetes
- **Version:** 1.24 or later
- **kubectl CLI** configured

### OpenShift
- **Version:** 4.10 or later  
- **oc CLI** configured

### Install kubectl

**On Ubuntu/Debian:**
```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
```

**On RHEL/Fedora:**
```bash
sudo dnf install kubernetes-client
```

## Docker (Optional)

Only needed if you want to build container images.

**Check if installed:**
```bash
docker --version
```

**Install Docker:**

**On Ubuntu:**
```bash
sudo apt-get install docker.io
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
```

**On RHEL/Fedora:**
```bash
sudo dnf install docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
```

## Quick Setup Script

Run this to check all prerequisites:

```bash
#!/bin/bash
echo "Checking prerequisites..."
echo ""

# Check Go
if command -v go &> /dev/null; then
    echo "✓ Go installed: $(go version)"
else
    echo "✗ Go not found. Install from: https://go.dev/dl/"
fi

# Check Make
if command -v make &> /dev/null; then
    echo "✓ Make installed: $(make --version | head -1)"
else
    echo "✗ Make not found. Install build-essential"
fi

# Check Git
if command -v git &> /dev/null; then
    echo "✓ Git installed: $(git --version)"
else
    echo "✗ Git not found. Install git"
fi

# Check Docker (optional)
if command -v docker &> /dev/null; then
    echo "✓ Docker installed: $(docker --version)"
else
    echo "○ Docker not found (optional)"
fi

# Check kubectl (optional)
if command -v kubectl &> /dev/null; then
    echo "✓ kubectl installed: $(kubectl version --client --short 2>/dev/null)"
else
    echo "○ kubectl not found (optional, only for operator)"
fi

echo ""
echo "Ready to build? Run: make build-cli"
```

Save this as `check-prereqs.sh` and run:
```bash
chmod +x check-prereqs.sh
./check-prereqs.sh
```

## Building the Project

Once prerequisites are installed:

### CLI Binary

```bash
# Clone repository
git clone https://github.com/ssonigra/prom2sar.git
cd prom2sar

# Build CLI (dependencies installed automatically)
make build-cli

# Install system-wide
sudo make install-cli

# Verify
prom2sar --version
```

### Operator

```bash
# Build operator
make build

# Build Docker image
make docker-build IMG=your-registry/prom2sar-operator:latest
```

## Troubleshooting

### "go: command not found"

**Solution:** Go is not installed or not in PATH.

```bash
# Check if Go is installed
which go

# If not in PATH, add it
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### "cannot find package"

**Solution:** Dependencies not downloaded.

```bash
# Download dependencies
make deps

# Or manually
go mod download
go mod tidy
```

### "permission denied" when installing

**Solution:** Need sudo for system-wide installation.

```bash
sudo make install-cli
```

### Go version too old

**Solution:** Update Go to 1.21+

```bash
# Remove old version
sudo rm -rf /usr/local/go

# Install new version (see Go installation above)
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
```

## Verification

After installing prerequisites, verify everything works:

```bash
# Check Go
go version
# Expected: go version go1.21.x linux/amd64

# Check Make  
make --version
# Expected: GNU Make 4.x

# Download dependencies
cd /path/to/prom2sar
make deps
# Expected: ✓ Dependencies ready

# Build
make build-cli
# Expected: Binary created at bin/prom2sar

# Test
./bin/prom2sar --version
# Expected: prom2sar version 1.0.0
```

## Minimum Specifications

### For Building
- **RAM:** 2GB minimum, 4GB recommended
- **Disk:** 1GB free space
- **CPU:** Any modern CPU

### For Running CLI
- **RAM:** 100MB minimum
- **Disk:** Depends on TSDB size (varies)
- **CPU:** Any CPU

### For Running Operator
- **Kubernetes cluster** with sufficient resources
- **RBAC permissions** to create CRDs

## Platform-Specific Notes

### Ubuntu 20.04+
All prerequisites available in default repos:
```bash
sudo apt-get update
sudo apt-get install golang make git build-essential
```

### RHEL/CentOS 8+
```bash
sudo dnf install golang make git
```

### Fedora
```bash
sudo dnf install golang make git
```

### macOS
```bash
brew install go make git
```

## Next Steps

After installing prerequisites:

1. **Build the CLI:**
   ```bash
   make build-cli
   ```

2. **Read documentation:**
   - [GETTING_STARTED.md](GETTING_STARTED.md) - Quick start guide
   - [CLI_GUIDE.md](CLI_GUIDE.md) - Complete CLI documentation

3. **Run your first conversion:**
   ```bash
   ./bin/prom2sar -tsdb /path/to/prometheus -output ./analysis
   ```

## Support

If you encounter issues with prerequisites:

1. Check this document thoroughly
2. Run the check-prereqs.sh script
3. See [TROUBLESHOOTING.md](CLI_GUIDE.md#troubleshooting)
4. Open an issue: https://github.com/ssonigra/prom2sar/issues

---

**Ready to build?** Run `make build-cli` and the dependencies will be installed automatically!
