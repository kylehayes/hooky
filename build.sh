#!/bin/bash
# Build script for cross-platform compilation

set -e

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="dist"

echo "🏗️  Building hooky v${VERSION}..."

# Clean previous builds
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# Build targets (OS/ARCH combinations)
declare -a targets=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
    "freebsd/amd64"
)

for target in "${targets[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$target"
    
    echo "  📦 Building for $GOOS/$GOARCH..."
    
    output_name="hooky"
    if [ "$GOOS" = "windows" ]; then
        output_name="hooky.exe"
    fi
    
    output_path="$BUILD_DIR/hooky-${VERSION}-${GOOS}-${GOARCH}"
    mkdir -p "$output_path"
    
    # Build binary
    env GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags "-X main.version=$VERSION -s -w" -o "$output_path/$output_name" .
    
    # Copy configuration and example files
    cp hooky.yaml "$output_path/"
    cp -r hooks "$output_path/"
    
    # Create archive
    if [ "$GOOS" = "windows" ]; then
        (cd "$BUILD_DIR" && zip -r "hooky-${VERSION}-${GOOS}-${GOARCH}.zip" "hooky-${VERSION}-${GOOS}-${GOARCH}/")
    else
        (cd "$BUILD_DIR" && tar -czf "hooky-${VERSION}-${GOOS}-${GOARCH}.tar.gz" "hooky-${VERSION}-${GOOS}-${GOARCH}/")
    fi
    
    # Remove the directory after archiving
    rm -rf "$output_path"
    
    echo "    ✅ Created archive for $GOOS/$GOARCH"
done

echo ""
echo "🎉 Build completed! Archives created in $BUILD_DIR:"
ls -la "$BUILD_DIR"

echo ""
echo "📝 To install, users should:"
echo "  1. Download the appropriate archive for their platform"
echo "  2. Extract it to a directory in their PATH"
echo "  3. Run 'hooky --install' in their git repository"