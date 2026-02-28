# BrandKit

A Go toolkit for processing brand icons — analyzing, verifying, and converting SVG files. Use it as a CLI tool or as a library in your own Go programs.

[![Build Status](https://github.com/grokify/brandkit/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/grokify/brandkit/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/grokify/brandkit)](https://goreportcard.com/report/github.com/grokify/brandkit)
[![Docs](https://pkg.go.dev/badge/github.com/grokify/brandkit)](https://pkg.go.dev/github.com/grokify/brandkit)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/grokify/brandkit/blob/master/LICENSE)

## Features

- **SVG Color Conversion** — Convert icons to white, any hex color, or preserve original colors
- **Background Removal** — Remove solid backgrounds while preserving transparency
- **Geometry Analysis** — Analyze centering, padding, and viewBox optimization
- **Pure Vector Verification** — Detect embedded raster data, base64, and data URIs
- **Security Scanning** — Detect XSS threats, malicious scripts, and unsafe elements
- **Sanitization** — Remove security threats while preserving valid SVG content
- **Brand Asset Library** — 55+ brand icons with standardized variants

## Quick Start

### CLI Installation

```bash
go install github.com/grokify/brandkit/cmd/svg@latest
```

### Library Installation

```bash
go get github.com/grokify/brandkit
```

### Basic Usage

```bash
# Create white icon on transparent background
brandkit white icon_orig.svg -o icon_white.svg

# Security scan SVG files
brandkit security-scan brands/

# Verify pure vector (no embedded raster)
brandkit verify icon.svg
```

## Documentation

- [Getting Started](getting-started.md) — Installation and first steps
- [CLI Reference](cli/index.md) — Complete command reference
- [Library API](library/index.md) — Go package documentation
- [Security Guide](security/index.md) — Threat detection and sanitization
- [Brand Assets](brands.md) — Available brand icons

## License

MIT License — See [LICENSE](https://github.com/grokify/brandkit/blob/master/LICENSE) for details.
