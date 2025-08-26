#!/bin/bash
# Pre-rebase safety check example hook script

echo "⚠️  Pre-rebase safety checks..."

# Get the upstream and branch being rebased
upstream="$1"
branch="$2"

# Check if we're on main/master branch (usually shouldn't rebase these)
current_branch=$(git rev-parse --abbrev-ref HEAD)
if [ "$current_branch" = "main" ] || [ "$current_branch" = "master" ]; then
    echo "❌ Attempting to rebase main/master branch"
    echo "This is generally not recommended. If you're sure, use --no-verify"
    exit 1
fi

# Check if there are uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "❌ You have uncommitted changes"
    echo "Please commit or stash your changes before rebasing"
    exit 1
fi

# Check if the upstream exists
if [ -n "$upstream" ] && ! git rev-parse --verify "$upstream" >/dev/null 2>&1; then
    echo "❌ Upstream '$upstream' does not exist"
    exit 1
fi

# Warn about public commits
if [ -n "$branch" ]; then
    # Check if any commits in the range have been pushed to origin
    commits_to_rebase=$(git rev-list --count "$upstream..$branch" 2>/dev/null || echo "0")
    if [ "$commits_to_rebase" -gt 0 ]; then
        echo "⚠️  About to rebase $commits_to_rebase commit(s)"
        
        # Check if any of these commits exist on the remote
        if git branch -r --contains "$branch" | grep -q "origin/"; then
            echo "❌ Some commits appear to have been pushed to remote"
            echo "Rebasing published commits can cause problems for other developers"
            echo "If you're sure, use --no-verify"
            exit 1
        fi
    fi
fi

echo "✅ Pre-rebase checks passed!"
exit 0