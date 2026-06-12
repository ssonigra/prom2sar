# Deployment Options

## Two Ways to Convert Prometheus to SAR

```
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│                   Prometheus TSDB Data                              │
│                   (Time Series Database)                            │
│                                                                     │
└────────────────────────┬────────────────────────────────────────────┘
                         │
                         │
        ┌────────────────┴─────────────────┐
        │                                  │
        │                                  │
        ▼                                  ▼
┌───────────────────┐            ┌──────────────────────┐
│   Option 1: CLI   │            │  Option 2: Operator  │
│   Standalone      │            │  Kubernetes          │
└───────────────────┘            └──────────────────────┘
        │                                  │
        │                                  │
        ▼                                  ▼
┌───────────────────┐            ┌──────────────────────┐
│ prom2sar          │            │ PrometheusDumpLoader │
│ (Single Binary)   │            │ (Custom Resource)    │
│                   │            │                      │
│ • No K8s needed   │            │ • Runs in cluster    │
│ • Run anywhere    │            │ • Automated          │
│ • Ad-hoc use      │            │ • Self-service       │
└───────────────────┘            └──────────────────────┘
        │                                  │
        │                                  │
        └────────────────┬─────────────────┘
                         │
                         ▼
             ┌───────────────────────┐
             │   SAR Format Files    │
             │                       │
             │ • sar-YYYYMMDD.txt    │
             │ • sar-summary-*.txt   │
             └───────────────────────┘
                         │
                         │
                         ▼
             ┌───────────────────────┐
             │  Kernel Engineers     │
             │                       │
             │ grep, awk, sed, cat   │
             │ Existing SAR scripts  │
             │ NO Prometheus needed! │
             └───────────────────────┘
```

---

## Option 1: Standalone CLI Binary

### When to Use
- ✅ Analyzing must-gather data
- ✅ Ad-hoc incident investigations
- ✅ Working on laptops/jumphosts
- ✅ No Kubernetes access
- ✅ Quick one-off conversions
- ✅ Scripting and automation
- ✅ Offline analysis

### Installation
```bash
# Build
make build-cli

# Install
sudo make install-cli

# Or just use directly
./bin/prom2sar --version
```

### Usage
```bash
# Basic
prom2sar -tsdb /prometheus

# With options
prom2sar \
  -tsdb /must-gather/prometheus/data \
  -start 2026-06-12T00:00:00Z \
  -end 2026-06-12T23:59:59Z \
  -profile cpu \
  -output ./analysis
```

### Pros
- 🚀 **Fast setup** - Just copy binary
- 🔓 **No dependencies** - Runs standalone
- 💻 **Works anywhere** - Any Linux system
- 📝 **Simple** - Single command
- 🎯 **Flexible** - Command-line control

### Cons
- ⚠️ Manual execution
- ⚠️ No automated scheduling
- ⚠️ Local file access only

### Best For
- Kernel engineers
- SREs doing incident analysis
- Must-gather processing
- Development/testing
- Scripts and automation

---

## Option 2: Kubernetes Operator

### When to Use
- ✅ Automated recurring conversions
- ✅ Team self-service portal
- ✅ Integration with cluster workflows
- ✅ Centralized conversion service
- ✅ Production monitoring
- ✅ Policy-based automation

### Installation
```bash
# Install CRDs
make install

# Deploy operator
make deploy
```

### Usage
```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: daily-sar-conversion
spec:
  sourcePath: /prometheus
  targetPath: /var/lib/dumps
  
  sarConversion:
    enabled: true
    outputPath: /var/lib/sar
    format: text
    interval: 60
    metricsProfile: all
```

```bash
oc apply -f conversion-job.yaml
```

### Pros
- ⚙️ **Automated** - Runs on schedule
- 🔄 **Declarative** - CRs define state
- 👥 **Multi-user** - Team access
- 🔐 **RBAC integration** - Cluster permissions
- 📊 **Status tracking** - Built-in monitoring

### Cons
- ⚠️ Requires Kubernetes
- ⚠️ More complex setup
- ⚠️ Cluster access needed

### Best For
- Production environments
- Operations teams
- Automated workflows
- Shared services
- Compliance/auditing

---

## Comparison Matrix

| Feature | CLI Binary | Operator |
|---------|-----------|----------|
| **Setup Time** | 1 minute | 10-15 minutes |
| **Requirements** | Binary file only | K8s cluster |
| **Execution** | Manual command | Automated CR |
| **Scheduling** | Via cron/scripts | Built-in |
| **Access Control** | File permissions | K8s RBAC |
| **Data Source** | Local filesystem | PVCs, mounts |
| **Output Location** | Local files | PVCs, ConfigMaps |
| **Multi-user** | No | Yes |
| **Portability** | 100% portable | Cluster-specific |
| **Use Cases** | Ad-hoc, scripts | Production, automation |
| **Learning Curve** | Minutes | Hours |
| **Best For** | Engineers | Operations |

---

## Use Case Scenarios

### Scenario 1: Must-Gather Analysis
**→ Use CLI Binary**

```bash
# Download must-gather
oc adm must-gather

# Convert locally
prom2sar -tsdb must-gather.local.*/monitoring/prometheus/*/data \
  -output ./analysis

# Analyze
cat analysis/sar-*.txt
```

**Why CLI?** One-time analysis, offline, no cluster needed.

---

### Scenario 2: Daily Performance Reports
**→ Use Operator**

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: daily-report
spec:
  sourcePath: /prometheus
  sarConversion:
    enabled: true
    outputPath: /var/reports
```

Add CronJob to create CR daily.

**Why Operator?** Automated, recurring, team access.

---

### Scenario 3: Incident Response
**→ Use CLI Binary**

```bash
# During incident, quick conversion
prom2sar -tsdb /prometheus \
  -start $(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ') \
  -end $(date -u '+%Y-%m-%dT%H:%M:%SZ') \
  -summary
```

**Why CLI?** Fast, flexible, no CR approval needed.

---

### Scenario 4: Customer Support
**→ Use CLI Binary**

```bash
# Customer provides Prometheus snapshot
scp customer@server:/prometheus-snapshot.tar.gz .
tar -xzf prometheus-snapshot.tar.gz

# Convert for analysis
prom2sar -tsdb ./prometheus-snapshot -output ./customer-analysis

# Share SAR files with kernel team
tar -czf analysis.tar.gz customer-analysis/
```

**Why CLI?** External data, offline, portable.

---

### Scenario 5: Continuous Monitoring
**→ Use Operator**

Deploy operator in monitoring namespace, create CRs for:
- Hourly CPU metrics
- Daily full reports
- Weekly summaries

**Why Operator?** Automated, consistent, team visibility.

---

## Hybrid Approach (Recommended)

**Use Both!**

```
┌─────────────────────────────────────────────┐
│                                             │
│  Production Cluster                         │
│                                             │
│  • Operator for automated conversions      │
│  • Scheduled daily/weekly reports          │
│  • Team self-service                       │
│                                             │
└─────────────────────────────────────────────┘

             ↓ Export data when needed

┌─────────────────────────────────────────────┐
│                                             │
│  Local Systems / Jumphosts                  │
│                                             │
│  • CLI for ad-hoc analysis                 │
│  • Must-gather processing                  │
│  • Incident investigations                 │
│  • Customer data analysis                  │
│                                             │
└─────────────────────────────────────────────┘
```

---

## Getting Started

### Start with CLI
1. Build: `make build-cli`
2. Test: `./bin/prom2sar -tsdb /test/data`
3. Use: Analyze real incidents

### Graduate to Operator
1. Deploy: `make install && make deploy`
2. Test: Create sample CR
3. Automate: Schedule regular conversions

---

## Distribution Recommendations

### For Kernel Engineering Team
**→ Provide CLI Binary**

```bash
# Create distribution package
./build.sh
cd dist/
tar -czf prom2sar-toolkit.tar.gz \
  prom2sar-linux-amd64-static \
  ../QUICK_REFERENCE.md \
  ../CLI_GUIDE.md

# Distribute
# Team extracts and uses prom2sar directly
```

### For Operations Team
**→ Deploy Operator**

```bash
# Deploy to cluster
make install deploy

# Provide examples
cp examples/*.yaml /team-docs/

# Train on CR creation
```

### For Both Teams
**→ Provide Documentation**
- CLI_README.md for CLI users
- SAR_CONVERSION_GUIDE.md for SAR format
- QUICK_REFERENCE.md for kernel engineers

---

## Summary

**You have options!**

| Need | Solution |
|------|----------|
| Quick analysis | CLI binary |
| Must-gather processing | CLI binary |
| Production automation | Operator |
| Team self-service | Operator |
| Scripts/cron jobs | CLI binary |
| Cluster integration | Operator |
| Offline analysis | CLI binary |
| Continuous monitoring | Operator |

**Both tools produce identical SAR output.** Choose based on your workflow!

---

## Files Reference

### CLI Binary
- Source: `cmd/prom2sar/main.go`
- Build: `make build-cli`
- Docs: `CLI_README.md`, `CLI_GUIDE.md`

### Operator
- Source: `cmd/main.go`, `pkg/`
- Deploy: `make install deploy`
- Docs: `README.md`, `SAR_CONVERSION_GUIDE.md`

### Shared
- TSDB Reader: `pkg/tsdb/`
- SAR Conversion: `pkg/sar/`
- Documentation: `QUICK_REFERENCE.md`
