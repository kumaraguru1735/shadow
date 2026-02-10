package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// AuthManager handles authentication generation and management
type AuthManager struct {
	homeDir string
}

// NewAuthManager creates a new authentication manager
func NewAuthManager() (*AuthManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	return &AuthManager{
		homeDir: home,
	}, nil
}

// OAuthCredentials represents Claude Code OAuth structure
type OAuthCredentials struct {
	AccessToken      string   `json:"accessToken"`
	RefreshToken     string   `json:"refreshToken"`
	ExpiresAt        int64    `json:"expiresAt"`
	Scopes           []string `json:"scopes"`
	SubscriptionType string   `json:"subscriptionType"`
	RateLimitTier    string   `json:"rateLimitTier"`
}

// ClaudeCredentials represents the full credentials file structure
type ClaudeCredentials struct {
	ClaudeAiOauth OAuthCredentials `json:"claudeAiOauth"`
}

// AuthStatus represents the current authentication status
type AuthStatus struct {
	HasOAuth       bool
	HasAPIKey      bool
	OAuthPath      string
	OAuthExpired   bool
	ExpiresIn      time.Duration
	Subscription   string
	RateLimitTier  string
	Scopes         []string
}

// GetAuthStatus checks the current authentication status
func (m *AuthManager) GetAuthStatus() (*AuthStatus, error) {
	status := &AuthStatus{}

	// Check for OAuth credentials
	claudeCredsPath := filepath.Join(m.homeDir, ".claude", ".credentials.json")
	if _, err := os.Stat(claudeCredsPath); err == nil {
		status.HasOAuth = true
		status.OAuthPath = claudeCredsPath

		// Read and parse credentials
		data, err := os.ReadFile(claudeCredsPath)
		if err == nil {
			var creds ClaudeCredentials
			if json.Unmarshal(data, &creds) == nil {
				// Check expiration
				expiresAt := time.Unix(creds.ClaudeAiOauth.ExpiresAt/1000, 0)
				now := time.Now()

				if expiresAt.Before(now) {
					status.OAuthExpired = true
				} else {
					status.ExpiresIn = expiresAt.Sub(now)
				}

				status.Subscription = creds.ClaudeAiOauth.SubscriptionType
				status.RateLimitTier = creds.ClaudeAiOauth.RateLimitTier
				status.Scopes = creds.ClaudeAiOauth.Scopes
			}
		}
	}

	// Check for API key
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		status.HasAPIKey = true
	}

	return status, nil
}

// ExtractOAuthToStandard extracts OAuth from Claude Code credentials to standard location
func (m *AuthManager) ExtractOAuthToStandard() error {
	claudeCredsPath := filepath.Join(m.homeDir, ".claude", ".credentials.json")
	oauthPath := filepath.Join(m.homeDir, ".claude", "oauth.json")

	// Read Claude Code credentials
	data, err := os.ReadFile(claudeCredsPath)
	if err != nil {
		return fmt.Errorf("failed to read credentials: %w", err)
	}

	// Parse credentials
	var creds ClaudeCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return fmt.Errorf("failed to parse credentials: %w", err)
	}

	// Write OAuth to standard location
	oauthData, err := json.MarshalIndent(creds.ClaudeAiOauth, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal OAuth: %w", err)
	}

	if err := os.WriteFile(oauthPath, oauthData, 0600); err != nil {
		return fmt.Errorf("failed to write OAuth file: %w", err)
	}

	return nil
}

// GenerateAPIKeyConfig creates a config file with API key placeholder
func (m *AuthManager) GenerateAPIKeyConfig() error {
	shadowDir := filepath.Join(m.homeDir, ".shadow")
	configPath := filepath.Join(shadowDir, "config.yaml")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(shadowDir, 0755); err != nil {
		return fmt.Errorf("failed to create shadow directory: %w", err)
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config already exists at %s", configPath)
	}

	config := `# Shadow Configuration
# Generated: ` + time.Now().Format(time.RFC3339) + `

# Anthropic Claude AI Settings
anthropic:
  api_key: ${ANTHROPIC_API_KEY}  # Set via environment variable
  model: claude-sonnet-4.5-20250929
  max_tokens: 4096

# Scanning Configuration
scanning:
  threads: 50
  timeout: 30s
  rate_limit: 100

# AI Analysis Configuration
ai:
  enabled: true
  auto_analyze: false
  retry_attempts: 3
  retry_delay: 15s
`

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// ValidateAuthentication tests if authentication works
func (m *AuthManager) ValidateAuthentication() error {
	// Try to initialize pi client
	analyzer, err := NewAdvancedClaudeAnalyzer()
	if err != nil {
		return fmt.Errorf("authentication validation failed: %w", err)
	}
	defer analyzer.Close()

	return nil
}

// RefreshOAuth attempts to refresh OAuth tokens using Claude Code
func (m *AuthManager) RefreshOAuth() error {
	// Check if Claude Code CLI is available
	if _, err := exec.LookPath("claude"); err != nil {
		return fmt.Errorf("Claude Code CLI not found - cannot refresh OAuth")
	}

	// Try to refresh using Claude Code
	cmd := exec.Command("claude", "auth", "refresh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("OAuth refresh failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// ShowOAuthToken displays OAuth token information (masked)
func (m *AuthManager) ShowOAuthToken() error {
	claudeCredsPath := filepath.Join(m.homeDir, ".claude", ".credentials.json")

	data, err := os.ReadFile(claudeCredsPath)
	if err != nil {
		return fmt.Errorf("failed to read credentials: %w", err)
	}

	var creds ClaudeCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return fmt.Errorf("failed to parse credentials: %w", err)
	}

	// Mask tokens
	accessToken := creds.ClaudeAiOauth.AccessToken
	refreshToken := creds.ClaudeAiOauth.RefreshToken

	if len(accessToken) > 20 {
		accessToken = accessToken[:20] + "..."
	}
	if len(refreshToken) > 20 {
		refreshToken = refreshToken[:20] + "..."
	}

	fmt.Println("OAuth Token Information:")
	fmt.Printf("  Access Token:  %s\n", accessToken)
	fmt.Printf("  Refresh Token: %s\n", refreshToken)
	fmt.Printf("  Expires At:    %s\n", time.Unix(creds.ClaudeAiOauth.ExpiresAt/1000, 0).Format(time.RFC3339))
	fmt.Printf("  Scopes:        %v\n", creds.ClaudeAiOauth.Scopes)
	fmt.Printf("  Subscription:  %s\n", creds.ClaudeAiOauth.SubscriptionType)
	fmt.Printf("  Rate Tier:     %s\n", creds.ClaudeAiOauth.RateLimitTier)

	return nil
}

// SetupAPIKey helps setup API key authentication
func (m *AuthManager) SetupAPIKey(apiKey string) error {
	shadowDir := filepath.Join(m.homeDir, ".shadow")
	envPath := filepath.Join(shadowDir, ".env")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(shadowDir, 0755); err != nil {
		return fmt.Errorf("failed to create shadow directory: %w", err)
	}

	// Create .env file
	envContent := fmt.Sprintf("# Shadow Environment Variables\n# Generated: %s\n\nANTHROPIC_API_KEY=%s\n",
		time.Now().Format(time.RFC3339), apiKey)

	if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	return nil
}

// BackupCredentials creates a backup of current credentials
func (m *AuthManager) BackupCredentials() (string, error) {
	claudeCredsPath := filepath.Join(m.homeDir, ".claude", ".credentials.json")
	backupDir := filepath.Join(m.homeDir, ".shadow", "backups")

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("credentials_backup_%s.json", timestamp))

	// Read and copy credentials
	data, err := os.ReadFile(claudeCredsPath)
	if err != nil {
		return "", fmt.Errorf("failed to read credentials: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return "", fmt.Errorf("failed to write backup: %w", err)
	}

	return backupPath, nil
}
