# Shadow 🕵️

> AI-augmented cybersecurity reconnaissance and analysis platform

Shadow combines the power of Go's performance with Claude AI's intelligence to deliver comprehensive, automated security assessments.

## Features

- 🔐 **Claude Code OAuth** - Automatic authentication, zero configuration needed
- 🔍 **Network Intelligence** - Port scanning, subdomain discovery, SSL/TLS analysis
- 🌐 **Web Security Analysis** - Crawling, header validation, vulnerability detection
- 🤖 **AI-Powered Analysis** - Claude AI for intelligent vulnerability assessment
- 📊 **Smart Reporting** - Executive summaries and technical reports
- 🔄 **CI/CD Integration** - Jenkins, GitLab CI, GitHub Actions support
- 🚀 **High Performance** - Go-powered concurrency, distributed scanning
- 🔌 **Extensible** - Plugin architecture for custom modules

## Quick Start

```bash
# Install
go install github.com/yourusername/shadow@latest

# Basic scan
shadow scan --target https://example.com

# With AI analysis
shadow scan --target https://example.com --ai-analysis

# Deep scan with all modules
shadow scan --target https://example.com --profile deep
```

## Installation

### From Source

```bash
git clone https://github.com/yourusername/shadow.git
cd shadow
go build -o shadow cmd/shadow/main.go
sudo mv shadow /usr/local/bin/
```

### Using Go

```bash
go install github.com/yourusername/shadow@latest
```

## Configuration

Create `~/.shadow/config.yaml`:

```yaml
anthropic:
  api_key: ${ANTHROPIC_API_KEY}
  model: claude-sonnet-4.5-20250929

scanning:
  threads: 50
  timeout: 30s
  rate_limit: 100  # requests per second

modules:
  - subdomain
  - port_scan
  - web_analysis
  - ssl_check
```

## Usage Examples

### Basic Reconnaissance

```bash
# Quick scan
shadow scan --target example.com

# Subdomain discovery
shadow subdomain --target example.com

# Port scanning
shadow portscan --target example.com --ports 1-65535

# SSL/TLS analysis
shadow ssl --target example.com
```

### AI-Powered Analysis

```bash
# Analyze findings with AI
shadow analyze --scan-id abc123

# Generate report
shadow report --scan-id abc123 --format pdf

# Ask questions about results
shadow query --scan-id abc123 "What are the critical vulnerabilities?"
```

### Advanced Features

```bash
# Distributed scanning
shadow scan --target example.com --workers 10

# CI/CD integration
shadow scan --target $CI_TARGET --format json --output scan.json

# Scheduled scanning
shadow schedule --target example.com --cron "0 2 * * *"

# Compare scans
shadow diff --baseline scan1.json --current scan2.json
```

## Architecture

```
shadow/
├── cmd/
│   └── shadow/          # CLI entry point
├── internal/
│   ├── scanner/         # Core scanning engine
│   ├── ai/              # Claude AI integration
│   ├── modules/         # Security modules
│   ├── report/          # Reporting engine
│   └── database/        # Data persistence
├── pkg/
│   ├── models/          # Data models
│   └── utils/           # Utilities
└── plugins/             # Plugin system
```

## Modules

### Network Intelligence
- **Subdomain Discovery** - Multiple techniques (DNS, CT logs, brute force)
- **Port Scanning** - Fast TCP/UDP scanning with service detection
- **SSL/TLS Analysis** - Certificate validation, vulnerability checks
- **DNS Enumeration** - Zone transfers, DNSSEC validation

### Web Security
- **Web Crawler** - JavaScript-aware crawling
- **Header Analysis** - Security header validation
- **Technology Detection** - Stack fingerprinting
- **Vulnerability Detection** - OWASP Top 10 checks

### AI Analysis
- **Risk Prioritization** - Intelligent vulnerability ranking
- **Attack Chain Detection** - Identify exploitation paths
- **Remediation Guidance** - Actionable fix recommendations
- **Natural Language Reporting** - Executive summaries

## Performance

- **Scan Speed**: ~10 minutes for medium-sized website
- **Concurrency**: 50+ simultaneous targets
- **Memory**: <500MB for most scans
- **Binary Size**: ~15MB (single binary, no dependencies)

## Security & Ethics

⚠️ **Authorization Required**: Only scan systems you own or have explicit permission to test.

Shadow includes built-in safeguards:
- Permission validation prompts
- Scope restriction enforcement
- Audit logging
- Rate limiting to prevent abuse

## Roadmap

- [ ] v0.1.0 - Core scanning + AI integration
- [ ] v0.2.0 - Plugin system
- [ ] v0.3.0 - Web UI dashboard
- [ ] v0.4.0 - Distributed scanning
- [ ] v0.5.0 - CI/CD integrations
- [ ] v1.0.0 - Production ready

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

Built with:
- [Claude AI](https://anthropic.com) - Intelligent analysis
- [Nuclei](https://github.com/projectdiscovery/nuclei) - Vulnerability templates
- [Subfinder](https://github.com/projectdiscovery/subfinder) - Subdomain discovery
- And many other amazing open-source tools

---

**Made with ❤️ for the security community**
