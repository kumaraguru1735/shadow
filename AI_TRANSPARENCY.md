# AI Transparency & Security

## Overview

Shadow provides **full transparency** about what AI is doing, similar to Claude Code. You can see exactly what prompts are sent, what tools are used, and what the AI is analyzing.

## 🔍 Why Transparency Matters

For security tools, you MUST know:
- ✅ What prompts are being sent to AI
- ✅ What data is being analyzed
- ✅ What tools/APIs are being called
- ✅ What the AI is actually doing
- ✅ No hidden operations

## 📋 What You See

### 1. AI Activity Log

Every time AI is called, you see:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔍 AI ACTIVITY LOG
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 Action: Sending prompt to Claude Opus 4.6
   🔧 Tool: Claude API (Anthropic)
   🎯 Model: claude-opus-4.6
   🧠 Thinking: HIGH (extended reasoning enabled)
   ⏱️  Timeout: 10 minutes
   📊 Input size: 2847 characters
   🔒 Security: Read-only analysis, no code execution
   📍 Endpoint: api.anthropic.com
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**What this tells you:**
- Exact API being called
- Which AI model
- Thinking depth
- Timeout limits
- Data size being sent
- Security constraints
- Where data is going

### 2. Prompt Preview

You see what Shadow asks the AI:

```
📤 PROMPT PREVIEW (First/Last Lines):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Initial Security Analysis

## Target
example.com

## Initial Findings
1. [Medium] Missing Security Headers
   Description: HSTS header not present
   Location: https://example.com

2. [Low] Directory Listing Enabled
   Description: /images/ directory is browsable
   Evidence: Index of /images/

## Your Task
Analyze these findings and THINK DEEPLY about:

1. **What do these findings tell us?**
   - What vulnerabilities might exist?
   - What attack vectors are possible?
   - What might the developers have overlooked?

... [Full prompt: 45 lines] ...

### NEXT STEPS
[Concrete actions to take]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**What this tells you:**
- Exact question being asked
- Data being analyzed
- Instructions given to AI
- Expected output format
- No hidden prompts

### 3. Real-Time Status Updates

You see what AI is doing in real-time:

```
⏳ WAITING FOR AI RESPONSE...
   The AI is now:
   • Reading and understanding the findings
   • Thinking critically about security implications
   • Generating hypotheses about vulnerabilities
   • Planning investigation strategies

   This typically takes 3-8 minutes with extended thinking.

🔄 API Call in Progress...
   Started: 14:23:45

   [... AI is working ...]

   Completed: 14:28:12

✅ AI Response Received
📊 Response size: 3542 characters
```

**What this tells you:**
- When API call started
- What AI is currently doing
- Expected duration
- When completed
- Response size received

### 4. AI Analysis Output

You see what AI found:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🧠 AI'S ANALYSIS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 AI's Critical Thinking:
"The missing HSTS header combined with directory listing
suggests insufficient security hardening. The application
appears to be in early deployment stage with default
configurations still active.

Key concerns:
1. HSTS absence enables MITM attacks
2. Directory listing could expose sensitive files
3. Likely other default configs remain unchanged"

### HYPOTHESES
1. Additional security headers likely missing
   - Indicators: HSTS absent, suggests poor security config
   - Severity: Medium
   - Investigation: Check CSP, X-Frame-Options, etc.

2. Sensitive files may be exposed
   - Indicators: Directory listing enabled
   - Severity: High if /backup or /config accessible
   - Investigation: Check common sensitive directories
```

**What this tells you:**
- AI's reasoning process
- Hypotheses generated
- Severity assessments
- Next investigation steps

## 🔒 Security Guarantees

### What Shadow Shows You

✅ **Every API call is logged**
- Tool name (Claude API)
- Endpoint (api.anthropic.com)
- Model used (claude-opus-4.6)

✅ **Every prompt is shown**
- Full context of what's asked
- Data being analyzed
- No hidden prompts

✅ **Security mode is explicit**
- "Read-only analysis, no code execution"
- Shown in every activity log
- Cannot be hidden

✅ **Timing is transparent**
- Start time logged
- Completion time logged
- Duration visible

### What Shadow DOESN'T Do

❌ **No hidden API calls**
- Every call is logged and visible
- No background operations
- No silent data transmission

❌ **No code execution by AI**
- AI only analyzes text
- Cannot run commands
- Cannot modify systems
- Explicitly stated in logs

❌ **No data persistence by default**
- API calls are stateless
- No automatic data storage
- You control all outputs

❌ **No hidden prompts**
- All instructions visible
- No secret system prompts
- Everything is transparent

## 📊 Example Full Session

```
🕵️  Shadow v0.1.0 - Autonomous Security Research
🎯 Target: example.com

🧠 Initializing Autonomous AI Security Researcher
   Model: Claude Opus 4.6 (most capable, extended thinking)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔍 Step 1: Running initial security scan...
✅ Initial scan complete: 2 findings

🤖 Step 2: Launching Autonomous AI Researcher...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🧠 Starting autonomous security research...
🎯 Target: example.com
📊 Initial findings: 2

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔬 ITERATION 1: Initial Analysis & Hypothesis Generation
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🤔 AI is thinking about what these findings might indicate...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 TRANSPARENCY: What AI is Being Asked
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔍 AI ACTIVITY LOG
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 Action: Sending prompt to Claude Opus 4.6
   🔧 Tool: Claude API (Anthropic)
   🎯 Model: claude-opus-4.6
   🧠 Thinking: HIGH (extended reasoning enabled)
   ⏱️  Timeout: 10 minutes
   📊 Input size: 2847 characters
   🔒 Security: Read-only analysis, no code execution
   📍 Endpoint: api.anthropic.com
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📤 PROMPT PREVIEW (First/Last Lines):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Initial Security Analysis

## Target
example.com

## Initial Findings
1. [Medium] Missing Security Headers
   Description: HSTS not enabled
   Location: https://example.com

2. [Low] Directory Listing
   Description: /images/ browsable
...

[Full prompt: 45 lines]

### NEXT STEPS
[Concrete actions to take]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

⏳ WAITING FOR AI RESPONSE...
   The AI is now:
   • Reading and understanding the findings
   • Thinking critically about security implications
   • Generating hypotheses about vulnerabilities
   • Planning investigation strategies

   This typically takes 3-8 minutes with extended thinking.

🔄 API Call in Progress...
   Started: 14:23:45
   Completed: 14:28:12

✅ AI Response Received
📊 Response size: 3542 characters

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🧠 AI'S ANALYSIS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Initial analysis complete

📋 AI's Critical Thinking:
[AI's detailed reasoning...]

[Continue with remaining iterations...]
```

## 🎯 Transparency Features

### For Every AI Call

| Feature | What You See | Why It Matters |
|---------|--------------|----------------|
| **Tool Used** | Claude API (Anthropic) | Know what service is called |
| **Model** | claude-opus-4.6 | Know which AI model |
| **Thinking Mode** | HIGH (extended reasoning) | Know reasoning depth |
| **Timeout** | 10 minutes | Know max wait time |
| **Input Size** | Character count | Know data sent |
| **Security Mode** | Read-only, no execution | Know safety constraints |
| **Endpoint** | api.anthropic.com | Know where data goes |
| **Prompt** | First/last lines | Know what's asked |
| **Timing** | Start/end timestamps | Track duration |
| **Response Size** | Character count | Know data received |

### Activity Log Types

1. **API Call Log**
   - Before each AI call
   - Shows all parameters
   - Includes security mode

2. **Prompt Preview**
   - Shows actual prompt
   - First 10 + last 5 lines
   - Total line count

3. **Status Updates**
   - Real-time progress
   - What AI is doing
   - Expected duration

4. **Result Log**
   - Response received
   - Size information
   - Completion time

## 🔐 Data Privacy

### What Gets Sent to Anthropic

✅ **Sent to API:**
- Target hostname (e.g., "example.com")
- Security findings from scan
- Severity levels
- Technical descriptions

❌ **NOT sent:**
- User credentials
- Session tokens
- Personal data (unless in scan results)
- System information
- File contents (unless explicitly scanned)

### How to Verify

1. **Check Activity Logs**
   - Every API call is logged
   - Prompt preview shows data

2. **Review Prompt**
   - Exact text sent is shown
   - No hidden data

3. **Monitor Network**
   - Only api.anthropic.com contacted
   - Can verify with netstat/wireshark

## 💡 Best Practices

### Review Before Scanning

1. **Check what will be scanned**
   ```bash
   # Review target first
   curl -I https://target.com
   ```

2. **Start with read-only analysis**
   ```bash
   # Research mode is read-only
   ./shadow research target.com
   ```

3. **Review activity logs**
   - Read each API call log
   - Verify prompt content
   - Check security mode

### During Scanning

1. **Watch for prompts**
   - Review prompt previews
   - Verify data being sent
   - Check security mode

2. **Monitor timing**
   - Note unusual delays
   - Check completion times
   - Verify timeouts

3. **Review findings**
   - Read AI's reasoning
   - Verify conclusions
   - Check hypotheses

## 🎓 Comparison to Other Tools

### Traditional Security Scanners
```
Runs scans → Generates report
[Black box - you don't see what's happening]
```

### Shadow with Transparency
```
Shows every action:
1. What tool is called
2. What prompt is sent
3. What AI is analyzing
4. What results come back
5. Timing of all operations

[Crystal clear - full visibility]
```

## 📊 Summary

Shadow provides **Claude Code-level transparency** for AI operations:

✅ **Every AI call is logged**
✅ **Every prompt is shown**
✅ **Security mode is explicit**
✅ **Timing is transparent**
✅ **No hidden operations**
✅ **Read-only by default**

You know **exactly** what the AI is doing at all times.

---

**Try it with full transparency:**
```bash
./shadow research yourtarget.com
```

Watch every step of the AI's analysis! 🔍👁️
