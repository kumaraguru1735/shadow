# Shadow Architecture

## Overview

Shadow is designed as a modular, high-performance security reconnaissance platform built in Go. The architecture emphasizes:

- **Performance**: Go's concurrency model for parallel scanning
- **Modularity**: Plugin-based architecture for extensibility
- **Intelligence**: Claude AI integration for smart analysis
- **Scalability**: Distributed scanning capabilities

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Interface                         │
│                    (cmd/shadow/main.go)                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    Core Scanner Engine                       │
│                 (internal/scanner/)                          │
│  ┌────────────┐  ┌────────────┐  ┌─────────────────────┐  │
│  │  Module    │  │  Module    │  │   Module            │  │
│  │  Manager   │  │  Executor  │  │   Registry          │  │
│  └────────────┘  └────────────┘  └─────────────────────┘  │
└──────────────┬────────────────┬──────────────┬─────────────┘
               │                │              │
    ┌──────────▼─────┐  ┌──────▼──────┐  ┌───▼──────────┐
    │   Scanning     │  │    AI       │  │   Database   │
    │   Modules      │  │   Engine    │  │              │
    │                │  │             │  │              │
    │ • Subdomain    │  │ • Analysis  │  │ • Results    │
    │ • Port Scan    │  │ • Query     │  │ • History    │
    │ • Web Crawl    │  │ • Report    │  │ • Cache      │
    │ • SSL Check    │  │ • Compare   │  │              │
    │ • Headers      │  │             │  │              │
    └────────────────┘  └─────────────┘  └──────────────┘
```

## Core Components

### 1. CLI Layer (`cmd/shadow/`)

**Responsibility**: Command-line interface and user interaction

```go
main.go
├── Command definitions (scan, analyze, report, etc.)
├── Flag parsing
├── User input validation
└── Authorization checks
```

**Key Features**:
- Cobra-based command structure
- Interactive permission confirmation
- Progress indicators
- Colored output

### 2. Scanner Engine (`internal/scanner/`)

**Responsibility**: Core scanning orchestration

```go
scanner/
├── scanner.go        // Main scanner logic
├── module.go         // Module interface
├── executor.go       // Parallel execution
├── profile.go        // Scan profiles
└── coordinator.go    // Module coordination
```

**Key Features**:
- Goroutine-based parallelism
- Profile-based module loading
- Error handling and retry logic
- Resource management

### 3. Scanning Modules (`internal/modules/`)

**Responsibility**: Individual security checks

```go
modules/
├── subdomain/
│   ├── dns.go           // DNS enumeration
│   ├── bruteforce.go    // Subdomain brute force
│   └── ct_logs.go       // Certificate transparency
├── portscan/
│   ├── tcp.go           // TCP scanning
│   ├── udp.go           // UDP scanning
│   └── service.go       // Service detection
├── web/
│   ├── crawler.go       // Web crawling
│   ├── headers.go       // Security headers
│   └── technology.go    // Stack fingerprinting
└── ssl/
    ├── certificate.go   // Certificate validation
    └── vulnerabilities.go // SSL/TLS vulns
```

**Module Interface**:
```go
type Module interface {
    Name() string
    Description() string
    Run(target string, opts Options) ([]Finding, error)
    Cleanup() error
}
```

### 4. AI Engine (`internal/ai/`)

**Responsibility**: Claude AI integration

```go
ai/
├── claude.go          // Claude API client
├── analyzer.go        // Vulnerability analysis
├── prompts.go         // Prompt templates
├── cache.go           // Response caching
└── optimizer.go       // API call optimization
```

**Key Features**:
- Intelligent prompt construction
- Response caching
- Batch processing
- Token management
- Error handling with retry

### 5. Data Layer (`internal/database/`)

**Responsibility**: Data persistence

```go
database/
├── database.go        // Database interface
├── sqlite.go          // SQLite implementation
├── postgres.go        // PostgreSQL implementation
├── models.go          // Database models
└── migrations.go      // Schema migrations
```

**Schema**:
```sql
scans (id, target, start_time, end_time, status, metadata)
findings (id, scan_id, type, severity, title, description)
ai_analysis (id, scan_id, summary, risk_score, timestamp)
```

### 6. Reporting (`internal/report/`)

**Responsibility**: Report generation

```go
report/
├── generator.go       // Report generator
├── html.go            // HTML reports
├── pdf.go             // PDF reports
├── json.go            // JSON exports
└── templates/         // Report templates
```

## Data Flow

### 1. Scan Execution Flow

```
User Command
    ↓
Permission Check
    ↓
Load Configuration
    ↓
Initialize Scanner
    ↓
Load Modules (based on profile)
    ↓
Execute Modules (parallel)
    │
    ├→ Module 1 → Findings → Database
    ├→ Module 2 → Findings → Database
    ├→ Module 3 → Findings → Database
    └→ Module N → Findings → Database
    ↓
Aggregate Results
    ↓
AI Analysis (if enabled)
    ↓
Generate Report
    ↓
Save & Display Results
```

### 2. AI Analysis Flow

```
Scan Results
    ↓
Load from Database
    ↓
Build Analysis Prompt
    ↓
Claude API Call
    ↓
Parse Response
    ↓
Extract:
    - Summary
    - Risk Score
    - Recommendations
    - Attack Chains
    ↓
Save to Database
    ↓
Return Analysis
```

## Concurrency Model

Shadow leverages Go's goroutines for efficient parallel processing:

```go
// Worker pool pattern
type WorkerPool struct {
    workers int
    jobs    chan Job
    results chan Result
}

// Execute jobs concurrently
func (wp *WorkerPool) Execute(jobs []Job) []Result {
    // Spawn workers
    for i := 0; i < wp.workers; i++ {
        go wp.worker()
    }

    // Distribute jobs
    for _, job := range jobs {
        wp.jobs <- job
    }

    // Collect results
    results := []Result{}
    for range jobs {
        results = append(results, <-wp.results)
    }

    return results
}
```

## Module System

### Module Lifecycle

1. **Registration**: Modules register with the registry
2. **Initialization**: Scanner initializes required modules
3. **Execution**: Modules run concurrently
4. **Results**: Modules return findings
5. **Cleanup**: Resources are released

### Adding a New Module

```go
// 1. Implement Module interface
type MyModule struct {
    config Config
}

func (m *MyModule) Name() string {
    return "My Custom Module"
}

func (m *MyModule) Run(target string, opts Options) ([]Finding, error) {
    // Scanning logic
    return findings, nil
}

// 2. Register module
func init() {
    registry.Register(&MyModule{})
}
```

## Performance Optimizations

### 1. Connection Pooling
- Reuse HTTP connections
- DNS caching
- TCP connection pooling

### 2. Resource Management
- Configurable worker pools
- Memory limits
- CPU throttling

### 3. Caching
- AI response caching (similar queries)
- DNS lookup caching
- HTTP response caching

### 4. Batch Processing
- Batch AI API calls
- Batch database writes
- Batch HTTP requests

## Security Considerations

### 1. Input Validation
- Target URL validation
- Scope checking
- Configuration validation

### 2. Authorization
- Permission confirmation required
- Scope enforcement
- Audit logging

### 3. Rate Limiting
- Respectful scanning rates
- Configurable delays
- Automatic throttling

### 4. Data Protection
- API key encryption
- Secure credential storage
- Audit log encryption

## Scalability

### Horizontal Scaling

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Worker 1 │     │ Worker 2 │     │ Worker N │
└────┬─────┘     └────┬─────┘     └────┬─────┘
     │                │                │
     └────────────────┴────────────────┘
                      │
            ┌─────────▼─────────┐
            │  Coordinator      │
            │  (distributes     │
            │   work)           │
            └─────────┬─────────┘
                      │
            ┌─────────▼─────────┐
            │  Shared Database  │
            └───────────────────┘
```

### Resource Management

```go
type ResourceManager struct {
    maxMemory    int64
    maxCPU       float64
    maxGoroutines int
}

func (rm *ResourceManager) CanAllocate() bool {
    currentMemory := getCurrentMemoryUsage()
    currentCPU := getCurrentCPUUsage()

    return currentMemory < rm.maxMemory &&
           currentCPU < rm.maxCPU
}
```

## Testing Strategy

### Unit Tests
- Individual module testing
- Mock interfaces
- Table-driven tests

### Integration Tests
- Module interaction testing
- Database integration
- AI API integration

### End-to-End Tests
- Full scan workflows
- Report generation
- Error scenarios

## Future Enhancements

1. **Plugin System**: External plugin loading
2. **Web UI**: Browser-based dashboard
3. **API Server**: RESTful API for integrations
4. **Distributed Scanning**: Multi-node scanning
5. **ML Models**: Custom vulnerability detection models
6. **Real-time Monitoring**: Continuous scanning

---

This architecture provides a solid foundation for Shadow while maintaining flexibility for future enhancements.
