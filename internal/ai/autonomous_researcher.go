package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	pi "github.com/joshp123/pi-golang"
	"github.com/kumaraguru1735/shadow/pkg/models"
)

// AutonomousSecurityResearcher conducts iterative security research
type AutonomousSecurityResearcher struct {
	client        *pi.OneShotClient
	findings      []models.Finding
	hypotheses    []SecurityHypothesis
	investigations []Investigation
	maxIterations int
}

// SecurityHypothesis represents AI's theory about potential vulnerabilities
type SecurityHypothesis struct {
	ID          string
	Description string
	Severity    string
	Indicators  []string
	NextSteps   []string
	Confidence  float64
}

// Investigation represents an AI-driven security investigation
type Investigation struct {
	ID           string
	Hypothesis   string
	Method       string
	Findings     string
	Conclusion   string
	FollowUp     []string
	Timestamp    time.Time
}

// NewAutonomousSecurityResearcher creates an autonomous AI security researcher
func NewAutonomousSecurityResearcher() (*AutonomousSecurityResearcher, error) {
	opts := pi.DefaultOneShotOptions()
	opts.AppName = "shadow-autonomous-researcher"
	opts.Mode = pi.ModeDragons
	opts.Dragons = pi.DragonsOptions{
		Provider: "anthropic",
		Model:    "claude-opus-4.6", // Use most capable model for deep thinking
		Thinking: "high",             // Maximum thinking depth
	}

	opts.SystemPrompt = `You are an elite autonomous security researcher and threat hunter.

Your capabilities:
- Deep critical thinking about security vulnerabilities
- Pattern recognition for backdoors and malicious code
- Hypothesis generation about potential threats
- Iterative investigation strategies
- Code analysis and reverse engineering mindset
- Advanced exploitation techniques understanding

Your approach:
1. THINK critically about what you observe
2. HYPOTHESIZE about potential vulnerabilities
3. INVESTIGATE specific attack vectors
4. ADAPT your strategy based on findings
5. ITERATE deeper into suspicious areas

You look for:
- Backdoors and hidden functionality
- Authentication bypasses
- Injection vulnerabilities
- Logic flaws and race conditions
- Misconfigurations
- Suspicious code patterns
- Hidden admin panels
- Information disclosure
- Supply chain vulnerabilities
- Zero-day potential

You reason like a real attacker:
- "What would I exploit here?"
- "What's the most likely attack path?"
- "What might the developers have missed?"
- "Where are the hidden entry points?"
- "What could go wrong in this implementation?"

Be thorough, creative, and think outside the box.`

	client, err := pi.StartOneShot(opts)
	if err != nil {
		return nil, err
	}

	return &AutonomousSecurityResearcher{
		client:        client,
		findings:      make([]models.Finding, 0),
		hypotheses:    make([]SecurityHypothesis, 0),
		investigations: make([]Investigation, 0),
		maxIterations: 5,
	}, nil
}

// ConductAutonomousResearch performs iterative AI-driven security research
func (asr *AutonomousSecurityResearcher) ConductAutonomousResearch(
	ctx context.Context,
	target string,
	initialFindings []models.Finding,
	progress ProgressCallback,
) (*AutonomousResearchReport, error) {
	asr.findings = initialFindings

	if progress != nil {
		progress("🧠 Starting autonomous security research...")
		progress(fmt.Sprintf("🎯 Target: %s", target))
		progress(fmt.Sprintf("📊 Initial findings: %d", len(initialFindings)))
		progress("")
	}

	report := &AutonomousResearchReport{
		Target:         target,
		StartTime:      time.Now(),
		Iterations:     make([]ResearchIteration, 0),
	}

	// Iteration 1: Initial Analysis & Hypothesis Generation
	if progress != nil {
		progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		progress("🔬 ITERATION 1: Initial Analysis & Hypothesis Generation")
		progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	iteration1, err := asr.initialAnalysis(ctx, target, initialFindings, progress)
	if err != nil {
		return nil, fmt.Errorf("iteration 1 failed: %w", err)
	}
	report.Iterations = append(report.Iterations, *iteration1)

	// Iteration 2: Backdoor & Hidden Threat Detection
	if progress != nil {
		progress("")
		progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		progress("🔬 ITERATION 2: Backdoor & Hidden Threat Detection")
		progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	iteration2, err := asr.backdoorDetection(ctx, target, iteration1.Findings, progress)
	if err != nil {
		return nil, fmt.Errorf("iteration 2 failed: %w", err)
	}
	report.Iterations = append(report.Iterations, *iteration2)

	// Iteration 3: Attack Path & Exploitation Analysis
	if progress != nil {
		progress("")
		progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		progress("🔬 ITERATION 3: Attack Path & Exploitation Analysis")
		progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	iteration3, err := asr.attackPathAnalysis(ctx, target, iteration2.Findings, progress)
	if err != nil {
		return nil, fmt.Errorf("iteration 3 failed: %w", err)
	}
	report.Iterations = append(report.Iterations, *iteration3)

	// Iteration 4: Deep Dive Investigations
	if progress != nil {
		progress("")
		progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		progress("🔬 ITERATION 4: Deep Dive Investigations")
		progress("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	iteration4, err := asr.deepDiveInvestigation(ctx, target, iteration3.NewHypotheses, progress)
	if err != nil {
		return nil, fmt.Errorf("iteration 4 failed: %w", err)
	}
	report.Iterations = append(report.Iterations, *iteration4)

	// Final Summary
	report.EndTime = time.Now()
	report.TotalDuration = report.EndTime.Sub(report.StartTime)
	report.FinalConclusions = asr.synthesizeFindings(report.Iterations)

	return report, nil
}

// initialAnalysis - AI thinks about initial findings
func (asr *AutonomousSecurityResearcher) initialAnalysis(
	ctx context.Context,
	target string,
	findings []models.Finding,
	progress ProgressCallback,
) (*ResearchIteration, error) {
	if progress != nil {
		progress("🤔 AI is thinking about what these findings might indicate...")
	}

	prompt := fmt.Sprintf(`# Initial Security Analysis

## Target
%s

## Initial Findings
%s

## Your Task
Analyze these findings and THINK DEEPLY about:

1. **What do these findings tell us?**
   - What vulnerabilities might exist?
   - What attack vectors are possible?
   - What might the developers have overlooked?

2. **Generate hypotheses about potential threats:**
   - Hidden backdoors or malicious code
   - Authentication/authorization bypasses
   - Injection vulnerabilities
   - Logic flaws
   - Configuration issues
   - Information disclosure

3. **Prioritize investigation areas:**
   - What's most critical to investigate?
   - What could lead to system compromise?
   - What needs deeper analysis?

4. **Recommend next steps:**
   - Specific tests to run
   - Endpoints to probe
   - Code patterns to search for

Think like an attacker. What would YOU target? What looks suspicious?

## Output Format
### CRITICAL THINKING
[Your reasoning about the findings]

### HYPOTHESES
1. [Hypothesis about potential vulnerability]
   - Indicators: [what suggests this]
   - Severity: [critical/high/medium/low]
   - Investigation plan: [how to verify]

### PRIORITY TARGETS
1. [Specific area to investigate]
   - Why: [reasoning]
   - How: [method]

### NEXT STEPS
[Concrete actions to take]`, target, formatFindingsDetailed(findings))

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	result, err := asr.client.Run(ctx, prompt)
	if err != nil {
		return nil, err
	}

	if progress != nil {
		progress("✅ Initial analysis complete")
		progress("")
		progress("📋 AI's Critical Thinking:")
		progress(extractSection(result.Text, "CRITICAL THINKING"))
		progress("")
	}

	iteration := &ResearchIteration{
		Number:        1,
		Phase:         "Initial Analysis",
		Findings:      result.Text,
		NewHypotheses: parseHypotheses(result.Text),
		NextSteps:     parseNextSteps(result.Text),
		Timestamp:     time.Now(),
	}

	return iteration, nil
}

// backdoorDetection - AI specifically looks for backdoors
func (asr *AutonomousSecurityResearcher) backdoorDetection(
	ctx context.Context,
	target string,
	previousFindings string,
	progress ProgressCallback,
) (*ResearchIteration, error) {
	if progress != nil {
		progress("🚪 AI is hunting for backdoors and hidden threats...")
	}

	prompt := fmt.Sprintf(`# Backdoor & Hidden Threat Detection

## Target
%s

## Previous Analysis
%s

## Your Mission
Think like an attacker who wants to hide malicious functionality.

Look for indicators of:

1. **Backdoors**
   - Hidden admin accounts
   - Secret URL parameters
   - Hardcoded credentials
   - Debug/test endpoints left in production
   - Undocumented API endpoints
   - Comment-based commands

2. **Hidden Functionality**
   - Easter eggs that could be exploited
   - Developer tools exposed in production
   - Hidden file upload locations
   - Unreferenced admin panels
   - Backup/old files (.bak, .old, ~)

3. **Suspicious Code Patterns**
   - Obfuscated JavaScript
   - Base64 encoded strings (could be backdoor code)
   - Suspicious redirects or iframes
   - Unusual cookie handling
   - Strange request headers being checked

4. **Supply Chain Indicators**
   - Suspicious third-party scripts
   - CDN resources from unusual domains
   - Outdated dependencies with known backdoors

5. **Data Exfiltration Paths**
   - External requests to unusual domains
   - Suspicious logging mechanisms
   - Unnecessary data being sent to third parties

## Think Creatively
- What would a sophisticated attacker hide here?
- What "features" could be abused as backdoors?
- What looks innocent but might be malicious?

## Output Format
### BACKDOOR ANALYSIS
[Your detailed analysis]

### SUSPECTED BACKDOORS
1. [Specific backdoor indicator]
   - Type: [backdoor type]
   - Location: [where found]
   - Exploitation: [how to use it]
   - Proof: [evidence]

### HIDDEN THREATS
[Other suspicious findings]

### VERIFICATION STEPS
[How to confirm these threats]`, target, previousFindings)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	result, err := asr.client.Run(ctx, prompt)
	if err != nil {
		return nil, err
	}

	if progress != nil {
		progress("✅ Backdoor detection complete")
		progress("")
		backdoors := extractSection(result.Text, "SUSPECTED BACKDOORS")
		if backdoors != "" {
			progress("🚨 AI found potential backdoors:")
			progress(backdoors)
		} else {
			progress("✓ No obvious backdoors detected")
		}
		progress("")
	}

	iteration := &ResearchIteration{
		Number:        2,
		Phase:         "Backdoor Detection",
		Findings:      result.Text,
		NewHypotheses: parseHypotheses(result.Text),
		NextSteps:     parseNextSteps(result.Text),
		Timestamp:     time.Now(),
	}

	return iteration, nil
}

// attackPathAnalysis - AI maps complete attack chains
func (asr *AutonomousSecurityResearcher) attackPathAnalysis(
	ctx context.Context,
	target string,
	previousFindings string,
	progress ProgressCallback,
) (*ResearchIteration, error) {
	if progress != nil {
		progress("🎯 AI is mapping complete attack paths...")
	}

	prompt := fmt.Sprintf(`# Attack Path & Exploitation Analysis

## Target
%s

## Findings So Far
%s

## Your Task
Map COMPLETE attack chains from initial access to full compromise.

Think step-by-step:
1. What's the weakest entry point?
2. After getting in, what's next?
3. How to escalate privileges?
4. How to maintain persistence?
5. What's the ultimate impact?

Analyze:

### ATTACK CHAINS
For each vulnerability found, map the COMPLETE path:
- Initial Access → Exploitation → Privilege Escalation → Impact

### CHAINED VULNERABILITIES
Can multiple vulnerabilities be chained together?
- Info disclosure + injection = RCE?
- CSRF + XSS = Account takeover?
- Multiple bugs combined for greater impact?

### EXPLOITATION PREREQUISITES
What does an attacker need?
- Specific user roles?
- Authenticated access?
- Specific timing/race conditions?
- External services?

### REAL-WORLD SCENARIOS
How would this be exploited in practice?
- What would a real attacker do?
- What's the most likely attack path?
- What's the highest impact path?

## Output Format
### ATTACK CHAIN 1: [Name]
Entry Point: [vulnerability]
↓
Step 1: [action]
↓
Step 2: [action]
↓
Final Impact: [compromise level]

Difficulty: [easy/medium/hard]
Detectability: [easy/medium/hard]
Impact: [low/medium/high/critical]

### MOST LIKELY ATTACK
[The path a real attacker would take]

### HIGHEST IMPACT ATTACK
[The path causing most damage]`, target, previousFindings)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	result, err := asr.client.Run(ctx, prompt)
	if err != nil {
		return nil, err
	}

	if progress != nil {
		progress("✅ Attack path analysis complete")
		progress("")
		mostLikely := extractSection(result.Text, "MOST LIKELY ATTACK")
		if mostLikely != "" {
			progress("⚡ Most likely attack path:")
			progress(mostLikely)
		}
		progress("")
	}

	iteration := &ResearchIteration{
		Number:        3,
		Phase:         "Attack Path Analysis",
		Findings:      result.Text,
		NewHypotheses: parseHypotheses(result.Text),
		NextSteps:     parseNextSteps(result.Text),
		Timestamp:     time.Now(),
	}

	return iteration, nil
}

// deepDiveInvestigation - AI conducts deep investigation of specific findings
func (asr *AutonomousSecurityResearcher) deepDiveInvestigation(
	ctx context.Context,
	target string,
	hypotheses []string,
	progress ProgressCallback,
) (*ResearchIteration, error) {
	if progress != nil {
		progress("🔬 AI is conducting deep dive investigations...")
	}

	prompt := fmt.Sprintf(`# Deep Dive Security Investigation

## Target
%s

## Hypotheses to Investigate
%s

## Your Mission
Go DEEP on the most critical findings.

For each critical issue:

1. **Root Cause Analysis**
   - WHY does this vulnerability exist?
   - What went wrong in the implementation?
   - Is this a pattern across the application?

2. **Exploitation Details**
   - Exact steps to exploit
   - Payloads that would work
   - Bypass techniques
   - Edge cases and variants

3. **Impact Analysis**
   - Worst-case scenario
   - Data that could be compromised
   - Systems that could be accessed
   - Business impact

4. **Similar Vulnerabilities**
   - Where else might this pattern exist?
   - Related vulnerabilities to check
   - Common variations

5. **Defense Evasion**
   - Would this bypass WAF/IDS?
   - How to avoid detection?
   - What logs would be generated?

Think at the DEEPEST level. Consider edge cases, race conditions, timing issues.

## Output Format
### DEEP DIVE: [Vulnerability Name]

**Root Cause:**
[Why it exists]

**Exploitation:**
[Step-by-step]

**Payloads:**
[Specific payloads]

**Impact:**
[Detailed impact]

**Detection:**
[How to detect exploitation]

**Related Issues:**
[Similar vulnerabilities]`, target, strings.Join(hypotheses, "\n"))

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	result, err := asr.client.Run(ctx, prompt)
	if err != nil {
		return nil, err
	}

	if progress != nil {
		progress("✅ Deep dive investigation complete")
	}

	iteration := &ResearchIteration{
		Number:        4,
		Phase:         "Deep Dive Investigation",
		Findings:      result.Text,
		NewHypotheses: make([]string, 0),
		NextSteps:     parseNextSteps(result.Text),
		Timestamp:     time.Now(),
	}

	return iteration, nil
}

// ResearchIteration represents one iteration of autonomous research
type ResearchIteration struct {
	Number        int
	Phase         string
	Findings      string
	NewHypotheses []string
	NextSteps     []string
	Timestamp     time.Time
}

// AutonomousResearchReport contains complete research results
type AutonomousResearchReport struct {
	Target           string
	StartTime        time.Time
	EndTime          time.Time
	TotalDuration    time.Duration
	Iterations       []ResearchIteration
	FinalConclusions string
}

// Helper functions

func formatFindingsDetailed(findings []models.Finding) string {
	if len(findings) == 0 {
		return "No initial findings provided."
	}

	var result strings.Builder
	for i, finding := range findings {
		result.WriteString(fmt.Sprintf("\n%d. [%s] %s\n",
			i+1, finding.Severity, finding.Title))
		if finding.Description != "" {
			result.WriteString(fmt.Sprintf("   Description: %s\n", finding.Description))
		}
		if finding.Evidence != "" {
			result.WriteString(fmt.Sprintf("   Evidence: %s\n", finding.Evidence))
		}
		if finding.Location != "" {
			result.WriteString(fmt.Sprintf("   Location: %s\n", finding.Location))
		}
	}
	return result.String()
}

func extractSection(text string, sectionName string) string {
	lines := strings.Split(text, "\n")
	var section strings.Builder
	inSection := false

	for _, line := range lines {
		if strings.Contains(line, sectionName) {
			inSection = true
			continue
		}
		if inSection && strings.HasPrefix(line, "###") {
			break
		}
		if inSection && strings.TrimSpace(line) != "" {
			section.WriteString(line + "\n")
		}
	}

	return strings.TrimSpace(section.String())
}

func parseHypotheses(text string) []string {
	hypotheses := make([]string, 0)
	lines := strings.Split(text, "\n")
	inHypotheses := false

	for _, line := range lines {
		if strings.Contains(line, "HYPOTHESES") || strings.Contains(line, "SUSPECTED") {
			inHypotheses = true
			continue
		}
		if inHypotheses && strings.HasPrefix(strings.TrimSpace(line), "-") {
			hypotheses = append(hypotheses, strings.TrimSpace(line))
		}
		if inHypotheses && strings.HasPrefix(line, "###") {
			break
		}
	}

	return hypotheses
}

func parseNextSteps(text string) []string {
	steps := make([]string, 0)
	lines := strings.Split(text, "\n")
	inNextSteps := false

	for _, line := range lines {
		if strings.Contains(line, "NEXT STEPS") || strings.Contains(line, "VERIFICATION") {
			inNextSteps = true
			continue
		}
		if inNextSteps && (strings.HasPrefix(strings.TrimSpace(line), "-") || strings.HasPrefix(strings.TrimSpace(line), "1.")) {
			steps = append(steps, strings.TrimSpace(line))
		}
	}

	return steps
}

func (asr *AutonomousSecurityResearcher) synthesizeFindings(iterations []ResearchIteration) string {
	var synthesis strings.Builder

	synthesis.WriteString("# Final Security Assessment\n\n")
	synthesis.WriteString("## Research Summary\n")
	synthesis.WriteString(fmt.Sprintf("Conducted %d iterations of autonomous security research.\n\n", len(iterations)))

	synthesis.WriteString("## Key Findings\n")
	for _, iteration := range iterations {
		synthesis.WriteString(fmt.Sprintf("### %s\n", iteration.Phase))
		synthesis.WriteString(extractSection(iteration.Findings, "CRITICAL") + "\n\n")
	}

	return synthesis.String()
}

// PrintReport displays the complete autonomous research report
func (report *AutonomousResearchReport) PrintReport() {
	fmt.Println("\n🧠 Autonomous Security Research Report")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\n🎯 Target: %s\n", report.Target)
	fmt.Printf("⏱️  Duration: %v\n", report.TotalDuration.Round(time.Second))
	fmt.Printf("🔬 Iterations: %d\n", len(report.Iterations))

	for _, iteration := range report.Iterations {
		fmt.Printf("\n━━ Iteration %d: %s ━━\n", iteration.Number, iteration.Phase)
		fmt.Println(iteration.Findings[:min(500, len(iteration.Findings))] + "...")
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Close closes the researcher
func (asr *AutonomousSecurityResearcher) Close() {
	if asr.client != nil {
		asr.client.Close()
	}
}
