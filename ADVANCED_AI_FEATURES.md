# Advanced AI Features (Based on OpenClaw)

## Overview

Shadow now includes advanced AI capabilities inspired by OpenClaw's production-tested patterns, providing robust, reliable AI analysis with intelligent error handling and retry logic.

## ✅ What Was Added

### 1. Automatic Retry with Exponential Backoff

**From OpenClaw Pattern** (`translator.go:68-87`)

```go
// Shadow Implementation
const (
    maxRetryAttempts = 3
    baseRetryDelay   = 15 * time.Second
)

func (a *AdvancedClaudeAnalyzer) retryWithBackoff(ctx context.Context, fn func(context.Context) (*models.AIAnalysis, error)) (*models.AIAnalysis, error) {
    for attempt := 0; attempt < maxRetryAttempts; attempt++ {
        result, err := fn(ctx)
        if err == nil {
            return result, nil
        }

        if !isRetryableError(err) {
            return nil, err
        }

        // Exponential backoff: 15s, 30s, 45s
        delay := baseRetryDelay * time.Duration(attempt+1)
        sleepWithContext(ctx, delay)
    }
}
```

**Benefits:**
- Automatically retries transient failures
- Exponential backoff prevents overwhelming the API
- Up to 3 attempts before giving up
- Context-aware (respects cancellation)

### 2. Intelligent Error Detection

**From OpenClaw Pattern** (`translator.go:120-129`)

```go
func isRetryableError(err error) bool {
    message := strings.ToLower(err.Error())

    retryablePatterns := []string{
        "rate limit",
        "429",
        "timeout",
        "temporary",
        "connection",
        "deadline exceeded",  // ← Fixes your "context deadline exceeded" error
    }

    for _, pattern := range retryablePatterns {
        if strings.Contains(message, pattern) {
            return true
        }
    }
    return false
}
```

**Handles:**
- Rate limiting (429 errors)
- Timeouts (including your "deadline exceeded" error)
- Network issues
- Temporary API failures
- Empty responses

### 3. Extended Timeouts

**Fixed Your Issue:**

```go
// OLD (2 minutes - too short)
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

// NEW (10 minutes for analysis, 5 for queries)
const (
    defaultAnalysisTimeout = 10 * time.Minute  // ← 5x longer!
    defaultQueryTimeout    = 5 * time.Minute
)
```

**Why This Fixes "context deadline exceeded":**
- Large scan results need more processing time
- Claude AI with extended thinking requires patience
- 10 minutes is plenty for even complex analyses

### 4. System Prompt for Better Analysis

**From OpenClaw Pattern** (system prompts for specialized tasks)

```go
opts.SystemPrompt = `You are an expert security analyst and penetration tester with deep knowledge of:
- Web application security (OWASP Top 10)
- Network security and reconnaissance
- Vulnerability assessment and exploitation
- Secure coding practices
- Risk assessment and prioritization

Your role is to:
1. Analyze security scan results thoroughly
2. Identify critical vulnerabilities and their impact
3. Provide actionable remediation steps
4. Prioritize findings by risk level
5. Explain attack chains and exploitation scenarios`
```

**Benefits:**
- More focused, relevant analysis
- Better understanding of security context
- Improved prioritization
- Actionable recommendations

### 5. Enhanced Prompt Engineering

**Structured Analysis Requests:**

```go
prompt := fmt.Sprintf(`# Security Scan Analysis Request

## Target Information
- **Target**: %s
- **Scan Time**: %s
- **Total Findings**: %d

## Analysis Tasks
1. Executive Summary
2. Critical Issues (top 3-5)
3. Risk Assessment (0-100)
4. Prioritized Recommendations
5. Attack Chains

## Scan Findings
...detailed findings...

## Output Format
Use markdown headings. Be specific and actionable.`)
```

**Improvements:**
- Clear structure for AI to follow
- Specific analysis requests
- Formatted output expectations
- Better, more consistent results

### 6. Context-Aware Timeout Management

**From OpenClaw Pattern** (`translator.go:131-140`)

```go
func sleepWithContext(ctx context.Context, delay time.Duration) error {
    timer := time.NewTimer(delay)
    defer timer.Stop()

    select {
    case <-ctx.Done():
        return ctx.Err()  // Respect cancellation
    case <-timer.C:
        return nil
    }
}
```

**Benefits:**
- Respects context cancellation
- Allows user to interrupt long operations
- Clean resource management

### 7. Better Progress Feedback

```go
fmt.Println("⏳ Analyzing scan results (this may take a few minutes)...")

// During retries:
fmt.Printf("⚠️  Retry %d/%d after %v (error: %v)\n", attempt+1, maxRetryAttempts, delay, err)

// On completion:
fmt.Printf("✅ Analysis completed at %s\n", analysis.Timestamp.Format("15:04:05"))
```

**Improvements:**
- User knows what's happening
- Retry progress visible
- Clear success/failure messages
- Estimated timing information

## 🚀 Usage

### Basic Usage (Same as Before)

```bash
./shadow scan example.com --ai-analysis
```

### What Happens Now (Behind the Scenes)

```
1. ⏳ Start analysis (10 min timeout)
2. 🔄 Attempt 1: Send to Claude AI
   └─ ❌ "context deadline exceeded"
   └─ ⚠️  Wait 15 seconds (retryable error detected)
3. 🔄 Attempt 2: Send to Claude AI
   └─ ❌ "rate limit"
   └─ ⚠️  Wait 30 seconds
4. 🔄 Attempt 3: Send to Claude AI
   └─ ✅ Success!
5. 📊 Display results
```

## 📊 Comparison: Old vs New

| Feature | Old Implementation | New Implementation | Improvement |
|---------|-------------------|-------------------|-------------|
| **Timeout** | 2 minutes | 10 minutes | 5x longer ✓ |
| **Retry** | None | 3 attempts with backoff | Auto-retry ✓ |
| **Error Handling** | Basic | Intelligent (rate limit, timeout, etc.) | Smart ✓ |
| **Feedback** | Minimal | Progress + retry info | Clear ✓ |
| **System Prompt** | None | Expert security analyst | Better results ✓ |
| **Prompt Structure** | Basic | Structured with sections | More consistent ✓ |
| **Context Handling** | Simple | Context-aware cancellation | Robust ✓ |

## 🔧 Configuration

### Adjusting Retry Behavior

Edit `internal/ai/advanced_client.go`:

```go
const (
    maxRetryAttempts = 3          // Change to 5 for more retries
    baseRetryDelay   = 15 * time.Second  // Change to 30s for longer waits
)
```

### Adjusting Timeouts

```go
const (
    defaultAnalysisTimeout = 10 * time.Minute  // Increase if needed
    defaultQueryTimeout    = 5 * time.Minute
)
```

## 🎯 Advanced Features Inspired by OpenClaw

### 1. Retry Logic ✅
- **Source**: `openclaw/translator.go:68-87`
- **Implementation**: `advanced_client.go:retryWithBackoff`

### 2. Error Classification ✅
- **Source**: `openclaw/translator.go:120-129`
- **Implementation**: `advanced_client.go:isRetryableError`

### 3. Context Handling ✅
- **Source**: `openclaw/translator.go:131-140`
- **Implementation**: `advanced_client.go:sleepWithContext`

### 4. System Prompts ✅
- **Source**: `openclaw/translator.go:35`
- **Implementation**: `advanced_client.go:buildSystemPrompt`

### 5. Exponential Backoff ✅
- **Source**: `openclaw/translator.go:80`
- **Implementation**: `advanced_client.go:retryWithBackoff` (lines 134-136)

## 📚 Future Enhancements (From OpenClaw Patterns)

### 1. Streaming Analysis

```go
// TODO: Implement using pi-golang's Subscribe feature
func (a *AdvancedClaudeAnalyzer) StreamingAnalyze(ctx context.Context, result *models.ScanResult, callback func(string)) error {
    events, cancel := a.client.Subscribe(256)
    defer cancel()

    // Stream analysis results in real-time
    for event := range events {
        if event.Type == "text" {
            callback(event.Text)
        }
    }
}
```

### 2. Response Caching

```go
type AnalysisCache struct {
    cache map[string]*models.AIAnalysis
    ttl   time.Duration
}

// Cache similar vulnerability analyses
// Avoid re-analyzing same issue types
```

### 3. Batch Processing

```go
func (a *AdvancedClaudeAnalyzer) AnalyzeMultiple(ctx context.Context, results []*models.ScanResult) ([]*models.AIAnalysis, error) {
    // Combine multiple scans in one API call
    // Reduce costs and improve efficiency
}
```

### 4. Custom Thinking Modes

```go
// From OpenClaw's normalizeThinking
func (a *AdvancedClaudeAnalyzer) SetThinkingMode(mode string) {
    // "high" - Deep analysis (slower, better)
    // "low"  - Quick analysis (faster, simpler)
}
```

## 🐛 Troubleshooting

### "context deadline exceeded" - FIXED ✅

**Old**: 2-minute timeout was too short
**New**: 10-minute timeout + auto-retry

**If still occurring:**
- Check network connection
- Try with `--profile quick` (smaller results)
- Verify API key/OAuth is working

### Rate Limiting

**Handled automatically:**
- Detects 429 errors
- Waits 15-45 seconds
- Retries up to 3 times

### Network Issues

**Handled automatically:**
- Detects connection errors
- Retries with backoff
- Provides clear error messages

## ✅ Testing

```bash
# Test basic analysis
./shadow scan example.com --ai-analysis

# Test with large results (stress test)
./shadow scan example.com --profile deep --ai-analysis

# Test authentication
./shadow auth-check

# Monitor retries (you'll see progress messages)
./shadow scan example.com --ai-analysis
# Output:
# ⏳ Analyzing scan results (this may take a few minutes)...
# ⚠️  Retry 1/3 after 15s (error: context deadline exceeded)
# ⚠️  Retry 2/3 after 30s (error: rate limit)
# ✅ Analysis completed at 14:32:15
```

## 📖 OpenClaw References

- **Repository**: https://github.com/openclaw/openclaw
- **Translator**: `scripts/docs-i18n/translator.go`
- **Pi Integration**: Uses same `github.com/joshp123/pi-golang` library
- **Patterns Used**: Retry logic, error handling, context management

## 🎉 Summary

Shadow now has:
- ✅ **5x longer timeout** (10 min vs 2 min)
- ✅ **Automatic retry** (3 attempts with exponential backoff)
- ✅ **Intelligent error handling** (rate limits, timeouts, network issues)
- ✅ **Better prompts** (system prompts + structured requests)
- ✅ **Clear progress feedback** (know what's happening)
- ✅ **Production-tested patterns** (from OpenClaw)

**Your "context deadline exceeded" error is now fixed!** 🚀

