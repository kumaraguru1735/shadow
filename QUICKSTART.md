# Shadow Quick Start 🚀

## ✅ What's Already Done

Shadow is configured and ready to use!

- ✅ Binary built: `/opt/lampp/htdocs/shadow/shadow` (5.8MB)
- ✅ Config created: `~/.shadow/config.yaml`
- ✅ All commands available

## 🎯 Run Your First Scan (30 seconds)

### Option 1: Quick Scan (No API Key Needed)

```bash
cd /opt/lampp/htdocs/shadow

# Basic scan
./shadow scan example.com
```

**What happens:**
1. Shadow asks for permission (type `yes`)
2. Runs security checks
3. Shows findings
4. Saves results

**Time:** ~2-5 minutes
**Cost:** $0 (no AI)

### Option 2: With AI Analysis (OAuth or API Key)

Shadow automatically uses Claude Code OAuth token when available!

```bash
# If using Claude Code - no setup needed, OAuth is automatic!
./shadow scan example.com --ai-analysis

# OR set API key manually
export ANTHROPIC_API_KEY='sk-ant-your-key-here'
./shadow scan example.com --ai-analysis
```

**What happens:**
1. Same as Option 1
2. Plus: Claude AI analyzes findings
3. Plus: Intelligent risk prioritization
4. Plus: Remediation recommendations

**Time:** ~3-8 minutes
**Cost:** $2-5 (AI analysis)

## 📋 All Commands

### Scanning

```bash
# Quick scan (2-3 minutes, basic checks)
./shadow scan example.com --profile quick

# Standard scan (10-15 minutes, recommended)
./shadow scan example.com --profile standard

# Deep scan (30-60 minutes, comprehensive)
./shadow scan example.com --profile deep

# With specific modules only
./shadow scan example.com --modules subdomain,portscan

# Save to file
./shadow scan example.com --output scan.json
```

### Discovery

```bash
# Find subdomains
./shadow subdomain example.com

# Scan ports
./shadow portscan example.com

# Scan all ports (slow!)
./shadow portscan example.com --ports 1-65535

# Fast port scan (top 100)
./shadow portscan example.com --fast
```

### Security Checks

```bash
# Check SSL/TLS configuration
./shadow ssl example.com
```

### AI Features (requires API key)

```bash
# Analyze scan results
./shadow analyze <scan-id>

# Generate report
./shadow report <scan-id> --format html
./shadow report <scan-id> --format pdf
./shadow report <scan-id> --format json

# Ask questions about results
./shadow query <scan-id> "What are the critical issues?"
./shadow query <scan-id> "How do I fix the SSL problems?"
```

## ⚙️ Configuration

### AI Authentication (Multiple Options)

Shadow supports multiple authentication methods with automatic detection:

**Option 1: Claude Code OAuth (Automatic, Preferred)**
- ✅ No configuration needed!
- ✅ Uses your existing Claude Code authentication
- ✅ OAuth token automatically detected from:
  - `~/.claude/.credentials.json` (Claude Code)
  - `~/.claude/oauth.json`
  - `~/.pi/agent/oauth.json`
  - `~/.config/anthropic/oauth.json`

**Option 2: Environment Variable**
```bash
export ANTHROPIC_API_KEY='sk-ant-your-key-here'
```

**Option 3: Config File**
Edit `~/.shadow/config.yaml`:
```yaml
anthropic:
  api_key: sk-ant-your-key-here
```

**Check Authentication Status:**
```bash
./shadow auth-check
```

### Other Settings

```yaml
scanning:
  threads: 50          # Concurrent threads
  timeout: 30s         # Request timeout
  rate_limit: 100      # Max requests/second

ai:
  enabled: true        # Enable AI features
  auto_analyze: false  # Auto-analyze after scan
```

## 📊 Example Output

```bash
$ ./shadow scan example.com --profile quick

🕵️  Shadow v0.1.0
🎯 Target: example.com
📋 Profile: quick
🧵 Threads: 50

⚠️  AUTHORIZATION REQUIRED
You are about to scan: example.com
Do you have explicit permission to test this target? (yes/no): yes

🔍 Starting reconnaissance...
  ▶ Running Basic Security module...
    ✓ Found 1 findings
  ▶ Running Security Headers module...
    ✓ Found 3 findings

✅ Scan completed in 2m 34s
📊 Scan ID: abc123-def456-789
🔍 Findings: 4

Results saved to: scan_results/example.com_20250210/
```

## 💰 Cost Breakdown

| Scan Type | Time | Cost (No AI) | Cost (With AI) |
|-----------|------|--------------|----------------|
| **Quick** | 2-5 min | **$0** | $1-2 |
| **Standard** | 10-15 min | **$0** | $2-5 |
| **Deep** | 30-60 min | **$0** | $5-10 |

**Pro Tip:** Run scans without AI first, then use `./shadow analyze <scan-id>` to analyze only when needed!

## 🎯 Real-World Examples

### Example 1: Quick Health Check

```bash
# Fast security check before deployment
./shadow scan staging.myapp.com --profile quick
```

**Use Case:** Pre-deployment security check
**Time:** 2-3 minutes
**Cost:** $0

### Example 2: Weekly Security Audit

```bash
# Comprehensive weekly scan
./shadow scan myapp.com --profile deep --ai-analysis
```

**Use Case:** Regular security monitoring
**Time:** 30 minutes
**Cost:** $5-10

### Example 3: Continuous Integration

```bash
# In your CI/CD pipeline
./shadow scan $TARGET --profile quick --output scan.json --format json

# Fail build if critical issues found
if [ $(jq '.findings[] | select(.severity=="critical") | length' scan.json) -gt 0 ]; then
  echo "❌ Critical vulnerabilities found!"
  exit 1
fi
```

**Use Case:** Automated security in CI/CD
**Time:** 2-5 minutes per build
**Cost:** $0

### Example 4: Bug Bounty Reconnaissance

```bash
# Discover subdomains
./shadow subdomain target.com > subdomains.txt

# Scan each subdomain
for sub in $(cat subdomains.txt); do
  ./shadow scan $sub --profile standard
done
```

**Use Case:** Bug bounty hunting
**Time:** Varies
**Cost:** $0-5 per target

## 🐛 Troubleshooting

### "Permission denied"

```bash
chmod +x ./shadow
```

### "API key not set" (when using --ai-analysis)

```bash
export ANTHROPIC_API_KEY='your-key-here'
```

### "Target unreachable"

Check:
- Target URL is correct
- You have internet connection
- Target isn't blocking your IP

### "Module failed"

Some modules require external tools (nmap, subfinder, etc.). Shadow will skip unavailable modules gracefully.

## 📚 What's Next?

1. ✅ **Run your first scan** (example.com)
2. 📖 **Read** [WHY_SHADOW.md](WHY_SHADOW.md) to understand benefits
3. 🏗️ **Learn architecture** in [ARCHITECTURE.md](ARCHITECTURE.md)
4. ⚙️ **Customize** `~/.shadow/config.yaml`
5. 🤖 **Try AI features** with your API key

## 🎉 You're Ready!

Shadow is configured and working. Try it now:

```bash
./shadow scan example.com
```

**Happy Hunting! 🕵️**

---

**Need Help?**
- Run `./shadow --help` for commands
- Check documentation files
- Review example configs in `configs/`
