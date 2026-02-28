# Library API

BrandKit provides Go packages for programmatic SVG processing.

## Installation

```bash
go get github.com/grokify/brandkit
```

## Packages

| Package | Import Path | Description |
|---------|-------------|-------------|
| [svg](svg.md) | `github.com/grokify/brandkit/svg` | Core types: BoundingBox, ViewBox, path parsing |
| [analyze](analyze.md) | `github.com/grokify/brandkit/svg/analyze` | Geometry analysis: centering, padding |
| [convert](convert.md) | `github.com/grokify/brandkit/svg/convert` | Color conversion, background removal |
| [verify](verify.md) | `github.com/grokify/brandkit/svg/verify` | Pure vector validation |
| [security](security.md) | `github.com/grokify/brandkit/svg/security` | Security scanning and sanitization |

## Quick Examples

### Color Conversion

```go
import "github.com/grokify/brandkit/svg/convert"

result, err := convert.SVG("input.svg", "output.svg", convert.Options{
    Color:            "ffffff",
    RemoveBackground: true,
    IncludeStroke:    true,
})
```

### Security Scanning

```go
import "github.com/grokify/brandkit/svg/security"

result, err := security.SVG("icon.svg")
if !result.IsSecure {
    for _, threat := range result.Threats {
        fmt.Printf("Threat: %s\n", threat.Description)
    }
}
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

## High-Level Processing

For common workflows, use the top-level processing functions:

```go
import "github.com/grokify/brandkit"

// Create white icon
result, err := brandkit.ProcessWhite("input.svg", "output.svg")

// Create color icon (preserving original colors)
result, err := brandkit.ProcessColor("input.svg", "output.svg")
```

These functions handle the full pipeline:

1. Remove background
2. Convert colors (white only for ProcessWhite)
3. Center content
4. Verify pure vector
5. Security scan

## Error Handling

All functions return errors following Go conventions:

```go
result, err := security.SVG("icon.svg")
if err != nil {
    // File read error, parse error, etc.
    log.Fatal(err)
}

// Result contains operation-specific status
if !result.IsSecure {
    // Handle security threats
}
```

## Documentation

Full API documentation available at [pkg.go.dev](https://pkg.go.dev/github.com/grokify/brandkit).
