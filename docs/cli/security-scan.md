# brandkit security-scan

Scan SVG files for security threats.

## Synopsis

```bash
brandkit security-scan [path] [flags]
brandkit security-scan-all [path] [flags]
```

## Description

Scan SVG files for security threats that could enable XSS attacks, session hijacking, or other malicious behavior. Detects:

| Severity | Threats |
|----------|---------|
| **Critical** | Script elements, `javascript:` URIs, event handlers |
| **High** | External references, foreignObject, XML entities |
| **Medium** | Animation elements, anchor links (strict mode) |
| **Low** | Style blocks (strict mode) |

## Commands

### security-scan

Scan a single file or directory (non-recursive):

```bash
brandkit security-scan icon.svg
brandkit security-scan brands/anthropic/
```

### security-scan-all

Recursively scan all SVG files. Designed for CI pipelines:

```bash
brandkit security-scan-all brands/
```

## Flags

| Flag | Description |
|------|-------------|
| `--strict` | Detect all threats including style blocks and animations (default: true) |
| `--report` | Output JSON report file path |
| `--project` | Project name for report (default: brandkit) |
| `--version` | Version for report (default: CLI version) |
| `-h, --help` | Help for security-scan |

## Examples

Scan a single file:

```bash
brandkit security-scan icon.svg
```

Scan a directory:

```bash
brandkit security-scan brands/
```

Generate JSON report:

```bash
brandkit security-scan brands/ --report=security-report.json
```

Standard mode (critical/high only):

```bash
brandkit security-scan brands/ --strict=false
```

## Scan Levels

### Strict Mode (default)

Detects all 7 threat types including low-severity issues:

- Style blocks (often benign but can contain malicious CSS)
- Animation elements (can trigger delayed attacks)
- Anchor links (unnecessary for static images)

### Standard Mode

Detects only critical and high severity threats:

- Scripts and JavaScript URIs
- Event handlers
- External references
- XML entities (XXE prevention)

## Output

### Console Output

```
Scanning brands/...

✗ brands/malicious/icon.svg
  CRITICAL: script element at line 5
  CRITICAL: onclick handler at line 12

✓ brands/react/icon.svg: Secure

Summary: 1 file(s) with threats, 1 secure
```

### JSON Report

The `--report` flag generates a JSON report following the multi-agent-spec team-report format:

```json
{
  "$schema": "...",
  "title": "SVG SECURITY SCAN REPORT",
  "project": "brandkit",
  "version": "0.4.0",
  "status": "NO-GO",
  "teams": [
    {
      "id": "script-detection",
      "name": "Script Detection",
      "status": "NO-GO",
      "tasks": [...]
    }
  ]
}
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | No threats detected (or warnings only in standard mode) |
| 1 | Threats detected |

## CI Integration

Add to your CI pipeline:

```yaml
- name: Security scan SVG icons
  run: brandkit security-scan-all brands/
```

Or using Make:

```bash
make security-scan-all
```

## See Also

- [sanitize](sanitize.md) — Remove security threats
- [Security Guide](../security/index.md) — Full security documentation
- [Threat Types](../security/threats.md) — Detailed threat descriptions
