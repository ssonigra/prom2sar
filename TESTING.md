# Testing Guide for prom2sar

Complete guide for testing the prom2sar CLI tool.

## Quick Start

```bash
# 1. Build the binary
make build-cli

# 2. Quick smoke test
./bin/prom2sar --version
./bin/prom2sar --help

# 3. Test with real data
./bin/prom2sar -tsdb /path/to/prometheus -output ./test-results -verbose
```

## Test with Docker (Easiest Method)

```bash
# Create test environment
mkdir -p /tmp/prom-test && cd /tmp/prom-test

# Create docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3'
services:
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    ports:
      - "9090:9090"
  node-exporter:
    image: prom/node-exporter:latest
    ports:
      - "9100:9100"
volumes:
  prometheus-data:
EOF

# Create prometheus config
cat > prometheus.yml << 'EOF'
global:
  scrape_interval: 15s
scrape_configs:
  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']
EOF

# Start and collect data
docker-compose up -d
echo "Waiting 5 minutes for data collection..."
sleep 300

# Extract data
VOLUME=$(docker volume ls -q | grep prom-test_prometheus-data)
docker run --rm -v ${VOLUME}:/source -v $(pwd)/data:/target alpine cp -r /source/* /target/

# Stop containers
docker-compose down

# Test prom2sar
cd /home/ssonigra/reprod/claude/prometheus-dump-operator
./bin/prom2sar -tsdb /tmp/prom-test/data -output ./test-results -verbose

# View results
cat test-results/sar-summary-*.txt
```

## Test Different Profiles

```bash
# CPU only
./bin/prom2sar -tsdb /path/to/prometheus -profile cpu -output ./cpu-test

# Memory only
./bin/prom2sar -tsdb /path/to/prometheus -profile memory -output ./memory-test

# Disk I/O
./bin/prom2sar -tsdb /path/to/prometheus -profile disk -output ./disk-test

# Network I/O
./bin/prom2sar -tsdb /path/to/prometheus -profile network -output ./network-test

# All metrics (default)
./bin/prom2sar -tsdb /path/to/prometheus -profile all -output ./all-test
```

## Test Time Ranges

```bash
# Last 24 hours (default)
./bin/prom2sar -tsdb /path/to/prometheus -output ./last-24h

# Specific range
./bin/prom2sar \
  -tsdb /path/to/prometheus \
  -start 2026-06-12T00:00:00Z \
  -end 2026-06-12T23:59:59Z \
  -output ./specific-range

# Last hour
START=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')
END=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
./bin/prom2sar -tsdb /path/to/prometheus -start "$START" -end "$END" -output ./last-hour
```

## Analyze Output

```bash
# View summary
cat test-results/sar-summary-*.txt

# Find CPU spikes
awk '/CPU/ && $2 > 80 {print $0}' test-results/sar-*.txt

# Find memory pressure
grep -A 5 "Memory" test-results/sar-*.txt | awk '$2 < 100000 {print}'

# Find disk bottlenecks
grep -A 10 "Disk" test-results/sar-*.txt | awk '$NF > 90 {print}'

# Grep specific time
grep "12:00:00" test-results/sar-*.txt
```

## Expected Output Format

### Summary
```
=== Prometheus to SAR Conversion Summary ===
Source: prometheus-tsdb
Time Range: 2026-06-11 12:00:00 to 2026-06-12 12:00:00
Duration: 24h0m0s
Data Points: 1440

CPU:
  Average User:   15.3%
  Average System: 8.2%
  Average Idle:   74.4%

Memory:
  Total:     16384 MB
  Used Avg:  8192 MB (50.0%)
```

### CPU Report
```
12:00:00    CPU    %user  %nice  %system  %iowait  %idle
12:01:00    all    15.32   0.00     8.21     2.14  74.33
```

## Validation Tests

```bash
# Test error handling
./bin/prom2sar                        # Should error: tsdb required
./bin/prom2sar -tsdb /nonexistent     # Should error: path not exists
./bin/prom2sar -tsdb /tmp -start bad  # Should error: invalid time

# Test summary only
./bin/prom2sar -tsdb /path/to/prometheus -summary -output ./summary-only
ls summary-only/  # Should only have summary file
```

## Real-World Example

```bash
# Investigate incident on June 12, 10am-12pm
./bin/prom2sar \
  -tsdb /var/lib/prometheus/data \
  -start 2026-06-12T10:00:00Z \
  -end 2026-06-12T12:00:00Z \
  -output ./incident-analysis \
  -profile all \
  -verbose

cd incident-analysis

# Check summary
cat sar-summary-*.txt

# Find issues
grep -A 2 "CPU" sar-*.txt | grep -v "all" | awk '{if ($2 > 80) print}'
grep "10:45" sar-*.txt | grep -A 5 "Memory"
```

## Success Criteria

- ✅ Binary runs without crashes
- ✅ Output files created in specified directory
- ✅ SAR format matches standard sar output
- ✅ Data values are reasonable
- ✅ Time ranges match request
- ✅ Profile filtering works
- ✅ Output works with grep/awk/sed
- ✅ Verbose mode shows progress

## Troubleshooting

```bash
# Check TSDB structure
ls -la /path/to/prometheus/
# Should see: 01ABCDEF123456/ (block dirs), chunks_head/, wal/

# Verify blocks
find /path/to/prometheus -name "meta.json" | head -3

# Debug mode
./bin/prom2sar -tsdb /path/to/prometheus -verbose 2>&1 | tee debug.log
```

## Performance Testing

```bash
# 1 hour of data
time ./bin/prom2sar -tsdb /path/to/prometheus -start "$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')" -output ./perf-1h

# 24 hours of data
time ./bin/prom2sar -tsdb /path/to/prometheus -output ./perf-24h
```

## Need Help?

- Use `-verbose` flag for debugging
- See `BUILD_SUCCESS.md` for build info
- See `CLI_GUIDE.md` for detailed usage
- Issues: https://github.com/ssonigra/prom2sar/issues
