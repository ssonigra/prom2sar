## Description

<!-- Provide a brief description of the changes in this PR -->

## Type of Change

<!-- Check all that apply -->

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring
- [ ] Performance improvement
- [ ] Test addition/update

## Related Issues

<!-- Link related issues here -->

Fixes #(issue number)
Related to #(issue number)

## Changes Made

<!-- Describe the changes in detail -->

- Change 1
- Change 2
- Change 3

## Testing Performed

<!-- Describe the testing you've done -->

### Unit Tests

```bash
# Commands used for testing
make test
go test ./pkg/sar/ -v
```

### Integration Tests

<!-- If applicable -->

```bash
# How you tested the full flow
prom2sar -tsdb testdata/prometheus -output /tmp/test
```

### Manual Testing

<!-- Describe manual testing steps -->

1. Built CLI: `make build-cli`
2. Tested with command: `prom2sar -tsdb /prometheus`
3. Verified output: `cat sar-output/sar-*.txt`
4. Results: Working as expected

## Output/Screenshots

<!-- If applicable, add screenshots or sample output -->

<details>
<summary>Sample Output</summary>

```
Paste output here
```

</details>

## Documentation

<!-- Check all that apply -->

- [ ] Documentation updated (if needed)
- [ ] README updated (if needed)
- [ ] CHANGELOG.md updated (if significant change)
- [ ] Code comments added for complex logic
- [ ] Examples added/updated (if applicable)

## Checklist

<!-- Ensure all items are completed -->

- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## Performance Impact

<!-- Describe any performance implications -->

- [ ] No performance impact
- [ ] Performance improved
- [ ] Performance degraded (with justification)

**Details:**

## Breaking Changes

<!-- If this is a breaking change, describe migration steps -->

**Migration Guide:**

1. Step 1
2. Step 2

## Additional Notes

<!-- Any additional information for reviewers -->

## Reviewer Checklist

<!-- For maintainers -->

- [ ] Code quality is acceptable
- [ ] Tests are adequate
- [ ] Documentation is sufficient
- [ ] No security concerns
- [ ] Ready to merge
