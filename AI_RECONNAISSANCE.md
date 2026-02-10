# AI-Driven Reconnaissance System

## Overview

Shadow's AI-Driven Reconnaissance is inspired by [OpenClaw's](https://github.com/openclaw/openclaw) autonomous AI workflow approach. Instead of running predefined scans, Shadow's AI **analyzes your target, plans the optimal reconnaissance strategy, and requests permissions as needed**.

## 🤖 How It Works

### Traditional Approach (Other Tools)
```
User → Hardcoded Scans → Results
```

### Shadow's AI Approach
```
User → AI Analyzes Target → AI Creates Custom Plan →
User Approves Plan → AI Requests Permissions → Execute → Results
```

## 🚀 Quick Start

### Basic Smart Scan
```bash
./shadow smart-scan example.com
```

**What happens:**
1. 🤖 AI analyzes the target
2. 📋 AI creates reconnaissance plan
3. ✅ You review and approve the plan
4. 🔐 AI requests root permissions if needed
5. 🔍 Reconnaissance executes phase by phase

### With Profile
```bash
./shadow smart-scan example.com --profile quick   # Fast
./shadow smart-scan example.com --profile standard # Balanced
./shadow smart-scan example.com --profile deep    # Thorough
```

## 📋 Example AI-Generated Plan

```
🎯 AI-Generated Reconnaissance Plan
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📍 Target: example.com

💭 Strategy:
This target appears to be a web application. I recommend a
three-phase approach: (1) DNS and subdomain enumeration to map
the attack surface, (2) port scanning to identify services,
(3) web technology fingerprinting for vulnerability research.

📋 Reconnaissance Phases (3):

1. ### PHASE 1: DNS & Subdomain Reconnaissance
   Priority: 🔴 Critical
   Map all DNS records and discover subdomains
   Tools needed:
      • dig (requires root: no) - DNS record enumeration
      • subfinder (requires root: no) - Subdomain discovery
   Expected outputs: DNS records, subdomain list

2. ### PHASE 2: Port & Service Discovery
   Priority: 🟠 High
   Identify open ports and running services
   Tools needed:
      • nmap [ROOT REQUIRED] - Fast SYN port scanning
   Expected outputs: Open ports, service versions

3. ### PHASE 3: Web Technology Fingerprinting
   Priority: 🟡 Medium
   Identify web technologies and frameworks
   Tools needed:
      • whatweb (requires root: no) - Technology detection
      • curl (requires root: no) - HTTP analysis
   Expected outputs: Tech stack, framework versions

🔐 Permissions:
   ⚠️  Root/sudo access required for some scans
   💡 Shadow will ask for permission before running privileged commands

🛠️  Required Tools:
   • dig
   • subfinder
   • nmap
   • whatweb
   • curl

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

❓ Execute this reconnaissance plan? (yes/no):
```

## 🔐 Permission Management

### Interactive Root Request

When AI determines root access is needed:

```
🔐 Root Permission Request
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 Tool: nmap
🎯 Purpose: Fast SYN port scanning (more stealthy than TCP connect)
💻 Command: sudo nmap -sS -p- example.com

⚠️  This command requires elevated privileges (root/sudo)
🔒 Shadow will ONLY run the specific command shown above
📊 This is needed for comprehensive security scanning

Allow this command? (yes/no/always):
```

### Response Options

- **yes** - Allow this specific command once
- **no** - Deny and skip to fallback
- **always** - Allow all future requests for this tool

### Security Features

✅ **Specific Command Approval**
- Shows exact command before running
- No blanket root access
- Each command requires approval

✅ **Fallback Strategies**
- If nmap with sudo denied → Use TCP connect
- If tool unavailable → Use Go-based alternative
- Always inform user of trade-offs

✅ **Alternative Methods**
- Suggests setcap instead of sudo
- Shows sudoers configuration options
- Explains security implications

## 🛡️ Permission Alternatives

### Option 1: Linux Capabilities (Recommended)

```bash
# One-time setup (more secure than sudo)
sudo setcap cap_net_raw,cap_net_admin,cap_net_bind_service+eip /usr/bin/nmap

# Now nmap works without sudo
nmap -sS example.com
```

**Benefits:**
- ✅ More secure than sudo
- ✅ Granular permissions
- ✅ No password prompts
- ✅ Audit trail

### Option 2: Sudoers Entry

```bash
# For persistent sudo access without password
echo 'yourusername ALL=(ALL) NOPASSWD: /usr/bin/nmap' | sudo tee /etc/sudoers.d/shadow-nmap
sudo chmod 440 /etc/sudoers.d/shadow-nmap
```

**Security Notes:**
- ⚠️ Only allow specific tools
- ⚠️ Use absolute paths
- ⚠️ Never use NOPASSWD for all commands

### Option 3: Non-Privileged Alternatives

```bash
# Shadow automatically falls back to:
nmap -sT example.com  # TCP connect (slower but no root)
# Or Go-based scanner (always available)
```

## 🎯 AI Decision Making

### What AI Analyzes

1. **Target Type**
   - Domain vs IP
   - Single host vs network range
   - Web app vs infrastructure

2. **Available Information**
   - DNS records
   - HTTP headers
   - Certificate information

3. **Reconnaissance Goals**
   - Attack surface mapping
   - Vulnerability identification
   - Technology profiling

4. **Tool Availability**
   - Checks what tools are installed
   - Plans around missing tools
   - Suggests alternatives

5. **Permission Requirements**
   - Determines if root needed
   - Explains why permissions required
   - Provides fallback options

### AI Reasoning Example

```
💭 AI Strategy:

"Target example.com appears to be a production web application
based on its HTTPS certificate and security headers. I recommend
starting with passive reconnaissance (DNS, subdomains) to avoid
detection, followed by targeted port scanning of discovered assets.

Root access would enable faster SYN scanning, but TCP connect
scans are acceptable for this target. I've prioritized stealth
over speed given the production environment indicators."
```

## 📊 Execution Flow

### Phase-by-Phase Execution

```
🚀 Executing reconnaissance plan...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📍 Phase 1/3: DNS & Subdomain Reconnaissance
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 Map all DNS records and discover subdomains

🔧 Running: dig
   Purpose: DNS record enumeration
   ✅ dig ready to execute

🔧 Running: subfinder
   Purpose: Subdomain discovery
   ✅ subfinder ready to execute

📍 Phase 2/3: Port & Service Discovery
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 Identify open ports and running services

🔧 Running: nmap
   Purpose: Fast SYN port scanning
   ⚠️  This tool requires root access

💡 Alternative: Use Linux Capabilities Instead of sudo
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📝 To allow nmap without sudo:
   sudo setcap cap_net_raw+eip /usr/bin/nmap

✅ Benefits:
   • More secure than sudo
   • No password prompts
   • Granular permissions

🔐 Root Permission Request
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 Tool: nmap
🎯 Purpose: Fast SYN port scanning
💻 Command: sudo nmap <args>

Allow this command? (yes/no/always): yes
   ✅ nmap ready to execute
```

## 🔄 Iterative Reconnaissance

Future enhancement: AI can adapt plan based on findings

```
Phase 1: DNS Enumeration
   ↓
   AI: "Found 15 subdomains, let's scan them"
   ↓
Phase 2: Subdomain Port Scanning
   ↓
   AI: "Found admin panel on port 8080, let's fingerprint"
   ↓
Phase 3: Targeted Web Analysis
```

## 📚 Use Cases

### Use Case 1: Bug Bounty
```bash
./shadow smart-scan target.bugcrowd.com --profile deep
```

**AI will:**
- Map entire attack surface
- Discover hidden subdomains
- Identify technology stack
- Find potential entry points

### Use Case 2: Penetration Testing
```bash
./shadow smart-scan client-app.com --profile deep
```

**AI will:**
- Create comprehensive recon plan
- Request necessary permissions
- Execute phased approach
- Provide detailed findings

### Use Case 3: Quick Assessment
```bash
./shadow smart-scan startup.io --profile quick
```

**AI will:**
- Focus on high-priority scans
- Skip root-requiring tools
- Provide rapid initial assessment

## 🔧 Configuration

### Customize AI Behavior

The AI planner uses Sonnet 4.5 with high thinking mode for optimal reconnaissance planning.

**Future**: Custom AI instructions and tool preferences

## 🎓 Comparison

### Traditional Scanner
```bash
nmap -A target.com
# Fixed scan, no adaptation
# Runs same scans for all targets
# No permission management
```

### Shadow AI Scanner
```bash
./shadow smart-scan target.com
# AI analyzes target first
# Creates custom plan
# Adapts to target type
# Manages permissions intelligently
```

## 🚧 Current Limitations

**Note**: Current implementation focuses on **planning**. Actual tool execution coming in future releases.

Current features:
- ✅ AI creates reconnaissance plans
- ✅ Permission management system
- ✅ Interactive approval workflow
- ✅ Fallback strategies

Coming soon:
- 🔜 Execute planned reconnaissance
- 🔜 Feed results back to AI
- 🔜 Adaptive multi-phase scanning
- 🔜 Learning from scan results

## 💡 Best Practices

1. **Review AI Plans**
   - Always review before executing
   - Understand what will be scanned
   - Check permission requirements

2. **Use Capabilities Over Sudo**
   - More secure
   - Granular permissions
   - No password prompts

3. **Start with Quick Profile**
   - Test on known targets first
   - Understand AI behavior
   - Then use deep profile

4. **Provide Authorization**
   - Only scan authorized targets
   - AI doesn't validate authorization
   - You are responsible for legal compliance

## 🔐 Security Considerations

### What Shadow Does

✅ Asks permission for each privileged command
✅ Shows exact command before execution
✅ Suggests secure alternatives (setcap)
✅ Falls back to non-privileged options
✅ Tracks and reports all permissions granted

### What Shadow Doesn't Do

❌ Never runs commands without approval
❌ No blanket root access
❌ No hidden privileged operations
❌ No automatic sudo password entry

## 🎯 Summary

Shadow's AI-Driven Reconnaissance represents a new approach to security scanning:

**Traditional Tools**: Fixed scans → Results
**Shadow**: AI thinks → Plans → Requests permissions → Executes → Adapts

**Key Features:**
- 🤖 AI analyzes targets and creates custom plans
- 🔐 Interactive permission management
- 🛡️ Security-first design
- 🔄 Graceful fallbacks
- 📊 Clear reasoning and explanations

**Inspired By:**
- OpenClaw's autonomous AI workflow
- Modern pentesting methodologies
- Security-first permission models

---

**Try it now:**
```bash
./shadow smart-scan yourtarget.com
```

Let AI plan your reconnaissance strategy! 🚀
