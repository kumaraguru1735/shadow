# Shadow Build Guide

## Prerequisites

### Required

- **Go 1.22+** - Shadow is written in Go
- **Git** - For version control
- **jq** - For OAuth setup (optional but recommended)

### Optional

- **make** - Build automation (recommended)
- **npm** - For pi CLI installation

## Quick Build

```bash
# Navigate to Shadow directory
cd /opt/lampp/htdocs/shadow

# Build using Makefile (recommended)
make build

# Or build directly with Go
go build -v -o shadow cmd/shadow/main.go
```

**Output**: `./shadow` binary (6.4MB)

## Step-by-Step Build

### 1. Install Dependencies

```bash
# Install Go dependencies
go mod download
go mod tidy
```

### 2. Build Binary

```bash
# Option A: Using Makefile
make build

# Option B: Using Go directly
go build -v -o shadow cmd/shadow/main.go

# Option C: With custom output location
go build -v -o /usr/local/bin/shadow cmd/shadow/main.go
```

### 3. Verify Build

```bash
# Check binary exists
ls -lh shadow

# Test version
./shadow --version

# Test help
./shadow --help
```

## Build Commands

### Using Makefile

```bash
# Build (default)
make build

# Install to /usr/local/bin
sudo make install

# Clean build artifacts
make clean

# Run tests
make test

# Build for production (optimized)
make build-release

# Build for multiple platforms
make build-all
```

### Using Go Commands

```bash
# Development build (fast, includes debug info)
go build -v -o shadow cmd/shadow/main.go

# Production build (optimized, smaller binary)
go build -v -ldflags="-s -w" -o shadow cmd/shadow/main.go

# Build with version info
VERSION=$(git describe --tags --always)
go build -ldflags="-X main.version=$VERSION" -o shadow cmd/shadow/main.go

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o shadow-linux-amd64 cmd/shadow/main.go
GOOS=darwin GOARCH=arm64 go build -o shadow-darwin-arm64 cmd/shadow/main.go
GOOS=windows GOARCH=amd64 go build -o shadow-windows-amd64.exe cmd/shadow/main.go
```

## Build Targets

### Development Build

```bash
# Fast build with debug info
go build -v -o shadow cmd/shadow/main.go
```

**Characteristics:**
- Fast compilation
- Includes debug symbols
- ~6.4MB binary
- Good for development/testing

### Production Build

```bash
# Optimized build
go build -v -ldflags="-s -w" -o shadow cmd/shadow/main.go
```

**Characteristics:**
- Stripped symbols (-s)
- Stripped DWARF debug info (-w)
- Smaller binary (~5MB)
- Good for distribution

### Static Build

```bash
# Fully static binary (no dynamic linking)
CGO_ENABLED=0 go build -v -ldflags="-s -w -extldflags '-static'" -o shadow cmd/shadow/main.go
```

**Characteristics:**
- No external dependencies
- Works on any Linux system
- Slightly larger (~6.5MB)
- Good for containers/distribution

## Cross-Platform Builds

### Linux

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o shadow-linux-amd64 cmd/shadow/main.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o shadow-linux-arm64 cmd/shadow/main.go
```

### macOS

```bash
# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o shadow-darwin-amd64 cmd/shadow/main.go

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o shadow-darwin-arm64 cmd/shadow/main.go
```

### Windows

```bash
# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o shadow-windows-amd64.exe cmd/shadow/main.go
```

## Post-Build Setup

### 1. Setup OAuth (Recommended)

```bash
# Extract OAuth from Claude Code
./setup_oauth.sh

# Or manually
jq '.claudeAiOauth' ~/.claude/.credentials.json > ~/.claude/oauth.json
chmod 600 ~/.claude/oauth.json
```

### 2. Verify Installation

```bash
# Check authentication
./shadow auth-check

# Should show:
# ✓ Claude Code OAuth token found
# ✅ AI client initialized successfully
```

### 3. Create Config (Optional)

```bash
# Create Shadow config directory
mkdir -p ~/.shadow

# Copy example config
cp configs/config.example.yaml ~/.shadow/config.yaml

# Edit as needed
nano ~/.shadow/config.yaml
```

## Build Troubleshooting

### "Command not found: go"

**Solution**: Install Go

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang-go

# macOS
brew install go

# Or download from https://go.dev/dl/
```

### "go.mod not found"

**Solution**: Run from correct directory

```bash
cd /opt/lampp/htdocs/shadow
ls go.mod  # Should exist
go build -v -o shadow cmd/shadow/main.go
```

### "Cannot find package"

**Solution**: Download dependencies

```bash
go mod download
go mod tidy
go build -v -o shadow cmd/shadow/main.go
```

### "Permission denied"

**Solution**: Make binary executable

```bash
chmod +x shadow
```

### Build takes too long

**Solution**: Use cached builds

```bash
# First build (slow)
go build -v -o shadow cmd/shadow/main.go

# Subsequent builds (fast, uses cache)
go build -o shadow cmd/shadow/main.go
```

## Complete Build Script

Save as `build.sh`:

```bash
#!/bin/bash
set -e

echo "🔨 Building Shadow..."

# Check prerequisites
command -v go >/dev/null 2>&1 || { echo "❌ Go not installed"; exit 1; }
command -v jq >/dev/null 2>&1 || { echo "⚠️  jq not installed (optional)"; }

# Clean previous builds
rm -f shadow

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

# Build
echo "🔧 Building binary..."
go build -v -ldflags="-s -w" -o shadow cmd/shadow/main.go

# Verify
if [ -f shadow ]; then
    chmod +x shadow
    SIZE=$(du -h shadow | cut -f1)
    echo "✅ Build complete: ./shadow ($SIZE)"

    # Test
    ./shadow --version

    # Setup OAuth if available
    if [ -f ~/.claude/.credentials.json ]; then
        echo "🔐 Setting up OAuth..."
        ./setup_oauth.sh
    fi

    echo ""
    echo "✅ Shadow is ready to use!"
    echo "   ./shadow scan example.com"
else
    echo "❌ Build failed"
    exit 1
fi
```

Make executable and run:

```bash
chmod +x build.sh
./build.sh
```

## Development Workflow

### 1. Make Changes

```bash
# Edit code
nano internal/ai/pi_client.go
```

### 2. Rebuild

```bash
# Quick rebuild (uses cache)
go build -o shadow cmd/shadow/main.go
```

### 3. Test

```bash
# Run tests
go test ./...

# Test binary
./shadow --help
./shadow auth-check
```

### 4. Format & Lint

```bash
# Format code
go fmt ./...

# Run linter (if installed)
golangci-lint run
```

## Distribution

### Create Release

```bash
# Build for all platforms
make build-all

# Or manually
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o shadow-linux-amd64 cmd/shadow/main.go
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o shadow-darwin-amd64 cmd/shadow/main.go
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o shadow-darwin-arm64 cmd/shadow/main.go
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o shadow-windows-amd64.exe cmd/shadow/main.go

# Create checksums
sha256sum shadow-* > SHA256SUMS
```

### Package

```bash
# Create tarball
tar -czf shadow-v0.1.0-linux-amd64.tar.gz shadow-linux-amd64 README.md LICENSE

# Create zip for Windows
zip shadow-v0.1.0-windows-amd64.zip shadow-windows-amd64.exe README.md LICENSE
```

## Summary

**Quick Build:**
```bash
cd /opt/lampp/htdocs/shadow
make build
./setup_oauth.sh
./shadow auth-check
```

**Binary Location:** `./shadow` (6.4MB)

**Next Steps:**
- Run `./shadow auth-check` to verify OAuth
- Read [QUICKSTART.md](QUICKSTART.md) for usage
- Read [CLAUDE_CODE_OAUTH.md](CLAUDE_CODE_OAUTH.md) for OAuth setup

---

**Built with Go. Powered by Claude AI.** 🚀
