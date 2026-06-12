---
name: Bug Report
about: Create a report to help us improve
title: '[BUG] '
labels: bug
assignees: ''
---

## Bug Description

A clear and concise description of what the bug is.

## Steps to Reproduce

Steps to reproduce the behavior:

1. Run command: `prom2sar -tsdb /prometheus ...`
2. Observe error: '...'
3. Check output: '...'

## Expected Behavior

A clear and concise description of what you expected to happen.

## Actual Behavior

A clear and concise description of what actually happened.

## Environment

**CLI Version (if applicable):**
```bash
prom2sar --version
```

**Operator Version (if applicable):**
```bash
kubectl get deployment -n prometheus-dump-operator
```

**System Information:**
- OS: [e.g., Ubuntu 22.04, RHEL 9]
- Go Version: [e.g., 1.21.3]
- Kubernetes Version: [e.g., 1.28.0] (if using operator)
- OpenShift Version: [e.g., 4.14] (if applicable)

## Logs

<details>
<summary>CLI Output / Operator Logs</summary>

```
Paste logs here
```

</details>

## TSDB Information

**TSDB Path:**
```
/path/to/tsdb
```

**Available Blocks:**
```bash
ls -la /path/to/tsdb/
```

**Time Range:**
- Start: `2026-06-12T00:00:00Z`
- End: `2026-06-12T23:59:59Z`

## Command Used

```bash
prom2sar -tsdb /prometheus -start ... -end ... -output ...
```

Or if using operator:

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: my-conversion
spec:
  ...
```

## Screenshots

If applicable, add screenshots to help explain your problem.

## Additional Context

Add any other context about the problem here. For example:
- Is this a regression? (worked in previous version)
- Does it only happen with specific TSDB data?
- Any workarounds you've found?

## Possible Solution

If you have ideas on how to fix this, please share them here.
