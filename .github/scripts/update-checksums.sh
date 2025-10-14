#!/bin/bash

# Script to update checksums in package manager files
# Usage: ./update-checksums.sh <version> <release-url>

set -e

VERSION=${1:-"1.0.0"}
RELEASE_URL=${2:-"https://github.com/cloud-ru/evo-ai-agents-cli/releases/download/v$VERSION"}

echo "🔄 Updating checksums for version $VERSION"

# Download release assets
mkdir -p temp_assets
cd temp_assets

echo "📥 Downloading release assets..."
wget -q "$RELEASE_URL/ai-agents-cli-darwin-amd64.tar.gz"
wget -q "$RELEASE_URL/ai-agents-cli-darwin-arm64.tar.gz"
wget -q "$RELEASE_URL/ai-agents-cli-linux-amd64.tar.gz"
wget -q "$RELEASE_URL/ai-agents-cli-windows-amd64.zip"

# Calculate checksums
DARWIN_AMD64_SHA=$(sha256sum ai-agents-cli-darwin-amd64.tar.gz | cut -d' ' -f1)
DARWIN_ARM64_SHA=$(sha256sum ai-agents-cli-darwin-arm64.tar.gz | cut -d' ' -f1)
LINUX_AMD64_SHA=$(sha256sum ai-agents-cli-linux-amd64.tar.gz | cut -d' ' -f1)
WINDOWS_AMD64_SHA=$(sha256sum ai-agents-cli-windows-amd64.zip | cut -d' ' -f1)

echo "📊 Calculated checksums:"
echo "  Darwin AMD64: $DARWIN_AMD64_SHA"
echo "  Darwin ARM64: $DARWIN_ARM64_SHA"
echo "  Linux AMD64:  $LINUX_AMD64_SHA"
echo "  Windows AMD64: $WINDOWS_AMD64_SHA"

cd ..

# Update Homebrew formula
echo "🍺 Updating Homebrew formula..."
sed -i "s/version \"[^\"]*\"/version \"$VERSION\"/" .github/package-managers/brew/ai-agents-cli.rb
sed -i "s/sha256 \"[^\"]*\"/sha256 \"$DARWIN_AMD64_SHA\"/" .github/package-managers/brew/ai-agents-cli.rb
sed -i "s/sha256 \"[^\"]*_ARM64\"/sha256 \"$DARWIN_ARM64_SHA\"/" .github/package-managers/brew/ai-agents-cli.rb
sed -i "s/sha256 \"[^\"]*_LINUX\"/sha256 \"$LINUX_AMD64_SHA\"/" .github/package-managers/brew/ai-agents-cli.rb

# Update URLs in Homebrew formula
sed -i "s|url \"https://github.com/cloud-ru/evo-ai-agents-cli/releases/download/[^\"]*\"|url \"$RELEASE_URL/ai-agents-cli-darwin-amd64.tar.gz\"|" .github/package-managers/brew/ai-agents-cli.rb

# Update Winget manifest
echo "🪟 Updating Winget manifest..."
sed -i "s/PackageVersion: [0-9.]*/PackageVersion: $VERSION/" .github/package-managers/winget/ai-agents-cli.yaml
sed -i "s/InstallerSha256: [a-f0-9]*/InstallerSha256: $WINDOWS_AMD64_SHA/" .github/package-managers/winget/ai-agents-cli.yaml
sed -i "s|InstallerUrl: https://github.com/cloud-ru/evo-ai-agents-cli/releases/download/[^\"]*|InstallerUrl: $RELEASE_URL/ai-agents-cli-windows-amd64.zip|" .github/package-managers/winget/ai-agents-cli.yaml

# Cleanup
rm -rf temp_assets

echo "✅ Checksums updated successfully!"
echo "📝 Don't forget to commit the changes:"
echo "   git add .github/package-managers/"
echo "   git commit -m \"Update package manager checksums to v$VERSION\""
