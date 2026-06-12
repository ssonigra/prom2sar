# prom2sar - Standalone CLI Guide

## Overview

`prom2sar` is a standalone command-line tool that converts Prometheus TSDB dumps to SAR format. **No Kubernetes required!** Use it on any Linux system with Prometheus data.

## Installation

### Option 1: Build from Source

```bash
# Clone the repo
cd prometheus-dump-operator

# Build the CLI binary
make build-cli

# Binary is created at: bin/prom2sar
./bin/prom2sar --version
```

### Option 2: Install System-Wide

```bash
# Build and install to /usr/local/bin
make install-cli

# Now use from anywhere
prom2sar --version
```

### Option 3: Manual Build

```bash
go build -o prom2sar cmd/prom2sar/main.go
sudo mv prom2sar /usr/local/bin/
```

### Option 4: Cross-Compile for Other Systems

```bash
# For Linux AMD64
GOOS=linux GOARCH=amd64 go build -o prom2sar-linux-amd64 cmd/prom2sar/main.go

# For Linux ARM64
GOOS=linux GOARCH=arm64 go build -o prom2sar-linux-arm64 cmd/prom2sar/main.go

# Create static binary (no dependencies)
CGO_ENABLED=0 GOOS=linux go build -o prom2sar-static cmd/prom2sar/main.go
```

## Quick Start

### Basic Usage

```bash
# Convert Prometheus data from the last 24 hours
prom2sar -tsdb /prometheus

# Output:
#   sar-output/sar-20260612.txt
#   sar-output/sar-summary-20260612.txt
```

### Specify Time Range

```bash
prom2sar -tsdb /prometheus \
  -start 2026-06-12T00:00:00Z \
  -end 2026-06-12T23:59:59Z
```

### Custom Output Directory

```bash
prom2sar -tsdb /prometheus -output /tmp/incident-analysis
```

### CPU Metrics Only

```bash
prom2sar -tsdb /prometheus -profile cpu -output ./cpu-analysis
```

### Quick Summary

```bash
prom2sar -tsdb /prometheus -summary
```

## Command-Line Options

```
Usage: prom2sar [options]

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
  -verbose
        Verbose output
  -summary
        Generate summary only (no full report)
  -version
        Show version
```

## Profiles

| Profile | Metrics Included | Use Case |
|---------|------------------|----------|
| `all` | CPU, memory, disk, network | Complete system analysis (default) |
| `cpu` | CPU utilization only | CPU spike investigation |
| `memory` | Memory utilization only | Memory leak analysis |
| `disk` | Disk I/O only | Disk bottleneck analysis |
| `network` | Network I/O only | Network issue diagnosis |

## Getting Prometheus TSDB Data

### From Running Prometheus

```bash
# Port-forward to Prometheus
kubectl port-forward -n openshift-monitoring prometheus-k8s-0 9090:9090

# Create snapshot
curl -XPOST http://localhost:9090/api/v1/admin/tsdb/snapshot

# Response: {"status":"success","data":{"name":"20260612T120000Z-..."}}

# Copy snapshot from pod
kubectl cp openshift-monitoring/prometheus-k8s-0:/prometheus/snapshots/20260612T120000Z-... ./prometheus-data
```

### From Must-Gather

```bash
# Extract must-gather
oc adm must-gather

# Find Prometheus data
find must-gather.local.* -path "*/monitoring/prometheus/*/data" -type d

# Use that path with prom2sar
prom2sar -tsdb must-gather.local.*/monitoring/prometheus/prometheus-k8s-0/data
```

### From Prometheus Backup

```bash
# If you have Prometheus backup/snapshot
tar -xzf prometheus-backup.tar.gz
prom2sar -tsdb ./prometheus-backup/data
```

## Usage Examples

### Example 1: Incident Investigation

```bash
# Analyze specific incident time window
prom2sar \
  -tsdb /prometheus \
  -start 2026-06-12T14:00:00Z \
  -end 2026-06-12T16:00:00Z \
  -output /tmp/incident-20260612 \
  -verbose

# Results in:
#   /tmp/incident-20260612/sar-20260612.txt
#   /tmp/incident-20260612/sar-summary-20260612.txt
```

### Example 2: CPU Performance Analysis

```bash
# High-resolution CPU metrics (30-second intervals)
prom2sar \
  -tsdb /prometheus \
  -profile cpu \
  -interval 30 \
  -output ./cpu-analysis \
  -start 2026-06-12T08:00:00Z \
  -end 2026-06-12T18:00:00Z
```

### Example 3: Memory Leak Detection

```bash
# Memory metrics for last 7 days
prom2sar \
  -tsdb /prometheus \
  -profile memory \
  -start $(date -u -d '7 days ago' '+%Y-%m-%dT%H:%M:%SZ') \
  -end $(date -u '+%Y-%m-%dT%H:%M:%SZ') \
  -output ./memory-trend
```

### Example 4: Quick Health Check

```bash
# Just show summary, no full report
prom2sar -tsdb /prometheus -summary
```

### Example 5: Batch Processing

```bash
#!/bin/bash
# Convert each day separately

for day in {01..07}; do
  prom2sar \
    -tsdb /prometheus \
    -start 2026-06-${day}T00:00:00Z \
    -end 2026-06-${day}T23:59:59Z \
    -output ./analysis/day-${day}
done
```

## Output Files

After running `prom2sar`, you get:

### sar-YYYYMMDD.txt
Main SAR report with all metrics. Use this for analysis.

```
Linux TSDB  06/12/2026

--------------------------------------------------------------------------------
CPU Utilization (sar -u)
--------------------------------------------------------------------------------

Timestamp    CPU      %user   %nice %system %iowait  %steal   %idle
09:00:00     all       12.50    0.00    3.20    0.50    0.00   83.80
09:01:00     all       15.30    0.00    3.80    0.60    0.00   80.30
...
```

### sar-summary-YYYYMMDD.txt
Quick overview and statistics.

```
=== Prometheus to SAR Conversion Summary ===

Hostname: prometheus-tsdb
Total samples: 1440
Time range: 2026-06-12T00:00:00Z to 2026-06-12T23:59:59Z

✓ CPU statistics available
✓ Memory statistics available
✓ Disk I/O statistics available
✓ Network statistics available

=== Quick Summary ===
CPU (avg):
  User:    12.50%
  System:   3.20%
  IOWait:   0.50%
  Idle:    83.80%
...
```

## Analyzing SAR Output

Once you have the SAR files, use standard Unix tools:

```bash
# View full report
cat sar-20260612.txt

# Find CPU spikes (>80% user)
awk '/^[0-9]{2}:[0-9]{2}/ && $4 > 80 {print}' sar-20260612.txt

# Extract specific time range
sed -n '/14:00/,/16:00/p' sar-20260612.txt

# Get memory stats
grep kbmemfree sar-20260612.txt

# Find disk bottlenecks
awk '/sda/ && $9 > 80 {print $1, "Util:", $9"%"}' sar-20260612.txt

# Network errors
awk '/eth0/ && ($7 > 0 || $8 > 0) {print}' sar-20260612.txt
```

## Troubleshooting

### Error: TSDB path does not exist

```bash
# Verify path exists
ls -la /prometheus

# Check for TSDB blocks
ls -la /prometheus/01*/
```

### Error: No data found in the specified time range

```bash
# Check available data range
prom2sar -tsdb /prometheus -verbose | grep "Block"

# Adjust your time range to match available data
```

### Error: Failed to open TSDB

```bash
# Check permissions
ls -la /prometheus

# Run with sudo if needed (be careful!)
sudo prom2sar -tsdb /prometheus
```

### Empty SAR output

```bash
# Use verbose mode to see what's happening
prom2sar -tsdb /prometheus -verbose

# Check if node-exporter metrics exist
# (The tool requires Prometheus node-exporter metrics)
```

## Performance Tips

### Large Time Ranges

For very large time ranges (months), process in chunks:

```bash
#!/bin/bash
for month in {01..12}; do
  prom2sar \
    -tsdb /prometheus \
    -start 2026-${month}-01T00:00:00Z \
    -end 2026-${month}-28T23:59:59Z \
    -output ./analysis/2026-${month}
done
```

### High-Frequency Sampling

Smaller intervals create larger files:

```bash
# 10-second intervals = 6x more data than 60-second
prom2sar -tsdb /prometheus -interval 10 -output ./detailed

# vs normal 60-second intervals
prom2sar -tsdb /prometheus -interval 60 -output ./normal
```

### Specific Metrics Only

Use profiles to reduce processing time:

```bash
# Faster - CPU only
prom2sar -tsdb /prometheus -profile cpu

# Slower - all metrics
prom2sar -tsdb /prometheus -profile all
```

## Integration Examples

### Cron Job

```bash
# Add to crontab for daily conversion
0 1 * * * /usr/local/bin/prom2sar -tsdb /prometheus -output /var/log/sar-$(date +\%Y\%m\%d)
```

### Script Integration

```bash
#!/bin/bash
# analyze-incident.sh

INCIDENT_TIME="$1"
OUTPUT_DIR="/tmp/incident-$$"

echo "Analyzing incident at $INCIDENT_TIME"

# Convert to SAR
prom2sar \
  -tsdb /prometheus \
  -start "${INCIDENT_TIME}" \
  -end "$(date -u -d "${INCIDENT_TIME} + 2 hours" '+%Y-%m-%dT%H:%M:%SZ')" \
  -output "$OUTPUT_DIR"

# Run analysis
echo "Top 10 CPU spikes:"
awk '/^[0-9]{2}:[0-9]{2}/ {print $1, $4}' "$OUTPUT_DIR"/sar-*.txt | sort -k2 -nr | head -10

echo "Memory usage:"
grep kbmemused "$OUTPUT_DIR"/sar-*.txt | tail -5

echo "Full report: $OUTPUT_DIR/sar-*.txt"
```

### CI/CD Pipeline

```yaml
# .gitlab-ci.yml
performance-analysis:
  script:
    - prom2sar -tsdb /prometheus -summary
    - if grep -q "User:.*9[0-9]\." sar-summary-*.txt; then exit 1; fi
  artifacts:
    paths:
      - sar-output/
```

## Comparison with Kubernetes Operator

| Feature | CLI (`prom2sar`) | Operator |
|---------|------------------|----------|
| **Deployment** | Standalone binary | Kubernetes cluster |
| **Usage** | Command-line | Custom Resource |
| **Requirements** | Just the binary | OpenShift/K8s |
| **Best For** | Ad-hoc analysis, scripts | Automated workflows |
| **Installation** | Copy binary | Deploy with `oc apply` |
| **Permissions** | File system access | Cluster RBAC |

**Use CLI when:**
- Analyzing must-gather data
- Working on a jumphost/bastion
- Scripting/automation outside K8s
- One-off incident analysis
- No cluster access available

**Use Operator when:**
- Regular automated conversions
- Integration with cluster workflows
- Team self-service portal
- Continuous monitoring

## Building Static Binary

For maximum portability (no dependencies):

```bash
# Build completely static binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags '-extldflags "-static"' \
  -o prom2sar-static \
  cmd/prom2sar/main.go

# Test it works
./prom2sar-static --version

# This binary runs on any Linux system!
```

## Version and Updates

```bash
# Check version
prom2sar --version

# Update (if installed from source)
cd prometheus-dump-operator
git pull
make build-cli
sudo cp bin/prom2sar /usr/local/bin/
```

## Quick Reference Card

```bash
# Basic conversion
prom2sar -tsdb /prometheus

# Specific time range
prom2sar -tsdb /prometheus -start 2026-06-12T00:00:00Z -end 2026-06-12T23:59:59Z

# CPU only
prom2sar -tsdb /prometheus -profile cpu

# Quick summary
prom2sar -tsdb /prometheus -summary

# Verbose mode
prom2sar -tsdb /prometheus -verbose

# Custom output location
prom2sar -tsdb /prometheus -output /tmp/analysis

# High resolution (30s intervals)
prom2sar -tsdb /prometheus -interval 30
```

## Getting Help

```bash
# Show help
prom2sar --help

# Show version
prom2sar --version
```

## Real-World Example

```bash
# Scenario: Investigating high CPU on June 12, 2026 at 14:30

# Step 1: Convert relevant time window
prom2sar \
  -tsdb /must-gather/prometheus/data \
  -start 2026-06-12T14:00:00Z \
  -end 2026-06-12T15:00:00Z \
  -output /tmp/incident-investigation \
  -verbose

# Output:
# === Prometheus to SAR Conversion ===
# TSDB Path:    /must-gather/prometheus/data
# Output Path:  /tmp/incident-investigation
# ...
# ✓ Generated summary: /tmp/incident-investigation/sar-summary-20260612.txt
# ✓ Generated report:  /tmp/incident-investigation/sar-20260612.txt
# ✓ Conversion completed successfully!

# Step 2: Quick summary
cat /tmp/incident-investigation/sar-summary-20260612.txt

# Step 3: Find exact spike time
grep "14:[23]" /tmp/incident-investigation/sar-20260612.txt | \
  awk '{print $1, "User:", $4"%", "System:", $6"%"}'

# Step 4: Check memory at that time
sed -n '/14:25/,/14:35/p' /tmp/incident-investigation/sar-20260612.txt | \
  grep kbmem

# Done! Standard sar analysis workflow.
```

---

**Bottom Line**: `prom2sar` gives you the power to analyze Prometheus data anywhere, using tools you already know. No Kubernetes required!
