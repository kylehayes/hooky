#!/bin/bash
# Dependency check example hook script (post-checkout)

echo "📦 Checking dependencies..."

# Get information about the checkout
prev_head="$1"
new_head="$2"
branch_checkout="$3"

# Skip if this is a file checkout, not a branch checkout
if [ "$branch_checkout" = "0" ]; then
    echo "✅ File checkout - skipping dependency check"
    exit 0
fi

dependency_files_changed=false

# Check if dependency files have changed
if [ "$prev_head" != "$new_head" ]; then
    # Check various dependency files
    for dep_file in "package.json" "package-lock.json" "yarn.lock" "go.mod" "go.sum" "requirements.txt" "Pipfile" "Pipfile.lock" "Cargo.toml" "Cargo.lock"; do
        if git diff --name-only "$prev_head" "$new_head" | grep -q "^$dep_file$"; then
            echo "  📄 $dep_file changed"
            dependency_files_changed=true
        fi
    done
fi

if [ "$dependency_files_changed" = false ]; then
    echo "✅ No dependency changes detected"
    exit 0
fi

echo "  🔄 Dependency files changed, checking if update needed..."

# Node.js projects
if [ -f "package.json" ]; then
    if [ -f "package-lock.json" ]; then
        echo "  📦 npm install might be needed"
        echo "  Run: npm install"
    elif [ -f "yarn.lock" ]; then
        echo "  📦 yarn install might be needed"
        echo "  Run: yarn install"
    fi
fi

# Go projects
if [ -f "go.mod" ]; then
    echo "  📦 Go dependencies might need updating"
    echo "  Run: go mod download && go mod tidy"
fi

# Python projects
if [ -f "requirements.txt" ]; then
    echo "  📦 Python dependencies might need updating"
    echo "  Run: pip install -r requirements.txt"
elif [ -f "Pipfile" ]; then
    echo "  📦 Pipenv dependencies might need updating"
    echo "  Run: pipenv install"
fi

# Rust projects
if [ -f "Cargo.toml" ]; then
    echo "  📦 Cargo dependencies might need updating"
    echo "  Run: cargo build"
fi

echo "⚠️  Consider updating dependencies after this checkout"
exit 0