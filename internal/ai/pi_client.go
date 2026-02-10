package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	pi "github.com/joshp123/pi-golang"
	"github.com/kumaraguru1735/shadow/pkg/models"
)

// PiClaudeAnalyzer provides AI-powered security analysis using pi-golang
// This automatically uses Claude Code OAuth token from ~/.claude/oauth.json
type PiClaudeAnalyzer struct {
	client *pi.OneShotClient
	model  string
}

// NewPiClaudeAnalyzer creates a new analyzer using pi-golang
// This automatically picks up Claude Code OAuth token
func NewPiClaudeAnalyzer() (*PiClaudeAnalyzer, error) {
	opts := pi.DefaultOneShotOptions()
	opts.AppName = "shadow"
	opts.Mode = pi.ModeDragons
	opts.Dragons = pi.DragonsOptions{
		Provider: "anthropic",
		Model:    "claude-sonnet-4.5-20250929",
		Thinking: "high",
	}

	client, err := pi.StartOneShot(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to start pi client: %w (ensure you have pi CLI installed or Claude Code OAuth configured)", err)
	}

	return &PiClaudeAnalyzer{
		client: client,
		model:  "claude-sonnet-4.5-20250929",
	}, nil
}

// AnalyzeScan performs AI analysis on scan results
func (a *PiClaudeAnalyzer) AnalyzeScan(ctx context.Context, result *models.ScanResult) (*models.AIAnalysis, error) {
	prompt := a.buildAnalysisPrompt(result)

	// Use the Run method which handles event parsing internally
	runResult, err := a.client.Run(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to run analysis: %w", err)
	}

	text := runResult.Text

	analysis := &models.AIAnalysis{
		ScanID:          result.ID,
		Summary:         parseAnalysisSummary(text),
		RiskScore:       parseRiskScore(text),
		CriticalIssues:  parseCriticalIssues(text),
		Recommendations: parseRecommendations(text),
	}

	return analysis, nil
}

// buildAnalysisPrompt constructs the analysis prompt for Claude
func (a *PiClaudeAnalyzer) buildAnalysisPrompt(result *models.ScanResult) string {
	prompt := fmt.Sprintf(`You are a senior security analyst reviewing scan results for %s.

Scan completed at: %s
Total findings: %d

Please analyze these security findings and provide:
1. Executive summary (2-3 sentences)
2. Critical issues that need immediate attention
3. Prioritized recommendations with implementation steps
4. Potential attack chains combining multiple vulnerabilities
5. Risk score (0-100)

Findings:
`, result.Target, result.EndTime, len(result.Findings))

	for _, finding := range result.Findings {
		prompt += fmt.Sprintf("\n- [%s] %s: %s", finding.Severity, finding.Title, finding.Description)
	}

	return prompt
}

// QueryResults allows natural language queries about scan results
func (a *PiClaudeAnalyzer) QueryResults(ctx context.Context, scanID string, question string) (string, error) {
	prompt := fmt.Sprintf("Scan ID: %s\nQuestion: %s", scanID, question)

	runResult, err := a.client.Run(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to run query: %w", err)
	}

	return runResult.Text, nil
}

// Close closes the pi client
func (a *PiClaudeAnalyzer) Close() {
	if a.client != nil {
		_ = a.client.Close()
	}
}

// Helper functions to parse Claude's response
func parseAnalysisSummary(text string) string {
	// Simple extraction - look for summary section
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), "summary") && i+1 < len(lines) {
			return strings.TrimSpace(lines[i+1])
		}
	}
	return strings.Split(text, "\n")[0] // fallback to first line
}

func parseRiskScore(text string) int {
	// Look for risk score in text
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "risk score") {
			// Extract number (basic parsing)
			for i := 0; i < len(line)-1; i++ {
				if line[i] >= '0' && line[i] <= '9' {
					score := 0
					for j := i; j < len(line) && line[j] >= '0' && line[j] <= '9'; j++ {
						score = score*10 + int(line[j]-'0')
					}
					if score >= 0 && score <= 100 {
						return score
					}
				}
			}
		}
	}
	return 50 // default medium risk
}

func parseCriticalIssues(text string) []string {
	issues := []string{}
	lines := strings.Split(text, "\n")
	inCritical := false

	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "critical") {
			inCritical = true
			continue
		}
		if inCritical && strings.TrimSpace(line) != "" {
			if strings.HasPrefix(strings.TrimSpace(line), "-") || strings.HasPrefix(strings.TrimSpace(line), "*") {
				issues = append(issues, strings.TrimSpace(line))
			}
		}
		if inCritical && strings.Contains(lower, "recommendation") {
			break
		}
	}

	return issues
}

func parseRecommendations(text string) []models.Recommendation {
	recommendations := []models.Recommendation{}
	lines := strings.Split(text, "\n")
	inRecommendations := false

	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "recommendation") {
			inRecommendations = true
			continue
		}
		if inRecommendations && strings.TrimSpace(line) != "" {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "-") ||
			   strings.HasPrefix(trimmed, "*") ||
			   strings.HasPrefix(trimmed, "1") {
				// Simple recommendation structure
				rec := models.Recommendation{
					Priority:    "medium",
					Title:       trimmed,
					Description: trimmed,
					Impact:      "unknown",
					Effort:      "medium",
				}
				recommendations = append(recommendations, rec)
			}
		}
	}

	return recommendations
}

// GetAuthenticationStatus checks what authentication method is available
func GetAuthenticationStatus() string {
	// Check for OAuth token (Claude Code)
	home, err := os.UserHomeDir()
	if err == nil {
		oauthPaths := []string{
			home + "/.claude/.credentials.json", // Claude Code credentials
			home + "/.claude/oauth.json",
			home + "/.config/claude/oauth.json",
			home + "/.config/anthropic/oauth.json",
			home + "/.pi/agent/oauth.json",
		}
		for _, path := range oauthPaths {
			if _, err := os.Stat(path); err == nil {
				return fmt.Sprintf("✓ Claude Code OAuth token found at %s", path)
			}
		}
	}

	// Check for API key
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		return "✓ ANTHROPIC_API_KEY environment variable set"
	}

	return "✗ No authentication found - set ANTHROPIC_API_KEY or use Claude Code OAuth"
}
