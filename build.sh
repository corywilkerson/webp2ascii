#!/bin/bash
# Build for multiple platforms

APP_NAME="webp2ascii"
VERSION="1.0.0"

# Clean previous builds
rm -rf dist/
mkdir -p dist/

# Build for different platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    OS="${PLATFORM%/*}"
    ARCH="${PLATFORM#*/}"
    OUTPUT="dist/${APP_NAME}-${VERSION}-${OS}-${ARCH}"
    
    if [ "$OS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi
    
    echo "Building for $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w" -o "$OUTPUT" main.go
done

# Create tarballs
cd dist/
for file in *; do
    if [[ ! "$file" == *.tar.gz ]]; then
        tar -czf "${file}.tar.gz" "$file"
        rm "$file"
    fi
done
cd ..

echo "Builds complete! Check the dist/ directory"