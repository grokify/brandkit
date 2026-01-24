# BrandKit

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

A Go toolkit for processing brand icons — analyzing, verifying, and converting SVG files. Use it as a CLI tool or as a library in your own Go programs.

## Installation

### CLI

```bash
go install github.com/grokify/brandkit/cmd/svg@latest
```

### Library

```bash
go get github.com/grokify/brandkit
```

## CLI Usage

```bash
# Convert icon colors to white on transparent background
brandkit white icon_orig.svg -o icon_white.svg

# Convert to a specific color
brandkit convert icon.svg -o output.svg --color ff5500

# Remove background and convert
brandkit convert icon.svg -o output.svg --color ffffff --remove-background --include-stroke

# Full pipeline: convert, center, verify
brandkit process icon.svg -o output.svg --color ffffff --center --strict

# Analyze SVG geometry (centering, padding)
brandkit analyze brands/

# Verify SVGs are pure vector (no embedded raster data)
brandkit verify brands/
```

## Library Usage

### Color Conversion

```go
import "github.com/grokify/brandkit/svg/convert"

result, err := convert.SVG("input.svg", "output.svg", convert.Options{
    Color:            "ffffff",
    RemoveBackground: true,
    IncludeStroke:    true,
    PreserveMasks:    true,
})
```

### Verification

```go
import "github.com/grokify/brandkit/svg/verify"

result, err := verify.SVG("icon.svg")
if result.IsSuccess() {
    fmt.Println("Pure vector SVG")
}
```

### Analysis

```go
import "github.com/grokify/brandkit/svg/analyze"

result, err := analyze.SVG("icon.svg")
fmt.Printf("Padding: L:%.1f%% R:%.1f%% T:%.1f%% B:%.1f%%\n",
    result.PaddingLeft, result.PaddingRight,
    result.PaddingTop, result.PaddingBottom)
```

### SVG Types

```go
import "github.com/grokify/brandkit/svg"

// Parse and work with SVG paths
commands := svg.ParsePath("M 0 0 L 100 0 L 100 100 Z")
bounds := svg.CalculatePathBounds("M 10 10 L 90 90")

// ViewBox parsing
vb, _ := svg.ParseViewBox("0 0 100 100")
```

## Packages

| Package | Description |
|---------|-------------|
| `svg` | Shared SVG types: `BoundingBox`, `ViewBox`, path parsing, file utilities |
| `svg/analyze` | Geometry analysis: centering, padding, viewBox suggestions |
| `svg/convert` | Color conversion, background removal, mask preservation |
| `svg/verify` | Pure vector validation, embedded data detection |

## Brand Assets

The `brands/` directory contains SVG icons for 36 brands with standardized variants:

- `icon_orig.svg` — Original source icon
- `icon_white.svg` — White on transparent background
- `icon_color.svg` — Color variant (select brands)

## Development

```bash
make build       # Build CLI binary
make test        # Run tests
make lint        # Run golangci-lint
make white       # Process all brands to white icons
make verify      # Verify all brand SVGs are pure vector
make analyze     # Analyze all brands for centering issues
make build-all   # Cross-platform builds
```

## License

MIT

 [build-status-svg]: https://github.com/grokify/brandkit/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/brandkit/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/brandkit/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/brandkit/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/brandkit
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/brandkit
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/brandkit
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/brandkit
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fbrandkit
 [loc-svg]: https://tokei.rs/b1/github/grokify/brandkit
 [repo-url]: https://github.com/grokify/brandkit
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/brandkit/blob/master/LICENSE
