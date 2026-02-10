#!/bin/bash

# Shadow OAuth Setup Script
# Extracts OAuth token from Claude Code credentials and places it where pi-golang can find it

set -e

CLAUDE_CREDS="$HOME/.claude/.credentials.json"
OAUTH_FILE="$HOME/.claude/oauth.json"

echo "╔════════════════════════════════════════════════════════════╗"
echo "║         Shadow OAuth Setup from Claude Code               ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check if Claude Code credentials exist
if [ ! -f "$CLAUDE_CREDS" ]; then
    echo "❌ Claude Code credentials not found at $CLAUDE_CREDS"
    echo ""
    echo "💡 Alternatives:"
    echo "   1. Run Claude Code to create credentials"
    echo "   2. Set ANTHROPIC_API_KEY environment variable"
    echo "   3. Install pi CLI: npm install -g @mariozechner/pi-coding-agent"
    exit 1
fi

echo "✓ Found Claude Code credentials at $CLAUDE_CREDS"

# Check if jq is available
if ! command -v jq &> /dev/null; then
    echo "❌ jq is required but not installed"
    echo "   Install: sudo apt-get install jq"
    exit 1
fi

echo "✓ jq is available"
echo ""

# Extract OAuth token
echo "🔐 Extracting OAuth token..."
if jq '.claudeAiOauth' "$CLAUDE_CREDS" > "$OAUTH_FILE.tmp"; then
    mv "$OAUTH_FILE.tmp" "$OAUTH_FILE"
    chmod 600 "$OAUTH_FILE"
    echo "✅ OAuth token extracted to $OAUTH_FILE"
else
    echo "❌ Failed to extract OAuth token"
    rm -f "$OAUTH_FILE.tmp"
    exit 1
fi

echo ""
echo "📊 Token Info:"
ACCESS_TOKEN=$(jq -r '.accessToken' "$OAUTH_FILE" | cut -c1-20)
EXPIRES_AT=$(jq -r '.expiresAt' "$OAUTH_FILE")
SCOPES=$(jq -r '.scopes | join(", ")' "$OAUTH_FILE")

echo "   Access Token: ${ACCESS_TOKEN}..."
echo "   Expires: $(date -d @$(($EXPIRES_AT / 1000)) '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo 'N/A')"
echo "   Scopes: $SCOPES"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "✅ OAuth Setup Complete!"
echo ""
echo "🧪 Testing authentication..."
./shadow auth-check | grep -E "(OAuth token found|AI client initialized)" || echo "⚠️  Auth check failed"
echo ""
echo "🚀 Ready to use:"
echo "   ./shadow scan example.com --ai-analysis"
echo ""
