# Security Guide

BrandKit includes comprehensive SVG security scanning and sanitization to protect against XSS attacks, session hijacking, and other malicious content.

## Overview

SVG files are XML documents that can contain executable content. While this enables powerful features like animation and interactivity, it also creates security risks when SVGs are:

- Uploaded by users
- Embedded in web pages
- Used as brand assets in applications

BrandKit's security features detect and remove these threats while preserving valid SVG content.

## Quick Start

### Scan for Threats

```bash
# Scan a single file
brandkit security-scan icon.svg

# Scan a directory
brandkit security-scan brands/

# Recursive scan with JSON report
brandkit security-scan-all brands/ --report=security-report.json
```

### Remove Threats

```bash
# Remove all threats
brandkit sanitize malicious.svg -o clean.svg

# Remove only critical threats
brandkit sanitize icon.svg -o clean.svg --remove-scripts --remove-event-handlers
```

### Integrated Pipeline

The `white` and `color` commands include security scanning by default:

```bash
# Fails if threats detected
brandkit white icon.svg -o icon_white.svg

# Warns but doesn't fail
brandkit white icon.svg -o icon_white.svg --insecure
```

## Threat Categories

| Category | Severity | Risk |
|----------|----------|------|
| [Scripts](threats.md#scripts) | Critical | XSS, session hijacking |
| [Event Handlers](threats.md#event-handlers) | Critical | XSS via user interaction |
| [External References](threats.md#external-references) | High | Data exfiltration, tracking |
| [XML Entities](threats.md#xml-entities) | High | XXE attacks, DoS |
| [Animation](threats.md#animation) | Medium | Delayed XSS, UI manipulation |
| [Links](threats.md#links) | Medium | Phishing, navigation hijacking |
| [Style Blocks](threats.md#style-blocks) | Low | CSS injection, UI manipulation |

## Scan Levels

### Strict Mode (Default)

Detects all 7 threat types. Use for maximum security:

```bash
brandkit security-scan icon.svg --strict
```

### Standard Mode

Detects only critical and high severity threats. Use when style blocks and animations are acceptable:

```bash
brandkit security-scan icon.svg --strict=false
```

## Documentation

- [Threat Types](threats.md) — Detailed threat descriptions
- [Scanning](scanning.md) — Scanning options and usage
- [Sanitization](sanitization.md) — Removing threats
- [Reports](reports.md) — JSON report format

## Library Usage

```go
import "github.com/grokify/brandkit/svg/security"

// Scan
result, _ := security.SVG("icon.svg")
if !result.IsSecure {
    for _, threat := range result.Threats {
        fmt.Printf("%s: %s\n", threat.Type.Severity(), threat.Description)
    }
}

// Sanitize
security.Sanitize("input.svg", "output.svg", security.SanitizeOptions{
    RemoveAll: true,
})
```

## CI Integration

Add security scanning to your CI pipeline:

```yaml
- name: Security scan SVG icons
  run: brandkit security-scan-all brands/
```

Or use Make:

```bash
make security-scan-all
```
