package ai

import (
	"fmt"
	"sync"
	"time"
)

// ModelPricing contains pricing information for Claude models (per million tokens)
type ModelPricing struct {
	InputCostPerMToken  float64
	OutputCostPerMToken float64
}

var modelPricing = map[string]ModelPricing{
	"claude-opus-4.6": {
		InputCostPerMToken:  15.00,
		OutputCostPerMToken: 75.00,
	},
	"claude-sonnet-4.5": {
		InputCostPerMToken:  3.00,
		OutputCostPerMToken: 15.00,
	},
	"claude-sonnet-4.5-20250929": {
		InputCostPerMToken:  3.00,
		OutputCostPerMToken: 15.00,
	},
	"claude-haiku-4.5": {
		InputCostPerMToken:  0.80,
		OutputCostPerMToken: 4.00,
	},
}

// UsageStats tracks model usage for a single operation
type UsageStats struct {
	Model        string
	Agent        string
	InputTokens  int64
	OutputTokens int64
	Duration     time.Duration
	StartTime    time.Time
	EndTime      time.Time
	Success      bool
	Error        string
}

// CalculateCost estimates the cost of this usage
func (u *UsageStats) CalculateCost() float64 {
	pricing, ok := modelPricing[u.Model]
	if !ok {
		return 0.0
	}

	inputCost := (float64(u.InputTokens) / 1_000_000.0) * pricing.InputCostPerMToken
	outputCost := (float64(u.OutputTokens) / 1_000_000.0) * pricing.OutputCostPerMToken

	return inputCost + outputCost
}

// UsageTracker tracks all model usage across agents
type UsageTracker struct {
	mu     sync.RWMutex
	usages []UsageStats
}

// NewUsageTracker creates a new usage tracker
func NewUsageTracker() *UsageTracker {
	return &UsageTracker{
		usages: make([]UsageStats, 0),
	}
}

// RecordUsage adds a usage record
func (t *UsageTracker) RecordUsage(stats UsageStats) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.usages = append(t.usages, stats)
}

// GetSummary returns a summary of all usage
func (t *UsageTracker) GetSummary() UsageSummary {
	t.mu.RLock()
	defer t.mu.RUnlock()

	summary := UsageSummary{
		ByAgent: make(map[string]AgentSummary),
		ByModel: make(map[string]ModelSummary),
	}

	for _, usage := range t.usages {
		// Overall totals
		summary.TotalInputTokens += usage.InputTokens
		summary.TotalOutputTokens += usage.OutputTokens
		summary.TotalCost += usage.CalculateCost()
		summary.TotalDuration += usage.Duration
		summary.TotalOperations++
		if usage.Success {
			summary.SuccessfulOperations++
		}

		// By agent
		agentSummary := summary.ByAgent[usage.Agent]
		agentSummary.Agent = usage.Agent
		agentSummary.Model = usage.Model
		agentSummary.InputTokens += usage.InputTokens
		agentSummary.OutputTokens += usage.OutputTokens
		agentSummary.Cost += usage.CalculateCost()
		agentSummary.Duration += usage.Duration
		agentSummary.Operations++
		if usage.Success {
			agentSummary.Successes++
		}
		summary.ByAgent[usage.Agent] = agentSummary

		// By model
		modelSummary := summary.ByModel[usage.Model]
		modelSummary.Model = usage.Model
		modelSummary.InputTokens += usage.InputTokens
		modelSummary.OutputTokens += usage.OutputTokens
		modelSummary.Cost += usage.CalculateCost()
		modelSummary.Operations++
		summary.ByModel[usage.Model] = modelSummary
	}

	return summary
}

// UsageSummary provides aggregated usage statistics
type UsageSummary struct {
	TotalInputTokens      int64
	TotalOutputTokens     int64
	TotalCost             float64
	TotalDuration         time.Duration
	TotalOperations       int
	SuccessfulOperations  int
	ByAgent               map[string]AgentSummary
	ByModel               map[string]ModelSummary
}

// AgentSummary provides per-agent statistics
type AgentSummary struct {
	Agent        string
	Model        string
	InputTokens  int64
	OutputTokens int64
	Cost         float64
	Duration     time.Duration
	Operations   int
	Successes    int
}

// ModelSummary provides per-model statistics
type ModelSummary struct {
	Model        string
	InputTokens  int64
	OutputTokens int64
	Cost         float64
	Operations   int
}

// PrintSummary prints a formatted summary of usage
func (s *UsageSummary) PrintSummary() {
	fmt.Println("\n📊 AI Model Usage Summary")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Overall stats
	fmt.Printf("\n📈 Overall Statistics:\n")
	fmt.Printf("   Operations: %d/%d successful\n", s.SuccessfulOperations, s.TotalOperations)
	fmt.Printf("   Total Tokens: %s input, %s output\n",
		formatTokens(s.TotalInputTokens),
		formatTokens(s.TotalOutputTokens))
	fmt.Printf("   Estimated Cost: $%.4f\n", s.TotalCost)
	fmt.Printf("   Total Duration: %v\n", s.TotalDuration.Round(time.Second))

	// By agent
	if len(s.ByAgent) > 0 {
		fmt.Printf("\n🤖 By Agent:\n")
		for _, agent := range s.ByAgent {
			fmt.Printf("   %s (using %s)\n", agent.Agent, getModelShortName(agent.Model))
			fmt.Printf("      Tokens: %s in, %s out\n",
				formatTokens(agent.InputTokens),
				formatTokens(agent.OutputTokens))
			fmt.Printf("      Cost: $%.4f | Duration: %v | Success: %d/%d\n",
				agent.Cost,
				agent.Duration.Round(time.Second),
				agent.Successes,
				agent.Operations)
		}
	}

	// By model
	if len(s.ByModel) > 0 {
		fmt.Printf("\n🎯 By Model:\n")
		for _, model := range s.ByModel {
			fmt.Printf("   %s\n", getModelDisplayName(model.Model))
			fmt.Printf("      Tokens: %s in, %s out\n",
				formatTokens(model.InputTokens),
				formatTokens(model.OutputTokens))
			fmt.Printf("      Cost: $%.4f | Operations: %d\n",
				model.Cost,
				model.Operations)
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// Helper functions

func formatTokens(tokens int64) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	}
	return fmt.Sprintf("%.1fK", float64(tokens)/1000.0)
}

func getModelShortName(model string) string {
	switch model {
	case "claude-opus-4.6":
		return "Opus 4.6"
	case "claude-sonnet-4.5", "claude-sonnet-4.5-20250929":
		return "Sonnet 4.5"
	case "claude-haiku-4.5":
		return "Haiku 4.5"
	default:
		return model
	}
}

func getModelDisplayName(model string) string {
	switch model {
	case "claude-opus-4.6":
		return "Claude Opus 4.6 (most capable)"
	case "claude-sonnet-4.5", "claude-sonnet-4.5-20250929":
		return "Claude Sonnet 4.5 (balanced)"
	case "claude-haiku-4.5":
		return "Claude Haiku 4.5 (fast & efficient)"
	default:
		return model
	}
}
