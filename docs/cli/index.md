# CLI Reference

BrandKit provides a command-line interface for SVG processing, analysis, and security scanning.

## Installation

```bash
go install github.com/grokify/brandkit/cmd/svg@latest
```

## Commands Overview

| Command | Description |
|---------|-------------|
| [`white`](white.md) | Create white icon on transparent background |
| [`color`](color.md) | Create centered color icon preserving original colors |
| [`convert`](convert.md) | Convert SVG colors with fine-grained control |
| [`process`](process.md) | Full pipeline: convert, center, verify |
| [`analyze`](analyze.md) | Analyze SVG geometry (centering, padding) |
| [`verify`](verify.md) | Verify SVG is pure vector |
| [`security-scan`](security-scan.md) | Scan for security threats |
| [`sanitize`](sanitize.md) | Remove security threats from SVG |

## Global Flags

```
-h, --help      Help for any command
-v, --version   Show version information
```

## Usage Pattern

```bash
brandkit <command> [arguments] [flags]
```

## Examples

```bash
# Create white icon
brandkit white icon.svg -o icon_white.svg

# Create color icon
brandkit color icon.svg -o icon_color.svg

# Security scan directory
brandkit security-scan brands/

# Verify all SVGs recursively
brandkit verify-all brands/

# Full processing pipeline
brandkit process icon.svg -o output.svg --color ffffff --center --strict
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (invalid input, processing failure, threats detected) |

## Environment

BrandKit reads no environment variables. All configuration is via command-line flags.
