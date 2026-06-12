#!/bin/bash
# Script to initialize GitHub repository

set -e

echo "=== Initializing GitHub Repository for Prometheus to SAR Converter ==="
echo ""

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo "Error: git is not installed"
    exit 1
fi

# Check if we're already in a git repository
if [ -d ".git" ]; then
    echo "Warning: This directory is already a git repository"
    read -p "Do you want to continue? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    # Initialize git repository
    echo "1. Initializing git repository..."
    git init
    echo "✓ Git repository initialized"
    echo ""
fi

# Create .gitignore if it doesn't exist
if [ ! -f ".gitignore" ]; then
    echo "2. Creating .gitignore..."
    cat > .gitignore <<EOF
# Binaries
bin/
dist/
*.exe
*.exe~
*.dll
*.so
*.dylib
prom2sar
prom2sar-*
prometheus-dump-operator

# Test artifacts
*.test
*.out
coverage.out
coverage.html

# Build artifacts
*.o
*.a

# Go workspace
go.work
go.work.sum

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Temporary files
tmp/
temp/
*.tmp

# SAR output (for testing)
sar-output/
*.sar

# Kubernetes
kubeconfig

# Local testing
testdata/prometheus/
must-gather.local.*/

# Environment
.env
.env.local
EOF
    echo "✓ .gitignore created"
    echo ""
else
    echo "2. .gitignore already exists, skipping..."
    echo ""
fi

# Copy GITHUB_README.md to README.md
echo "3. Setting up README.md..."
if [ -f "GITHUB_README.md" ]; then
    cp GITHUB_README.md README_GITHUB.md
    echo "✓ Copied GITHUB_README.md to README_GITHUB.md"
    echo "  (You can replace README.md with this for GitHub)"
else
    echo "⚠ GITHUB_README.md not found, skipping..."
fi
echo ""

# Stage all files
echo "4. Staging files..."
git add .
echo "✓ Files staged"
echo ""

# Create initial commit
echo "5. Creating initial commit..."
if git diff --cached --quiet; then
    echo "⚠ No changes to commit"
else
    git commit -m "Initial commit: Prometheus to SAR Converter

- Standalone CLI binary (prom2sar)
- Kubernetes Operator for automated conversions
- Complete SAR format output (CPU, memory, disk, network)
- Comprehensive documentation (10+ guides)
- Build and deployment tooling
- Examples and test cases
- GitHub workflows for CI/CD
- Contributing guide and issue templates"
    echo "✓ Initial commit created"
fi
echo ""

# Instructions for GitHub
echo "=== Next Steps ==="
echo ""
echo "6. Create a GitHub repository:"
echo "   - Go to https://github.com/new"
echo "   - Repository name: prometheus-dump-operator"
echo "   - Description: Convert Prometheus TSDB to SAR format - Enable kernel engineers to analyze Prometheus metrics using familiar Unix tools"
echo "   - Public or Private: Your choice"
echo "   - Do NOT initialize with README, .gitignore, or license (we have them)"
echo ""
echo "7. Add GitHub remote and push:"
echo ""
echo "   git remote add origin https://github.com/YOUR-USERNAME/prometheus-dump-operator.git"
echo "   git branch -M main"
echo "   git push -u origin main"
echo ""
echo "8. Configure GitHub repository settings:"
echo "   - Add topics: prometheus, sar, monitoring, kubernetes, golang, sysadmin"
echo "   - Add description from step 6"
echo "   - Enable Issues"
echo "   - Enable Discussions (optional but recommended)"
echo "   - Enable Wiki (optional)"
echo ""
echo "9. Create first release:"
echo "   - Build binaries: ./build.sh"
echo "   - Create tag: git tag -a v1.0.0 -m 'First release'"
echo "   - Push tag: git push origin v1.0.0"
echo "   - GitHub Actions will create release automatically"
echo ""
echo "10. Optional: Set up branch protection"
echo "    - Go to Settings > Branches"
echo "    - Add rule for 'main' branch"
echo "    - Require pull request reviews"
echo "    - Require status checks to pass"
echo ""
echo "=== Repository Contents ==="
echo ""
echo "Source Code:"
echo "  - cmd/prom2sar/        CLI binary"
echo "  - cmd/main.go          Operator"
echo "  - pkg/tsdb/            TSDB reader"
echo "  - pkg/sar/             SAR conversion"
echo ""
echo "Documentation:"
echo "  - README.md            Main overview"
echo "  - GITHUB_README.md     GitHub-specific README"
echo "  - CLI_GUIDE.md         CLI documentation"
echo "  - GETTING_STARTED.md   Quick start guide"
echo "  - CONTRIBUTING.md      Contribution guide"
echo "  + 6 more guides"
echo ""
echo "GitHub Configuration:"
echo "  - .github/workflows/   CI/CD pipelines"
echo "  - .github/ISSUE_TEMPLATE/  Issue templates"
echo "  - LICENSE              Apache 2.0"
echo "  - CHANGELOG.md         Version history"
echo ""
echo "=== All Set! ==="
echo ""
echo "Your repository is ready to push to GitHub!"
echo ""
