# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Binary SAR format support (sadc-compatible)
- Per-node metric breakdowns
- Additional Prometheus exporters support
- Real-time conversion from live Prometheus
- Performance optimizations for large datasets

## [1.0.0] - 2026-06-12

### Added
- **Standalone CLI binary (`prom2sar`)** for Linux systems
  - Convert Prometheus TSDB to SAR format without Kubernetes
  - Command-line options for time range, profiles, output paths
  - Static binary build support for portability
  - Verbose mode for debugging
  - Summary-only mode for quick health checks
- **Kubernetes Operator** for automated conversions
  - Custom Resource Definition (PrometheusDumpLoader)
  - Controller for reconciling conversion requests
  - Status tracking and reporting
  - Time range filtering
  - Multiple metric profiles (all, cpu, memory, disk, network)
- **TSDB Reader** (`pkg/tsdb/`)
  - Read Prometheus Time Series Database blocks
  - Query metrics by name and labels
  - Block discovery and metadata extraction
- **Metrics Mapper** (`pkg/sar/mapper.go`)
  - Map Prometheus node-exporter metrics to SAR equivalents
  - CPU utilization metrics (user, system, iowait, idle, steal)
  - Memory utilization metrics (free, used, cached, buffers)
  - Disk I/O metrics (tps, read/write rates, utilization)
  - Network metrics (rx/tx packets, bytes, errors)
- **SAR Generator** (`pkg/sar/generator.go`)
  - Generate text-based SAR reports
  - CPU utilization report (sar -u equivalent)
  - Memory utilization report (sar -r equivalent)
  - Disk I/O report (sar -d equivalent)
  - Network statistics report (sar -n DEV equivalent)
  - Summary reports with averages
  - Compact reports for quick viewing
- **Documentation**
  - Comprehensive README with quick start
  - CLI usage guide (CLI_GUIDE.md)
  - SAR conversion guide (SAR_CONVERSION_GUIDE.md)
  - Quick reference for kernel engineers (QUICK_REFERENCE.md)
  - Testing guide (TESTING.md)
  - Implementation summary (IMPLEMENTATION_SUMMARY.md)
  - Project overview with diagrams (PROJECT_OVERVIEW.md)
  - Deployment options guide (DEPLOYMENT_OPTIONS.md)
  - Contributing guide (CONTRIBUTING.md)
- **Build Tools**
  - Makefile with targets for CLI and operator
  - Build script for multi-platform binaries
  - Installation script
  - Docker build support
- **Examples**
  - Basic SAR conversion example
  - CPU-only analysis example
  - Custom metrics mapping example

### Technical Details
- Go 1.21+ support
- Prometheus TSDB library integration
- Kubernetes controller-runtime for operator
- Support for OpenShift 4.10+
- Support for Kubernetes 1.24+

### Metrics Supported
- CPU: %user, %nice, %system, %iowait, %steal, %idle
- Memory: kbmemfree, kbmemused, kbbuffers, kbcached, kbswpfree, %memused
- Disk: tps, rd_sec/s, wr_sec/s, avgrq-sz, avgqu-sz, %util
- Network: rxpck/s, txpck/s, rxkB/s, txkB/s, rxerr/s, txerr/s

## [0.1.0] - Initial Development

### Added
- Basic project structure
- Initial operator skeleton
- CRD definitions for PrometheusDumpLoader

---

## Version History

| Version | Release Date | Highlights |
|---------|-------------|------------|
| 1.0.0   | 2026-06-12  | First production release with CLI and operator |
| 0.1.0   | 2026-06-01  | Initial development version |

## Upgrade Guide

### From 0.1.0 to 1.0.0

This is a major release with significant new features:

**Breaking Changes:**
- None (first production release)

**New Features:**
- Standalone CLI binary
- SAR conversion capability
- Complete metric mapping

**Migration Steps:**
1. Update CRDs: `kubectl apply -f deploy/crds/`
2. Update operator deployment: `kubectl apply -f deploy/`
3. No changes needed to existing CRs (backward compatible)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

## Links

- [GitHub Repository](https://github.com/yourusername/prometheus-dump-operator)
- [Issue Tracker](https://github.com/yourusername/prometheus-dump-operator/issues)
- [Releases](https://github.com/yourusername/prometheus-dump-operator/releases)
