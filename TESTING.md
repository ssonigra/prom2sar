# Testing Guide

## Prerequisites

1. **OpenShift cluster** with cluster-admin access
2. **Prometheus TSDB blocks** available (from a Prometheus snapshot)
3. **kubectl/oc CLI** configured

## Quick Test

### 1. Build and Deploy the Operator

```bash
# Build the operator image
make docker-build IMG=quay.io/<your-repo>/prometheus-dump-operator:latest

# Push to registry
make docker-push IMG=quay.io/<your-repo>/prometheus-dump-operator:latest

# Install CRDs
make install

# Deploy operator
make deploy
```

### 2. Prepare Test Data

If you don't have Prometheus TSDB blocks, you can get them from:

**Option A: From a running Prometheus**
```bash
# Port-forward to Prometheus
oc port-forward -n openshift-monitoring prometheus-k8s-0 9090:9090

# Create a snapshot
curl -XPOST http://localhost:9090/api/v1/admin/tsdb/snapshot

# The snapshot is created in /prometheus/snapshots/<timestamp>
```

**Option B: From must-gather**
```bash
# Extract must-gather
oc adm must-gather

# Prometheus data is typically in:
# must-gather.local.*/monitoring/prometheus/*/data
```

### 3. Create a PV/PVC for TSDB Data

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: prometheus-tsdb-pvc
  namespace: default
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

Copy your TSDB blocks to this PVC:
```bash
# Create a temporary pod to copy data
oc run tsdb-loader --image=busybox --command -- sleep 3600
oc set volume pod/tsdb-loader --add --name=tsdb --claim-name=prometheus-tsdb-pvc --mount-path=/tsdb

# Copy TSDB blocks
oc cp /path/to/your/tsdb/blocks tsdb-loader:/tsdb/
```

### 4. Create a Test CR

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: test-sar-conversion
  namespace: default
spec:
  sourcePath: /tsdb
  targetPath: /tmp/prometheus-dumps
  
  timeRange:
    start: "2026-06-12T00:00:00Z"
    end: "2026-06-12T12:00:00Z"
  
  sarConversion:
    enabled: true
    outputPath: /tmp/sar-output
    format: text
    interval: 60
    metricsProfile: all
```

Apply:
```bash
oc apply -f test-cr.yaml
```

### 5. Monitor Progress

```bash
# Watch the CR status
oc get promethedusdumploader test-sar-conversion -w

# Check operator logs
make logs

# Detailed status
oc get promethedusdumploader test-sar-conversion -o yaml
```

Expected status:
```yaml
status:
  phase: Completed
  message: Successfully copied 15 files (1234567 bytes)
  filesCopied: 15
  bytesCopied: 1234567
  sarConversionStatus:
    phase: Completed
    metricsConverted: 720
    sarFilesGenerated: 2
    outputLocation: /tmp/sar-output
    timeRangeCovered:
      start: "2026-06-12T00:00:00Z"
      end: "2026-06-12T12:00:00Z"
```

### 6. Retrieve SAR Files

```bash
# Get operator pod name
OPERATOR_POD=$(oc get pods -n prometheus-dump-operator -l app=prometheus-dump-operator -o jsonpath='{.items[0].metadata.name}')

# Copy SAR report
oc cp prometheus-dump-operator/$OPERATOR_POD:/tmp/sar-output/sar-20260612.txt ./sar-report.txt

# Copy summary
oc cp prometheus-dump-operator/$OPERATOR_POD:/tmp/sar-output/sar-summary-20260612.txt ./sar-summary.txt

# View the report
cat sar-report.txt
```

### 7. Verify Output

The sar report should look like:

```
Linux TSDB	06/12/2026

--------------------------------------------------------------------------------
CPU Utilization (sar -u)
--------------------------------------------------------------------------------

Timestamp    CPU      %user   %nice %system %iowait  %steal   %idle
00:00:00     all       12.50    0.00    3.20    0.50    0.00   83.80
00:01:00     all       15.30    0.00    3.80    0.60    0.00   80.30
...
```

## Testing Different Profiles

### CPU-Only Test

```yaml
spec:
  sarConversion:
    enabled: true
    outputPath: /tmp/sar-output/cpu
    format: text
    interval: 30
    metricsProfile: cpu
```

Verify only CPU section appears in output.

### Memory-Only Test

```yaml
spec:
  sarConversion:
    enabled: true
    outputPath: /tmp/sar-output/memory
    format: text
    interval: 60
    metricsProfile: memory
```

Verify only Memory section appears in output.

## Testing Edge Cases

### Empty Time Range

```yaml
timeRange:
  start: "2026-06-12T00:00:00Z"
  end: "2026-06-12T00:01:00Z"
```

Should generate minimal data or gracefully handle no data.

### Future Time Range

```yaml
timeRange:
  start: "2027-01-01T00:00:00Z"
  end: "2027-01-02T00:00:00Z"
```

Should report no data available.

### Invalid TSDB Path

```yaml
spec:
  sourcePath: /nonexistent
```

Should fail with clear error message.

## Manual Testing with Real Data

If you have access to a cluster with real Prometheus data:

```bash
# Create CR pointing to live Prometheus data
cat <<EOF | oc apply -f -
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: live-prometheus-sar
  namespace: openshift-monitoring
spec:
  sourcePath: /prometheus
  targetPath: /tmp/prometheus-dumps
  
  timeRange:
    start: "$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')"
    end: "$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  
  sarConversion:
    enabled: true
    outputPath: /tmp/sar-output
    format: text
    interval: 60
    metricsProfile: all
EOF
```

## Validation Checklist

- [ ] Operator deploys successfully
- [ ] CRD is registered
- [ ] CR is accepted
- [ ] TSDB blocks are read
- [ ] Metrics are extracted
- [ ] SAR files are generated
- [ ] SAR format matches standard sar output
- [ ] CPU metrics are accurate
- [ ] Memory metrics are accurate
- [ ] Disk metrics are present
- [ ] Network metrics are present
- [ ] Time range filtering works
- [ ] Different profiles work (cpu, memory, disk, network, all)
- [ ] Status updates correctly
- [ ] Error handling works

## Troubleshooting Tests

### No SAR Files Generated

```bash
# Check operator logs
oc logs -n prometheus-dump-operator -l app=prometheus-dump-operator

# Check CR status
oc describe promethedusdumploader test-sar-conversion

# Verify TSDB path
oc exec <operator-pod> -- ls -la /tsdb
```

### SAR Files Empty

```bash
# Verify time range has data
oc exec <operator-pod> -- ls -la /tsdb/01*/

# Check TSDB block metadata
oc exec <operator-pod> -- cat /tsdb/01*/meta.json
```

### Metrics Missing

```bash
# List available metrics in TSDB
# This requires promtool
oc exec <operator-pod> -- promtool tsdb dump /tsdb | head -100

# Verify node-exporter metrics exist
oc exec <operator-pod> -- promtool tsdb dump /tsdb | grep node_cpu
```

## Performance Testing

### Large Time Range

```yaml
timeRange:
  start: "2026-06-01T00:00:00Z"
  end: "2026-06-12T23:59:59Z"  # 12 days
```

Monitor:
- Conversion time
- Memory usage
- Output file size

### High-Frequency Sampling

```yaml
sarConversion:
  interval: 10  # 10 second intervals
```

Monitor:
- Processing time
- File size
- Data accuracy

## Cleanup

```bash
# Delete test CRs
oc delete promethedusdumploader --all

# Remove operator
make undeploy

# Remove CRDs
make uninstall
```

## CI/CD Integration

Example GitHub Actions workflow:

```yaml
name: Test Operator
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Run tests
        run: make test
      
      - name: Build
        run: make build
```
