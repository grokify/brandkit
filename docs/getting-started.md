# Getting Started

This guide covers installation and basic usage of BrandKit.

## Installation

### CLI Tool

Install the CLI using Go:

```bash
go install github.com/grokify/brandkit/cmd/svg@latest
```

The binary is named `brandkit`. Verify installation:

```bash
brandkit --version
```

### Go Library

Add BrandKit to your Go project:

```bash
go get github.com/grokify/brandkit
```

## Quick Start

### Creating White Icons

The most common workflow is converting a color icon to white on a transparent background:

```bash
brandkit white icon_orig.svg -o icon_white.svg
```

This command:

1. Removes any solid background
2. Converts all colors to white (`#ffffff`)
3. Centers the content in the viewBox
4. Verifies the output is pure vector
5. Scans for security threats

### Creating Color Icons

To create a centered icon preserving original colors:

```bash
brandkit color icon_orig.svg -o icon_color.svg
```

### Security Scanning

Scan SVG files for potential security threats:

```bash
# Scan a single file
brandkit security-scan icon.svg

# Scan a directory
brandkit security-scan brands/

# Generate JSON report
brandkit security-scan brands/ --report=security-report.json
```

### Verification

Verify that SVG files are pure vector (no embedded raster data):

```bash
# Single file
brandkit verify icon.svg

# Recursive directory scan
brandkit verify-all brands/
```

## Using as a Library

### Color Conversion

```go
package main

import (
    "fmt"
    "github.com/grokify/brandkit/svg/convert"
)

func main() {
    result, err := convert.SVG("input.svg", "output.svg", convert.Options{
        Color:            "ffffff",
        RemoveBackground: true,
        IncludeStroke:    true,
        PreserveMasks:    true,
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("Converted: %s\n", result.OutputPath)
}
```

### Security Scanning

```go
package main

import (
    "fmt"
    "github.com/grokify/brandkit/svg/security"
)

func main() {
    result, err := security.SVG("icon.svg")
    if err != nil {
        panic(err)
    }

    if result.IsSecure {
        fmt.Println("No threats detected")
    } else {
        for _, threat := range result.Threats {
            fmt.Printf("Threat: %s (%s)\n", threat.Description, threat.Type)
        }
    }
}
```

### Verification

```go
package main

import (
    "fmt"
    "github.com/grokify/brandkit/svg/verify"
)

func main() {
    result, err := verify.SVG("icon.svg")
    if err != nil {
        panic(err)
    }

    if result.IsSuccess() {
        fmt.Println("Pure vector SVG")
    } else {
        for _, issue := range result.Issues {
            fmt.Printf("Issue: %s\n", issue)
        }
    }
}
```

## Next Steps

- [CLI Reference](cli/index.md) — Learn all available commands
- [Library API](library/index.md) — Explore the Go packages
- [Security Guide](security/index.md) — Understand threat detection
