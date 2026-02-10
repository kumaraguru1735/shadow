# Getting Started with Shadow 🕵️

## Quick Start (5 minutes)

### 1. Build Shadow

```bash
cd /opt/lampp/htdocs/shadow
make build
```

### 2. Set up API Key

```bash
export ANTHROPIC_API_KEY="your-key-here"
```

### 3. Run Your First Scan

```bash
./shadow scan example.com
```

That's it! Shadow will:
- ✅ Check for authorization
- ✅ Perform security reconnaissance
- ✅ Generate findings
- ✅ Save results

## Basic Commands

### Scan a Target

```bash
# Quick scan (fast, essential checks)
./shadow scan --target example.com --profile quick

# Standard scan (recommended)
./shadow scan --target example.com --profile standard

# Deep scan (comprehensive, slower)
./shadow scan --target example.com --profile deep
```

### Scan with AI Analysis

```bash
# Enable AI-powered vulnerability analysis
./shadow scan --target example.com --ai-analysis
```

### Subdomain Discovery

```bash
./shadow subdomain example.com
```

### Port Scanning

```bash
# Scan common ports
./shadow portscan example.com

# Scan all ports
./shadow portscan example.com --ports 1-65535

# Fast scan (top 100 ports)
./shadow portscan example.com --fast
```

### SSL/TLS Analysis

```bash
./shadow ssl example.com
```

## Scan Profiles

| Profile | Speed | Depth | Use Case |
|---------|-------|-------|----------|
| `quick` | ⚡ Fast (2-5 min) | Basic | Quick health check |
| `standard` | 🚀 Medium (10-15 min) | Good | Regular assessment |
| `deep` | 🐢 Slow (30-60 min) | Comprehensive | Thorough audit |

## Configuration

### Create Config File

```bash
mkdir -p ~/.shadow
cp configs/config.example.yaml ~/.shadow/config.yaml
```

### Edit Configuration

```yaml
# ~/.shadow/config.yaml
anthropic:
  api_key: ${ANTHROPIC_API_KEY}

scanning:
  threads: 50
  timeout: 30s
  rate_limit: 100
```

## Output

Shadow saves results in the current directory:

```
scan_results/
└── example.com_20250210_143022/
    ├── scan.json           # Raw scan data
    ├── findings.json       # Discovered vulnerabilities
    ├── ai_analysis.json    # AI analysis (if enabled)
    └── report.html         # HTML report
```

## AI-Powered Features

### Analyze Scan Results

```bash
# After running a scan, analyze with AI
./shadow analyze <scan-id>
```

### Ask Questions

```bash
# Natural language queries about your scan
./shadow query <scan-id> "What are the critical vulnerabilities?"
./shadow query <scan-id> "How do I fix the SSL issues?"
```

### Generate Reports

```bash
# HTML report (default)
./shadow report <scan-id>

# PDF report
./shadow report <scan-id> --format pdf

# JSON export
./shadow report <scan-id> --format json
```

## Advanced Usage

### Custom Modules

```bash
# Run specific modules only
./shadow scan --target example.com --modules subdomain,portscan
```

### Multi-Threading

```bash
# Increase concurrent threads for faster scanning
./shadow scan --target example.com --threads 100
```

### Output to File

```bash
# Save results to specific file
./shadow scan --target example.com --output results.json
```

## Troubleshooting

### "ANTHROPIC_API_KEY not set"

```bash
export ANTHROPIC_API_KEY="sk-ant-your-key"
```

Or add to `~/.shadow/config.yaml`

### "Permission denied"

Make Shadow executable:
```bash
chmod +x ./shadow
```

Or install system-wide:
```bash
sudo make install
```

### "Target unreachable"

Check:
1. Target is accessible from your network
2. Firewall isn't blocking outbound connections
3. Target URL is correct (include `https://` if needed)

## Best Practices

### 1. Always Get Permission ⚠️

**NEVER** scan systems without authorization. Shadow includes permission checks, but **you** are responsible for ethical use.

### 2. Start with Quick Scans

Use `--profile quick` first to get rapid feedback, then run deeper scans if needed.

### 3. Use AI Analysis Selectively

AI analysis costs API credits. Run it on important scans, not every quick check.

### 4. Respect Rate Limits

Configure appropriate rate limiting to avoid overwhelming target systems:

```yaml
scanning:
  rate_limit: 10  # requests per second
```

### 5. Save Results

Always use `--output` to save results for future reference and comparison.

## Next Steps

1. ✅ Run your first scan
2. 📖 Read [ARCHITECTURE.md](ARCHITECTURE.md) to understand how Shadow works
3. 🔧 Customize `~/.shadow/config.yaml` for your needs
4. 🧪 Try different scan profiles
5. 🤖 Experiment with AI analysis features
6. 📊 Generate reports and share with your team

## Getting Help

```bash
# Command help
./shadow --help
./shadow scan --help

# Check version
./shadow --version
```

## What's Next?

Once you're comfortable with basic usage:

- Explore [advanced configuration options](configs/config.example.yaml)
- Learn about [Shadow's architecture](ARCHITECTURE.md)
- Understand [why Shadow beats other tools](WHY_SHADOW.md)
- Contribute to the project!

---

**Happy Hunting! 🎯**
