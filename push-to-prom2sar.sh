#!/bin/bash
# Push to your prom2sar repository

set -e

echo "════════════════════════════════════════════════════════════════════"
echo "  Pushing to your prom2sar repository"
echo "════════════════════════════════════════════════════════════════════"
echo ""

# Get GitHub username
echo "What is your GitHub username?"
read -p "Username: " GITHUB_USER

if [ -z "$GITHUB_USER" ]; then
    echo "❌ Error: Username cannot be empty"
    exit 1
fi

echo ""
echo "Repository URL: https://github.com/${GITHUB_USER}/prom2sar"
echo ""

# Check if remote exists
if git remote | grep -q "^origin$"; then
    echo "⚠️  Remote 'origin' already exists. Updating URL..."
    git remote set-url origin "https://github.com/${GITHUB_USER}/prom2sar.git"
else
    echo "Adding remote 'origin'..."
    git remote add origin "https://github.com/${GITHUB_USER}/prom2sar.git"
fi

echo "✅ Remote configured"
echo ""

# Show what will be pushed
echo "Commits to push:"
git log --oneline
echo ""

# Confirm
read -p "Push to https://github.com/${GITHUB_USER}/prom2sar? (y/n) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

echo ""
echo "Pushing to GitHub..."
echo "(You may be prompted for credentials)"
echo ""

# Push
git push -u origin main

echo ""
echo "════════════════════════════════════════════════════════════════════"
echo "  ✅ SUCCESS! Repository pushed to GitHub!"
echo "════════════════════════════════════════════════════════════════════"
echo ""
echo "View at: https://github.com/${GITHUB_USER}/prom2sar"
echo ""
