#!/bin/bash
# Script to push repository to GitHub

set -e

echo "╔══════════════════════════════════════════════════════════════════════╗"
echo "║         Push Prometheus to SAR Converter to GitHub                  ║"
echo "╚══════════════════════════════════════════════════════════════════════╝"
echo ""

# Check if we're in the right directory
if [ ! -f "README.md" ] || [ ! -d ".git" ]; then
    echo "❌ Error: Please run this script from the prometheus-dump-operator directory"
    exit 1
fi

# Check if git is clean
if ! git diff-index --quiet HEAD --; then
    echo "⚠️  Warning: You have uncommitted changes"
    git status
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Get GitHub username
echo "Step 1: Enter your GitHub username"
echo "────────────────────────────────────────"
read -p "GitHub username: " GITHUB_USER

if [ -z "$GITHUB_USER" ]; then
    echo "❌ Error: GitHub username cannot be empty"
    exit 1
fi

echo ""
echo "Step 2: Checking if you created the repository on GitHub..."
echo "────────────────────────────────────────"
echo "Have you created the repository on GitHub?"
echo "  → https://github.com/new"
echo "  → Repository name: prometheus-dump-operator"
echo "  → Public repository"
echo "  → DO NOT add README, .gitignore, or license"
echo ""
read -p "Repository created on GitHub? (y/n) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "Please create the repository first:"
    echo "1. Go to https://github.com/new"
    echo "2. Repository name: prometheus-dump-operator"
    echo "3. Make it Public"
    echo "4. DO NOT check any initialization options"
    echo "5. Click 'Create repository'"
    echo ""
    echo "Then run this script again!"
    exit 0
fi

echo ""
echo "Step 3: Adding GitHub remote..."
echo "────────────────────────────────────────"

REPO_URL="https://github.com/${GITHUB_USER}/prometheus-dump-operator.git"
echo "Repository URL: $REPO_URL"

# Check if remote already exists
if git remote | grep -q "^origin$"; then
    echo "⚠️  Remote 'origin' already exists"
    EXISTING_URL=$(git remote get-url origin)
    echo "   Current URL: $EXISTING_URL"

    if [ "$EXISTING_URL" != "$REPO_URL" ]; then
        read -p "Update remote URL? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            git remote set-url origin "$REPO_URL"
            echo "✅ Remote URL updated"
        fi
    else
        echo "✅ Remote already points to correct repository"
    fi
else
    git remote add origin "$REPO_URL"
    echo "✅ Remote 'origin' added"
fi

echo ""
echo "Step 4: Verifying remote..."
echo "────────────────────────────────────────"
git remote -v

echo ""
echo "Step 5: Pushing to GitHub..."
echo "────────────────────────────────────────"
echo "This will push to: $REPO_URL"
read -p "Continue with push? (y/n) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Push cancelled. You can push manually with:"
    echo "  git push -u origin main"
    exit 0
fi

echo ""
echo "Pushing to GitHub..."
echo "(You may be prompted for GitHub credentials or token)"
echo ""

# Push to GitHub
git push -u origin main

echo ""
echo "╔══════════════════════════════════════════════════════════════════════╗"
echo "║                        🎉 SUCCESS! 🎉                                ║"
echo "╚══════════════════════════════════════════════════════════════════════╝"
echo ""
echo "✅ Repository pushed to GitHub!"
echo ""
echo "View your repository at:"
echo "   https://github.com/${GITHUB_USER}/prometheus-dump-operator"
echo ""
echo "Next steps:"
echo "1. Configure repository settings (see GITHUB_SETUP_GUIDE.md)"
echo "2. Add topics/tags to your repository"
echo "3. Enable GitHub Discussions (optional)"
echo "4. Create first release: git tag -a v1.0.0 -m \"First release\""
echo "                        git push origin v1.0.0"
echo ""
echo "Share your project:"
echo "  - Reddit: r/golang, r/kubernetes, r/sysadmin"
echo "  - Hacker News: Show HN"
echo "  - Twitter/LinkedIn with #golang #kubernetes tags"
echo ""
echo "Congratulations on your open source project! 🚀"
echo ""
