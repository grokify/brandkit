# Security Scanning

BrandKit provides multiple ways to scan SVG files for security threats.

## CLI Commands

### security-scan

Scan a single file or directory (non-recursive).

```bash
# Single file
brandkit security-scan icon.svg

# Directory
brandkit security-scan brands/anthropic/
```

### security-scan-all

Recursively scan all SVG files in a directory tree. Designed for CI pipelines.

```bash
brandkit security-scan-all brands/
```

## Command Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--strict` | true | Detect all threats (false = critical/high only) |
| `--report` | "" | Output JSON report file path |
| `--project` | "brandkit" | Project name for report |
| `--version` | CLI version | Version for report |

## Scan Levels

### Strict Mode (Default)

Detects all 7 threat types:

- Scripts (critical)
- Event handlers (critical)
- External references (high)
- XML entities (high)
- Animation (medium)
- Links (medium)
- Style blocks (low)

```bash
brandkit security-scan icon.svg --strict
```

### Standard Mode

Detects only critical and high severity threats:

- Scripts (critical)
- Event handlers (critical)
- External references (high)
- XML entities (high)

```bash
brandkit security-scan icon.svg --strict=false
```

Use standard mode when:

- Style blocks are acceptable (CSS classes for styling)
- Animations are intentional
- Links are acceptable

## Output

### Console Output

```
Scanning brands/...

✗ brands/malicious/icon.svg
  CRITICAL: script element
  CRITICAL: onclick handler

✓ brands/react/icon.svg: Secure

Summary: 1 file(s) with threats, 1 secure
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | No threats detected |
| 1 | Threats detected |

## JSON Reports

Generate detailed JSON reports for CI integration:

```bash
brandkit security-scan-all brands/ --report=security-report.json
```

See [Reports](reports.md) for format details.

## Library Usage

### Scan Single File

```go
import "github.com/grokify/brandkit/svg/security"

result, err := security.SVG("icon.svg")
if err != nil {
    log.Fatal(err)
}

if result.IsSecure {
    fmt.Println("No threats detected")
} else {
    for _, threat := range result.Threats {
        fmt.Printf("[%s] %s: %s\n",
            threat.Type.Severity(),
            threat.Description,
            threat.Match)
    }
}
```

### Scan with Level

```go
// Standard mode (critical/high only)
result, err := security.SVGWithLevel("icon.svg", security.ScanLevelStandard)

// Strict mode (all threats)
result, err := security.SVGWithLevel("icon.svg", security.ScanLevelStrict)
```

### Scan Directory

```go
// Non-recursive
results, err := security.Directory("brands/anthropic/")

// Recursive
results, err := security.DirectoryRecursive("brands/")

for _, r := range results {
    if !r.IsSecure {
        fmt.Printf("%s: %d threats\n", r.FilePath, len(r.Threats))
    }
}
```

### Scan Content String

```go
content := `<svg><script>alert(1)</script></svg>`

result := &security.Result{
    Threats:      []security.Threat{},
    ThreatCounts: make(map[security.ThreatType]int),
}

result = security.ScanContent(content, result)
fmt.Printf("Found %d threats\n", len(result.Threats))
```

## CI Integration

### GitHub Actions

```yaml
- name: Security scan SVG icons
  run: |
    brandkit security-scan-all brands/ --report=security-report.json

- name: Upload security report
  uses: actions/upload-artifact@v3
  if: always()
  with:
    name: security-report
    path: security-report.json
```

### Makefile

```makefile
security-scan-all: build
	$(BUILD_DIR)/$(BINARY_NAME) security-scan-all brands/
```

Then in CI:

```yaml
- name: Security scan
  run: make security-scan-all
```
