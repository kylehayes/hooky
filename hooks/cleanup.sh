#!/bin/bash
# Post-merge cleanup example hook script

echo "🧹 Post-merge cleanup..."

# Get merge information
squash_merge="$1"

# Clean up common temporary files
echo "  🗑️  Cleaning temporary files..."

# Remove common editor temporary files
find . -name ".DS_Store" -delete 2>/dev/null || true
find . -name "*.tmp" -delete 2>/dev/null || true
find . -name "*.temp" -delete 2>/dev/null || true
find . -name "*~" -delete 2>/dev/null || true

# Clean up common build artifacts that might not be in .gitignore
find . -name "*.pyc" -delete 2>/dev/null || true
find . -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true
find . -name ".pytest_cache" -type d -exec rm -rf {} + 2>/dev/null || true

# Node.js cleanup
if [ -d "node_modules" ]; then
    echo "  📦 Checking node_modules..."
    # Only suggest cleanup, don't do it automatically
    echo "  💡 Consider running 'npm prune' to clean unused packages"
fi

# Check for merge conflict markers that might have been accidentally committed
echo "  🔍 Checking for leftover conflict markers..."
if grep -r "^<<<<<<< \|^=======$\|^>>>>>>> " . --exclude-dir=.git 2>/dev/null; then
    echo "  ⚠️  Found potential conflict markers in files"
    echo "  Please review and clean up any remaining conflict markers"
fi

# Update file permissions if needed
echo "  🔐 Ensuring script permissions..."
find hooks -name "*.sh" -type f -exec chmod +x {} \; 2>/dev/null || true

echo "✅ Post-merge cleanup completed!"
exit 0