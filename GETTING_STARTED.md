## Getting Started with Prometheus to SAR Converter

This guide will help you get up and running with the Prometheus to SAR converter in under 10 minutes.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Start with CLI](#quick-start-with-cli)
3. [Quick Start with Operator](#quick-start-with-operator)
4. [Your First Conversion](#your-first-conversion)
5. [Analyzing the Output](#analyzing-the-output)
6. [Next Steps](#next-steps)

## Prerequisites

### For CLI Usage

- **Linux system** (or any OS with Go installed)
- **Access to Prometheus TSDB data**
  - From a Prometheus snapshot
  - From must-gather
  - From backup/archive
- **Go 1.21+** (only if building from source)

### For Operator Usage

- **Kubernetes 1.24+** or **OpenShift 4.10+**
- **Cluster admin access** (for CRD installation)
- **kubectl or oc CLI** configured
- **Prometheus with node-exporter** running in cluster

## Quick Start with CLI

### Step 1: Get the Binary

**Option A: Download Pre-built Binary** (Easiest)

```bash
# Download latest release
wget https://github.com/yourusername/prometheus-dump-operator/releases/latest/download/prom2sar-linux-amd64-static

# Make executable
chmod +x prom2sar-linux-amd64-static

# Move to PATH
sudo mv prom2sar-linux-amd64-static /usr/local/bin/prom2sar

# Verify
prom2sar --version
```

**Option B: Build from Source**

```bash
# Clone repository
git clone https://github.com/yourusername/prometheus-dump-operator.git
cd prometheus-dump-operator

# Build CLI
make build-cli

# Install system-wide (optional)
sudo make install-cli

# Or just use from bin/
./bin/prom2sar --version
```

### Step 2: Get Prometheus Data

You need Prometheus TSDB blocks. Here's how to get them:

**From Running Prometheus:**

```bash
# Port-forward to Prometheus
kubectl port-forward -n openshift-monitoring prometheus-k8s-0 9090:9090

# Create snapshot (in another terminal)
curl -XPOST http://localhost:9090/api/v1/admin/tsdb/snapshot

# Response shows snapshot name:
# {"status":"success","data":{"name":"20260612T120000Z-..."}}

# Copy snapshot from pod
kubectl cp openshift-monitoring/prometheus-k8s-0:/prometheus/snapshots/20260612T120000Z-... ./prometheus-data
```

**From Must-Gather:**

```bash
# Collect must-gather
oc adm must-gather

# Find Prometheus data
find must-gather.local.* -path "*/prometheus/*/data" -type d

# Note the path for next step
# Example: must-gather.local.5555/monitoring/prometheus/prometheus-k8s-0/data
```

### Step 3: Run Your First Conversion

```bash
# Basic conversion (uses last 24 hours)
prom2sar -tsdb ./prometheus-data

# Or specify exact time range
prom2sar \
  -tsdb ./prometheus-data \
  -start 2026-06-12T00:00:00Z \
  -end 2026-06-12T23:59:59Z \
  -output ./my-analysis \
  -verbose
```

**Output:**

```
=== Prometheus to SAR Conversion ===
TSDB Path:    ./prometheus-data
Output Path:  ./my-analysis
Time Range:   2026-06-12T00:00:00Z to 2026-06-12T23:59:59Z
Interval:     60 seconds
Profile:      all

Opening TSDB at ./prometheus-data...
Found 24 TSDB blocks
Extracting all metrics...
Extracted 1440 data points

✓ Generated summary: ./my-analysis/sar-summary-20260612.txt
✓ Generated report:  ./my-analysis/sar-20260612.txt
✓ Conversion completed successfully!
```

### Step 4: View Results

```bash
# Quick summary
cat my-analysis/sar-summary-20260612.txt

# Full report
cat my-analysis/sar-20260612.txt

# Find CPU spikes
grep "%" my-analysis/sar-20260612.txt | awk '$4 > 80'
```

**Congratulations!** You've successfully converted Prometheus data to SAR format.

## Quick Start with Operator

### Step 1: Install the Operator

```bash
# Clone repository
git clone https://github.com/yourusername/prometheus-dump-operator.git
cd prometheus-dump-operator

# Install CRDs
make install

# Deploy operator
make deploy
```

Verify deployment:

```bash
kubectl get pods -n prometheus-dump-operator

# Expected output:
# NAME                                          READY   STATUS    
# prometheus-dump-operator-xxxxx                1/1     Running
```

### Step 2: Create Your First Conversion

```bash
# Create a conversion job
cat <<EOF | kubectl apply -f -
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: my-first-conversion
  namespace: default
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
EOF
```

### Step 3: Monitor Progress

```bash
# Watch status
kubectl get promethedusdumploader my-first-conversion -w

# View detailed status
kubectl get promethedusdumploader my-first-conversion -o yaml
```

### Step 4: Retrieve SAR Files

```bash
# Get operator pod name
OPERATOR_POD=$(kubectl get pods -n prometheus-dump-operator -l app=prometheus-dump-operator -o jsonpath='{.items[0].metadata.name}')

# Copy SAR files
kubectl cp prometheus-dump-operator/$OPERATOR_POD:/var/lib/sar-output/sar-20260612.txt ./sar-report.txt

# View report
cat sar-report.txt
```

## Your First Conversion

Let's do a practical example: analyzing high CPU during an incident.

### Scenario

Your application had high CPU usage on June 12, 2026 between 14:00 and 16:00 UTC. Let's analyze it.

### Step-by-Step

**1. Get the data:**

```bash
# Copy Prometheus snapshot
kubectl cp openshift-monitoring/prometheus-k8s-0:/prometheus ./prom-data
```

**2. Convert to SAR format:**

```bash
prom2sar \
  -tsdb ./prom-data \
  -start 2026-06-12T14:00:00Z \
  -end 2026-06-12T16:00:00Z \
  -profile cpu \
  -output ./incident-analysis \
  -verbose
```

**3. Analyze CPU spikes:**

```bash
# Find times when CPU user% > 80%
awk '/^[0-9]{2}:[0-9]{2}/ && $4 > 80 {print $1, "User CPU:", $4"%"}' \
  incident-analysis/sar-20260612.txt

# Output:
# 14:23:00 User CPU: 85.3%
# 14:24:00 User CPU: 92.1%
# 14:25:00 User CPU: 88.7%
```

**4. Check other metrics at those times:**

```bash
# Memory usage at 14:23-14:25
sed -n '/14:2[345]/p' incident-analysis/sar-20260612.txt | grep kbmem

# System CPU at those times
sed -n '/14:2[345]/p' incident-analysis/sar-20260612.txt | awk '{print $1, "System:", $6"%"}'
```

**5. Share findings:**

```bash
# Package analysis for team
tar -czf incident-20260612.tar.gz incident-analysis/

# Email or share with kernel engineers
# They can analyze using their existing tools!
```

## Analyzing the Output

### Understanding SAR Format

The output follows standard `sar` format:

```
Linux TSDB  06/12/2026

--------------------------------------------------------------------------------
CPU Utilization (sar -u)
--------------------------------------------------------------------------------

Timestamp    CPU      %user   %nice %system %iowait  %steal   %idle
14:00:00     all       12.50    0.00    3.20    0.50    0.00   83.80
14:01:00     all       15.30    0.00    3.80    0.60    0.00   80.30
```

**Columns:**
- `Timestamp`: Time (HH:MM:SS)
- `CPU`: "all" for aggregated across all CPUs
- `%user`: User space CPU time
- `%system`: Kernel CPU time
- `%iowait`: Waiting for I/O
- `%idle`: Idle time

### Common Analysis Commands

```bash
# Find CPU spikes
awk '/^[0-9]{2}:[0-9]{2}/ && $4 > 80' sar-20260612.txt

# Extract specific time range
sed -n '/14:00/,/15:00/p' sar-20260612.txt

# Memory usage trend
grep kbmemused sar-20260612.txt | awk '{print $1, $3}'

# Disk with highest utilization
awk '/sda|sdb|nvme/ {print $2, $9}' sar-20260612.txt | sort -k2 -nr | head -1

# Network interfaces with errors
awk '$7 > 0 || $8 > 0 {print $1, $2, "RxErr:", $7, "TxErr:", $8}' sar-20260612.txt
```

### Using Existing SAR Scripts

If you have existing scripts that analyze SAR data, they should work as-is:

```bash
# Your existing script
./analyze-sar.sh sar-20260612.txt

# Your analysis pipeline
cat sar-20260612.txt | ./parse-sar.sh | ./generate-report.sh
```

## Next Steps

### Learn More

- **[CLI Guide](CLI_GUIDE.md)** - Complete CLI documentation
- **[SAR Conversion Guide](SAR_CONVERSION_GUIDE.md)** - How metrics are mapped
- **[Quick Reference](QUICK_REFERENCE.md)** - One-page cheat sheet

### Try Different Profiles

```bash
# CPU only (faster)
prom2sar -tsdb ./prom-data -profile cpu

# Memory only
prom2sar -tsdb ./prom-data -profile memory

# Network only
prom2sar -tsdb ./prom-data -profile network
```

### Automate It

Create a script for regular conversions:

```bash
#!/bin/bash
# daily-sar-convert.sh

DATE=$(date -u -d yesterday '+%Y-%m-%d')
OUTPUT="/var/reports/sar-${DATE}"

prom2sar \
  -tsdb /prometheus \
  -start "${DATE}T00:00:00Z" \
  -end "${DATE}T23:59:59Z" \
  -output "$OUTPUT"

echo "Daily SAR report generated: ${OUTPUT}"
```

Add to cron:

```bash
# Run daily at 1 AM
0 1 * * * /usr/local/bin/daily-sar-convert.sh
```

### Join the Community

- **Report issues**: [GitHub Issues](https://github.com/yourusername/prometheus-dump-operator/issues)
- **Ask questions**: [GitHub Discussions](https://github.com/yourusername/prometheus-dump-operator/discussions)
- **Contribute**: See [CONTRIBUTING.md](CONTRIBUTING.md)

## Troubleshooting

### Problem: "TSDB path does not exist"

```bash
# Verify path
ls -la /path/to/prometheus

# Check for TSDB blocks
ls -la /path/to/prometheus/01*/
```

### Problem: "No data found in time range"

```bash
# Check available data with verbose mode
prom2sar -tsdb /prometheus -verbose

# Look for block time ranges
# Adjust your -start and -end to match
```

### Problem: "Permission denied"

```bash
# Check permissions
ls -la /prometheus

# Run with appropriate user
sudo prom2sar -tsdb /prometheus
```

### Still stuck?

- Check [CLI_GUIDE.md](CLI_GUIDE.md#troubleshooting) for more solutions
- Open an issue with the `question` label
- Search existing issues: [Issues](https://github.com/yourusername/prometheus-dump-operator/issues)

## Summary

You've learned:

✅ How to install the CLI binary  
✅ How to get Prometheus TSDB data  
✅ How to run your first conversion  
✅ How to analyze SAR output  
✅ Common analysis commands  

**Next:** Dive deeper with [CLI_GUIDE.md](CLI_GUIDE.md) or try the operator with [README.md](README.md).

---

**Questions?** Open an issue or check our [documentation](CLI_GUIDE.md)!
