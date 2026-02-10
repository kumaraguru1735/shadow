# Shadow + Claude Code OAuth Integration

## Overview

Shadow now automatically uses Claude Code's OAuth credentials stored at `~/.claude/.credentials.json`.

## 🎯 How It Works

### Claude Code Storage

Claude Code stores OAuth credentials at:
```
~/.claude/.credentials.json
```

Structure:
```json
{
  "claudeAiOauth": {
    "accessToken": "sk-ant-oat01-...",
    "refreshToken": "sk-ant-ort01-...",
    "expiresAt": 1770758958832,
    "scopes": [
      "user:inference",
      "user:mcp_servers",
      "user:profile",
      "user:sessions:claude_code"
    ],
    "subscriptionType": "max",
    "rateLimitTier": "default_claude_max_20x"
  }
}
```

### Shadow Integration

Shadow reads from the same location and extracts the OAuth token for use with pi-golang.

## 🚀 Quick Setup

### Option 1: Automatic Setup (Recommended)

```bash
# Run the setup script (one-time)
./setup_oauth.sh

# Verify
./shadow auth-check

# Start using Shadow
./shadow scan example.com --ai-analysis
```

### Option 2: Manual Setup

```bash
# Extract OAuth token from Claude Code credentials
jq '.claudeAiOauth' ~/.claude/.credentials.json > ~/.claude/oauth.json
chmod 600 ~/.claude/oauth.json

# Verify
./shadow auth-check
```

## 🔍 Verification

Check authentication status:

```bash
$ ./shadow auth-check

🔐 Claude AI Authentication Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ Claude Code OAuth token found at ~/.claude/.credentials.json

🧪 Testing AI connection...
✅ AI client initialized successfully!
✅ Shadow can use Claude AI for analysis
```

## 📋 OAuth Token Details

View your OAuth token info:

```bash
# Access token (truncated)
jq -r '.claudeAiOauth.accessToken' ~/.claude/.credentials.json | cut -c1-20

# Expiration date
jq -r '.claudeAiOauth.expiresAt' ~/.claude/.credentials.json

# Scopes
jq -r '.claudeAiOauth.scopes' ~/.claude/.credentials.json

# Subscription type
jq -r '.claudeAiOauth.subscriptionType' ~/.claude/.credentials.json
```

## 🔄 Token Refresh

Claude Code automatically refreshes the OAuth token. To update Shadow:

```bash
# Re-run setup script to sync the latest token
./setup_oauth.sh
```

Or manually:

```bash
jq '.claudeAiOauth' ~/.claude/.credentials.json > ~/.claude/oauth.json
```

## 🏗️ How Shadow Uses OAuth

```
Shadow CLI
    ↓
internal/ai/pi_client.go
    ↓
Detect OAuth from:
  1. ~/.claude/.credentials.json (Claude Code) ← Primary
  2. ~/.claude/oauth.json (extracted)
    ↓
pi-golang library
    ↓
Copy to ~/.shadow/pi-agent/oauth.json
    ↓
Spawn pi CLI with OAuth
    ↓
Claude AI API
```

## 📊 OAuth Detection Priority

Shadow checks these locations in order:

1. **`~/.claude/.credentials.json`** ← Claude Code (detected!)
2. `~/.claude/oauth.json` ← Extracted token
3. `~/.pi/agent/oauth.json` ← pi CLI
4. `~/.config/claude/oauth.json`
5. `~/.config/anthropic/oauth.json`
6. `ANTHROPIC_OAUTH_TOKEN` env var
7. `ANTHROPIC_API_KEY` env var

## 🔐 Security

- OAuth tokens are stored with `600` permissions (owner read/write only)
- Tokens automatically expire and refresh
- Shadow never exposes tokens in logs or output
- All communication uses HTTPS

## 💡 Troubleshooting

### "No authentication found"

**Solution 1: Run setup script**
```bash
./setup_oauth.sh
```

**Solution 2: Check Claude Code credentials**
```bash
ls -la ~/.claude/.credentials.json
jq '.claudeAiOauth' ~/.claude/.credentials.json
```

**Solution 3: Use API key instead**
```bash
export ANTHROPIC_API_KEY='your-key'
```

### "Token expired"

Claude Code automatically refreshes tokens. Re-run setup:
```bash
./setup_oauth.sh
```

### "jq not found"

Install jq:
```bash
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq
```

## 🎯 Usage Examples

### Basic Scan with OAuth

```bash
# OAuth is used automatically
./shadow scan example.com --ai-analysis
```

### AI-Powered Queries

```bash
./shadow query <scan-id> "What are the most critical vulnerabilities?"
```

### Generate Reports

```bash
./shadow report <scan-id> --format html
```

## ✅ Verification Checklist

- [ ] Claude Code credentials exist at `~/.claude/.credentials.json`
- [ ] Run `./setup_oauth.sh` successfully
- [ ] `./shadow auth-check` shows OAuth detected
- [ ] AI client initializes successfully
- [ ] Can run scans with `--ai-analysis`

## 📚 Related Documentation

- [OAUTH_SUPPORT.md](OAUTH_SUPPORT.md) - Complete OAuth guide
- [QUICKSTART.md](QUICKSTART.md) - Quick start guide
- [README.md](README.md) - Main documentation

## 🎉 Benefits

- ✅ **Zero Manual Configuration** - Uses existing Claude Code auth
- ✅ **Automatic Token Management** - Claude Code handles refresh
- ✅ **Secure** - Proper file permissions and HTTPS
- ✅ **Simple Setup** - One command to configure
- ✅ **Always Current** - Uses latest token from Claude Code

---

**Shadow uses the exact same OAuth credentials as Claude Code!** 🚀

