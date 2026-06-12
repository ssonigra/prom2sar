# Prometheus to SAR Converter

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

Convert Prometheus Time Series Database (TSDB) dumps to standard SAR (System Activity Report) format. Enable kernel engineers and system administrators to analyze Prometheus metrics using familiar Unix tools—**no Prometheus expertise required!**

## 🎯 Problem Statement

Prometheus is powerful but requires specialized knowledge. System administrators and kernel engineers are experts in traditional tools like `sar`, `grep`, `awk`, and `sed`. This project bridges that gap by converting Prometheus data into the familiar SAR format.

## ✨ Features

- 🚀 **Two deployment modes**: Standalone CLI or Kubernetes Operator
- 📊 **Complete metrics**: CPU, memory, disk I/O, network statistics
- 🔧 **Standard SAR format**: Works with existing tools and scripts
- ⚡ **Fast processing**: Handles gigabytes of TSDB data efficiently
- 📝 **Zero learning curve**: Use grep, awk, sed on the output
- 🎯 **Flexible profiles**: Extract all metrics or focus on specific categories
- 🔄 **Time range filtering**: Analyze specific incidents or time windows

## 📦 Two Ways to Use

### Option 1: Standalone CLI Binary (Recommended for Getting Started)

Perfect for ad-hoc analysis, must-gather processing, and incident investigations.

```bash
# Build the CLI
make build-cli

# Convert Prometheus data to SAR format
./bin/prom2sar -tsdb /prometheus -output ./sar-analysis

# Analyze with standard Unix tools
cat sar-analysis/sar-20260612.txt
grep "14:30" sar-analysis/sar-20260612.txt
awk '$4 > 80' sar-analysis/sar-20260612.txt
```

**[→ CLI Documentation](CLI_README.md)**

### Option 2: Kubernetes Operator

Perfect for automated, recurring conversions in production clusters.

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: daily-sar-conversion
spec:
  sourcePath: /prometheus
  targetPath: /var/lib/dumps
  sarConversion:
    enabled: true
    outputPath: /var/lib/sar
    metricsProfile: all
```

**[→ Operator Documentation](README.md)**

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later (for building from source)
- Access to Prometheus TSDB data
- Linux system (for CLI usage)
- Kubernetes/OpenShift cluster (for operator deployment)

### Install CLI

```bash
# Clone the repository
git clone https://github.com/yourusername/prometheus-dump-operator.git
cd prometheus-dump-operator

# Build the CLI binary
make build-cli

# Install system-wide (optional)
sudo make install-cli

# Verify installation
prom2sar --version
```

### First Conversion

```bash
# Basic conversion (last 24 hours)
prom2sar -tsdb /prometheus

# With specific time range
prom2sar \
  -tsdb /prometheus \
  -start 2026-06-12T00:00:00Z \
  -end 2026-06-12T23:59:59Z \
  -output ./my-analysis

# CPU metrics only
prom2sar -tsdb /prometheus -profile cpu
```

## 📊 Sample Output

The tool generates standard SAR-formatted reports:

```
Linux TSDB  06/12/2026

--------------------------------------------------------------------------------
CPU Utilization (sar -u)
--------------------------------------------------------------------------------

Timestamp    CPU      %user   %nice %system %iowait  %steal   %idle
09:00:00     all       12.50    0.00    3.20    0.50    0.00   83.80
09:01:00     all       15.30    0.00    3.80    0.60    0.00   80.30
09:02:00     all       85.20    0.00   10.50    2.30    0.00    2.00

--------------------------------------------------------------------------------
Memory Utilization (sar -r)
--------------------------------------------------------------------------------

Timestamp     kbmemfree    kbmemused   kbbuffers    kbcached  kbswpfree %memused
09:00:00       2048000      6144000      512000     1536000    4096000    75.00
09:01:00       2040000      6152000      512000     1540000    4096000    75.10
```

Analyze with standard tools:

```bash
# Find CPU spikes over 80%
awk '/^[0-9]{2}:[0-9]{2}/ && $4 > 80' sar-20260612.txt

# Extract memory stats for specific time
sed -n '/14:00/,/15:00/p' sar-20260612.txt | grep kbmem

# Network errors
awk '/eth0/ && ($7 > 0 || $8 > 0)' sar-20260612.txt
```

## 📖 Documentation

### Getting Started
- **[CLI Quick Start](CLI_README.md)** - Standalone binary usage
- **[Operator Guide](README.md)** - Kubernetes deployment
- **[Quick Reference](QUICK_REFERENCE.md)** - One-page cheat sheet for kernel engineers

### In-Depth Guides
- **[Complete CLI Guide](CLI_GUIDE.md)** - All CLI options and examples
- **[SAR Conversion Guide](SAR_CONVERSION_GUIDE.md)** - How metrics are mapped
- **[Deployment Options](DEPLOYMENT_OPTIONS.md)** - When to use CLI vs Operator
- **[Testing Guide](TESTING.md)** - How to test the tool

### Technical Documentation
- **[Implementation Summary](IMPLEMENTATION_SUMMARY.md)** - Architecture and design
- **[Project Overview](PROJECT_OVERVIEW.md)** - Visual architecture diagrams
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute

## 🎯 Common Use Cases

### Incident Investigation
```bash
prom2sar \
  -tsdb /must-gather/prometheus/data \
  -start 2026-06-12T14:00:00Z \
  -end 2026-06-12T16:00:00Z \
  -output /tmp/incident-analysis \
  -verbose
```

### Must-Gather Analysis
```bash
# Extract must-gather
oc adm must-gather

# Convert Prometheus data
prom2sar -tsdb must-gather.local.*/monitoring/prometheus/*/data

# Share SAR files with kernel team
tar -czf sar-analysis.tar.gz sar-output/
```

### Daily Performance Reports
```bash
# Automated via cron
0 1 * * * /usr/local/bin/prom2sar -tsdb /prometheus -output /var/reports/$(date +\%Y-\%m-\%d)
```

### CPU Performance Analysis
```bash
prom2sar \
  -tsdb /prometheus \
  -profile cpu \
  -interval 30 \
  -output ./cpu-analysis
```

## 🔧 Supported Metrics

| SAR Category | Metrics | Prometheus Source |
|--------------|---------|-------------------|
| **CPU** (`sar -u`) | %user, %nice, %system, %iowait, %steal, %idle | `node_cpu_seconds_total{mode=*}` |
| **Memory** (`sar -r`) | kbmemfree, kbmemused, kbbuffers, kbcached, %memused | `node_memory_*_bytes` |
| **Disk I/O** (`sar -d`) | tps, rd_sec/s, wr_sec/s, avgrq-sz, %util | `node_disk_*` |
| **Network** (`sar -n DEV`) | rxpck/s, txpck/s, rxkB/s, txkB/s, rxerr/s, txerr/s | `node_network_*` |

All metrics are sourced from Prometheus node-exporter.

## 🏗️ Architecture

```
Prometheus TSDB Blocks
         ↓
    TSDB Reader
         ↓
   Metrics Mapper (Prometheus → SAR)
         ↓
    SAR Generator
         ↓
   SAR Format Files
         ↓
  Standard Unix Tools (grep, awk, sed)
```

**Key Components:**
- **TSDB Reader** (`pkg/tsdb/`) - Reads Prometheus TSDB blocks
- **Metrics Mapper** (`pkg/sar/mapper.go`) - Maps Prometheus metrics to SAR fields
- **SAR Generator** (`pkg/sar/generator.go`) - Generates SAR-formatted output
- **CLI Binary** (`cmd/prom2sar/`) - Standalone executable
- **Operator** (`cmd/main.go`, `pkg/controller/`) - Kubernetes controller

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Run tests (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Areas for Contribution

- 🐛 **Bug fixes** - Found a bug? Fix it and submit a PR!
- 📚 **Documentation** - Improve guides, add examples, fix typos
- ✨ **Features** - Binary SAR format, additional metrics, performance improvements
- 🧪 **Testing** - Add test coverage, integration tests
- 🌍 **Internationalization** - Translate documentation
- 📦 **Packaging** - RPM/DEB packages, container images

See [open issues](https://github.com/yourusername/prometheus-dump-operator/issues) for ideas.

## 📋 Requirements

### For CLI Usage
- Linux system (or any OS with Go support)
- Access to Prometheus TSDB files
- No runtime dependencies (static binary available)

### For Operator Usage
- Kubernetes 1.24+ or OpenShift 4.10+
- Cluster admin access (for CRD installation)
- Prometheus with node-exporter metrics

### For Building from Source
- Go 1.21 or later
- make
- git

## 🔨 Building

### Build Everything
```bash
# Build both CLI and operator
make build-all
```

### Build CLI Only
```bash
# Current platform
make build-cli

# Static binary (no dependencies)
CGO_ENABLED=0 go build -o prom2sar cmd/prom2sar/main.go

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o prom2sar-linux cmd/prom2sar/main.go
```

### Build Operator
```bash
make build
```

### Build Docker Image
```bash
make docker-build IMG=quay.io/yourorg/prometheus-dump-operator:latest
```

## 🧪 Testing

```bash
# Run unit tests
make test

# Run with coverage
go test ./... -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out
```

See [TESTING.md](TESTING.md) for comprehensive testing guide.

## 📦 Installation Methods

### Method 1: Pre-built Binary
```bash
# Download from releases
wget https://github.com/yourusername/prometheus-dump-operator/releases/download/v1.0.0/prom2sar-linux-amd64

# Make executable
chmod +x prom2sar-linux-amd64

# Move to PATH
sudo mv prom2sar-linux-amd64 /usr/local/bin/prom2sar
```

### Method 2: Build from Source
```bash
git clone https://github.com/yourusername/prometheus-dump-operator.git
cd prometheus-dump-operator
make build-cli
sudo cp bin/prom2sar /usr/local/bin/
```

### Method 3: Go Install
```bash
go install github.com/yourusername/prometheus-dump-operator/cmd/prom2sar@latest
```

### Method 4: Operator Deployment
```bash
# Install CRDs
make install

# Deploy operator
make deploy NAMESPACE=prometheus-dump-operator
```

## 🌟 Examples

See the [examples/](examples/) directory for:
- Basic SAR conversion
- CPU-only analysis
- Custom metrics mapping
- Time range filtering
- Batch processing scripts

## 🐛 Troubleshooting

### Common Issues

**Problem:** No data found in specified time range
```bash
# Solution: Check available data range
prom2sar -tsdb /prometheus -verbose
```

**Problem:** Permission denied accessing TSDB
```bash
# Solution: Check permissions or use sudo
ls -la /prometheus
sudo prom2sar -tsdb /prometheus
```

**Problem:** Empty SAR output
```bash
# Solution: Verify node-exporter metrics exist
# The tool requires Prometheus node-exporter metrics
```

See [Troubleshooting Guide](CLI_GUIDE.md#troubleshooting) for more solutions.

## 📜 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Prometheus project for the TSDB library
- SAR/sysstat for the output format inspiration
- Kubernetes community for operator patterns
- All contributors and users of this project

## 📞 Support

- **Documentation**: Check our comprehensive [guides](CLI_GUIDE.md)
- **Issues**: [GitHub Issues](https://github.com/yourusername/prometheus-dump-operator/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/prometheus-dump-operator/discussions)
- **Questions**: Open an issue with the `question` label

## 🗺️ Roadmap

- [ ] Binary SAR format support (sadc-compatible)
- [ ] Per-node metric breakdowns
- [ ] Additional metric sources beyond node-exporter
- [ ] Real-time conversion from live Prometheus
- [ ] GUI for easier configuration
- [ ] Pre-built packages (RPM, DEB)
- [ ] Container images on multiple registries
- [ ] Prometheus-compatible query interface

## 📈 Project Stats

- **Language**: Go 1.21+
- **Lines of Code**: 2000+
- **Documentation**: 10 comprehensive guides
- **License**: Apache 2.0
- **Status**: Production Ready ✅

## 🔗 Related Projects

- [Prometheus](https://prometheus.io/) - Monitoring system and time series database
- [sysstat](https://github.com/sysstat/sysstat) - System performance tools including sar
- [node-exporter](https://github.com/prometheus/node_exporter) - Hardware and OS metrics exporter

---

**Made with ❤️ for kernel engineers and system administrators who prefer traditional Unix tools**

⭐ Star this repo if you find it useful!
