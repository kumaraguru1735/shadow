package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	pi "github.com/joshp123/pi-golang"
	"github.com/kumaraguru1735/shadow/pkg/models"
)

// AgentManager orchestrates multiple specialized AI agents
type AgentManager struct {
	agents  map[models.AgentType]*Agent
	tracker *UsageTracker
}

// Agent represents a specialized AI agent
type Agent struct {
	config *models.AgentConfig
	client *pi.OneShotClient
}

// NewAgentManager creates a new multi-agent manager
func NewAgentManager() (*AgentManager, error) {
	manager := &AgentManager{
		agents:  make(map[models.AgentType]*Agent),
		tracker: NewUsageTracker(),
	}

	// Initialize all default agents
	configs := models.GetDefaultAgents()
	for i := range configs {
		agent, err := manager.createAgent(&configs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to create agent %s: %w", configs[i].Name, err)
		}
		manager.agents[configs[i].Type] = agent
	}

	return manager, nil
}

// createAgent creates a new agent with the given configuration
func (m *AgentManager) createAgent(config *models.AgentConfig) (*Agent, error) {
	opts := pi.DefaultOneShotOptions()
	opts.AppName = "shadow"
	opts.Mode = pi.ModeDragons
	opts.Dragons = pi.DragonsOptions{
		Provider: "anthropic",
		Model:    config.Model,
		Thinking: normalizeThinking(config.Thinking),
	}

	// Set agent-specific system prompt
	opts.SystemPrompt = m.buildSystemPrompt(config)

	client, err := pi.StartOneShot(opts)
	if err != nil {
		return nil, err
	}

	return &Agent{
		config: config,
		client: client,
	}, nil
}

// buildSystemPrompt creates a system prompt for the agent
func (m *AgentManager) buildSystemPrompt(config *models.AgentConfig) string {
	basePrompt := `You are an expert security analyst and penetration tester.`

	var rolePrompt string
	switch config.Type {
	case models.AgentTypeQuickScan:
		rolePrompt = `
Your role: QUICK SCANNER
- Perform rapid triage of scan results
- Identify obvious vulnerabilities quickly
- Flag critical issues for deeper analysis
- Keep analysis concise and actionable`

	case models.AgentTypeRecon:
		rolePrompt = `
Your role: RECONNAISSANCE ANALYST
- Analyze target attack surface comprehensively
- Identify exposed services and technologies
- Map potential entry points
- Assess configuration weaknesses
- Provide detailed reconnaissance insights`

	case models.AgentTypeVulnerability:
		rolePrompt = `
Your role: VULNERABILITY RESEARCHER
- Deep analysis of security vulnerabilities
- OWASP Top 10 assessment
- CVE research and correlation
- Risk scoring and prioritization
- Detailed exploitation prerequisites`

	case models.AgentTypeExploitation:
		rolePrompt = `
Your role: EXPLOITATION SPECIALIST
- Develop complete attack chains
- Identify exploitation paths
- Analyze exploitation feasibility
- Provide proof-of-concept guidance
- Assess real-world impact`

	case models.AgentTypeReport:
		rolePrompt = `
Your role: SECURITY REPORTER
- Generate executive summaries
- Create technical reports
- Develop remediation roadmaps
- Prioritize actions by business impact
- Communicate clearly to both technical and non-technical audiences`
	}

	return basePrompt + rolePrompt
}

// AnalyzeWithAgent performs analysis using a specific agent
func (m *AgentManager) AnalyzeWithAgent(
	ctx context.Context,
	agentType models.AgentType,
	prompt string,
	progress ProgressCallback,
) (string, error) {
	agent, ok := m.agents[agentType]
	if !ok {
		return "", fmt.Errorf("agent type %s not found", agentType)
	}

	if progress != nil {
		progress(fmt.Sprintf("🤖 Using %s (%s)",
			agent.config.Name,
			getModelShortName(agent.config.Model)))
		progress(fmt.Sprintf("📋 Task: %s", agent.config.Description))
	}

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultAnalysisTimeout)
	defer cancel()

	startTime := time.Now()

	// Show detailed progress updates
	done := make(chan bool)
	if progress != nil {
		go func() {
			ticker := time.NewTicker(15 * time.Second)
			defer ticker.Stop()

			stages := []string{
				"🔍 Analyzing security findings",
				"📊 Evaluating risk levels",
				"🎯 Identifying attack vectors",
				"🔗 Mapping attack chains",
				"📝 Generating recommendations",
				"✅ Finalizing analysis",
			}
			stageIdx := 0

			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					elapsed := time.Since(startTime)
					if stageIdx < len(stages) {
						progress(fmt.Sprintf("   %s (%.0fs)", stages[stageIdx], elapsed.Seconds()))
						stageIdx++
					} else {
						progress(fmt.Sprintf("   ⏱️  %s completing analysis... (%.0fs elapsed)",
							agent.config.Name, elapsed.Seconds()))
					}
				}
			}
		}()
	}

	// Run analysis
	result, err := agent.client.Run(timeoutCtx, prompt)
	close(done)

	duration := time.Since(startTime)

	// Record usage stats (note: pi-golang doesn't expose token counts, so we estimate)
	stats := UsageStats{
		Model:     agent.config.Model,
		Agent:     agent.config.Name,
		Duration:  duration,
		StartTime: startTime,
		EndTime:   time.Now(),
		Success:   err == nil,
	}

	if err != nil {
		stats.Error = err.Error()
	} else {
		// Estimate tokens (rough approximation: 1 token ≈ 4 characters)
		stats.InputTokens = int64(len(prompt) / 4)
		stats.OutputTokens = int64(len(result.Text) / 4)
	}

	m.tracker.RecordUsage(stats)

	if err != nil {
		return "", fmt.Errorf("analysis failed: %w", err)
	}

	return result.Text, nil
}

// AnalyzeScanWithAgents performs multi-agent analysis of scan results
func (m *AgentManager) AnalyzeScanWithAgents(
	ctx context.Context,
	result *models.ScanResult,
	profile string,
	progress ProgressCallback,
) (*models.AIAnalysis, error) {
	if progress != nil {
		progress("🚀 Starting multi-agent analysis...")
	}

	var analysis *models.AIAnalysis
	var err error

	switch profile {
	case "quick":
		// Quick scan: Use only Haiku for fast analysis
		analysis, err = m.runQuickAnalysis(ctx, result, progress)

	case "standard":
		// Standard: Use Sonnet for balanced analysis
		analysis, err = m.runStandardAnalysis(ctx, result, progress)

	case "deep":
		// Deep: Use multiple agents (Sonnet + Opus)
		analysis, err = m.runDeepAnalysis(ctx, result, progress)

	default:
		analysis, err = m.runStandardAnalysis(ctx, result, progress)
	}

	return analysis, err
}

// runQuickAnalysis uses Haiku for fast analysis
func (m *AgentManager) runQuickAnalysis(
	ctx context.Context,
	result *models.ScanResult,
	progress ProgressCallback,
) (*models.AIAnalysis, error) {
	prompt := buildAnalysisPrompt(result)

	text, err := m.AnalyzeWithAgent(ctx, models.AgentTypeQuickScan, prompt, progress)
	if err != nil {
		return nil, err
	}

	return parseAnalysisResponse(text, result.ID), nil
}

// runStandardAnalysis uses Sonnet for balanced analysis
func (m *AgentManager) runStandardAnalysis(
	ctx context.Context,
	result *models.ScanResult,
	progress ProgressCallback,
) (*models.AIAnalysis, error) {
	if progress != nil {
		progress(fmt.Sprintf("📋 Preparing analysis of %d findings from %s",
			len(result.Findings), result.Target))

		// Show what we're analyzing
		if len(result.Findings) > 0 {
			for i, finding := range result.Findings {
				if i < 3 { // Show first 3 findings
					progress(fmt.Sprintf("   • [%s] %s", finding.Severity, finding.Title))
				}
			}
			if len(result.Findings) > 3 {
				progress(fmt.Sprintf("   • ... and %d more", len(result.Findings)-3))
			}
		}
		progress("")
	}

	// Use vulnerability agent for standard analysis
	prompt := buildAnalysisPrompt(result)

	text, err := m.AnalyzeWithAgent(ctx, models.AgentTypeVulnerability, prompt, progress)
	if err != nil {
		return nil, err
	}

	return parseAnalysisResponse(text, result.ID), nil
}

// runDeepAnalysis uses multiple agents for comprehensive analysis
func (m *AgentManager) runDeepAnalysis(
	ctx context.Context,
	result *models.ScanResult,
	progress ProgressCallback,
) (*models.AIAnalysis, error) {
	if progress != nil {
		progress("🔬 Deep analysis mode: Using multiple specialized agents")
	}

	// Stage 1: Reconnaissance
	if progress != nil {
		progress("\n📍 Stage 1/3: Reconnaissance Analysis")
	}

	reconPrompt := buildReconPrompt(result)
	reconResult, err := m.AnalyzeWithAgent(ctx, models.AgentTypeRecon, reconPrompt, progress)
	if err != nil {
		return nil, fmt.Errorf("recon stage failed: %w", err)
	}

	// Stage 2: Vulnerability Analysis
	if progress != nil {
		progress("\n🔍 Stage 2/3: Vulnerability Analysis")
	}

	vulnPrompt := buildVulnPrompt(result, reconResult)
	vulnResult, err := m.AnalyzeWithAgent(ctx, models.AgentTypeVulnerability, vulnPrompt, progress)
	if err != nil {
		return nil, fmt.Errorf("vulnerability stage failed: %w", err)
	}

	// Stage 3: Exploitation Analysis (if critical vulns found)
	if progress != nil {
		progress("\n💥 Stage 3/3: Exploitation Analysis")
	}

	exploitPrompt := buildExploitPrompt(result, reconResult, vulnResult)
	exploitResult, err := m.AnalyzeWithAgent(ctx, models.AgentTypeExploitation, exploitPrompt, progress)
	if err != nil {
		// Don't fail the whole analysis if exploitation stage fails
		if progress != nil {
			progress(fmt.Sprintf("⚠️  Exploitation analysis unavailable: %v", err))
		}
		exploitResult = "Exploitation analysis not available."
	}

	// Combine results
	combinedText := fmt.Sprintf(`# Reconnaissance Findings
%s

# Vulnerability Analysis
%s

# Exploitation Assessment
%s`, reconResult, vulnResult, exploitResult)

	return parseAnalysisResponse(combinedText, result.ID), nil
}

// GetUsageSummary returns usage statistics
func (m *AgentManager) GetUsageSummary() UsageSummary {
	return m.tracker.GetSummary()
}

// Close closes all agents
func (m *AgentManager) Close() {
	for _, agent := range m.agents {
		if agent.client != nil {
			agent.client.Close()
		}
	}
}

// Helper functions

func normalizeThinking(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "low", "high":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "high"
	}
}

func buildAnalysisPrompt(result *models.ScanResult) string {
	return fmt.Sprintf(`# Security Scan Analysis Request

## Target Information
- **Target**: %s
- **Scan Time**: %s
- **Total Findings**: %d

## Analysis Tasks
1. Executive Summary (2-3 sentences)
2. Critical Issues (top 3-5 vulnerabilities)
3. Risk Assessment (score 0-100)
4. Prioritized Recommendations
5. Potential Attack Chains

## Scan Findings
%s

## Output Format
Use markdown headings. Be specific and actionable.`,
		result.Target,
		result.StartTime.Format(time.RFC3339),
		len(result.Findings),
		formatFindings(result.Findings))
}

func buildReconPrompt(result *models.ScanResult) string {
	return fmt.Sprintf(`# Reconnaissance Analysis

Analyze the target's attack surface:
- Exposed services and ports
- Technology stack and versions
- Configuration weaknesses
- Potential entry points

Target: %s
Findings: %s

Provide detailed reconnaissance insights.`,
		result.Target,
		formatFindings(result.Findings))
}

func buildVulnPrompt(result *models.ScanResult, reconData string) string {
	return fmt.Sprintf(`# Vulnerability Analysis

Based on reconnaissance findings, perform deep vulnerability analysis:
- OWASP Top 10 vulnerabilities
- Known CVEs
- Configuration issues
- Security misconfigurations

Target: %s
Reconnaissance Data:
%s

Scan Findings:
%s

Provide comprehensive vulnerability assessment with risk scores.`,
		result.Target,
		reconData,
		formatFindings(result.Findings))
}

func buildExploitPrompt(result *models.ScanResult, reconData, vulnData string) string {
	return fmt.Sprintf(`# Exploitation Analysis

Analyze exploitation possibilities:
- Complete attack chains
- Exploitation prerequisites
- Proof-of-concept approaches
- Real-world impact

Target: %s

Reconnaissance:
%s

Vulnerabilities:
%s

Provide detailed exploitation assessment.`,
		result.Target,
		reconData,
		vulnData)
}

func formatFindings(findings []models.Finding) string {
	if len(findings) == 0 {
		return "No findings detected"
	}

	var result strings.Builder
	for i, finding := range findings {
		result.WriteString(fmt.Sprintf("\n%d. [%s] %s",
			i+1,
			finding.Severity,
			finding.Title))
		if finding.Description != "" {
			result.WriteString(fmt.Sprintf("\n   Details: %s", finding.Description))
		}
	}

	return result.String()
}

func parseAnalysisResponse(text string, scanID string) *models.AIAnalysis {
	return &models.AIAnalysis{
		ScanID:          scanID,
		Summary:         parseAnalysisSummary(text),
		RiskScore:       parseRiskScore(text),
		CriticalIssues:  parseCriticalIssues(text),
		Recommendations: parseRecommendations(text),
		Timestamp:       time.Now(),
	}
}
