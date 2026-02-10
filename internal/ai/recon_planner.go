package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	pi "github.com/joshp123/pi-golang"
)

// ReconPlanner uses AI to plan reconnaissance strategy
type ReconPlanner struct {
	client *pi.OneShotClient
}

// ReconPlan represents the AI's reconnaissance strategy
type ReconPlan struct {
	Target          string
	Phases          []ReconPhase
	RequiresRoot    bool
	RequiredTools   []string
	EstimatedTime   string
	Reasoning       string
}

// ReconPhase represents a single phase of reconnaissance
type ReconPhase struct {
	Name            string
	Description     string
	Tools           []ToolRequirement
	Priority        string // "critical", "high", "medium", "low"
	ExpectedOutputs []string
}

// ToolRequirement defines what a tool needs to run
type ToolRequirement struct {
	Name         string
	Command      string
	RequiresRoot bool
	Flags        []string
	Purpose      string
	Fallback     string // Alternative if tool unavailable
}

// NewReconPlanner creates a new reconnaissance planner
func NewReconPlanner() (*ReconPlanner, error) {
	opts := pi.DefaultOneShotOptions()
	opts.AppName = "shadow-recon-planner"
	opts.Mode = pi.ModeDragons
	opts.Dragons = pi.DragonsOptions{
		Provider: "anthropic",
		Model:    "claude-sonnet-4.5-20250929",
		Thinking: "high",
	}

	opts.SystemPrompt = `You are an expert penetration tester and reconnaissance specialist.

Your role is to:
1. Analyze a target URL/domain
2. Plan comprehensive reconnaissance strategy
3. Determine what tools and scans are needed
4. Identify permission requirements (root access, etc.)
5. Prioritize reconnaissance phases
6. Suggest fallback options when tools are unavailable

Consider:
- Target type (web app, API, infrastructure)
- Available tools (nmap, subfinder, whatweb, curl, dig, etc.)
- Permission requirements (root for SYN scans, etc.)
- Information gathering priorities
- OSINT opportunities
- Legal and ethical boundaries

Provide structured, executable reconnaissance plans.`

	client, err := pi.StartOneShot(opts)
	if err != nil {
		return nil, err
	}

	return &ReconPlanner{client: client}, nil
}

// PlanReconnaissance asks AI to create a reconnaissance plan
func (rp *ReconPlanner) PlanReconnaissance(ctx context.Context, target string, mode string) (*ReconPlan, error) {
	prompt := fmt.Sprintf(`# Reconnaissance Planning Request

## Target
%s

## Mode
%s (quick/standard/deep)

## Available Tools
The following tools may be available:
- nmap (port scanning - requires root for SYN scans, falls back to TCP connect)
- subfinder (subdomain enumeration)
- whatweb (web technology detection)
- curl/wget (HTTP requests)
- dig/nslookup (DNS queries)
- whois (domain information)
- openssl (SSL/TLS analysis)
- Go-based HTTP scanner (always available, no root needed)

## Task
Create a comprehensive reconnaissance plan for this target.

## Output Format
Provide your plan in the following format:

### OVERVIEW
Brief description of target and reconnaissance approach

### PHASE 1: [Phase Name]
Priority: [critical/high/medium/low]
Description: [What this phase accomplishes]
Tools needed:
- [tool name] (requires root: yes/no) - [purpose]
- [tool name] (requires root: yes/no) - [purpose]
Expected outputs: [what we'll learn]

### PHASE 2: [Phase Name]
[Same format...]

### PERMISSIONS REQUIRED
- Root access: [yes/no and why]
- Sudo for specific commands: [list if needed]

### FALLBACK OPTIONS
If root not available: [alternative approach]
If tool X not available: [alternative]

### ESTIMATED TIME
[time estimate for full reconnaissance]

### REASONING
[Why this approach is optimal for this target]

Be specific about commands and explain your reasoning.`, target, mode)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	result, err := rp.client.Run(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to create recon plan: %w", err)
	}

	// Parse the AI's response into a structured plan
	plan := rp.parseReconPlan(result.Text, target)
	return plan, nil
}

// parseReconPlan converts AI's text response into structured ReconPlan
func (rp *ReconPlanner) parseReconPlan(response string, target string) *ReconPlan {
	plan := &ReconPlan{
		Target: target,
		Phases: make([]ReconPhase, 0),
	}

	lines := strings.Split(response, "\n")
	var currentPhase *ReconPhase
	var inPermissions, inReasoning bool

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse phases
		if strings.HasPrefix(line, "### PHASE") {
			if currentPhase != nil {
				plan.Phases = append(plan.Phases, *currentPhase)
			}
			currentPhase = &ReconPhase{
				Name:            line,
				Tools:           make([]ToolRequirement, 0),
				ExpectedOutputs: make([]string, 0),
			}
			inPermissions = false
			inReasoning = false
		}

		// Parse permissions section
		if strings.HasPrefix(line, "### PERMISSIONS REQUIRED") {
			inPermissions = true
			inReasoning = false
			continue
		}

		// Parse reasoning section
		if strings.HasPrefix(line, "### REASONING") {
			inReasoning = true
			inPermissions = false
			if currentPhase != nil {
				plan.Phases = append(plan.Phases, *currentPhase)
				currentPhase = nil
			}
			continue
		}

		// Extract root requirement
		if inPermissions && strings.Contains(strings.ToLower(line), "root access: yes") {
			plan.RequiresRoot = true
		}

		// Extract reasoning
		if inReasoning && line != "" && !strings.HasPrefix(line, "#") {
			plan.Reasoning += line + " "
		}

		// Parse phase details
		if currentPhase != nil {
			if strings.HasPrefix(line, "Priority:") {
				priority := strings.TrimPrefix(line, "Priority:")
				currentPhase.Priority = strings.TrimSpace(priority)
			} else if strings.HasPrefix(line, "Description:") {
				desc := strings.TrimPrefix(line, "Description:")
				currentPhase.Description = strings.TrimSpace(desc)
			} else if strings.HasPrefix(line, "- ") && strings.Contains(line, "(requires root:") {
				// Parse tool requirement
				tool := rp.parseToolRequirement(line)
				if tool.Name != "" {
					currentPhase.Tools = append(currentPhase.Tools, tool)
					plan.RequiredTools = append(plan.RequiredTools, tool.Name)
				}
			}
		}
	}

	// Add last phase
	if currentPhase != nil {
		plan.Phases = append(plan.Phases, *currentPhase)
	}

	// Store full reasoning
	plan.Reasoning = strings.TrimSpace(plan.Reasoning)

	return plan
}

// parseToolRequirement extracts tool details from a line
func (rp *ReconPlanner) parseToolRequirement(line string) ToolRequirement {
	// Example: "- nmap (requires root: yes) - Port scanning"
	tool := ToolRequirement{}

	// Extract tool name
	if idx := strings.Index(line, "(requires root:"); idx > 0 {
		tool.Name = strings.TrimSpace(strings.TrimPrefix(line[:idx], "- "))
	}

	// Check if requires root
	if strings.Contains(strings.ToLower(line), "requires root: yes") {
		tool.RequiresRoot = true
	}

	// Extract purpose (after the dash following root requirement)
	if parts := strings.Split(line, "-"); len(parts) >= 3 {
		tool.Purpose = strings.TrimSpace(parts[len(parts)-1])
	}

	return tool
}

// PrintPlan displays the reconnaissance plan to the user
func (plan *ReconPlan) PrintPlan() {
	fmt.Println("\n🎯 AI-Generated Reconnaissance Plan")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\n📍 Target: %s\n", plan.Target)

	if plan.Reasoning != "" {
		fmt.Printf("\n💭 Strategy:\n%s\n", wrapText(plan.Reasoning, 70))
	}

	// Show phases
	fmt.Printf("\n📋 Reconnaissance Phases (%d):\n", len(plan.Phases))
	for i, phase := range plan.Phases {
		fmt.Printf("\n%d. %s\n", i+1, phase.Name)
		if phase.Priority != "" {
			fmt.Printf("   Priority: %s\n", getPriorityEmoji(phase.Priority))
		}
		if phase.Description != "" {
			fmt.Printf("   %s\n", phase.Description)
		}
		if len(phase.Tools) > 0 {
			fmt.Println("   Tools needed:")
			for _, tool := range phase.Tools {
				rootBadge := ""
				if tool.RequiresRoot {
					rootBadge = " [ROOT REQUIRED]"
				}
				fmt.Printf("      • %s%s - %s\n", tool.Name, rootBadge, tool.Purpose)
			}
		}
	}

	// Show permission requirements
	fmt.Println("\n🔐 Permissions:")
	if plan.RequiresRoot {
		fmt.Println("   ⚠️  Root/sudo access required for some scans")
		fmt.Println("   💡 Shadow will ask for permission before running privileged commands")
	} else {
		fmt.Println("   ✓ No elevated permissions needed")
	}

	// Show required tools
	if len(plan.RequiredTools) > 0 {
		fmt.Println("\n🛠️  Required Tools:")
		for _, tool := range plan.RequiredTools {
			fmt.Printf("   • %s\n", tool)
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// Helper functions

func getPriorityEmoji(priority string) string {
	switch strings.ToLower(priority) {
	case "critical":
		return "🔴 Critical"
	case "high":
		return "🟠 High"
	case "medium":
		return "🟡 Medium"
	case "low":
		return "🟢 Low"
	default:
		return priority
	}
}

func wrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine)+len(word)+1 > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// Close closes the planner
func (rp *ReconPlanner) Close() {
	if rp.client != nil {
		rp.client.Close()
	}
}
