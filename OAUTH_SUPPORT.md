# Shadow OAuth Authentication Support

## Overview

Shadow seamlessly integrates with Claude Code's OAuth authentication, eliminating the need for manual API key management in most cases.

## 🎯 Key Features

- ✅ **Automatic OAuth Detection** - Works out-of-the-box with Claude Code
- ✅ **Multiple Auth Methods** - OAuth, API Key, or Config File
- ✅ **Zero Configuration** - Uses existing Claude Code authentication
- ✅ **Same as OpenClaw** - Uses proven `pi-golang` library
- ✅ **Seamless Fallback** - Automatically tries multiple auth methods

## 🔐 How It Works

Shadow uses the [pi-golang](https://github.com/joshp123/pi-golang) library (same as OpenClaw) which automatically searches for authentication in this order:

1. **OAuth Tokens** (preferred):
   - `~/.claude/.credentials.json` (Claude Code - detected!)
   - `~/.claude/oauth.json` (Claude Code alternative)
   - `~/.pi/agent/oauth.json` (pi CLI)
   - `~/.config/claude/oauth.json`
   - `~/.config/anthropic/oauth.json`
   - `~/.clawdis/credentials/oauth.json`

2. **Environment Variables**:
   - `ANTHROPIC_OAUTH_TOKEN`
   - `ANTHROPIC_API_KEY`
   - `ANTHROPIC_TOKEN_FILE`

3. **pi CLI** (if installed):
   - Falls back to pi CLI's built-in authentication

## 📋 Setup Methods

### Method 1: Use Claude Code (Recommended)

If you're using Claude Code, you're already set up! Just run:

```bash
./shadow scan example.com --ai-analysis
```

Shadow automatically uses Claude Code's OAuth token.

### Method 2: Install pi CLI

```bash
# Install pi CLI globally
npm install -g @mariozechner/pi-coding-agent

# Authenticate with pi
pi auth

# Shadow will now use pi's authentication
./shadow scan example.com --ai-analysis
```

### Method 3: Manual API Key

```bash
# Set API key
export ANTHROPIC_API_KEY='sk-ant-your-key-here'

# Run Shadow
./shadow scan example.com --ai-analysis
```

## 🧪 Testing Authentication

Check your authentication status:

```bash
./shadow auth-check
```

Output shows:
- Which authentication method is detected
- Whether AI client initializes successfully
- Available authentication options

## 🔄 Comparison with OpenClaw

Shadow uses **exactly the same authentication mechanism** as OpenClaw:

| Feature | OpenClaw | Shadow |
|---------|----------|--------|
| **Library** | pi-golang | pi-golang ✓ |
| **OAuth Detection** | Automatic | Automatic ✓ |
| **Fallback Methods** | Multiple | Multiple ✓ |
| **Configuration** | Zero-config | Zero-config ✓ |

## 💡 Common Scenarios

### Scenario 1: Using Claude Code
```bash
# No setup needed - OAuth is automatic!
./shadow scan example.com --ai-analysis
./shadow query <scan-id> "What are the critical issues?"
```

### Scenario 2: Using pi CLI
```bash
# Install pi (one-time)
npm install -g @mariozechner/pi-coding-agent

# Authenticate (one-time)
pi auth

# Use Shadow
./shadow scan example.com --ai-analysis
```

### Scenario 3: Using API Key
```bash
# Export API key
export ANTHROPIC_API_KEY='sk-ant-...'

# Use Shadow
./shadow scan example.com --ai-analysis
```

### Scenario 4: CI/CD Pipeline
```bash
# Set API key in CI environment
export ANTHROPIC_API_KEY=${{ secrets.ANTHROPIC_API_KEY }}

# Run Shadow in pipeline
./shadow scan $TARGET --ai-analysis --output report.json
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────┐
│                    Shadow                        │
│  ┌────────────────────────────────────────┐    │
│  │      internal/ai/pi_client.go          │    │
│  │  Uses pi-golang for authentication     │    │
│  └─────────────────┬────────────────────── ┘    │
│                    │                             │
│                    ▼                             │
│  ┌────────────────────────────────────────┐    │
│  │       github.com/joshp123/pi-golang    │    │
│  │   Manages OAuth, API keys, pi CLI      │    │
│  └─────────────────┬────────────────────── ┘    │
│                    │                             │
└────────────────────┼─────────────────────────────┘
                     │
                     ▼
       ┌─────────────────────────────┐
       │  Authentication Sources     │
       ├─────────────────────────────┤
       │ • ~/.claude/oauth.json      │
       │ • ANTHROPIC_API_KEY         │
       │ • pi CLI authentication     │
       │ • ~/.pi/agent/oauth.json    │
       └─────────────────────────────┘
```

## 🚀 Benefits

1. **Zero Configuration** - Works immediately with Claude Code
2. **Secure** - OAuth tokens never exposed in commands
3. **Flexible** - Multiple authentication methods
4. **Proven** - Same library as OpenClaw
5. **Automatic** - No manual token management

## 🔧 Troubleshooting

### "AI client initialization failed"

**Solution 1: Install pi CLI**
```bash
npm install -g @mariozechner/pi-coding-agent
```

**Solution 2: Use API Key**
```bash
export ANTHROPIC_API_KEY='your-key'
```

**Solution 3: Check OAuth Token**
```bash
ls -la ~/.claude/oauth.json
```

### "No authentication found"

This warning can be ignored if:
- pi CLI is installed (Shadow will use it automatically)
- `./shadow auth-check` shows "AI client initialized successfully!"

The warning only indicates OAuth file not found at specific paths, but pi CLI handles auth internally.

## 📚 References

- [pi-golang Library](https://github.com/joshp123/pi-golang)
- [OpenClaw](https://github.com/openclaw/openclaw) - Reference implementation
- [pi Coding Agent](https://github.com/badlogic/pi-mono)
- [Anthropic API Docs](https://docs.anthropic.com)

## ✅ Verification

To verify OAuth support is working:

```bash
# Step 1: Check authentication
./shadow auth-check

# Step 2: Run a quick test
./shadow --version

# Step 3: Test AI initialization (dry run)
./shadow auth-check | grep "AI client initialized successfully"
```

If you see "✅ AI client initialized successfully!" - you're all set!

---

**Shadow + Claude Code OAuth = Seamless AI-Powered Security Scanning** 🚀
