# Multi-Agent AI System

## Overview

Shadow uses a sophisticated multi-agent AI system inspired by [OpenClaw](https://github.com/openclaw/openclaw) to provide specialized, cost-effective security analysis. Each agent is optimized for specific tasks using the most appropriate Claude model.

## 🤖 Available Agents

### 1. Quick Scanner
- **Model**: Claude Haiku 4.5
- **Thinking**: Low
- **Cost**: $0.80-4.00 per million tokens
- **Use Case**: Fast triage and basic vulnerability identification
- **Best For**: Quick scans, initial assessment, rapid triage

### 2. Reconnaissance Analyst
- **Model**: Claude Sonnet 4.5
- **Thinking**: High
- **Cost**: $3-15 per million tokens
- **Use Case**: Deep reconnaissance and attack surface analysis
- **Tasks**:
  - Technology stack identification
  - Service enumeration
  - Attack surface mapping
  - Configuration analysis

### 3. Vulnerability Researcher
- **Model**: Claude Sonnet 4.5
- **Thinking**: High
- **Cost**: $3-15 per million tokens
- **Use Case**: Comprehensive vulnerability analysis
- **Tasks**:
  - OWASP Top 10 assessment
  - CVE research and correlation
  - Vulnerability prioritization
  - Risk scoring

### 4. Exploitation Specialist
- **Model**: Claude Opus 4.6
- **Thinking**: High
- **Cost**: $15-75 per million tokens
- **Use Case**: Advanced exploitation path analysis
- **Tasks**:
  - Attack chain development
  - Exploitation techniques
  - Proof-of-concept guidance
  - Real-world impact assessment

### 5. Security Reporter
- **Model**: Claude Sonnet 4.5
- **Thinking**: High
- **Cost**: $3-15 per million tokens
- **Use Case**: Executive and technical report generation
- **Tasks**:
  - Risk assessment
  - Executive summaries
  - Remediation roadmaps
  - Business impact analysis

## 📊 Scan Profiles

### Quick Profile
```bash
./shadow scan example.com --ai-analysis --profile quick
```

**Agent Used**: Quick Scanner (Haiku 4.5)
**Analysis Time**: ~30-60 seconds
**Estimated Cost**: $0.01-0.05 per scan
**Best For**:
- Rapid vulnerability triage
- CI/CD integration
- High-frequency scanning
- Cost-sensitive operations

### Standard Profile (Recommended)
```bash
./shadow scan example.com --ai-analysis --profile standard
```

**Agent Used**: Vulnerability Researcher (Sonnet 4.5)
**Analysis Time**: ~2-5 minutes
**Estimated Cost**: $0.05-0.20 per scan
**Best For**:
- Comprehensive vulnerability assessment
- Balanced speed and depth
- Most security assessments
- Production use

### Deep Profile
```bash
./shadow scan example.com --ai-analysis --profile deep
```

**Agents Used**:
1. Reconnaissance Analyst (Sonnet 4.5)
2. Vulnerability Researcher (Sonnet 4.5)
3. Exploitation Specialist (Opus 4.6)

**Analysis Time**: ~5-15 minutes
**Estimated Cost**: $0.50-2.00 per scan
**Best For**:
- Critical systems
- Pre-production security audits
- Penetration testing
- Complex applications

**Deep Analysis Stages**:
```
Stage 1: Reconnaissance Analysis
   └─ Maps attack surface and identifies entry points

Stage 2: Vulnerability Analysis
   └─ Deep vulnerability assessment with CVE research

Stage 3: Exploitation Analysis
   └─ Advanced attack chain and exploitation path analysis
```

## 💰 Cost Optimization

### Model Pricing (per million tokens)

| Model | Input | Output | When to Use |
|-------|-------|--------|-------------|
| **Haiku 4.5** | $0.80 | $4.00 | Quick scans, triage, simple analysis |
| **Sonnet 4.5** | $3.00 | $15.00 | Standard analysis, balanced approach |
| **Opus 4.6** | $15.00 | $75.00 | Complex exploitation, critical systems |

### Cost Estimation Examples

**Quick Scan**:
- Input: ~2K tokens
- Output: ~1K tokens
- Cost: ~$0.01-0.02

**Standard Scan**:
- Input: ~5K tokens
- Output: ~3K tokens
- Cost: ~$0.06-0.10

**Deep Scan**:
- Input: ~15K tokens (across 3 agents)
- Output: ~10K tokens
- Cost: ~$0.50-1.00

## 📈 Usage Tracking

Shadow automatically tracks and displays usage statistics:

```
📊 AI Model Usage Summary
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📈 Overall Statistics:
   Operations: 3/3 successful
   Total Tokens: 14.5K input, 9.2K output
   Estimated Cost: $0.65
   Total Duration: 8m 32s

🤖 By Agent:
   Reconnaissance Analyst (using Sonnet 4.5)
      Tokens: 5.0K in, 3.0K out
      Cost: $0.15 | Duration: 2m 15s | Success: 1/1

   Vulnerability Researcher (using Sonnet 4.5)
      Tokens: 5.5K in, 3.5K out
      Cost: $0.18 | Duration: 3m 10s | Success: 1/1

   Exploitation Specialist (using Opus 4.6)
      Tokens: 4.0K in, 2.7K out
      Cost: $0.32 | Duration: 3m 7s | Success: 1/1

🎯 By Model:
   Claude Sonnet 4.5 (balanced)
      Tokens: 10.5K in, 6.5K out
      Cost: $0.33 | Operations: 2

   Claude Opus 4.6 (most capable)
      Tokens: 4.0K in, 2.7K out
      Cost: $0.32 | Operations: 1
```

## 🎯 Agent Selection Strategy

Shadow automatically selects agents based on:

1. **Scan Profile**: User-specified profile (quick/standard/deep)
2. **Task Complexity**: Simple triage vs complex exploitation
3. **Cost Constraints**: Budget-aware model selection
4. **Performance Requirements**: Speed vs thoroughness trade-off

### Decision Flow

```
User selects profile
    ↓
Quick? → Use Haiku (fast, cheap)
    ↓
Standard? → Use Sonnet (balanced)
    ↓
Deep? → Use multi-agent (thorough)
    ├─ Recon: Sonnet
    ├─ Vuln: Sonnet
    └─ Exploit: Opus
```

## 🔧 Configuration

### Default Agent Configuration

Agents are pre-configured with optimal settings. Configuration in `pkg/models/agent.go`:

```go
{
    Name:        "Quick Scanner",
    Type:        AgentTypeQuickScan,
    Model:       "claude-haiku-4.5",
    Thinking:    "low",
    Description: "Fast initial scan analysis",
}
```

### Custom System Prompts

Each agent has a specialized system prompt optimized for its role. See `internal/ai/agent_manager.go:buildSystemPrompt()`.

## 📚 Advanced Usage

### List Available Agents

```bash
./shadow agents
```

Shows all agents, models, pricing, and use cases.

### View Real-Time Progress

During analysis, Shadow shows which agent is active:

```
🤖 Using Vulnerability Researcher (Sonnet 4.5)
📋 Task: Comprehensive vulnerability analysis
🧠 Sending to Claude AI (with extended thinking)...
⏱️  Vulnerability Researcher analyzing... (45s elapsed)
```

### Monitor Usage Patterns

Usage statistics are displayed after every scan with `--ai-analysis`:

- Total tokens consumed
- Cost per agent
- Duration per agent
- Success rates
- Model distribution

## 🚀 Performance Tips

### For Speed
```bash
# Use quick profile with Haiku
./shadow scan target.com --ai-analysis --profile quick
```

### For Cost Optimization
```bash
# Standard profile provides best value
./shadow scan target.com --ai-analysis --profile standard
```

### For Maximum Thoroughness
```bash
# Deep analysis with multi-agent pipeline
./shadow scan target.com --ai-analysis --profile deep
```

## 🔍 Under the Hood

### Agent Manager Architecture

```
AgentManager
├── Agent Pool (5 specialized agents)
│   ├── Quick Scanner (Haiku)
│   ├── Recon Analyst (Sonnet)
│   ├── Vuln Researcher (Sonnet)
│   ├── Exploit Specialist (Opus)
│   └── Reporter (Sonnet)
├── Usage Tracker
│   ├── Per-agent statistics
│   ├── Per-model statistics
│   └── Cost estimation
└── Analysis Pipeline
    ├── Profile selection
    ├── Agent orchestration
    └── Result aggregation
```

### Key Files

- `internal/ai/agent_manager.go` - Multi-agent orchestration
- `internal/ai/usage_tracker.go` - Token and cost tracking
- `pkg/models/agent.go` - Agent definitions and configuration
- `cmd/shadow/main.go` - CLI integration

## 🎓 Inspired By

This multi-agent system is inspired by [OpenClaw's](https://github.com/openclaw/openclaw) production-tested patterns:
- Model selection based on task complexity
- Usage tracking and cost optimization
- Thinking mode configuration
- Specialized system prompts

## 📊 Example Outputs

### Quick Analysis Output
```
🤖 Using Quick Scanner (Haiku 4.5)
📋 Task: Fast initial scan analysis

📊 AI Analysis Results:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📝 Summary:
Target shows 2 medium-severity findings requiring attention.

🎯 Risk Score: 35/100

💡 Top Recommendations:
  1. [High] Update security headers
  2. [Medium] Review SSL/TLS configuration

📊 AI Model Usage Summary
   Operations: 1/1 successful
   Estimated Cost: $0.02
   Total Duration: 45s
```

### Deep Analysis Output
```
🤖 Starting Multi-Agent AI Analysis...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📍 Stage 1/3: Reconnaissance Analysis
   🤖 Using Reconnaissance Analyst (Sonnet 4.5)
   ✅ Reconnaissance complete (2m 15s)

🔍 Stage 2/3: Vulnerability Analysis
   🤖 Using Vulnerability Researcher (Sonnet 4.5)
   ✅ Vulnerability analysis complete (3m 10s)

💥 Stage 3/3: Exploitation Analysis
   🤖 Using Exploitation Specialist (Opus 4.6)
   ✅ Exploitation assessment complete (3m 7s)

[Detailed analysis results...]

📊 AI Model Usage Summary
   Operations: 3/3 successful
   Estimated Cost: $0.65
   Total Duration: 8m 32s
```

## 🛠️ Troubleshooting

### High Costs
- Use `--profile quick` for routine scans
- Reserve deep analysis for critical systems
- Monitor usage with `shadow agents` output

### Slow Analysis
- Quick profile completes in <1 minute
- Standard profile: 2-5 minutes
- Deep profile: 5-15 minutes (3 agents)

### Agent Failures
- Check authentication: `shadow auth-check`
- Verify network connectivity
- Review usage limits

## 📝 Best Practices

1. **Use Quick for CI/CD**: Fast feedback in pipelines
2. **Use Standard for Regular Scans**: Best value/performance
3. **Use Deep for Critical Systems**: Thorough pre-production audit
4. **Monitor Costs**: Review usage summaries regularly
5. **Profile Selection**: Match profile to risk level

## 🔮 Future Enhancements

- [ ] Custom agent configurations
- [ ] Agent result caching
- [ ] Streaming analysis progress
- [ ] Custom model selection per agent
- [ ] Cost budgeting and alerts
- [ ] Historical usage analytics

---

**Cost-Aware. Production-Ready. OpenClaw-Inspired.** 🚀
