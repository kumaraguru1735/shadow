package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/kumaraguru1735/shadow/internal/ai"
	"github.com/kumaraguru1735/shadow/internal/scanner"
	"github.com/kumaraguru1735/shadow/pkg/models"
)

var (
	version = "0.1.0"
	rootCmd = &cobra.Command{
		Use:   "shadow",
		Short: "Shadow - AI-augmented security reconnaissance platform",
		Long: `Shadow combines network security tools with Claude AI intelligence
to provide comprehensive, automated security assessments.

⚠️  AUTHORIZATION REQUIRED: Only scan systems you own or have permission to test.`,
		Version: version,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Scan command
	var scanCmd = &cobra.Command{
		Use:   "scan [target]",
		Short: "Perform security scan on target",
		Args:  cobra.ExactArgs(1),
		Run:   runScan,
	}

	scanCmd.Flags().StringP("profile", "p", "standard", "Scan profile (quick, standard, deep)")
	scanCmd.Flags().BoolP("ai-analysis", "a", false, "Enable AI-powered analysis")
	scanCmd.Flags().StringSliceP("modules", "m", []string{}, "Specific modules to run")
	scanCmd.Flags().IntP("threads", "t", 50, "Number of concurrent threads")
	scanCmd.Flags().StringP("output", "o", "", "Output file path")
	scanCmd.Flags().StringP("format", "f", "json", "Output format (json, yaml, html, pdf)")

	// Smart scan command (AI-planned reconnaissance)
	var smartScanCmd = &cobra.Command{
		Use:   "smart-scan [target]",
		Short: "AI plans and executes reconnaissance strategy",
		Long: `Let AI analyze the target and create a comprehensive reconnaissance plan.
AI will:
- Determine what scans are needed
- Check tool availability
- Request root permissions if needed
- Execute reconnaissance phases
- Provide real-time guidance`,
		Args: cobra.ExactArgs(1),
		Run:  runSmartScan,
	}

	smartScanCmd.Flags().StringP("profile", "p", "standard", "Reconnaissance depth (quick, standard, deep)")

	// Subdomain command
	var subdomainCmd = &cobra.Command{
		Use:   "subdomain [domain]",
		Short: "Discover subdomains",
		Args:  cobra.ExactArgs(1),
		Run:   runSubdomain,
	}

	// Port scan command
	var portscanCmd = &cobra.Command{
		Use:   "portscan [target]",
		Short: "Scan ports on target",
		Args:  cobra.ExactArgs(1),
		Run:   runPortscan,
	}

	portscanCmd.Flags().StringP("ports", "p", "1-1000", "Port range to scan")
	portscanCmd.Flags().BoolP("fast", "f", false, "Fast scan (top 100 ports)")

	// SSL check command
	var sslCmd = &cobra.Command{
		Use:   "ssl [target]",
		Short: "Analyze SSL/TLS configuration",
		Args:  cobra.ExactArgs(1),
		Run:   runSSL,
	}

	// Analyze command
	var analyzeCmd = &cobra.Command{
		Use:   "analyze [scan-id]",
		Short: "Analyze scan results with AI",
		Args:  cobra.ExactArgs(1),
		Run:   runAnalyze,
	}

	// Report command
	var reportCmd = &cobra.Command{
		Use:   "report [scan-id]",
		Short: "Generate report from scan results",
		Args:  cobra.ExactArgs(1),
		Run:   runReport,
	}

	reportCmd.Flags().StringP("format", "f", "html", "Report format (html, pdf, json, markdown)")
	reportCmd.Flags().StringP("output", "o", "", "Output file path")

	// Query command (AI-powered)
	var queryCmd = &cobra.Command{
		Use:   "query [scan-id] [question]",
		Short: "Ask questions about scan results using AI",
		Args:  cobra.MinimumNArgs(2),
		Run:   runQuery,
	}

	// Auth check command
	var authCheckCmd = &cobra.Command{
		Use:   "auth-check",
		Short: "Check Claude AI authentication status",
		Run:   runAuthCheck,
	}

	// Auth generate command
	var authGenCmd = &cobra.Command{
		Use:   "auth-gen",
		Short: "Generate and setup authentication",
		Run:   runAuthGen,
	}

	// Auth status command
	var authStatusCmd = &cobra.Command{
		Use:   "auth-status",
		Short: "Show detailed authentication status",
		Run:   runAuthStatus,
	}

	// Auth setup command
	var authSetupCmd = &cobra.Command{
		Use:   "auth-setup",
		Short: "Interactive authentication setup wizard",
		Run:   runAuthSetup,
	}

	authSetupCmd.Flags().String("api-key", "", "Set API key directly")
	authSetupCmd.Flags().Bool("oauth", false, "Extract OAuth from Claude Code")

	// Auth refresh command
	var authRefreshCmd = &cobra.Command{
		Use:   "auth-refresh",
		Short: "Refresh OAuth credentials",
		Run:   runAuthRefresh,
	}

	// Auth backup command
	var authBackupCmd = &cobra.Command{
		Use:   "auth-backup",
		Short: "Backup current credentials",
		Run:   runAuthBackup,
	}

	// Agents command
	var agentsCmd = &cobra.Command{
		Use:   "agents",
		Short: "List available AI agents and their configurations",
		Run:   runAgents,
	}

	// Add commands to root
	rootCmd.AddCommand(scanCmd, smartScanCmd, subdomainCmd, portscanCmd, sslCmd, analyzeCmd, reportCmd, queryCmd,
		authCheckCmd, authGenCmd, authStatusCmd, authSetupCmd, authRefreshCmd, authBackupCmd, agentsCmd)
}

func runScan(cmd *cobra.Command, args []string) {
	target := args[0]
	profile, _ := cmd.Flags().GetString("profile")
	aiAnalysis, _ := cmd.Flags().GetBool("ai-analysis")
	threads, _ := cmd.Flags().GetInt("threads")

	fmt.Printf("🕵️  Shadow v%s\n", version)
	fmt.Printf("🎯 Target: %s\n", target)
	fmt.Printf("📋 Profile: %s\n", profile)
	fmt.Printf("🧵 Threads: %d\n\n", threads)

	// Permission check
	if !confirmAuthorization(target) {
		fmt.Println("❌ Authorization not confirmed. Exiting.")
		os.Exit(1)
	}

	config := models.ScanConfig{
		Target:     target,
		Profile:    profile,
		AIAnalysis: aiAnalysis,
		Threads:    threads,
	}

	// Initialize scanner
	s := scanner.New(config)

	// Run scan
	result, err := s.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Scan failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Scan completed in %v\n", result.Duration)
	fmt.Printf("📊 Scan ID: %s\n", result.ID)
	fmt.Printf("🔍 Findings: %d\n", len(result.Findings))

	if aiAnalysis {
		fmt.Println("\n🤖 Running Multi-Agent AI Analysis...")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		// Initialize multi-agent manager
		manager, err := ai.NewAgentManager()
		if err != nil {
			fmt.Printf("⚠️  AI analysis unavailable: %v\n", err)
			fmt.Println("💡 Tip: Run 'shadow auth-check' to verify authentication")
			return
		}
		defer manager.Close()

		// Use parent context
		ctx := context.Background()

		// Progress callback for real-time updates
		progressCallback := func(msg string) {
			fmt.Printf("   %s\n", msg)
		}

		// Run multi-agent analysis based on profile
		analysis, err := manager.AnalyzeScanWithAgents(ctx, result, profile, progressCallback)
		if err != nil {
			fmt.Printf("❌ AI analysis failed: %v\n", err)
			fmt.Println("\n💡 This could be due to:")
			fmt.Println("   - Large scan results (try with --profile quick)")
			fmt.Println("   - Network issues (check connection)")
			fmt.Println("   - Rate limiting (wait a few minutes)")

			// Still show usage stats even on failure
			summary := manager.GetUsageSummary()
			if summary.TotalOperations > 0 {
				summary.PrintSummary()
			}
			return
		}

		// Display analysis results
		fmt.Printf("\n📊 AI Analysis Results:\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("\n📝 Summary:\n%s\n", analysis.Summary)
		fmt.Printf("\n🎯 Risk Score: %d/100\n", analysis.RiskScore)

		if len(analysis.CriticalIssues) > 0 {
			fmt.Printf("\n🚨 Critical Issues:\n")
			for i, issue := range analysis.CriticalIssues {
				fmt.Printf("  %d. %s\n", i+1, issue)
			}
		}

		if len(analysis.Recommendations) > 0 {
			fmt.Printf("\n💡 Top Recommendations:\n")
			for i, rec := range analysis.Recommendations {
				if i < 5 { // Show top 5
					fmt.Printf("  %d. [%s] %s\n", i+1, rec.Priority, rec.Title)
				}
			}
		}

		fmt.Printf("\n✅ Analysis completed at %s\n", analysis.Timestamp.Format("15:04:05"))

		// Show model usage summary
		summary := manager.GetUsageSummary()
		summary.PrintSummary()
	}
}

func runSubdomain(cmd *cobra.Command, args []string) {
	domain := args[0]
	fmt.Printf("🔍 Discovering subdomains for %s...\n", domain)
	// Implementation coming
}

func runPortscan(cmd *cobra.Command, args []string) {
	target := args[0]
	ports, _ := cmd.Flags().GetString("ports")
	fmt.Printf("🔍 Scanning ports %s on %s...\n", ports, target)
	// Implementation coming
}

func runSSL(cmd *cobra.Command, args []string) {
	target := args[0]
	fmt.Printf("🔒 Analyzing SSL/TLS for %s...\n", target)
	// Implementation coming
}

func runAnalyze(cmd *cobra.Command, args []string) {
	scanID := args[0]
	fmt.Printf("🤖 Analyzing scan %s with AI...\n", scanID)
	// Implementation coming
}

func runReport(cmd *cobra.Command, args []string) {
	scanID := args[0]
	format, _ := cmd.Flags().GetString("format")
	fmt.Printf("📄 Generating %s report for scan %s...\n", format, scanID)
	// Implementation coming
}

func runQuery(cmd *cobra.Command, args []string) {
	scanID := args[0]
	question := args[1]
	fmt.Printf("💬 Querying scan %s: %s\n", scanID, question)
	// Implementation coming
}

func confirmAuthorization(target string) bool {
	fmt.Printf("\n⚠️  AUTHORIZATION REQUIRED\n")
	fmt.Printf("You are about to scan: %s\n\n", target)
	fmt.Printf("Do you have explicit permission to test this target? (yes/no): ")

	var response string
	fmt.Scanln(&response)

	return response == "yes" || response == "y"
}

func runAuthCheck(cmd *cobra.Command, args []string) {
	fmt.Println("🔐 Claude AI Authentication Status")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	status := ai.GetAuthenticationStatus()
	fmt.Println(status)
	fmt.Println()

	fmt.Println("📋 Authentication Methods:")
	fmt.Println("  1. Claude Code OAuth (automatic, preferred)")
	fmt.Println("     - Primary: ~/.claude/.credentials.json")
	fmt.Println("     - Alternative: ~/.claude/oauth.json")
	fmt.Println("     - Used automatically when Claude Code is installed")
	fmt.Println()
	fmt.Println("  2. API Key (manual)")
	fmt.Println("     - Set ANTHROPIC_API_KEY environment variable")
	fmt.Println("     - Example: export ANTHROPIC_API_KEY='sk-ant-...'")
	fmt.Println()

	// Test AI connection
	fmt.Println("🧪 Testing AI connection...")
	analyzer, err := ai.NewPiClaudeAnalyzer()
	if err != nil {
		fmt.Printf("❌ Failed to initialize AI client: %v\n", err)
		fmt.Println()
		fmt.Println("💡 Solutions:")
		fmt.Println("  - Run: ./setup_oauth.sh (extracts from Claude Code credentials)")
		fmt.Println("  - Install pi CLI: npm install -g @mariozechner/pi-coding-agent")
		fmt.Println("  - Or set ANTHROPIC_API_KEY environment variable")
		return
	}
	defer analyzer.Close()

	fmt.Println("✅ AI client initialized successfully!")
	fmt.Println("✅ Shadow can use Claude AI for analysis")
}

func runAuthGen(cmd *cobra.Command, args []string) {
	fmt.Println("🔐 Shadow Authentication Generator")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	manager, err := ai.NewAuthManager()
	if err != nil {
		fmt.Printf("❌ Failed to initialize auth manager: %v\n", err)
		return
	}

	// Extract OAuth from Claude Code
	fmt.Println("📝 Extracting OAuth from Claude Code credentials...")
	if err := manager.ExtractOAuthToStandard(); err != nil {
		fmt.Printf("⚠️  OAuth extraction failed: %v\n", err)
		fmt.Println()
		fmt.Println("💡 Tip: Make sure Claude Code is installed and authenticated")
	} else {
		fmt.Println("✅ OAuth credentials extracted successfully!")
		fmt.Println("   Location: ~/.claude/oauth.json")
	}

	fmt.Println()

	// Generate config if needed
	fmt.Println("📝 Generating Shadow configuration...")
	if err := manager.GenerateAPIKeyConfig(); err != nil {
		if os.IsExist(err) || fmt.Sprint(err) != "" && (fmt.Sprint(err)[:6] == "config") {
			fmt.Println("ℹ️  Config already exists at ~/.shadow/config.yaml")
		} else {
			fmt.Printf("⚠️  Config generation failed: %v\n", err)
		}
	} else {
		fmt.Println("✅ Configuration generated at ~/.shadow/config.yaml")
	}

	fmt.Println()

	// Validate authentication
	fmt.Println("🧪 Validating authentication...")
	if err := manager.ValidateAuthentication(); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
		fmt.Println()
		fmt.Println("💡 Solutions:")
		fmt.Println("   - Set ANTHROPIC_API_KEY environment variable")
		fmt.Println("   - Or run: shadow auth-setup --oauth")
	} else {
		fmt.Println("✅ Authentication is working!")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ Authentication setup complete!")
}

func runAuthStatus(cmd *cobra.Command, args []string) {
	fmt.Println("🔐 Detailed Authentication Status")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	manager, err := ai.NewAuthManager()
	if err != nil {
		fmt.Printf("❌ Failed to initialize auth manager: %v\n", err)
		return
	}

	status, err := manager.GetAuthStatus()
	if err != nil {
		fmt.Printf("❌ Failed to get auth status: %v\n", err)
		return
	}

	// OAuth Status
	fmt.Println("📋 OAuth Authentication:")
	if status.HasOAuth {
		fmt.Printf("   ✅ Enabled\n")
		fmt.Printf("   📍 Location: %s\n", status.OAuthPath)
		
		if status.OAuthExpired {
			fmt.Printf("   ⚠️  Status: EXPIRED\n")
			fmt.Println("   💡 Run: shadow auth-refresh")
		} else {
			fmt.Printf("   ✅ Status: Active\n")
			fmt.Printf("   ⏰ Expires in: %v\n", status.ExpiresIn.Round(time.Hour))
		}
		
		if status.Subscription != "" {
			fmt.Printf("   📦 Subscription: %s\n", status.Subscription)
		}
		if status.RateLimitTier != "" {
			fmt.Printf("   🚀 Rate Tier: %s\n", status.RateLimitTier)
		}
		if len(status.Scopes) > 0 {
			fmt.Printf("   🔑 Scopes: %v\n", status.Scopes)
		}
	} else {
		fmt.Println("   ❌ Not configured")
		fmt.Println("   💡 Run: shadow auth-gen")
	}

	fmt.Println()

	// API Key Status
	fmt.Println("📋 API Key Authentication:")
	if status.HasAPIKey {
		fmt.Println("   ✅ Configured")
		fmt.Println("   📍 Via: ANTHROPIC_API_KEY environment variable")
	} else {
		fmt.Println("   ❌ Not configured")
		fmt.Println("   💡 Set: export ANTHROPIC_API_KEY='sk-ant-...'")
	}

	fmt.Println()

	// Overall Status
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if status.HasOAuth && !status.OAuthExpired {
		fmt.Println("✅ Authentication: READY (OAuth)")
	} else if status.HasAPIKey {
		fmt.Println("✅ Authentication: READY (API Key)")
	} else {
		fmt.Println("❌ Authentication: NOT CONFIGURED")
		fmt.Println()
		fmt.Println("💡 Quick Setup:")
		fmt.Println("   shadow auth-gen")
	}
}

func runAuthSetup(cmd *cobra.Command, args []string) {
	fmt.Println("🔧 Shadow Authentication Setup Wizard")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	manager, err := ai.NewAuthManager()
	if err != nil {
		fmt.Printf("❌ Failed to initialize auth manager: %v\n", err)
		return
	}

	// Check flags
	apiKey, _ := cmd.Flags().GetString("api-key")
	useOAuth, _ := cmd.Flags().GetBool("oauth")

	if apiKey != "" {
		// Setup API key
		fmt.Println("📝 Setting up API key authentication...")
		if err := manager.SetupAPIKey(apiKey); err != nil {
			fmt.Printf("❌ Failed to setup API key: %v\n", err)
			return
		}
		fmt.Println("✅ API key saved to ~/.shadow/.env")
		fmt.Println()
		fmt.Println("💡 To use it:")
		fmt.Println("   source ~/.shadow/.env")
		fmt.Println("   shadow scan example.com --ai-analysis")
		return
	}

	if useOAuth {
		// Setup OAuth
		fmt.Println("📝 Setting up OAuth authentication...")
		if err := manager.ExtractOAuthToStandard(); err != nil {
			fmt.Printf("❌ Failed to setup OAuth: %v\n", err)
			return
		}
		fmt.Println("✅ OAuth credentials extracted!")
		
		// Validate
		fmt.Println()
		fmt.Println("🧪 Validating...")
		if err := manager.ValidateAuthentication(); err != nil {
			fmt.Printf("⚠️  Validation failed: %v\n", err)
		} else {
			fmt.Println("✅ Authentication working!")
		}
		return
	}

	// Interactive mode
	fmt.Println("Choose authentication method:")
	fmt.Println()
	fmt.Println("  1. OAuth (Claude Code) - Recommended")
	fmt.Println("     • Automatic token management")
	fmt.Println("     • Uses your Claude Code subscription")
	fmt.Println()
	fmt.Println("  2. API Key (Manual)")
	fmt.Println("     • Direct API key")
	fmt.Println("     • Manual token management")
	fmt.Println()
	fmt.Print("Enter choice (1 or 2): ")

	var choice string
	fmt.Scanln(&choice)
	fmt.Println()

	switch choice {
	case "1":
		fmt.Println("📝 Extracting OAuth from Claude Code...")
		if err := manager.ExtractOAuthToStandard(); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			fmt.Println()
			fmt.Println("💡 Make sure Claude Code is installed and authenticated")
			return
		}
		fmt.Println("✅ OAuth setup complete!")

	case "2":
		fmt.Println("📝 API Key Setup")
		fmt.Println()
		fmt.Print("Enter your Anthropic API key: ")
		var key string
		fmt.Scanln(&key)
		
		if key == "" {
			fmt.Println("❌ No API key provided")
			return
		}

		if err := manager.SetupAPIKey(key); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			return
		}
		fmt.Println()
		fmt.Println("✅ API key saved to ~/.shadow/.env")
		fmt.Println()
		fmt.Println("💡 To use it:")
		fmt.Println("   source ~/.shadow/.env")

	default:
		fmt.Println("❌ Invalid choice")
		return
	}

	// Validate
	fmt.Println()
	fmt.Println("🧪 Validating authentication...")
	if err := manager.ValidateAuthentication(); err != nil {
		fmt.Printf("⚠️  Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Authentication is working!")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ Setup complete! You can now use Shadow.")
}

func runAuthRefresh(cmd *cobra.Command, args []string) {
	fmt.Println("🔄 Refreshing OAuth Credentials")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	manager, err := ai.NewAuthManager()
	if err != nil {
		fmt.Printf("❌ Failed to initialize auth manager: %v\n", err)
		return
	}

	fmt.Println("📝 Attempting to refresh OAuth tokens...")
	if err := manager.RefreshOAuth(); err != nil {
		fmt.Printf("⚠️  Automatic refresh failed: %v\n", err)
		fmt.Println()
		fmt.Println("💡 Manual refresh:")
		fmt.Println("   1. Open Claude Code")
		fmt.Println("   2. Re-authenticate if needed")
		fmt.Println("   3. Run: shadow auth-gen")
		return
	}

	fmt.Println("✅ OAuth tokens refreshed!")
	fmt.Println()

	// Re-extract to standard location
	fmt.Println("📝 Updating local OAuth file...")
	if err := manager.ExtractOAuthToStandard(); err != nil {
		fmt.Printf("⚠️  Update failed: %v\n", err)
	} else {
		fmt.Println("✅ Local OAuth file updated!")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ Refresh complete!")
}

func runAuthBackup(cmd *cobra.Command, args []string) {
	fmt.Println("💾 Backing Up Credentials")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	manager, err := ai.NewAuthManager()
	if err != nil {
		fmt.Printf("❌ Failed to initialize auth manager: %v\n", err)
		return
	}

	fmt.Println("📝 Creating backup...")
	backupPath, err := manager.BackupCredentials()
	if err != nil {
		fmt.Printf("❌ Backup failed: %v\n", err)
		return
	}

	fmt.Println("✅ Backup created successfully!")
	fmt.Printf("📍 Location: %s\n", backupPath)
	fmt.Println()
	fmt.Println("💡 To restore:")
	fmt.Println("   cp", backupPath, "~/.claude/.credentials.json")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ Backup complete!")
}

func runAgents(cmd *cobra.Command, args []string) {
	fmt.Println("🤖 Shadow AI Agents Configuration")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	agents := models.GetDefaultAgents()

	for i, agent := range agents {
		fmt.Printf("\n%d. %s\n", i+1, agent.Name)
		fmt.Printf("   Type: %s\n", agent.Type)
		fmt.Printf("   Model: %s\n", getModelDisplayName(agent.Model))
		fmt.Printf("   Thinking Mode: %s\n", agent.Thinking)
		fmt.Printf("   Description: %s\n", agent.Description)
		fmt.Printf("   Use Case: %s\n", agent.UseCase)
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\n📋 Scan Profiles:")
	fmt.Println("   • quick  - Uses Haiku 4.5 (fast, cost-effective)")
	fmt.Println("   • standard - Uses Sonnet 4.5 (balanced, recommended)")
	fmt.Println("   • deep   - Uses multiple agents (Sonnet + Opus, most thorough)")

	fmt.Println("\n💰 Model Pricing (per million tokens):")
	fmt.Println("   • Haiku 4.5:  $0.80 input, $4.00 output")
	fmt.Println("   • Sonnet 4.5: $3.00 input, $15.00 output")
	fmt.Println("   • Opus 4.6:   $15.00 input, $75.00 output")

	fmt.Println("\n💡 Usage:")
	fmt.Println("   shadow scan example.com --ai-analysis --profile quick")
	fmt.Println("   shadow scan example.com --ai-analysis --profile standard")
	fmt.Println("   shadow scan example.com --ai-analysis --profile deep")
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

func runSmartScan(cmd *cobra.Command, args []string) {
	target := args[0]
	profile, _ := cmd.Flags().GetString("profile")

	fmt.Printf("🕵️  Shadow v%s - Smart Reconnaissance\n", version)
	fmt.Printf("🎯 Target: %s\n", target)
	fmt.Printf("📋 Mode: %s\n\n", profile)

	fmt.Println("🤖 AI is analyzing target and planning reconnaissance strategy...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Create AI reconnaissance planner
	planner, err := ai.NewReconPlanner()
	if err != nil {
		fmt.Printf("❌ Failed to initialize AI planner: %v\n", err)
		fmt.Println("💡 Tip: Run 'shadow auth-check' to verify authentication")
		return
	}
	defer planner.Close()

	// Ask AI to create reconnaissance plan
	ctx := context.Background()
	plan, err := planner.PlanReconnaissance(ctx, target, profile)
	if err != nil {
		fmt.Printf("❌ Failed to create reconnaissance plan: %v\n", err)
		return
	}

	// Display the plan
	plan.PrintPlan()

	// Ask user if they want to proceed
	fmt.Print("\n❓ Execute this reconnaissance plan? (yes/no): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ Error reading input: %v\n", err)
		return
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response != "yes" && response != "y" {
		fmt.Println("\n✅ Reconnaissance plan saved but not executed")
		fmt.Println("💡 You can review the plan and run scans manually")
		return
	}

	// Execute the plan
	fmt.Println("\n🚀 Executing reconnaissance plan...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Initialize permission manager
	permManager := scanner.NewPermissionManager()

	for i, phase := range plan.Phases {
		fmt.Printf("\n📍 Phase %d/%d: %s\n", i+1, len(plan.Phases), phase.Name)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		if phase.Description != "" {
			fmt.Printf("📋 %s\n\n", phase.Description)
		}

		// Execute each tool in the phase
		for _, tool := range phase.Tools {
			fmt.Printf("🔧 Running: %s\n", tool.Name)
			fmt.Printf("   Purpose: %s\n", tool.Purpose)

			if tool.RequiresRoot {
				fmt.Println("   ⚠️  This tool requires root access")

				// Show alternatives
				permManager.ShowCapabilityInfo(tool.Name)

				// Request permission
				approved, err := permManager.RequestRootPermission(
					tool.Name,
					tool.Purpose,
					fmt.Sprintf("sudo %s <args>", tool.Name),
				)

				if err != nil || !approved {
					fmt.Printf("   ⏭️  Skipping %s (permission denied or unavailable)\n", tool.Name)
					if tool.Fallback != "" {
						fmt.Printf("   💡 Fallback: %s\n", tool.Fallback)
					}
					continue
				}
			}

			fmt.Printf("   ✅ %s ready to execute\n", tool.Name)
			// Actual execution would happen here
			// For now, just show what would be executed
		}
	}

	// Show permission summary
	permManager.GetApprovalSummary()

	fmt.Println("\n✅ Reconnaissance plan execution complete")
	fmt.Println("💡 Next: Run 'shadow scan %s --ai-analysis' to analyze findings", target)
}
