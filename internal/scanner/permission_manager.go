package scanner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PermissionManager handles permission requests and sudo access
type PermissionManager struct {
	sudoAvailable bool
	sudoTested    bool
	userApproved  map[string]bool // Track which commands user approved
}

// NewPermissionManager creates a new permission manager
func NewPermissionManager() *PermissionManager {
	return &PermissionManager{
		userApproved: make(map[string]bool),
	}
}

// CheckSudoAvailable tests if sudo is available
func (pm *PermissionManager) CheckSudoAvailable() bool {
	if pm.sudoTested {
		return pm.sudoAvailable
	}

	// Check if sudo command exists
	_, err := exec.LookPath("sudo")
	if err != nil {
		pm.sudoAvailable = false
		pm.sudoTested = true
		return false
	}

	// Test if we can run sudo (might require password)
	cmd := exec.Command("sudo", "-n", "true")
	err = cmd.Run()
	pm.sudoAvailable = (err == nil)
	pm.sudoTested = true

	return pm.sudoAvailable
}

// RequestRootPermission asks user for permission to run privileged commands
func (pm *PermissionManager) RequestRootPermission(tool string, purpose string, command string) (bool, error) {
	// Check if already approved
	cacheKey := fmt.Sprintf("%s:%s", tool, command)
	if approved, exists := pm.userApproved[cacheKey]; exists {
		return approved, nil
	}

	fmt.Println("\n🔐 Root Permission Request")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\n📋 Tool: %s\n", tool)
	fmt.Printf("🎯 Purpose: %s\n", purpose)
	fmt.Printf("💻 Command: %s\n", command)

	if !pm.CheckSudoAvailable() {
		fmt.Println("\n⚠️  sudo is not available or not configured")
		fmt.Println("💡 Options:")
		fmt.Println("   1. Configure sudo access")
		fmt.Println("   2. Run Shadow as root (not recommended)")
		fmt.Println("   3. Skip this scan and use alternatives")
		return false, fmt.Errorf("sudo not available")
	}

	fmt.Println("\n⚠️  This command requires elevated privileges (root/sudo)")
	fmt.Println("🔒 Shadow will ONLY run the specific command shown above")
	fmt.Println("📊 This is needed for comprehensive security scanning")

	fmt.Print("\nAllow this command? (yes/no/always): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.ToLower(strings.TrimSpace(response))

	switch response {
	case "yes", "y":
		pm.userApproved[cacheKey] = true
		return true, nil
	case "always", "a":
		// Approve all future requests for this tool
		pm.userApproved[tool+":*"] = true
		pm.userApproved[cacheKey] = true
		return true, nil
	case "no", "n":
		pm.userApproved[cacheKey] = false
		return false, nil
	default:
		fmt.Println("⚠️  Invalid response, treating as 'no'")
		pm.userApproved[cacheKey] = false
		return false, nil
	}
}

// RunWithSudo executes a command with sudo after getting permission
func (pm *PermissionManager) RunWithSudo(tool string, purpose string, args ...string) ([]byte, error) {
	command := fmt.Sprintf("sudo %s %s", tool, strings.Join(args, " "))

	// Request permission
	approved, err := pm.RequestRootPermission(tool, purpose, command)
	if err != nil {
		return nil, err
	}

	if !approved {
		return nil, fmt.Errorf("user denied permission")
	}

	// Execute with sudo
	fmt.Printf("\n🔧 Executing: %s\n", command)

	cmdArgs := append([]string{tool}, args...)
	cmd := exec.Command("sudo", cmdArgs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}

	return output, nil
}

// RunWithFallback tries to run with sudo, falls back to non-root version
func (pm *PermissionManager) RunWithFallback(
	tool string,
	purpose string,
	rootArgs []string,
	fallbackArgs []string,
) ([]byte, bool, error) {
	// Try with sudo first
	if pm.CheckSudoAvailable() {
		command := fmt.Sprintf("sudo %s %s", tool, strings.Join(rootArgs, " "))
		approved, err := pm.RequestRootPermission(tool, purpose, command)

		if err == nil && approved {
			fmt.Printf("\n🔧 Executing privileged scan: %s\n", command)
			cmdArgs := append([]string{tool}, rootArgs...)
			cmd := exec.Command("sudo", cmdArgs...)
			output, err := cmd.CombinedOutput()

			if err == nil {
				fmt.Println("✓ Privileged scan completed")
				return output, true, nil
			}

			fmt.Printf("⚠️  Privileged scan failed: %v\n", err)
			fmt.Println("💡 Falling back to non-privileged scan...")
		}
	}

	// Fallback to non-root version
	fmt.Printf("🔧 Running non-privileged scan: %s %s\n", tool, strings.Join(fallbackArgs, " "))

	cmd := exec.Command(tool, fallbackArgs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, false, fmt.Errorf("fallback scan failed: %w", err)
	}

	fmt.Println("✓ Non-privileged scan completed")
	return output, false, nil
}

// ShowCapabilityInfo displays information about setcap as an alternative to sudo
func (pm *PermissionManager) ShowCapabilityInfo(tool string) {
	fmt.Println("\n💡 Alternative: Use Linux Capabilities Instead of sudo")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	switch tool {
	case "nmap":
		fmt.Println("\n📝 To allow nmap without sudo:")
		fmt.Println("   sudo setcap cap_net_raw,cap_net_admin,cap_net_bind_service+eip /usr/bin/nmap")
		fmt.Println("\n✅ Benefits:")
		fmt.Println("   • More secure than sudo")
		fmt.Println("   • No password prompts")
		fmt.Println("   • Granular permissions")
		fmt.Println("\n⚠️  Note: You'll need sudo once to set capabilities")

	default:
		fmt.Println("\n📝 Check if %s supports Linux capabilities", tool)
		fmt.Println("   man capabilities")
	}

	fmt.Println()
}

// SuggestSudoersEntry suggests a sudoers configuration for the tool
func (pm *PermissionManager) SuggestSudoersEntry(tool string) {
	fmt.Println("\n💡 Persistent sudo Configuration")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\n📝 To avoid repeated password prompts:")

	username := os.Getenv("USER")
	if username == "" {
		username = "yourusername"
	}

	switch tool {
	case "nmap":
		fmt.Printf("\n   echo '%s ALL=(ALL) NOPASSWD: /usr/bin/nmap' | sudo tee /etc/sudoers.d/shadow-nmap\n", username)
		fmt.Println("   sudo chmod 440 /etc/sudoers.d/shadow-nmap")

	default:
		toolPath, _ := exec.LookPath(tool)
		if toolPath == "" {
			toolPath = "/usr/bin/" + tool
		}
		fmt.Printf("\n   echo '%s ALL=(ALL) NOPASSWD: %s' | sudo tee /etc/sudoers.d/shadow-%s\n",
			username, toolPath, tool)
		fmt.Printf("   sudo chmod 440 /etc/sudoers.d/shadow-%s\n", tool)
	}

	fmt.Println("\n⚠️  Security Note:")
	fmt.Println("   • Only allow specific tools, not ALL commands")
	fmt.Println("   • Limit to absolute paths")
	fmt.Println("   • Review sudoers entries regularly")
	fmt.Println()
}

// GetApprovalSummary returns summary of approved commands
func (pm *PermissionManager) GetApprovalSummary() {
	if len(pm.userApproved) == 0 {
		return
	}

	fmt.Println("\n📊 Permission Summary")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	approvedCount := 0
	deniedCount := 0

	for _, approved := range pm.userApproved {
		if approved {
			approvedCount++
		} else {
			deniedCount++
		}
	}

	fmt.Printf("✅ Approved: %d commands\n", approvedCount)
	fmt.Printf("❌ Denied: %d commands\n", deniedCount)
	fmt.Println()
}
