#!/bin/bash
# Testing example hook script

echo "🧪 Running tests..."

# Example: Go tests
if find . -name "*.go" -not -path "./.git/*" | grep -q .; then
    if command -v go >/dev/null 2>&1; then
        echo "  Running Go tests..."
        if ! go test ./...; then
            echo "❌ Go tests failed"
            exit 1
        fi
    fi
fi

# Example: JavaScript/TypeScript tests
if [ -f "package.json" ]; then
    if command -v npm >/dev/null 2>&1; then
        echo "  Running npm tests..."
        if ! npm test; then
            echo "❌ npm tests failed"
            exit 1
        fi
    elif command -v yarn >/dev/null 2>&1; then
        echo "  Running yarn tests..."
        if ! yarn test; then
            echo "❌ yarn tests failed"
            exit 1
        fi
    fi
fi

# Example: Python tests
if find . -name "*.py" | grep -q .; then
    if command -v pytest >/dev/null 2>&1; then
        echo "  Running pytest..."
        if ! pytest; then
            echo "❌ pytest failed"
            exit 1
        fi
    elif command -v python >/dev/null 2>&1 && [ -f "manage.py" ]; then
        echo "  Running Django tests..."
        if ! python manage.py test; then
            echo "❌ Django tests failed"
            exit 1
        fi
    fi
fi

echo "✅ All tests passed!"
exit 0