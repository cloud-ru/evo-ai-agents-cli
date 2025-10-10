#!/bin/bash

# Release script for AI Agents CLI
# This script builds binaries for all platforms and creates GitHub releases

set -e

VERSION=${1:-"1.0.0"}
GITHUB_TOKEN=${GITHUB_TOKEN:-""}
REPO="cloudru/ai-agents-cli"

echo "🚀 Creating release v$VERSION for $REPO"

# Check if GitHub token is provided
if [ -z "$GITHUB_TOKEN" ]; then
    echo "❌ GITHUB_TOKEN environment variable is required"
    exit 1
fi

# Build for all platforms
echo "📦 Building binaries for all platforms..."

# Create release directory
mkdir -p release

# Build for different platforms
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r os arch <<< "$platform"
    echo "Building for $os/$arch..."
    
    output_name="ai-agents-cli"
    if [ "$os" = "windows" ]; then
        output_name+=".exe"
    fi
    
    GOOS=$os GOARCH=$arch go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$(date -u +%Y-%m-%d_%H:%M:%S) -X main.GitCommit=$(git rev-parse --short HEAD)" -o "release/ai-agents-cli-$os-$arch$([ "$os" = "windows" ] && echo ".exe")" .
done

# Create archives
echo "📁 Creating archives..."
cd release

for platform in "${platforms[@]}"; do
    IFS='/' read -r os arch <<< "$platform"
    archive_name="ai-agents-cli-$os-$arch"
    
    if [ "$os" = "windows" ]; then
        zip -r "$archive_name.zip" "ai-agents-cli-$os-$arch.exe"
    else
        tar -czf "$archive_name.tar.gz" "ai-agents-cli-$os-$arch"
    fi
done

cd ..

# Create GitHub release
echo "🏷️ Creating GitHub release..."

# Create release notes
cat > release_notes.md << EOF
# AI Agents CLI v$VERSION

## 🚀 What's New

- Initial release of AI Agents CLI
- Full support for MCP servers management
- AI agents and agent systems management
- CI/CD integration with validation and monitoring
- Beautiful colored interface with emojis
- Support for YAML and JSON configurations
- Multi-language support (Russian, English)

## 📦 Installation

### Using winget (Windows)
\`\`\`bash
winget install CloudRu.AIAgentsCLI
\`\`\`

### Using Homebrew (macOS/Linux)
\`\`\`bash
brew install cloudru/ai-agents-cli/ai-agents-cli
\`\`\`

### Manual Installation
Download the appropriate binary for your platform from the assets below.

## 🔧 Quick Start

1. Set your API key: \`export API_KEY="your-api-key"\`
2. Set your project ID: \`export PROJECT_ID="your-project-id"\`
3. Run: \`ai-agents-cli --help\`

## 📚 Documentation

- [README](https://github.com/$REPO/blob/main/README.md)
- [Examples](https://github.com/$REPO/tree/main/examples)
- [CI/CD Integration](https://github.com/$REPO/tree/main/examples/.github/workflows)

## 🐛 Bug Reports

Please report bugs in the [Issues](https://github.com/$REPO/issues) section.

## 📄 License

MIT License - see [LICENSE](https://github.com/$REPO/blob/main/LICENSE) for details.
EOF

# Create GitHub release
gh release create "v$VERSION" \
    --title "AI Agents CLI v$VERSION" \
    --notes-file release_notes.md \
    release/*.tar.gz \
    release/*.zip

echo "✅ Release v$VERSION created successfully!"
echo "🔗 View release: https://github.com/$REPO/releases/tag/v$VERSION"

# Cleanup
rm -rf release
rm release_notes.md

echo "🎉 Release process completed!"
