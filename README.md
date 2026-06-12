# prom2sar - Prometheus to SAR Converter

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/ssonigra/prom2sar)
[![Version](https://img.shields.io/badge/version-1.0.0-blue)](https://github.com/ssonigra/prom2sar/releases)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE)

Convert Prometheus TSDB dumps to SAR-compatible format for kernel engineers and system administrators.

**Perfect for**: Kernel engineers who need to analyze Prometheus metrics using familiar SAR tools (grep, awk, sed) without learning PromQL.

---

## 🚀 Quick Start

```bash
# 1. Clone repository
git clone https://github.com/ssonigra/prom2sar.git
cd prom2sar

# 2. Build CLI
make build-cli

# 3. Convert Prometheus data to SAR format
./bin/prom2sar -tsdb /path/to/prometheus -output ./sar-results -verbose

# 4. Analyze with standard Unix tools
cat sar-results/sar-summary-*.txt
grep "12:00:00" sar-results/sar-*.txt
awk '/CPU/ && $2 > 80 {print}' sar-results/sar-*.txt
```

---

## 📋 Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Usage](#-usage)
  - [CLI Tool](#cli-tool-recommended)
  - [Kubernetes Operator](#kubernetes-operator)
- [Testing](#-testing)
- [Output Format](#-output-format)
- [Use Cases](#-use-cases)
- [Documentation](#-documentation)
- [Contributing](#-contributing)

---

## ✨ Features

### 🎯 Core Capabilities

- **TSDB Reading** - Reads Prometheus Time Series Database blocks directly
- **SAR Conversion** - Converts Prometheus metrics to standard SAR format
- **Multiple Profiles** - CPU, memory, disk, network metrics (individually or all together)
- **Time Range Filtering** - Extract specific time windows for analysis
- **Familiar Output** - Standard SAR format that kernel teams already know
- **No PromQL Required** - Analyze metrics without learning Prometheus query language

### 📊 Supported Metrics

| Profile | Metrics | SAR Equivalent |
|---------|---------|----------------|
| **CPU** | User, System, IOWait, Idle | `sar -u` |
| **Memory** | Total, Used, Free, Cached, Buffers, Swap | `sar -r` |
| **Disk** | TPS, Read/Write KB/s, Utilization | `sar -d` |
| **Network** | RX/TX packets/s, RX/TX KB/s, Errors | `sar -n DEV` |
| **All** | All of the above | Combined report |

---

## 🔧 Installation

### Prerequisites

- **Go 1.21+** (for building)
- **Make**
- **Git**
- **Prometheus TSDB data** (for testing)

### Option 1: Build from Source

```bash
# Clone repository
git clone https://github.com/ssonigra/prom2sar.git
cd prom2sar

# Check prerequisites
./check-prereqs.sh

# Build CLI binary
make build-cli

# Install system-wide (optional)
sudo make install-cli
```

### Option 2: Quick Test Build

```bash
# Just build without installation
make build-cli

# Binary will be at: bin/prom2sar
./bin/prom2sar --version
```

For detailed prerequisites, see **[PREREQUISITES.md](PREREQUISITES.md)**

---

## 📖 Usage

### CLI Tool (Recommended)

The standalone CLI tool works on any Linux system - no Kubernetes required.

#### Basic Usage

```bash
# Convert last 24 hours of data
./bin/prom2sar -tsdb /var/lib/prometheus/data -output ./results

# Specific time range
./bin/prom2sar \
  -tsdb /var/lib/prometheus/data \
  -start 2026-06-12T00:00:00Z \
  -end 2026-06-12T23:59:59Z \
  -output ./incident-analysis

# CPU metrics only
./bin/prom2sar -tsdb /prometheus -profile cpu -output ./cpu-analysis

# With verbose output
./bin/prom2sar -tsdb /prometheus -output ./results -verbose
```

#### All CLI Options

```bash
./bin/prom2sar [options]

Options:
  -tsdb string
      Path to Prometheus TSDB directory (required)
  -output string
      Output directory for SAR files (default "./sar-output")
  -start string
      Start time (RFC3339 format, e.g., 2026-06-12T00:00:00Z)
  -end string
      End time (RFC3339 format, e.g., 2026-06-12T23:59:59Z)
  -interval int
      Sampling interval in seconds (default 60)
  -profile string
      Metrics profile: all, cpu, memory, disk, network (default "all")
  -summary
      Generate summary only (no full report)
  -verbose
      Verbose output
  -version
      Show version
```

For detailed CLI usage, see **[CLI_GUIDE.md](CLI_GUIDE.md)**

### Kubernetes Operator

Deploy as an OpenShift/Kubernetes operator for automated conversions.

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: prometheus-to-sar
spec:
  sourcePath: /prometheus
  targetPath: /var/lib/prometheus-dumps
  timeRange:
    start: "2026-06-12T00:00:00Z"
    end: "2026-06-12T23:59:59Z"
  
  sarConversion:
    enabled: true
    outputPath: /var/lib/sar-output
    format: text
    interval: 60
    metricsProfile: all
```

---

## 🧪 Testing

### Quick Smoke Test (30 seconds)

```bash
./QUICK_TEST.sh
```

Output:
```
✓ Test 1: Binary exists... PASS
✓ Test 2: Version check... PASS (prom2sar version 1.0.0)
✓ Test 3: Help output... PASS
✓ Test 4: Error handling (no tsdb)... PASS
✓ Test 5: Invalid path handling... PASS
✓ Test 6: Invalid time format... PASS
✓ Test 7: Invalid profile... PASS

All Tests Passed!
```

### Docker Test with Sample Data (10 minutes)

See **[HOW_TO_TEST.txt](HOW_TO_TEST.txt)** for complete Docker-based testing setup.

### With Real Prometheus Data

```bash
# From local Prometheus
./bin/prom2sar -tsdb /var/lib/prometheus/data -output ./results -verbose

# From OpenShift (copy data first)
oc rsync openshift-monitoring/prometheus-k8s-0:/prometheus ./prom-data
./bin/prom2sar -tsdb ./prom-data -output ./results -verbose
```

For comprehensive testing guide, see **[TESTING.md](TESTING.md)**

---

## 📄 Output Format

### Summary File

```
=== Prometheus to SAR Conversion Summary ===
Source: prometheus-tsdb
Time Range: 2026-06-11 12:00:00 to 2026-06-12 12:00:00
Duration: 24h0m0s
Data Points: 1440

=== System Statistics ===
CPU:
  Average User:   15.3%
  Average System: 8.2%
  Average IOWait: 2.1%
  Average Idle:   74.4%

Memory:
  Total:     16384 MB
  Used Avg:  8192 MB (50.0%)
  Free Avg:  4096 MB
```

### Full SAR Report

```
12:00:00    CPU    %user  %nice  %system  %iowait  %steal  %idle
12:01:00    all    15.32   0.00     8.21     2.14    0.00  74.33
12:02:00    all    16.45   0.00     7.89     1.98    0.00  73.68

12:00:00    kbmemfree  kbmemused  %memused  kbcached  kbbuffers
12:01:00      4194304    8388608     66.67   2097152    1048576
12:02:00      4128768    8454144     67.19   2105344    1052672
```

The output is **identical to standard `sar` command output** - use your existing SAR analysis tools!

---

## 🎯 Use Cases

### For Kernel Engineers

Analyze Prometheus data without learning PromQL:

```bash
# Find CPU spikes
awk '/CPU/ && $2 > 80 {print $0}' sar-*.txt

# Memory pressure analysis
grep -A 5 "Memory" sar-*.txt | awk '$2 < 100000 {print}'

# Disk bottlenecks
grep -A 10 "Disk" sar-*.txt | awk '$NF > 90 {print}'

# Network errors
grep -A 10 "Network" sar-*.txt | awk '$5 > 0 || $6 > 0 {print}'
```

### Incident Investigation

```bash
# 1. Convert incident timeframe
./bin/prom2sar \
  -tsdb /var/lib/prometheus/data \
  -start 2026-06-12T10:00:00Z \
  -end 2026-06-12T12:00:00Z \
  -output ./incident-june12 \
  -verbose

# 2. Quick summary
cat incident-june12/sar-summary-*.txt

# 3. Find issues at specific time
grep "10:45" incident-june12/sar-*.txt

# 4. Package for team
tar czf incident-analysis.tar.gz incident-june12/
```

### Performance Analysis

- Convert historical Prometheus data to SAR format
- Use existing SAR-based analysis scripts
- Integrate with performance monitoring workflows
- Share with teams unfamiliar with Prometheus

---

## 📚 Documentation

| Document | Description |
|----------|-------------|
| **[HOW_TO_TEST.txt](HOW_TO_TEST.txt)** | Complete step-by-step testing guide |
| **[QUICK_TEST.sh](QUICK_TEST.sh)** | Automated smoke test script |
| **[TESTING.md](TESTING.md)** | Comprehensive testing scenarios |
| **[CLI_GUIDE.md](CLI_GUIDE.md)** | Detailed CLI usage reference |
| **[SAR_CONVERSION_GUIDE.md](SAR_CONVERSION_GUIDE.md)** | SAR format and conversion details |
| **[BUILD_SUCCESS.md](BUILD_SUCCESS.md)** | Build verification report |
| **[PREREQUISITES.md](PREREQUISITES.md)** | Installation requirements |
| **[CONTRIBUTING.md](CONTRIBUTING.md)** | Contribution guidelines |

---

## 🏗️ Project Structure

```
prom2sar/
├── bin/
│   └── prom2sar                      # CLI binary (after build)
├── cmd/
│   ├── main.go                       # Operator entrypoint
│   └── prom2sar/
│       └── main.go                   # CLI entrypoint
├── pkg/
│   ├── apis/prometheus/v1alpha1/     # Custom Resource definitions
│   ├── controller/                   # Kubernetes controller
│   ├── loader/                       # TSDB dump loader
│   ├── tsdb/                         # TSDB block reader
│   │   └── reader.go                 # Prometheus TSDB API
│   └── sar/                          # SAR conversion engine
│       ├── mapper.go                 # Prometheus → SAR mapping
│       ├── generator.go              # SAR format output
│       └── converter.go              # Conversion orchestration
├── examples/                         # Example Custom Resources
├── deploy/                           # Kubernetes manifests
├── Makefile                          # Build automation
├── HOW_TO_TEST.txt                   # Testing walkthrough
├── QUICK_TEST.sh                     # Automated tests
└── README.md                         # This file
```

---

## 🔨 Development

### Build

```bash
# Build CLI binary
make build-cli

# Build operator binary
make build

# Build both
make all

# Clean build artifacts
make clean
```

### Run Locally

```bash
# Run CLI directly
go run cmd/prom2sar/main.go -tsdb /path/to/prometheus -output ./results

# Run operator
go run cmd/main.go
```

### Docker

```bash
# Build operator image
make docker-build IMG=quay.io/youruser/prom2sar:v1.0.0

# Push to registry
make docker-push IMG=quay.io/youruser/prom2sar:v1.0.0

# Deploy to cluster
make deploy IMG=quay.io/youruser/prom2sar:v1.0.0
```

---

## 🤝 Contributing

We welcome contributions! Please see **[CONTRIBUTING.md](CONTRIBUTING.md)** for guidelines.

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `./QUICK_TEST.sh`
5. Commit: `git commit -m 'Add amazing feature'`
6. Push: `git push origin feature/amazing-feature`
7. Open a Pull Request

---

## 📊 Build Status

- ✅ **Build**: Passing
- ✅ **CLI Binary**: 29MB
- ✅ **Operator Binary**: 62MB
- ✅ **Tests**: 7/7 passing
- ✅ **Go Version**: 1.21+
- ✅ **Dependencies**: All verified

See **[BUILD_SUCCESS.md](BUILD_SUCCESS.md)** for detailed build information.

---

## 🐛 Troubleshooting

### Common Issues

**Error: "TSDB path does not exist"**
- Verify the path is correct
- Ensure you have read permissions
- Check it's a valid Prometheus data directory

**Error: "No data found in time range"**
- Use `-verbose` to see available TSDB blocks
- Adjust start/end times to match available data

**Binary not found**
- Run: `make build-cli`
- Verify: `ls -lh bin/prom2sar`

For more help, see **[TESTING.md](TESTING.md#troubleshooting)**

---

## 📜 License

Apache License 2.0 - see [LICENSE](LICENSE) for details

---

## 🔗 Links

- **GitHub Repository**: https://github.com/ssonigra/prom2sar
- **Issue Tracker**: https://github.com/ssonigra/prom2sar/issues
- **Releases**: https://github.com/ssonigra/prom2sar/releases

---

## 🌟 Star History

If you find this project useful, please consider giving it a star! ⭐

---

## 📞 Support

- **Issues**: https://github.com/ssonigra/prom2sar/issues
- **Discussions**: https://github.com/ssonigra/prom2sar/discussions

---

**Made with ❤️ for kernel engineers who love SAR and need to work with Prometheus data**
