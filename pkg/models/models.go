package models

import "time"

// ScanConfig represents scan configuration
type ScanConfig struct {
	Target     string
	Profile    string
	AIAnalysis bool
	Threads    int
	Modules    []string
}

// ScanResult represents the output of a security scan
type ScanResult struct {
	ID        string        `json:"id"`
	Target    string        `json:"target"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Status    string        `json:"status"`
	Findings  []Finding     `json:"findings"`
	Metadata  ScanMetadata  `json:"metadata"`
}

// Finding represents a security finding
type Finding struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Severity    string            `json:"severity"` // critical, high, medium, low, info
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Evidence    string            `json:"evidence"`
	Location    string            `json:"location"`
	CVE         string            `json:"cve,omitempty"`
	CVSS        float64           `json:"cvss,omitempty"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	Timestamp   time.Time         `json:"timestamp"`
}

// ScanMetadata contains metadata about the scan
type ScanMetadata struct {
	Version    string    `json:"version"`
	Profile    string    `json:"profile"`
	Modules    []string  `json:"modules"`
	Threads    int       `json:"threads"`
	AIAnalyzed bool      `json:"ai_analyzed"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

// SubdomainResult represents discovered subdomains
type SubdomainResult struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
	Count      int      `json:"count"`
	Timestamp  time.Time `json:"timestamp"`
}

// PortScanResult represents port scan findings
type PortScanResult struct {
	Target    string       `json:"target"`
	Ports     []OpenPort   `json:"ports"`
	Count     int          `json:"count"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time    `json:"timestamp"`
}

// OpenPort represents an open port
type OpenPort struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"` // tcp, udp
	Service  string `json:"service"`
	Version  string `json:"version,omitempty"`
	State    string `json:"state"`
}

// SSLResult represents SSL/TLS analysis
type SSLResult struct {
	Target       string        `json:"target"`
	Valid        bool          `json:"valid"`
	Issuer       string        `json:"issuer"`
	Subject      string        `json:"subject"`
	NotBefore    time.Time     `json:"not_before"`
	NotAfter     time.Time     `json:"not_after"`
	DaysToExpiry int           `json:"days_to_expiry"`
	Version      string        `json:"version"`
	Cipher       string        `json:"cipher"`
	Issues       []string      `json:"issues"`
	Grade        string        `json:"grade"` // A+, A, B, C, D, F
}

// AIAnalysis represents AI-powered analysis results
type AIAnalysis struct {
	ScanID          string             `json:"scan_id"`
	Summary         string             `json:"summary"`
	CriticalIssues  []string           `json:"critical_issues"`
	Recommendations []Recommendation   `json:"recommendations"`
	AttackChains    []AttackChain      `json:"attack_chains"`
	RiskScore       int                `json:"risk_score"` // 0-100
	Timestamp       time.Time          `json:"timestamp"`
}

// Recommendation represents an AI-generated recommendation
type Recommendation struct {
	Priority    string `json:"priority"` // critical, high, medium, low
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Effort      string `json:"effort"` // low, medium, high
	Steps       []string `json:"steps"`
}

// AttackChain represents a potential attack path
type AttackChain struct {
	ID          string   `json:"id"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	Impact      string   `json:"impact"`
	Likelihood  string   `json:"likelihood"`
}
