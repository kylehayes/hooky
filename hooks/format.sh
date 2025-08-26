#!/bin/bash
# Format check example hook script

echo "🎨 Checking code formatting..."

# Example: Check if we have any Go files to format
if find . -name "*.go" -not -path "./.git/*" | grep -q .; then
    echo "  Checking Go files..."
    unformatted=$(gofmt -l . | grep -v vendor || true)
    if [ -n "$unformatted" ]; then
        echo "❌ The following files need formatting:"
        echo "$unformatted"
        echo "Run: gofmt -w ."
        exit 1
    fi
fi

# Example: Check JavaScript/TypeScript files
if find . -name "*.js" -o -name "*.ts" -o -name "*.jsx" -o -name "*.tsx" | grep -q .; then
    if command -v prettier >/dev/null 2>&1; then
        echo "  Checking JS/TS files with Prettier..."
        if ! prettier --check . 2>/dev/null; then
            echo "❌ Some files need formatting. Run: prettier --write ."
            exit 1
        fi
    fi
fi

# Example: Check Python files
if find . -name "*.py" | grep -q .; then
    if command -v black >/dev/null 2>&1; then
        echo "  Checking Python files with Black..."
        if ! black --check . 2>/dev/null; then
            echo "❌ Some Python files need formatting. Run: black ."
            exit 1
        fi
    fi
fi

echo "✅ Code formatting looks good!"
exit 0