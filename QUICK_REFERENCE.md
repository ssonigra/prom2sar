# Quick Reference - For Kernel Engineers

## What is This?

This operator converts **Prometheus monitoring data** into **sar format** so you can analyze it with tools you already know.

## Why Do I Care?

Instead of learning Prometheus/PromQL, you can analyze cluster metrics using:
- `cat`, `grep`, `awk`, `sed`
- Your existing sar analysis scripts
- Standard text viewing tools

## How Do I Use It?

### Step 1: Someone Creates a Conversion Job

An admin or SRE creates this YAML:

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: my-analysis
spec:
  sourcePath: /prometheus
  targetPath: /var/lib/dumps
  
  timeRange:
    start: "2026-06-12T08:00:00Z"  # Incident start
    end: "2026-06-12T10:00:00Z"    # Incident end
  
  sarConversion:
    enabled: true
    outputPath: /var/lib/sar
    metricsProfile: all  # or: cpu, memory, disk, network
```

### Step 2: Get the SAR Files

```bash
# Someone gives you the SAR files, or you copy them:
oc cp <namespace>/<pod>:/var/lib/sar/sar-20260612.txt ./analysis/
```

### Step 3: Analyze Like Normal SAR Output

```bash
# View the full report
cat sar-20260612.txt

# Find peak CPU time
grep "%" sar-20260612.txt | sort -k4 -nr | head -5

# Extract memory during incident
sed -n '/08:00/,/10:00/p' sar-20260612.txt | grep kbmemfree

# Check disk I/O
grep sda sar-20260612.txt

# Network traffic
grep eth0 sar-20260612.txt
```

## What Metrics Are Available?

### CPU (like `sar -u`)
```
Timestamp    CPU      %user   %nice %system %iowait  %steal   %idle
09:00:00     all       12.50    0.00    3.20    0.50    0.00   83.80
```

Fields: %user, %nice, %system, %iowait, %steal, %idle

### Memory (like `sar -r`)
```
Timestamp     kbmemfree    kbmemused   kbbuffers    kbcached  kbswpfree %memused
09:00:00       2048000      6144000      512000     1536000    4096000    75.00
```

Fields: kbmemfree, kbmemused, kbbuffers, kbcached, kbswpfree, %memused

### Disk I/O (like `sar -d`)
```
Timestamp    DEV          tps     rd_sec/s     wr_sec/s  avgrq-sz  avgqu-sz   %util
09:00:00     sda        50.00      1024.00      2048.00     61.44      0.50   45.00
```

Fields: tps, rd_sec/s, wr_sec/s, avgrq-sz, avgqu-sz, %util

### Network (like `sar -n DEV`)
```
Timestamp    IFACE      rxpck/s   txpck/s      rxkB/s      txkB/s  rxerr/s  txerr/s
09:00:00     eth0      1500.00   1200.00     1536.00      768.00     0.00     0.00
```

Fields: rxpck/s, txpck/s, rxkB/s, txkB/s, rxerr/s, txerr/s

## Common Analysis Tasks

### Find CPU Spikes
```bash
# Show times when CPU user% > 80%
awk '/^[0-9]{2}:[0-9]{2}/ && $4 > 80 {print}' sar-20260612.txt
```

### Memory Pressure
```bash
# Show when memory used > 90%
awk '/kbmemfree/ {next} /^[0-9]{2}:[0-9]{2}/ && $8 > 90 {print}' sar-20260612.txt
```

### Disk Bottlenecks
```bash
# Show when disk util > 80%
awk '/DEV/ {next} /^[0-9]{2}:[0-9]{2}/ && $9 > 80 {print}' sar-20260612.txt
```

### Network Issues
```bash
# Show network errors
awk '/IFACE/ {next} /^[0-9]{2}:[0-9]{2}/ && ($7 > 0 || $8 > 0) {print}' sar-20260612.txt
```

## File Structure

After conversion, you get:

```
/var/lib/sar/
├── sar-20260612.txt          # Main report (use this)
└── sar-summary-20260612.txt  # Quick summary
```

**Use the main report (`sar-20260612.txt`) for analysis.**

## Differences from Real SAR

### Same ✅
- Output format identical
- All standard fields present
- Time-based organization
- Human-readable text

### Different ⚠️
- Source is Prometheus, not kernel
- Data is aggregated across nodes (not per-node)
- Historical snapshots only (not real-time)

**For most analysis purposes, these differences don't matter!**

## Tips

1. **Start with the summary**: Check `sar-summary-20260612.txt` first for overview
2. **Use grep liberally**: Standard text, so grep/awk/sed all work
3. **Time format**: Times are in HH:MM:SS format (24-hour)
4. **Averages**: Metrics are averaged across the cluster
5. **Missing data**: If a section is missing, those metrics weren't collected

## Example Workflow

```bash
# 1. Get the files
scp user@server:/var/lib/sar/sar-20260612.txt .

# 2. Quick overview
head -50 sar-20260612.txt

# 3. Find the incident time
grep "14:30" sar-20260612.txt

# 4. Extract 30 minutes around incident
sed -n '/14:15/,/14:45/p' sar-20260612.txt > incident.txt

# 5. Analyze
cat incident.txt
```

## Need Help?

- **Full documentation**: See [SAR_CONVERSION_GUIDE.md](SAR_CONVERSION_GUIDE.md)
- **Examples**: Check [examples/](examples/) directory
- **Testing**: See [TESTING.md](TESTING.md)

## One-Page Cheat Sheet

| Task | Command |
|------|---------|
| View full report | `cat sar-20260612.txt` |
| Find high CPU | `grep "%" sar-*.txt \| awk '$4 > 80'` |
| Memory during time | `sed -n '/09:00/,/10:00/p' sar-*.txt \| grep kbmem` |
| Disk I/O peaks | `grep sda sar-*.txt \| sort -k4 -nr \| head` |
| Network errors | `awk '/eth0/ && ($7 > 0 \|\| $8 > 0)' sar-*.txt` |
| Extract time range | `sed -n '/START/,/END/p' sar-*.txt` |
| Get summary | `cat sar-summary-*.txt` |

## Bottom Line

**You don't need to learn Prometheus.** Just analyze the sar files like you would any other sar output. All your existing knowledge and scripts work!
