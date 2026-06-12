# SAR Conversion Guide

## Overview

The Prometheus Dump Loader Operator can convert Prometheus TSDB (Time Series Database) dumps into sar-compatible format, enabling kernel engineers and system administrators to analyze Prometheus metrics using familiar `sar` tools and workflows.

## Why SAR Conversion?

- **Familiar Interface**: Kernel engineers can use `sar` commands they already know
- **No Prometheus Knowledge Required**: Team members don't need to learn PromQL or Prometheus internals
- **Standard Analysis Tools**: Leverage existing sar-based analysis scripts and workflows
- **Text Format**: Human-readable output for quick inspection and debugging

## How It Works

1. **TSDB Reading**: The operator reads Prometheus TSDB blocks from the specified path
2. **Metric Mapping**: Prometheus node-exporter metrics are mapped to sar-equivalent fields
3. **Data Extraction**: Time series data is extracted and aggregated at specified intervals
4. **SAR Generation**: Output files are generated in sar-compatible text format

## Supported Metrics

### CPU Metrics (sar -u equivalent)
- `%user`: User CPU time
- `%nice`: Nice CPU time  
- `%system`: System CPU time
- `%iowait`: I/O wait time
- `%steal`: Steal time (virtualized environments)
- `%idle`: Idle time

**Source**: `node_cpu_seconds_total` with mode labels

### Memory Metrics (sar -r equivalent)
- `kbmemfree`: Free memory in KB
- `kbmemused`: Used memory in KB
- `kbbuffers`: Buffer cache in KB
- `kbcached`: Page cache in KB
- `kbswpfree`: Free swap in KB
- `%memused`: Memory utilization percentage

**Source**: `node_memory_*_bytes` metrics

### Disk I/O Metrics (sar -d equivalent)
- `tps`: Transfers per second
- `rd_sec/s`: Read sectors per second (KB/s * 2)
- `wr_sec/s`: Write sectors per second (KB/s * 2)
- `avgrq-sz`: Average request size
- `%util`: Device utilization percentage

**Source**: `node_disk_*` metrics

### Network Metrics (sar -n DEV equivalent)
- `rxpck/s`: Received packets per second
- `txpck/s`: Transmitted packets per second
- `rxkB/s`: Received kilobytes per second
- `txkB/s`: Transmitted kilobytes per second
- `rxerr/s`: Receive errors per second
- `txerr/s`: Transmit errors per second

**Source**: `node_network_*` metrics

## Usage Examples

### Basic Usage - All Metrics

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: prometheus-to-sar
spec:
  sourcePath: /prometheus
  targetPath: /var/lib/prometheus-dumps
  
  sarConversion:
    enabled: true
    outputPath: /var/lib/sar-output
    format: text
    interval: 60
    metricsProfile: all
```

Apply the CR:
```bash
oc apply -f examples/basic-sar-conversion.yaml
```

### CPU-Only Analysis

For focused CPU performance analysis:

```yaml
sarConversion:
  enabled: true
  outputPath: /var/lib/sar-output/cpu
  format: text
  interval: 30  # Higher resolution for CPU spikes
  metricsProfile: cpu
```

### Memory-Only Analysis

```yaml
sarConversion:
  enabled: true
  outputPath: /var/lib/sar-output/memory
  format: text
  interval: 60
  metricsProfile: memory
```

### Time Range Filtering

```yaml
spec:
  sourcePath: /prometheus
  targetPath: /var/lib/prometheus-dumps
  
  timeRange:
    start: "2026-06-12T08:00:00Z"
    end: "2026-06-12T18:00:00Z"  # Business hours only
  
  sarConversion:
    enabled: true
    outputPath: /var/lib/sar-output
    format: text
    interval: 60
    metricsProfile: all
```

## Output Files

The operator generates the following files in the output directory:

### 1. Main SAR Report (`sar-YYYYMMDD.txt`)

Contains detailed sar-formatted output with all requested metrics:

```
Linux TSDB	06/12/2026

--------------------------------------------------------------------------------
CPU Utilization (sar -u)
--------------------------------------------------------------------------------

Timestamp    CPU      %user   %nice %system %iowait  %steal   %idle
09:00:00     all       12.50    0.00    3.20    0.50    0.00   83.80
09:01:00     all       15.30    0.00    3.80    0.60    0.00   80.30
...

--------------------------------------------------------------------------------
Memory Utilization (sar -r)
--------------------------------------------------------------------------------

Timestamp     kbmemfree    kbmemused   kbbuffers    kbcached  kbswpfree %memused
09:00:00       2048000      6144000      512000     1536000    4096000    75.00
...
```

### 2. Summary Report (`sar-summary-YYYYMMDD.txt`)

Contains summary statistics and quick overview:

```
=== Prometheus to SAR Conversion Summary ===

Hostname: prometheus-tsdb
Total samples: 1440
Time range: 2026-06-12T00:00:00Z to 2026-06-12T23:59:59Z
Duration: 23h59m0s

✓ CPU statistics available
✓ Memory statistics available
✓ Disk I/O statistics available
✓ Network statistics available

=== Quick Summary ===
Samples: 1440

CPU (avg):
  User:    12.50%
  System:   3.20%
  IOWait:   0.50%
  Idle:    83.80%

Memory (avg):
  Used:    75.00% (   6144000 KB)
  Free:       2048000 KB
  Cached:     1536000 KB
```

## Reading SAR Output

### For Kernel Engineers

The output format is identical to standard `sar` command output. You can:

1. **View the main report**:
   ```bash
   cat /var/lib/sar-output/sar-20260612.txt
   ```

2. **Grep for specific times**:
   ```bash
   grep "09:00" /var/lib/sar-output/sar-20260612.txt
   ```

3. **Extract CPU data only**:
   ```bash
   sed -n '/CPU Utilization/,/Memory Utilization/p' /var/lib/sar-output/sar-20260612.txt
   ```

4. **Use existing sar analysis scripts**:
   Most scripts expecting sar text output will work with minimal or no modification.

## Monitoring Conversion Status

Check the status of your conversion:

```bash
oc get promethedusdumploader prometheus-to-sar -o yaml
```

Look for the `sarConversionStatus` field:

```yaml
status:
  phase: Completed
  sarConversionStatus:
    phase: Completed
    metricsConverted: 1440
    sarFilesGenerated: 2
    outputLocation: /var/lib/sar-output
    timeRangeCovered:
      start: "2026-06-12T00:00:00Z"
      end: "2026-06-12T23:59:59Z"
```

## Troubleshooting

### No Data in SAR Files

**Problem**: SAR files are generated but contain no metrics

**Solutions**:
1. Verify Prometheus TSDB blocks exist at `sourcePath`
2. Check that time range overlaps with available data
3. Ensure node-exporter metrics are present in the TSDB
4. Verify the operator has read access to the TSDB path

### Missing Metrics

**Problem**: Some metric categories are missing (e.g., no disk stats)

**Solutions**:
1. Check if node-exporter collected those metrics
2. Verify label selectors if using custom metrics
3. Some nodes may not have certain devices (e.g., swap disabled)

### Time Range Issues

**Problem**: Conversion fails with time range errors

**Solutions**:
1. Ensure start time is before end time
2. Check that requested time range has data in TSDB
3. Use `oc logs` to see operator logs for details

## Accessing SAR Files from Operator Pod

The SAR files are written to the output path specified in the CR. To access them:

1. **Find the operator pod**:
   ```bash
   oc get pods -n prometheus-dump-operator
   ```

2. **Copy files locally**:
   ```bash
   oc cp prometheus-dump-operator/<pod-name>:/var/lib/sar-output/sar-20260612.txt ./sar-report.txt
   ```

3. **Or exec into the pod**:
   ```bash
   oc exec -it <pod-name> -n prometheus-dump-operator -- cat /var/lib/sar-output/sar-20260612.txt
   ```

## Advanced: Custom Metrics Mapping

For specialized metrics not covered by standard profiles:

```yaml
sarConversion:
  enabled: true
  outputPath: /var/lib/sar-output/custom
  metricsProfile: custom
  
  customMetrics:
    - prometheusMetric: "container_cpu_usage_seconds_total"
      sarField: "custom_container_cpu"
      labelSelector:
        namespace: "monitoring"
      aggregation: "avg"
```

This feature is for advanced users who understand Prometheus metrics and want to create custom sar-like reports.

## Performance Considerations

- **Interval**: Lower intervals (30s, 60s) provide more detail but larger files
- **Time Range**: Longer ranges produce more data; consider splitting into days
- **Metrics Profile**: Use specific profiles (cpu, memory) instead of "all" when possible
- **TSDB Size**: Large TSDB blocks may take time to process

## Best Practices

1. **Start Small**: Begin with a short time range (1-2 hours) to verify output
2. **Use Specific Profiles**: Only extract metrics you need
3. **Reasonable Intervals**: 60-second intervals are usually sufficient
4. **Archive SAR Files**: Keep SAR files for historical analysis
5. **Automate**: Create CRs programmatically for regular conversions

## Integration with Existing Workflows

The SAR text format allows integration with:

- **Custom analysis scripts** expecting sar output
- **Performance monitoring dashboards** that parse sar files
- **Automated alerting** based on sar threshold checks
- **Capacity planning tools** using historical sar data

## Limitations

- **Text Format Only**: Binary sar format not currently supported
- **Node-Exporter Required**: Metrics must come from Prometheus node-exporter
- **Aggregation**: Multi-node metrics are averaged (no per-node breakdown in current version)
- **Historical Data Only**: Reads from TSDB snapshots, not live Prometheus queries

## Future Enhancements

Planned features:
- Binary sar format support (sadc-compatible)
- Per-node breakdowns
- Additional metric sources beyond node-exporter
- Real-time conversion from live Prometheus
