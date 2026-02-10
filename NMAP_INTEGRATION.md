# Nmap Integration & Root Permissions

## Status

Nmap integration is **planned** but not yet implemented in Shadow.

## Root Permission Handling

When implemented, Shadow will handle nmap's root requirements using one of these strategies:

### Option 1: sudo with NOPASSWD (Recommended for Security Tools)

```bash
# Add to /etc/sudoers.d/shadow
yourusername ALL=(ALL) NOPASSWD: /usr/bin/nmap
```

Shadow would then run:
```bash
sudo nmap -sS target.com
```

### Option 2: setcap (Linux Only)

Grant nmap raw socket capabilities:
```bash
sudo setcap cap_net_raw,cap_net_admin,cap_net_bind_service+eip /usr/bin/nmap
```

Nmap can then run without sudo:
```bash
nmap -sS target.com  # Works without sudo
```

### Option 3: Non-Privileged Scans

Use TCP connect scans that don't require root:
```bash
nmap -sT target.com  # No root needed
```

**Trade-off**: Slower and more detectable than SYN scans.

## Planned Implementation

```go
// internal/scanner/nmap.go (future)

type NmapScanner struct {
    useSudo     bool
    useSetcap   bool
    scanType    string  // "SYN", "Connect", "UDP"
}

func (n *NmapScanner) Scan(target string) error {
    var cmd *exec.Cmd

    if n.useSudo {
        cmd = exec.Command("sudo", "nmap", "-sS", target)
    } else if n.scanType == "Connect" {
        cmd = exec.Command("nmap", "-sT", target)
    } else {
        // Try setcap first, fallback to TCP connect
        cmd = exec.Command("nmap", "-sS", target)
    }

    output, err := cmd.CombinedOutput()
    if err != nil {
        // Fallback logic
    }

    return parseNmapOutput(output)
}
```

## User Experience

When nmap is integrated, Shadow will:

1. **Detect available permissions**:
   ```
   🔍 Checking nmap capabilities...
   ✓ nmap found at /usr/bin/nmap
   ⚠️  Root access not available
   ℹ️  Will use TCP connect scans (slower but no root needed)
   ```

2. **Offer setup guidance**:
   ```
   💡 For faster SYN scans, run:
      sudo setcap cap_net_raw+eip /usr/bin/nmap

   Or add to sudoers:
      echo "$(whoami) ALL=(ALL) NOPASSWD: /usr/bin/nmap" | sudo tee /etc/sudoers.d/shadow
   ```

3. **Graceful degradation**:
   - Try SYN scan with setcap
   - Fall back to TCP connect if no permissions
   - Inform user of scan type being used

## Security Considerations

### Why Root is Needed

Nmap requires raw socket access for:
- **SYN scanning** (-sS): Fastest, stealthiest
- **UDP scanning** (-sU): Service discovery
- **OS detection** (-O): Fingerprinting

### Security Best Practices

1. **Use setcap over sudo**: More granular permissions
2. **Limit nmap sudoers entry**: Only allow nmap, not all commands
3. **Audit nmap usage**: Log all privileged scans
4. **Validate targets**: Prevent scanning unauthorized hosts

## Current Workaround

Since nmap isn't integrated yet, Shadow uses:
- **Go-based port scanner**: No root needed
- **HTTP headers analysis**: Standard permissions
- **SSL/TLS checks**: Standard sockets

## Timeline

Nmap integration is planned for **v0.2.0** with:
- Auto-detection of permissions
- Graceful fallback strategies
- Clear user guidance
- Security-first implementation

## References

- [Nmap Privileges Documentation](https://nmap.org/book/inst-linux.html)
- [Linux Capabilities](https://man7.org/linux/man-pages/man7/capabilities.7.html)
- [Nmap Scan Types](https://nmap.org/book/scan-methods.html)

---

**Note**: Shadow is designed to work without root where possible. Nmap is optional and will have fallback strategies.
