package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	pi "github.com/joshp123/pi-golang"
	"github.com/kumaraguru1735/shadow/pkg/models"
)

const (
	// Retry configuration (from openclaw)
	maxRetryAttempts = 3
	baseRetryDelay   = 15 * time.Second

	// Timeout configuration (increased from 2min to handle long scans)
	defaultAnalysisTimeout = 10 * time.Minute
	defaultQueryTimeout    = 5 * time.Minute
)

var (
	errEmptyResponse      = errors.New("empty AI response")
	errRateLimitExceeded  = errors.New("rate limit exceeded")
)

// AdvancedClaudeAnalyzer provides advanced AI analysis with retry logic and better error handling
type AdvancedClaudeAnalyzer struct {
	client *pi.OneShotClient
	model  string
}

// NewAdvancedClaudeAnalyzer creates an advanced analyzer with openclaw-style features
func NewAdvancedClaudeAnalyzer() (*AdvancedClaudeAnalyzer, error) {
	opts := pi.DefaultOneShotOptions()
	opts.AppName = "shadow"
	opts.Mode = pi.ModeDragons
	opts.Dragons = pi.DragonsOptions{
		Provider: "anthropic",
		Model:    "claude-sonnet-4.5-20250929",
		Thinking: "high", // High thinking mode for better analysis
	}

	// Set system prompt for security analysis
	opts.SystemPrompt = buildSystemPrompt()

	client, err := pi.StartOneShot(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to start pi client: %w", err)
	}

	return &AdvancedClaudeAnalyzer{
		client: client,
		model:  "claude-sonnet-4.5-20250929",
	}, nil
}

// buildSystemPrompt creates a comprehensive system prompt for security analysis
func buildSystemPrompt() string {
	return `You are an expert security analyst and penetration tester with deep knowledge of:
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
5. Explain attack chains and exploitation scenarios

Always provide:
- Clear, concise analysis
- Specific remediation steps
- Risk scores with justification
- Practical security recommendations

Be direct and technical. Focus on actionable insights.`
}

// AnalyzeScanWithRetry performs AI analysis with automatic retry logic (openclaw pattern)
func (a *AdvancedClaudeAnalyzer) AnalyzeScanWithRetry(ctx context.Context, result *models.ScanResult) (*models.AIAnalysis, error) {
	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, defaultAnalysisTimeout)
	defer cancel()

	// Retry with exponential backoff
	return a.retryWithBackoff(ctx, func(ctx context.Context) (*models.AIAnalysis, error) {
		return a.analyzeScanOnce(ctx, result)
	})
}

// analyzeScanOnce performs a single analysis attempt
func (a *AdvancedClaudeAnalyzer) analyzeScanOnce(ctx context.Context, result *models.ScanResult) (*models.AIAnalysis, error) {
	prompt := a.buildAnalysisPrompt(result)

	runResult, err := a.client.Run(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to run analysis: %w", err)
	}

	text := runResult.Text
	if strings.TrimSpace(text) == "" {
		return nil, errEmptyResponse
	}

	analysis := &models.AIAnalysis{
		ScanID:          result.ID,
		Summary:         parseAnalysisSummary(text),
		RiskScore:       parseRiskScore(text),
		CriticalIssues:  parseCriticalIssues(text),
		Recommendations: parseRecommendations(text),
		Timestamp:       time.Now(),
	}

	return analysis, nil
}

// retryWithBackoff implements openclaw's retry pattern
func (a *AdvancedClaudeAnalyzer) retryWithBackoff(ctx context.Context, fn func(context.Context) (*models.AIAnalysis, error)) (*models.AIAnalysis, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetryAttempts; attempt++ {
		// Check context before attempting
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}

		// Check if error is retryable
		if !isRetryableError(err) {
			return nil, err
		}

		lastErr = err

		// Calculate backoff delay (exponential)
		if attempt+1 < maxRetryAttempts {
			delay := baseRetryDelay * time.Duration(attempt+1)
			fmt.Printf("⚠️  Retry %d/%d after %v (error: %v)\n", attempt+1, maxRetryAttempts, delay, err)

			if err := sleepWithContext(ctx, delay); err != nil {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// isRetryableError determines if an error should trigger a retry (openclaw pattern)
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific error types
	if errors.Is(err, errEmptyResponse) {
		return true
	}

	if errors.Is(err, errRateLimitExceeded) {
		return true
	}

	// Check error message for retryable conditions
	message := strings.ToLower(err.Error())
	retryablePatterns := []string{
		"rate limit",
		"429",
		"timeout",
		"temporary",
		"connection",
		"deadline exceeded",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(message, pattern) {
			return true
		}
	}

	return false
}

// sleepWithContext sleeps with context cancellation support (openclaw pattern)
func sleepWithContext(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// QueryWithRetry performs AI queries with retry logic
func (a *AdvancedClaudeAnalyzer) QueryWithRetry(ctx context.Context, scanID string, question string) (string, error) {
	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	return a.retryStringWithBackoff(ctx, func(ctx context.Context) (string, error) {
		prompt := fmt.Sprintf("Scan ID: %s\nQuestion: %s", scanID, question)

		runResult, err := a.client.Run(ctx, prompt)
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(runResult.Text) == "" {
			return "", errEmptyResponse
		}

		return runResult.Text, nil
	})
}

// retryStringWithBackoff implements retry for string-returning functions
func (a *AdvancedClaudeAnalyzer) retryStringWithBackoff(ctx context.Context, fn func(context.Context) (string, error)) (string, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetryAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}

		if !isRetryableError(err) {
			return "", err
		}

		lastErr = err

		if attempt+1 < maxRetryAttempts {
			delay := baseRetryDelay * time.Duration(attempt+1)
			fmt.Printf("⚠️  Retry %d/%d after %v\n", attempt+1, maxRetryAttempts, delay)

			if err := sleepWithContext(ctx, delay); err != nil {
				return "", err
			}
		}
	}

	return "", fmt.Errorf("max retries exceeded: %w", lastErr)
}

// buildAnalysisPrompt constructs an enhanced analysis prompt
func (a *AdvancedClaudeAnalyzer) buildAnalysisPrompt(result *models.ScanResult) string {
	prompt := fmt.Sprintf(`# Security Scan Analysis Request

## Target Information
- **Target**: %s
- **Scan Time**: %s
- **Total Findings**: %d
- **Scan ID**: %s

## Analysis Tasks

Please analyze these security findings and provide:

1. **Executive Summary** (2-3 sentences)
   - Overall security posture
   - Most critical concerns
   - Business impact

2. **Critical Issues** (list top 3-5)
   - Issue description
   - Severity justification
   - Immediate risk

3. **Risk Assessment**
   - Overall risk score (0-100)
   - Scoring rationale
   - Risk factors

4. **Prioritized Recommendations**
   - Quick wins (low effort, high impact)
   - Critical fixes (immediate attention)
   - Long-term improvements

5. **Attack Chains** (if applicable)
   - How vulnerabilities could be chained
   - Potential exploitation scenarios

## Scan Findings

`, result.Target, result.EndTime.Format(time.RFC3339), len(result.Findings), result.ID)

	for i, finding := range result.Findings {
		prompt += fmt.Sprintf("\n### Finding %d: [%s] %s\n", i+1, finding.Severity, finding.Title)
		prompt += fmt.Sprintf("- **Type**: %s\n", finding.Type)
		prompt += fmt.Sprintf("- **Description**: %s\n", finding.Description)
		if finding.Evidence != "" {
			prompt += fmt.Sprintf("- **Evidence**: %s\n", finding.Evidence)
		}
		if finding.Location != "" {
			prompt += fmt.Sprintf("- **Location**: %s\n", finding.Location)
		}
	}

	prompt += `

## Output Format

Please structure your response clearly with markdown headings for each section.
Be specific, technical, and actionable.`

	return prompt
}

// Close closes the AI client
func (a *AdvancedClaudeAnalyzer) Close() {
	if a.client != nil {
		_ = a.client.Close()
	}
}

// StreamingAnalyze provides streaming analysis (future enhancement)
func (a *AdvancedClaudeAnalyzer) StreamingAnalyze(ctx context.Context, result *models.ScanResult, callback func(string)) error {
	// TODO: Implement streaming analysis using pi-golang's Subscribe feature
	// This would provide real-time feedback during analysis
	return fmt.Errorf("streaming analysis not yet implemented")
}
