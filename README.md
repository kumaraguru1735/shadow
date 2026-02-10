# Shadow 🕵️

> AI-augmented cybersecurity reconnaissance and analysis platform

Shadow combines the power of Go's performance with Claude AI's intelligence to deliver comprehensive, automated security assessments with advanced retry logic and intelligent error handling.

## Features

- 🔐 **Advanced Authentication System** - OAuth + API key support with automatic management
- 🔍 **Network Intelligence** - Port scanning, subdomain discovery, SSL/TLS analysis
- 🌐 **Web Security Analysis** - Crawling, header validation, vulnerability detection
- 🤖 **AI-Powered Analysis** - Claude AI with intelligent retry and exponential backoff
- 📊 **Smart Reporting** - Executive summaries and technical reports
- 🔄 **Robust Error Handling** - Auto-retry for rate limits, timeouts, and transient failures
- 🚀 **High Performance** - Go-powered concurrency, 10-minute analysis timeout
- 🔌 **Production-Ready** - Patterns from OpenClaw's proven implementation

## Quick Start

```bash
# Build from source
git clone https://github.com/kumaraguru1735/shadow.git
cd shadow
make build

# Setup authentication (auto-detects Claude Code OAuth)
./shadow auth-gen

# Check authentication status
./shadow auth-status

# Basic scan with AI analysis
./shadow scan example.com --ai-analysis
```

## Installation

### From Source

```bash
git clone https://github.com/kumaraguru1735/shadow.git
cd shadow
make build
sudo make install  # Optional: installs to /usr/local/bin
```

### Prerequisites

- Go 1.22 or higher
- Claude Code installed (for OAuth) OR Anthropic API key
- External tools (optional): nmap, subfinder, whatweb

## Authentication

Shadow supports two authentication methods:

### 1. Claude Code OAuth (Recommended)

If you have Claude Code installed, Shadow automatically detects and uses your OAuth credentials:

```bash
# Check authentication status
./shadow auth-status

# Generate/extract OAuth credentials
./shadow auth-gen

# Refresh expired tokens
./shadow auth-refresh

# Create backup of credentials
./shadow auth-backup
```

### 2. API Key Authentication

```bash
# Interactive setup wizard
./shadow auth-setup --api-key sk-ant-your-key

# Or set environment variable
export ANTHROPIC_API_KEY='sk-ant-your-key'

# Verify authentication works
./shadow auth-check
```

### Authentication Commands

| Command | Description |
|---------|-------------|
| `auth-check` | Quick authentication verification |
| `auth-status` | Detailed status with expiration times |
| `auth-gen` | Auto-generate authentication setup |
| `auth-setup` | Interactive setup wizard |
| `auth-refresh` | Refresh OAuth tokens |
| `auth-backup` | Create timestamped credential backups |

## Usage

### Security Scanning

```bash
# Basic scan
./shadow scan example.com

# Scan with AI analysis
./shadow scan example.com --ai-analysis

# Deep scan with specific profile
./shadow scan example.com --profile deep --ai-analysis

# Custom output format
./shadow scan example.com --format yaml --output results.yaml

# Adjust thread count
./shadow scan example.com --threads 100
```

### Subdomain Discovery

```bash
# Discover subdomains
./shadow subdomain example.com

# Save results
./shadow subdomain example.com --output subdomains.txt
```

## Configuration

Shadow can be configured via `~/.shadow/config.yaml`:

```yaml
# Anthropic Claude AI Settings
anthropic:
  api_key: ${ANTHROPIC_API_KEY}  # Set via environment variable
  model: claude-sonnet-4.5-20250929
  max_tokens: 4096

# Scanning Configuration
scanning:
  threads: 50
  timeout: 30s
  rate_limit: 100

# AI Analysis Configuration
ai:
  enabled: true
  auto_analyze: false
  retry_attempts: 3     # Auto-retry on failures
  retry_delay: 15s      # Base delay (exponential backoff)
```

Generate a default config:

```bash
./shadow auth-gen  # Creates config automatically
```

## Advanced AI Features

Shadow includes production-tested AI patterns from [OpenClaw](https://github.com/openclaw/openclaw):

### Intelligent Retry Logic
- **3 automatic retries** with exponential backoff (15s, 30s, 45s)
- Detects and handles: rate limits, timeouts, network issues
- Context-aware cancellation support

### Extended Timeouts
- **10-minute analysis timeout** (vs 2 min previously)
- **5-minute query timeout** for complex questions
- Handles large scan results with extended thinking

### Smart Error Detection
Automatically retries on:
- Rate limiting (429 errors)
- Context deadline exceeded
- Temporary network issues
- Connection failures

### System Prompt Engineering
- Expert security analyst persona
- Structured analysis requests
- Consistent, actionable output

See [ADVANCED_AI_FEATURES.md](ADVANCED_AI_FEATURES.md) for detailed documentation.

## Architecture

```
shadow/
├── cmd/
│   └── shadow/          # CLI entry point with auth commands
├── internal/
│   ├── scanner/         # Core scanning engine
│   ├── ai/              # Claude AI integration
│   │   ├── pi_client.go           # Basic OAuth client
│   │   ├── advanced_client.go     # Advanced retry/error handling
│   │   └── auth_manager.go        # Authentication lifecycle
│   └── modules/         # Security modules
├── pkg/
│   └── models/          # Data models
└── docs/               # Documentation
```

## Performance

- **Scan Speed**: ~5-10 minutes for medium-sized website
- **AI Analysis**: Up to 10 minutes for complex scans (with auto-retry)
- **Concurrency**: 50+ simultaneous scan threads
- **Memory**: <500MB for most scans
- **Binary Size**: ~15MB (single binary)

## Troubleshooting

### "context deadline exceeded" error
**Fixed!** Shadow now uses 10-minute timeouts and automatic retry with exponential backoff.

### Authentication Issues

```bash
# Check detailed status
./shadow auth-status

# Regenerate authentication
./shadow auth-gen

# Refresh expired tokens
./shadow auth-refresh

# Validate authentication works
./shadow auth-check
```

### Rate Limiting
Automatically handled with exponential backoff. No manual intervention needed.

### Missing External Tools
Shadow will skip unavailable tools gracefully. Install for better results:

```bash
# Ubuntu/Debian
sudo apt install nmap

# Install subfinder
go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest

# Install whatweb
sudo apt install whatweb
```

## Security & Ethics

⚠️ **Authorization Required**: Only scan systems you own or have explicit permission to test.

Shadow includes built-in safeguards:
- Scope restriction enforcement
- Rate limiting to prevent abuse
- Credential protection (see .gitignore)
- Audit logging

## Development

### Building

```bash
# Build binary
make build

# Run tests
make test

# Clean build artifacts
make clean

# Install to /usr/local/bin
sudo make install
```

### Project Structure

See [BUILD.md](BUILD.md) for detailed build instructions.

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Submit a pull request

Ensure no sensitive data (OAuth tokens, API keys) is committed.

## Roadmap

- [x] v0.1.0 - Core scanning + Basic AI integration
- [x] v0.1.1 - Advanced AI with retry logic (OpenClaw patterns)
- [x] v0.1.2 - Authentication management system
- [ ] v0.2.0 - Additional security modules
- [ ] v0.3.0 - Web UI dashboard
- [ ] v1.0.0 - Production ready

## Acknowledgments

Built with:
- [Claude AI](https://anthropic.com) - Intelligent analysis with extended thinking
- [pi-golang](https://github.com/joshp123/pi-golang) - Go wrapper for Claude integration
- [OpenClaw](https://github.com/openclaw/openclaw) - Production-tested AI patterns
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- Security community tools (nmap, subfinder, whatweb)

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**Made with ❤️ for the security community**

*Defensive security tool - use responsibly with proper authorization*
