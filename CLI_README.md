# prom2sar - Standalone CLI Tool

Convert Prometheus TSDB dumps to SAR format. **No Kubernetes required!**

## Quick Install

```bash
# Build the CLI
cd prometheus-dump-operator
make build-cli

# Install system-wide (optional)
sudo make install-cli

# Or just use from bin/
./bin/prom2sar --version
```

## Usage

```bash
# Basic usage
prom2sar -tsdb /prometheus

# With time range
prom2sar -tsdb /prometheus \
  -start 2026-06-12T00:00:00Z \
  -end 2026-06-12T23:59:59Z

# CPU metrics only
prom2sar -tsdb /prometheus -profile cpu
```

## What You Get

```
sar-output/
├── sar-20260612.txt          # Main SAR report
└── sar-summary-20260612.txt  # Quick summary
```

Use standard Unix tools to analyze:

```bash
cat sar-20260612.txt
grep "14:30" sar-20260612.txt
awk '$4 > 80' sar-20260612.txt
```

## Features

- ✅ **Standalone binary** - No dependencies, runs anywhere
- ✅ **Standard SAR format** - Use grep, awk, sed
- ✅ **All metrics** - CPU, memory, disk, network
- ✅ **Time range filtering** - Analyze specific incidents
- ✅ **Multiple profiles** - Focus on what matters
- ✅ **Fast** - Processes gigabytes of data quickly

## Output Example

```
Linux TSDB  06/12/2026

CPU Utilization (sar -u)
Timestamp    CPU      %user   %nice %system %iowait  %steal   %idle
09:00:00     all       12.50    0.00    3.20    0.50    0.00   83.80
09:01:00     all       15.30    0.00    3.80    0.60    0.00   80.30
...
```

## Documentation

- **[CLI_GUIDE.md](CLI_GUIDE.md)** - Comprehensive guide with examples
- **[SAR_CONVERSION_GUIDE.md](SAR_CONVERSION_GUIDE.md)** - SAR format details
- **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - One-page reference

## Common Use Cases

### Incident Analysis
```bash
prom2sar -tsdb /must-gather/prometheus/data \
  -start 2026-06-12T14:00:00Z \
  -end 2026-06-12T16:00:00Z \
  -output /tmp/incident
```

### Performance Baseline
```bash
prom2sar -tsdb /prometheus -profile all -output ./baseline
```

### Quick Health Check
```bash
prom2sar -tsdb /prometheus -summary
```

## Building from Source

```bash
# Simple build
go build -o prom2sar cmd/prom2sar/main.go

# Static binary (no dependencies)
CGO_ENABLED=0 go build -o prom2sar-static cmd/prom2sar/main.go

# Cross-compile for other systems
GOOS=linux GOARCH=amd64 go build -o prom2sar-linux cmd/prom2sar/main.go
```

## All Available Options

```
  -tsdb string
        Path to Prometheus TSDB directory (required)
  -output string
        Output directory (default "./sar-output")
  -start string
        Start time (RFC3339 format)
  -end string
        End time (RFC3339 format)
  -interval int
        Sampling interval in seconds (default 60)
  -profile string
        Metrics profile: all, cpu, memory, disk, network (default "all")
  -verbose
        Verbose output
  -summary
        Summary only (no full report)
  -version
        Show version
```

## Where to Get Prometheus Data

### From Kubernetes/OpenShift
```bash
kubectl port-forward prometheus-0 9090:9090
curl -XPOST http://localhost:9090/api/v1/admin/tsdb/snapshot
kubectl cp prometheus-0:/prometheus/snapshots/... ./prometheus-data
```

### From Must-Gather
```bash
oc adm must-gather
# Data in: must-gather.local.*/monitoring/prometheus/*/data
```

### From Backup
```bash
tar -xzf prometheus-backup.tar.gz
prom2sar -tsdb ./prometheus-backup/data
```

## Requirements

- Linux system (or any OS with Go support)
- Access to Prometheus TSDB files
- Go 1.21+ (for building from source)
- No runtime dependencies (static binary)

## Troubleshooting

**No data found:**
```bash
# Check TSDB blocks exist
ls -la /prometheus/01*/

# Use verbose mode
prom2sar -tsdb /prometheus -verbose
```

**Permission denied:**
```bash
# Check file permissions
ls -la /prometheus

# Run with appropriate permissions
sudo prom2sar -tsdb /prometheus
```

## License

Part of the Prometheus Dump Loader Operator project.

## Support

- See [CLI_GUIDE.md](CLI_GUIDE.md) for detailed documentation
- Check [examples/](examples/) for sample use cases
- Read [QUICK_REFERENCE.md](QUICK_REFERENCE.md) for kernel engineers

---

**For Kernel Engineers**: This tool converts Prometheus data to the familiar SAR format. Use it exactly like you would use regular `sar` output files. No Prometheus knowledge needed!
