## GitHub Repository Setup Guide

This guide walks you through creating a GitHub repository for the Prometheus to SAR Converter project and setting it up for community contributions.

## Table of Contents

1. [Quick Setup](#quick-setup)
2. [Detailed Setup](#detailed-setup)
3. [Repository Configuration](#repository-configuration)
4. [First Release](#first-release)
5. [Community Guidelines](#community-guidelines)
6. [Maintenance](#maintenance)

---

## Quick Setup

### Option 1: Using the Script (Recommended)

```bash
# Run the initialization script
./init-github-repo.sh

# Follow the on-screen instructions
```

### Option 2: Manual Setup

```bash
# 1. Initialize git
git init

# 2. Add all files
git add .

# 3. Create initial commit
git commit -m "Initial commit: Prometheus to SAR Converter"

# 4. Create GitHub repo (do this on GitHub.com)
# Then add remote:
git remote add origin https://github.com/YOUR-USERNAME/prometheus-dump-operator.git

# 5. Push to GitHub
git branch -M main
git push -u origin main
```

---

## Detailed Setup

### Step 1: Prepare the Repository

Before creating the GitHub repository, ensure all files are ready:

```bash
# Check that key files exist
ls -la README.md
ls -la GITHUB_README.md
ls -la LICENSE
ls -la CONTRIBUTING.md
ls -la .github/

# Build binaries to verify everything works
make build-cli
make build

# Run tests
make test
```

### Step 2: Create GitHub Repository

1. **Go to GitHub** and create a new repository:
   - URL: https://github.com/new
   - **Repository name**: `prometheus-dump-operator`
   - **Description**: `Convert Prometheus TSDB to SAR format - Enable kernel engineers to analyze Prometheus metrics using familiar Unix tools`
   - **Visibility**: Public (recommended for open source)
   - **DO NOT** check any initialization options:
     - ❌ Add README
     - ❌ Add .gitignore
     - ❌ Choose license
   
   (We already have these files!)

2. **Click "Create repository"**

### Step 3: Push to GitHub

```bash
# Add the remote
git remote add origin https://github.com/YOUR-USERNAME/prometheus-dump-operator.git

# Rename branch to main (if not already)
git branch -M main

# Push
git push -u origin main
```

### Step 4: Verify Upload

```bash
# Check that everything is there
git log
git remote -v
```

Visit your repository: `https://github.com/YOUR-USERNAME/prometheus-dump-operator`

You should see:
- All source code files
- Documentation files
- .github/ directory with workflows and templates
- LICENSE file

---

## Repository Configuration

### Basic Settings

1. **Go to Settings** in your repository

2. **General Settings**:
   - Description: "Convert Prometheus TSDB to SAR format - Enable kernel engineers to analyze Prometheus metrics using familiar Unix tools"
   - Website: (your docs site if you have one)
   - Topics: Add these tags:
     - `prometheus`
     - `sar`
     - `monitoring`
     - `kubernetes`
     - `openshift`
     - `golang`
     - `sysadmin`
     - `system-administration`
     - `metrics`
     - `observability`

3. **Features**:
   - ✅ Issues
   - ✅ Discussions (highly recommended)
   - ❌ Sponsorships (optional)
   - ❌ Projects (optional)
   - ✅ Preserve this repository (recommended)
   - ❌ Wiki (we use Markdown docs)

### Branch Protection

Protect the `main` branch:

1. Go to **Settings > Branches**
2. Click "Add rule"
3. Branch name pattern: `main`
4. Configure protection:
   - ✅ Require pull request reviews before merging
   - ✅ Require status checks to pass before merging
     - Select: `test`, `build-cli`, `lint`
   - ✅ Require linear history
   - ✅ Include administrators (recommended)

### GitHub Actions

GitHub Actions should work automatically. Verify:

1. Go to **Actions** tab
2. You should see workflows:
   - ✅ CI (runs on push/PR)
   - ✅ Release (runs on tags)

3. First push triggers CI workflow. Check that it passes.

### Issue and PR Templates

Templates are already in `.github/ISSUE_TEMPLATE/` and `.github/pull_request_template.md`.

Verify they work:
1. Click **Issues > New Issue**
2. You should see templates:
   - Bug Report
   - Feature Request
   - Question

### Discussions

Enable Discussions for community Q&A:

1. Go to **Settings > Features**
2. Enable **Discussions**
3. Set up categories:
   - 💡 Ideas - Feature requests and enhancements
   - 🙏 Q&A - Questions from the community
   - 📣 Announcements - Project updates
   - 💬 General - General discussion

---

## First Release

### Create v1.0.0 Release

#### Step 1: Build Binaries

```bash
# Use the build script
./build.sh

# Or manually
make build-all

# Verify binaries work
./bin/prom2sar --version
./bin/prometheus-dump-operator --version
```

#### Step 2: Tag the Release

```bash
# Create annotated tag
git tag -a v1.0.0 -m "First stable release

- Standalone CLI binary (prom2sar)
- Kubernetes Operator
- Complete SAR format support
- CPU, memory, disk, network metrics
- Comprehensive documentation"

# Push tag to GitHub
git push origin v1.0.0
```

#### Step 3: Wait for GitHub Actions

The `release.yml` workflow will automatically:
- Build binaries for multiple platforms
- Create GitHub release
- Upload release artifacts
- Generate release notes from CHANGELOG.md

#### Step 4: Edit Release (Optional)

1. Go to **Releases**
2. Find the v1.0.0 release
3. Click "Edit release"
4. Add/modify release notes if needed
5. Upload any additional files
6. Click "Update release"

---

## Community Guidelines

### README Badges

Add badges to your README for better visibility. Update `GITHUB_README.md`:

```markdown
[![Go Version](https://img.shields.io/github/go-mod/go-version/YOUR-USERNAME/prometheus-dump-operator)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/YOUR-USERNAME/prometheus-dump-operator)](https://github.com/YOUR-USERNAME/prometheus-dump-operator/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![CI](https://github.com/YOUR-USERNAME/prometheus-dump-operator/workflows/CI/badge.svg)](https://github.com/YOUR-USERNAME/prometheus-dump-operator/actions)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
```

### Social Preview

1. Go to **Settings > Options**
2. Scroll to **Social preview**
3. Click "Edit"
4. Upload an image (1280x640px recommended)
   - Create a simple banner with project name and tagline
   - Or use GitHub's auto-generated preview

### About Section

In repository main page:
1. Click ⚙️ next to "About"
2. Fill in:
   - **Description**: "Convert Prometheus TSDB to SAR format"
   - **Website**: Link to docs
   - **Topics**: (see above)
   - ✅ Releases
   - ✅ Packages

### Star and Watch

Encourage users to:
- ⭐ **Star** the repository
- 👁️ **Watch** for updates
- 🍴 **Fork** to contribute

Add to README:

```markdown
⭐ **Star this repo** if you find it useful!

Want updates? Click **Watch** > **Custom** > **Releases**
```

---

## Maintenance

### Regular Tasks

#### Weekly
- Check and respond to new issues
- Review and merge PRs
- Update documentation if needed

#### Monthly
- Update dependencies: `go get -u ./...`
- Check for security updates: `go list -m -u all`
- Review and close stale issues

#### Per Release
- Update CHANGELOG.md
- Create release tag
- Verify GitHub Actions pass
- Test release artifacts
- Announce on Discussions

### Community Engagement

1. **Respond to Issues**:
   - Acknowledge within 48 hours
   - Label appropriately
   - Close resolved issues

2. **Review Pull Requests**:
   - Thank contributors
   - Provide constructive feedback
   - Merge when CI passes and approved

3. **Discussions**:
   - Answer questions
   - Share updates
   - Recognize contributors

### Documentation Updates

Keep docs current:

```bash
# When adding features
- Update relevant .md files
- Add examples if needed
- Update CHANGELOG.md

# When fixing bugs
- Document the fix if non-obvious
- Update troubleshooting guides

# For releases
- Update version numbers
- Update CHANGELOG.md
- Tag release
```

---

## Analytics and Insights

### View Repository Stats

GitHub provides insights:

1. **Insights > Traffic**:
   - Views, Clones
   - Referrers

2. **Insights > Community**:
   - Community health
   - Checklist completion

3. **Insights > Contributors**:
   - Who contributed
   - Contribution graph

### Optional: External Tools

- **CodeCov**: Code coverage reporting
- **Go Report Card**: Go code quality
- **pkg.go.dev**: Automatic Go documentation

---

## Promotion

### Share Your Project

1. **Reddit**:
   - r/golang
   - r/kubernetes
   - r/sysadmin
   - r/devops

2. **Hacker News**:
   - Show HN: Prometheus to SAR Converter

3. **Twitter/X**:
   - Tweet with hashtags: #golang #kubernetes #prometheus #sysadmin

4. **LinkedIn**:
   - Share in relevant groups

5. **Dev.to / Medium**:
   - Write blog post explaining the problem/solution

### Example Announcement

```
🎉 Introducing prom2sar - Convert Prometheus to SAR Format

Ever wished your kernel team could analyze Prometheus data without learning Prometheus?

prom2sar bridges the gap by converting Prometheus TSDB dumps to standard SAR format.

✅ Standalone CLI - works anywhere
✅ Kubernetes Operator - automated workflows  
✅ Zero learning curve - use grep, awk, sed
✅ Complete metrics - CPU, memory, disk, network

Perfect for must-gather analysis, incident investigations, and historical analysis.

https://github.com/YOUR-USERNAME/prometheus-dump-operator

#golang #kubernetes #prometheus #sysadmin #devops
```

---

## Checklist

Before making repository public:

- [ ] All sensitive data removed (no tokens, passwords, etc.)
- [ ] LICENSE file present
- [ ] README.md is clear and complete
- [ ] CONTRIBUTING.md explains how to contribute
- [ ] .gitignore is comprehensive
- [ ] CI/CD workflows configured
- [ ] Issue templates work
- [ ] PR template works
- [ ] Examples work
- [ ] Documentation is accurate
- [ ] Code builds successfully
- [ ] Tests pass

After making repository public:

- [ ] Topics/tags added
- [ ] Description set
- [ ] Social preview image uploaded
- [ ] Branch protection enabled
- [ ] Discussions enabled
- [ ] First release created
- [ ] Announced to community
- [ ] Starred your own repo! ⭐

---

## Support

If you need help setting up the repository:

1. Check GitHub's documentation: https://docs.github.com
2. Open an issue in your repository
3. Ask in GitHub Discussions
4. Check our CONTRIBUTING.md

---

## Summary

Your repository is ready for the community! 🎉

**What you have:**
- ✅ Complete source code
- ✅ Comprehensive documentation  
- ✅ CI/CD pipelines
- ✅ Issue/PR templates
- ✅ Contributing guidelines
- ✅ Apache 2.0 license

**Next steps:**
1. Push to GitHub
2. Configure repository settings
3. Create first release
4. Share with community
5. Accept contributions!

**Good luck with your open source project!** 🚀
