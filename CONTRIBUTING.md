# Contributing to Prometheus to SAR Converter

First off, thank you for considering contributing! This project aims to bridge Prometheus monitoring and traditional system administration tools, and your contributions help achieve that goal.

## 🎯 Our Mission

Make Prometheus data accessible to kernel engineers and system administrators using tools they already know—**no Prometheus expertise required**.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Coding Guidelines](#coding-guidelines)
- [Submitting Changes](#submitting-changes)
- [Testing](#testing)
- [Documentation](#documentation)

## 📜 Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inspiring community for all. Please be respectful and constructive in your interactions.

### Our Standards

**Positive behavior includes:**
- Using welcoming and inclusive language
- Being respectful of differing viewpoints
- Accepting constructive criticism gracefully
- Focusing on what's best for the community

**Unacceptable behavior includes:**
- Harassment, trolling, or discriminatory comments
- Publishing others' private information
- Other unprofessional or unethical conduct

## 🤝 How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates.

**When submitting a bug report, include:**

- **Clear title**: Describe the issue concisely
- **Description**: Detailed explanation of the problem
- **Steps to reproduce**: Numbered steps to recreate the issue
- **Expected behavior**: What you expected to happen
- **Actual behavior**: What actually happened
- **Environment**: OS, Go version, Kubernetes version (if applicable)
- **Logs**: Relevant error messages or logs
- **Screenshots**: If applicable

**Example:**

```markdown
## Bug: Conversion fails with time range error

**Description:**
When converting TSDB data with a specific time range, the tool fails with "time range exceeds block range" error.

**Steps to Reproduce:**
1. Run: `prom2sar -tsdb /prometheus -start 2026-06-01T00:00:00Z -end 2026-06-15T00:00:00Z`
2. Observe error

**Expected Behavior:**
Conversion should process all available blocks within the range.

**Actual Behavior:**
Error: "time range exceeds block range"

**Environment:**
- OS: Ubuntu 22.04
- Go: 1.21.3
- prom2sar: v1.0.0

**Logs:**
```
[error logs here]
```
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues.

**When suggesting enhancements, include:**

- **Use case**: Why is this enhancement needed?
- **Current behavior**: What happens now?
- **Proposed behavior**: What should happen?
- **Alternatives**: Other solutions you've considered
- **Additional context**: Screenshots, examples, etc.

**Enhancement ideas:**
- Binary SAR format support
- Additional metric sources
- Performance improvements
- New CLI options
- Operator features
- Documentation improvements

### Contributing Code

We welcome code contributions! Here are some areas:

#### 🐛 Bug Fixes
- Fix reported issues
- Improve error handling
- Edge case handling

#### ✨ New Features
- Binary SAR format (sadc-compatible)
- Per-node metric breakdowns
- Additional Prometheus exporters support
- Performance optimizations
- New CLI flags/options

#### 📚 Documentation
- Improve existing guides
- Add examples
- Fix typos
- Translate documentation
- Create tutorials

#### 🧪 Testing
- Add unit tests
- Integration tests
- Performance benchmarks
- Test coverage improvements

#### 🎨 Code Quality
- Refactoring
- Code organization
- Dependency updates
- CI/CD improvements

## 🛠️ Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- make
- Docker (for operator development)
- Kubernetes cluster (for operator testing, optional)

### Fork and Clone

```bash
# Fork the repository on GitHub, then:

git clone https://github.com/YOUR-USERNAME/prometheus-dump-operator.git
cd prometheus-dump-operator

# Add upstream remote
git remote add upstream https://github.com/ORIGINAL-OWNER/prometheus-dump-operator.git
```

### Build and Test

```bash
# Get dependencies
go mod download

# Build CLI
make build-cli

# Build operator
make build

# Run tests
make test

# Run formatting
make fmt

# Run linting
make vet
```

### Run Locally

```bash
# Run CLI
./bin/prom2sar -tsdb /path/to/prometheus/data -verbose

# Run operator locally (requires kubeconfig)
go run cmd/main.go
```

## 📝 Coding Guidelines

### Go Code Style

Follow standard Go conventions:

```go
// Good: Clear function name, documentation, error handling
// ExtractCPUStats extracts CPU utilization metrics from Prometheus data.
// Returns nil if no CPU metrics are found.
func (m *MetricsMapper) ExtractCPUStats(ctx context.Context, startTime, endTime int64) (*CPUStats, error) {
    if startTime >= endTime {
        return nil, fmt.Errorf("invalid time range: start (%d) >= end (%d)", startTime, endTime)
    }
    
    // Implementation
}

// Bad: No documentation, unclear naming, missing error context
func (m *MetricsMapper) getCPU(s, e int64) (*CPUStats, error) {
    // Implementation
}
```

### Code Organization

```
pkg/
├── tsdb/        # TSDB reading - keep focused on TSDB operations
├── sar/         # SAR conversion - metrics mapping and formatting
├── controller/  # Kubernetes controller logic
└── apis/        # API types and schemas
```

### Comments

```go
// Good: Explains WHY, not WHAT
// We query each CPU mode separately because aggregating across modes
// can lead to race conditions in the TSDB query layer
for _, mode := range cpuModes {
    // ...
}

// Bad: Explains obvious WHAT
// Loop through CPU modes
for _, mode := range cpuModes {
    // ...
}
```

### Error Handling

```go
// Good: Context in errors
if err != nil {
    return fmt.Errorf("failed to query metric %s for time range %d-%d: %w", 
        metricName, startTime, endTime, err)
}

// Bad: Generic errors
if err != nil {
    return err
}
```

### Testing

```go
// Good: Table-driven tests with clear cases
func TestExtractCPUStats(t *testing.T) {
    tests := []struct {
        name      string
        startTime int64
        endTime   int64
        want      *CPUStats
        wantErr   bool
    }{
        {
            name:      "valid time range",
            startTime: 1000,
            endTime:   2000,
            want:      &CPUStats{User: 12.5},
            wantErr:   false,
        },
        {
            name:      "invalid time range",
            startTime: 2000,
            endTime:   1000,
            want:      nil,
            wantErr:   true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := extractCPUStats(tt.startTime, tt.endTime)
            if (err != nil) != tt.wantErr {
                t.Errorf("extractCPUStats() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // ... assertions
        })
    }
}
```

### Formatting

```bash
# Format code
go fmt ./...

# Run goimports (if installed)
goimports -w .

# Or use make
make fmt
```

### Linting

```bash
# Run go vet
go vet ./...

# Or use make
make vet

# Optional: Use golangci-lint for comprehensive checks
golangci-lint run
```

## 🚀 Submitting Changes

### Branch Naming

Use descriptive branch names:

```bash
# Features
git checkout -b feature/binary-sar-format
git checkout -b feature/per-node-metrics

# Bug fixes
git checkout -b fix/time-range-validation
git checkout -b fix/memory-leak-mapper

# Documentation
git checkout -b docs/cli-examples
git checkout -b docs/fix-typos

# Refactoring
git checkout -b refactor/metrics-mapper
```

### Commit Messages

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

**Examples:**

```
feat(sar): add binary SAR format support

Implement sadc-compatible binary output format. This allows
the generated files to be read by standard sar tools.

Closes #42
```

```
fix(tsdb): handle empty TSDB blocks gracefully

Previously, empty blocks caused panic. Now we skip them
with a warning message.

Fixes #38
```

```
docs(cli): add batch processing examples

Add examples showing how to process multiple days of data
in a loop. Useful for historical analysis.
```

### Pull Request Process

1. **Update documentation** if needed
2. **Add tests** for new functionality
3. **Update CHANGELOG.md** (if major change)
4. **Ensure CI passes** (tests, linting)
5. **Request review** from maintainers

### Pull Request Template

When creating a PR, include:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring
- [ ] Performance improvement

## Related Issues
Fixes #123, Related to #456

## Testing
How was this tested?

## Checklist
- [ ] Code follows project style guidelines
- [ ] Comments added for complex logic
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests pass locally
- [ ] No new warnings introduced
```

## 🧪 Testing

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./pkg/sar/

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Verbose output
go test -v ./...

# Run specific test
go test -v ./pkg/sar/ -run TestExtractCPUStats
```

### Writing Tests

**Unit Tests:**

```go
// pkg/sar/mapper_test.go
func TestExtractCPUStats(t *testing.T) {
    // Create test TSDB data
    reader := createTestReader(t)
    mapper := NewMetricsMapper(reader)
    
    stats, err := mapper.extractCPUStats(context.Background(), 1000, 2000)
    
    assert.NoError(t, err)
    assert.NotNil(t, stats)
    assert.Greater(t, stats.User, 0.0)
}
```

**Integration Tests:**

```go
// test/integration/conversion_test.go
func TestFullConversion(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Setup test TSDB
    tmpDir := t.TempDir()
    setupTestTSDB(t, tmpDir)
    
    // Run conversion
    converter := NewConverter(...)
    status, err := converter.Convert(...)
    
    assert.NoError(t, err)
    assert.Greater(t, status.MetricsConverted, 0)
}
```

### Testing CLI

```bash
# Build CLI
make build-cli

# Test basic conversion
./bin/prom2sar -tsdb testdata/prometheus -output /tmp/test

# Verify output
cat /tmp/test/sar-*.txt
```

### Testing Operator

```bash
# Deploy to test cluster
make deploy

# Apply test CR
kubectl apply -f examples/basic-sar-conversion.yaml

# Check status
kubectl get promethedusdumploader -w

# View logs
kubectl logs -n prometheus-dump-operator -l app=prometheus-dump-operator
```

## 📚 Documentation

### Documentation Standards

- **Clear and concise**: Use simple language
- **Examples**: Provide working examples
- **Up-to-date**: Keep in sync with code
- **Tested**: Verify examples work

### Documentation Types

**Code Documentation:**
```go
// MetricsMapper maps Prometheus metrics to SAR data structures.
// It queries the TSDB reader and transforms time series data into
// SAR-compatible format for CPU, memory, disk, and network metrics.
type MetricsMapper struct {
    reader *tsdb.Reader
}
```

**User Documentation:**
- README files
- CLI guides
- Tutorials
- Examples

**API Documentation:**
- CRD documentation
- API reference
- Schema descriptions

### Updating Documentation

When making changes that affect:
- **CLI flags**: Update `CLI_GUIDE.md`
- **Operator CRs**: Update `SAR_CONVERSION_GUIDE.md`
- **New features**: Update main `README.md`
- **Breaking changes**: Update `CHANGELOG.md`

## 🔍 Code Review Process

### For Contributors

- Be open to feedback
- Respond to review comments
- Make requested changes promptly
- Ask questions if unclear

### For Reviewers

- Be respectful and constructive
- Explain the "why" behind suggestions
- Approve when ready
- Use GitHub's suggestion feature for small changes

### Review Checklist

- [ ] Code is clear and maintainable
- [ ] Tests cover new functionality
- [ ] Documentation is updated
- [ ] No unnecessary complexity
- [ ] Error handling is appropriate
- [ ] Performance is acceptable
- [ ] Security considerations addressed

## 🎨 Style Guide

### File Naming

```
good: metrics_mapper.go
bad:  metricsMapper.go, MetricsMapper.go

good: cpu_stats_test.go
bad:  cpuStatsTest.go
```

### Package Naming

```
good: package sar
bad:  package sarConverter
```

### Variable Naming

```go
// Good
var metricsConverted int
var cpuStats *CPUStats

// Bad  
var mc int
var cs *CPUStats
```

## 📊 Performance Guidelines

- Profile before optimizing
- Benchmark significant changes
- Consider memory allocations
- Avoid premature optimization

```bash
# Run benchmarks
go test -bench=. ./pkg/sar/

# Profile
go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./pkg/sar/
go tool pprof cpu.prof
```

## 🏷️ Issue Labels

We use labels to categorize issues:

- `bug`: Something isn't working
- `enhancement`: New feature or request
- `documentation`: Documentation improvements
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention needed
- `question`: Further information requested
- `wontfix`: This will not be worked on
- `duplicate`: Issue already exists

## 🎓 Learning Resources

### Go Resources
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Kubernetes Resources
- [Kubernetes API Conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
- [Operator SDK](https://sdk.operatorframework.io/)

### Project-Specific
- [SAR Format Documentation](SAR_CONVERSION_GUIDE.md)
- [Prometheus TSDB](https://github.com/prometheus/prometheus/tree/main/tsdb)

## 💬 Community

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and community support
- **Pull Requests**: Code contributions

## 🙏 Recognition

Contributors are recognized in:
- `CONTRIBUTORS.md` file
- Release notes
- Project README

Thank you for contributing! Every contribution, no matter how small, is valuable.

---

**Questions?** Open an issue with the `question` label or start a discussion!
