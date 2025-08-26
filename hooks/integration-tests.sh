#!/bin/bash
# Integration tests example hook script

echo "🔗 Running integration tests..."

# Example: Check if we should skip integration tests
if [ "$SKIP_INTEGRATION_TESTS" = "true" ]; then
    echo "⏭️  Skipping integration tests (SKIP_INTEGRATION_TESTS=true)"
    exit 0
fi

# Example: Docker-based integration tests
if [ -f "docker-compose.yml" ] || [ -f "docker-compose.yaml" ]; then
    if command -v docker-compose >/dev/null 2>&1; then
        echo "  Running docker-compose integration tests..."
        if ! docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from test; then
            echo "❌ Integration tests failed"
            docker-compose -f docker-compose.test.yml down
            exit 1
        fi
        docker-compose -f docker-compose.test.yml down
    fi
fi

# Example: API integration tests
if [ -f "integration_test.sh" ]; then
    echo "  Running custom integration tests..."
    if ! ./integration_test.sh; then
        echo "❌ Custom integration tests failed"
        exit 1
    fi
fi

# Example: Go integration tests with build tags
if find . -name "*_integration_test.go" | grep -q .; then
    if command -v go >/dev/null 2>&1; then
        echo "  Running Go integration tests..."
        if ! go test -tags=integration ./...; then
            echo "❌ Go integration tests failed"
            exit 1
        fi
    fi
fi

echo "✅ Integration tests passed!"
exit 0