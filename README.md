# Prometheus to SAR Converter

Convert Prometheus TSDB dumps to sar-compatible format for kernel engineers and system administrators.

## Two Ways to Use

### 1. Standalone CLI Tool (`prom2sar`) ⭐ NEW

**Use anywhere - no Kubernetes required!**

```bash
# Build and install
make build-cli
sudo make install-cli

# Convert Prometheus data to SAR format
prom2sar -tsdb /prometheus -output ./sar-analysis
```

Perfect for:
- Analyzing must-gather data on your laptop
- Ad-hoc incident investigations
- Scripts and automation
- Systems without Kubernetes

**[→ CLI Documentation](CLI_README.md)**

### 2. OpenShift Operator

Automated conversions in Kubernetes/OpenShift clusters.

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: sar-conversion
spec:
  sourcePath: /prometheus
  sarConversion:
    enabled: true
```

Perfect for:
- Automated recurring conversions
- Team self-service workflows
- Integration with cluster monitoring

## Overview

This project provides:
1. **TSDB Reading** - Loads Prometheus Time Series Database blocks
2. **SAR Conversion** - Converts metrics to sar-compatible format
3. **Human-Readable Output** - Standard sar format that kernel teams already know

## Key Features

### Dump Loading
- Automatically copies Prometheus dump data from source to target directories
- Supports filtering by time ranges
- Validates dump integrity
- Reports status via CR status field

### SAR Conversion ⭐ NEW
- **Converts Prometheus metrics to sar format** - no Prometheus expertise required
- **Familiar sar output** - CPU, memory, disk, network stats in standard sar format
- **Multiple profiles** - all, cpu, memory, disk, network
- **Configurable intervals** - adjust sampling rate for your needs
- **Text format** - human-readable reports for quick analysis

Perfect for kernel engineers who need to analyze Prometheus data using familiar system tools!

## Prerequisites

- OpenShift 4.x cluster
- Cluster admin privileges for deployment
- Prometheus dumps available at `/prometheus`

## Installation

1. Deploy the CRD:
```bash
oc apply -f deploy/crds/prometheus.openshift.io_promethedusdumploaders_crd.yaml
```

2. Create the operator namespace:
```bash
oc create namespace prometheus-dump-operator
```

3. Deploy the operator:
```bash
oc apply -f deploy/operator.yaml
```

## Quick Start - SAR Conversion

Create a PrometheusDumpLoader CR with sar conversion:

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: prometheus-to-sar
  namespace: default
spec:
  sourcePath: /prometheus
  targetPath: /var/lib/prometheus-dumps
  timeRange:
    start: "2026-06-12T00:00:00Z"
    end: "2026-06-12T23:59:59Z"
  
  # Enable SAR conversion
  sarConversion:
    enabled: true
    outputPath: /var/lib/sar-output
    format: text
    interval: 60
    metricsProfile: all  # Options: all, cpu, memory, disk, network
```

Apply it:
```bash
oc apply -f examples/basic-sar-conversion.yaml
```

Access the sar files:
```bash
oc exec -it <operator-pod> -- cat /var/lib/sar-output/sar-20260612.txt
```

## SAR Output Example

The operator generates standard sar-formatted output:

```
Linux TSDB	06/12/2026

--------------------------------------------------------------------------------
CPU Utilization (sar -u)
--------------------------------------------------------------------------------

Timestamp    CPU      %user   %nice %system %iowait  %steal   %idle
09:00:00     all       12.50    0.00    3.20    0.50    0.00   83.80
09:01:00     all       15.30    0.00    3.80    0.60    0.00   80.30

--------------------------------------------------------------------------------
Memory Utilization (sar -r)
--------------------------------------------------------------------------------

Timestamp     kbmemfree    kbmemused   kbbuffers    kbcached  kbswpfree %memused
09:00:00       2048000      6144000      512000     1536000    4096000    75.00
```

**For kernel engineers**: The output format is identical to standard `sar` command output!

## Usage - Basic Dump Loading (without SAR conversion)

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: my-dump-loader
  namespace: default
spec:
  sourcePath: /prometheus
  targetPath: /var/lib/prometheus-dumps
  timeRange:
    start: "2026-06-01T00:00:00Z"
    end: "2026-06-12T23:59:59Z"
```

## Architecture

- **Controller**: Reconciles PrometheusDumpLoader resources
- **Loader**: Handles file copying and validation
- **Status Reporter**: Updates CR status with progress

## Documentation

- **[SAR Conversion Guide](SAR_CONVERSION_GUIDE.md)** - Comprehensive guide for kernel engineers
- [Examples](examples/) - Sample CRs for different use cases

## Example Use Cases

### For Kernel Engineers
Use sar format to analyze system performance without learning Prometheus:
```bash
# View CPU utilization during incident
grep "14:00" /var/lib/sar-output/sar-20260612.txt

# Extract memory stats
sed -n '/Memory Utilization/,/Disk I\/O/p' /var/lib/sar-output/sar-20260612.txt
```

### For Performance Analysis
- Convert historical Prometheus data to sar format
- Use existing sar-based analysis scripts
- Integrate with performance monitoring workflows
- No Prometheus or PromQL knowledge required

## Development

Build the operator:
```bash
docker build -t prometheus-dump-operator:latest .
```

Run locally:
```bash
go run main.go
```

### Project Structure
```
prometheus-dump-operator/
├── pkg/
│   ├── apis/prometheus/v1alpha1/  # CR definitions with SAR spec
│   ├── controller/                # Reconciler with SAR integration
│   ├── loader/                    # TSDB dump loader
│   ├── tsdb/                      # TSDB block reader
│   └── sar/                       # SAR conversion engine
│       ├── mapper.go              # Prometheus → SAR metrics mapping
│       ├── generator.go           # SAR format output generator
│       └── converter.go           # Conversion orchestration
├── examples/                      # Example CRs
└── SAR_CONVERSION_GUIDE.md       # Detailed SAR usage guide
```
