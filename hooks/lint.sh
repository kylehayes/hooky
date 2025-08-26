#!/bin/bash
# Linting example hook script

echo "🔍 Running linting checks..."

# Example: Go linting
if find . -name "*.go" -not -path "./.git/*" | grep -q .; then
    if command -v golangci-lint >/dev/null 2>&1; then
        echo "  Running golangci-lint..."
        if ! golangci-lint run; then
            echo "❌ Go linting failed"
            exit 1
        fi
    elif command -v go >/dev/null 2>&1; then
        echo "  Running go vet..."
        if ! go vet ./...; then
            echo "❌ Go vet failed"
            exit 1
        fi
    fi
fi

# Example: JavaScript/TypeScript linting
if find . -name "*.js" -o -name "*.ts" -o -name "*.jsx" -o -name "*.tsx" | grep -q .; then
    if command -v eslint >/dev/null 2>&1; then
        echo "  Running ESLint..."
        if ! eslint . --ext .js,.ts,.jsx,.tsx; then
            echo "❌ ESLint failed"
            exit 1
        fi
    fi
fi

# Example: Python linting
if find . -name "*.py" | grep -q .; then
    if command -v flake8 >/dev/null 2>&1; then
        echo "  Running flake8..."
        if ! flake8 .; then
            echo "❌ Python linting failed"
            exit 1
        fi
    elif command -v pylint >/dev/null 2>&1; then
        echo "  Running pylint..."
        if ! pylint **/*.py; then
            echo "❌ Pylint failed"
            exit 1
        fi
    fi
fi

echo "✅ All linting checks passed!"
exit 0