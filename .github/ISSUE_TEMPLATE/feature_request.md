---
name: Feature Request
about: Suggest an idea for this project
title: '[FEATURE] '
labels: enhancement
assignees: ''
---

## Is your feature request related to a problem?

A clear and concise description of what the problem is. Ex. I'm always frustrated when [...]

## Proposed Solution

A clear and concise description of what you want to happen.

## Use Case

Describe your use case and how this feature would help. For example:

- **As a** kernel engineer
- **I want to** analyze per-node metrics separately
- **So that** I can identify which specific nodes are having issues

## Example Usage

### CLI Example (if applicable)

```bash
# How you would use this feature
prom2sar -tsdb /prometheus --new-flag value
```

### Operator Example (if applicable)

```yaml
apiVersion: prometheus.openshift.io/v1alpha1
kind: PrometheusDumpLoader
metadata:
  name: example
spec:
  newFeature:
    enabled: true
    option: value
```

## Expected Output

What would the output look like?

```
Example output here
```

## Alternatives Considered

A clear and concise description of any alternative solutions or features you've considered.

**Alternative 1:**
- Description
- Pros/Cons

**Alternative 2:**
- Description  
- Pros/Cons

## Implementation Considerations

If you have technical insights:

- **Complexity**: Low / Medium / High
- **Breaking Changes**: Yes / No
- **Dependencies**: Any new libraries needed?
- **Performance Impact**: Expected impact on processing time

## Additional Context

Add any other context, screenshots, or examples about the feature request here.

## Willingness to Contribute

- [ ] I'm willing to submit a PR to implement this feature
- [ ] I'm willing to help test this feature
- [ ] I'm willing to help with documentation

## Related Issues

Link any related issues or PRs:
- #123
- #456
