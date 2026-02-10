package models

// AgentType represents different types of AI agents
type AgentType string

const (
	AgentTypeRecon        AgentType = "reconnaissance"
	AgentTypeVulnerability AgentType = "vulnerability"
	AgentTypeExploitation  AgentType = "exploitation"
	AgentTypeReport        AgentType = "report"
	AgentTypeQuickScan     AgentType = "quick-scan"
)

// AgentConfig defines configuration for an AI agent
type AgentConfig struct {
	Name         string
	Type         AgentType
	Model        string
	Thinking     string // "low", "high"
	SystemPrompt string
	Description  string
	UseCase      string
}

// GetDefaultAgents returns the default agent configurations
func GetDefaultAgents() []AgentConfig {
	return []AgentConfig{
		{
			Name:        "Quick Scanner",
			Type:        AgentTypeQuickScan,
			Model:       "claude-haiku-4.5",
			Thinking:    "low",
			Description: "Fast initial scan analysis",
			UseCase:     "Quick triage and basic vulnerability identification",
		},
		{
			Name:        "Reconnaissance Analyst",
			Type:        AgentTypeRecon,
			Model:       "claude-sonnet-4.5-20250929",
			Thinking:    "high",
			Description: "Deep reconnaissance and attack surface analysis",
			UseCase:     "Technology identification, service enumeration, attack surface mapping",
		},
		{
			Name:        "Vulnerability Researcher",
			Type:        AgentTypeVulnerability,
			Model:       "claude-sonnet-4.5-20250929",
			Thinking:    "high",
			Description: "Comprehensive vulnerability analysis",
			UseCase:     "OWASP Top 10, CVE research, vulnerability prioritization",
		},
		{
			Name:        "Exploitation Specialist",
			Type:        AgentTypeExploitation,
			Model:       "claude-opus-4.6",
			Thinking:    "high",
			Description: "Advanced exploitation path analysis",
			UseCase:     "Attack chain development, exploitation techniques, proof-of-concept",
		},
		{
			Name:        "Security Reporter",
			Type:        AgentTypeReport,
			Model:       "claude-sonnet-4.5-20250929",
			Thinking:    "high",
			Description: "Executive and technical report generation",
			UseCase:     "Risk assessment, executive summaries, remediation roadmaps",
		},
	}
}

// GetAgentByType returns the agent configuration for a given type
func GetAgentByType(agentType AgentType) *AgentConfig {
	agents := GetDefaultAgents()
	for i := range agents {
		if agents[i].Type == agentType {
			return &agents[i]
		}
	}
	return nil
}
