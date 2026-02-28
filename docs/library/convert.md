# svg/convert Package

```go
import "github.com/grokify/brandkit/svg/convert"
```

Provides SVG color conversion and background removal.

## Types

### Options

Configures the color conversion behavior.

```go
type Options struct {
    Color            string // Target color (hex or named)
    IncludeStroke    bool   // Also convert stroke colors
    PreserveMasks    bool   // Don't modify colors in mask/clipPath
    RemoveBackground bool   // Remove background rect/circle elements
}
```

| Field | Default | Description |
|-------|---------|-------------|
| `Color` | "" | Target color (empty = no conversion) |
| `IncludeStroke` | false | Convert stroke colors too |
| `PreserveMasks` | false | Preserve mask/clipPath colors |
| `RemoveBackground` | false | Remove full-bleed backgrounds |

### Result

Contains the result of a color conversion.

```go
type Result struct {
    InputPath         string
    OutputPath        string
    OriginalColor     string
    TargetColor       string
    Converted         bool
    BackgroundRemoved bool
    Error             error
}
```

## Functions

### SVG

Converts colors in an SVG file.

```go
func SVG(inputPath, outputPath string, opts Options) (*Result, error)
```

**Example:**

```go
result, err := convert.SVG("input.svg", "output.svg", convert.Options{
    Color:            "ffffff",
    RemoveBackground: true,
    IncludeStroke:    true,
    PreserveMasks:    true,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Converted to %s\n", result.TargetColor)
if result.BackgroundRemoved {
    fmt.Println("Background removed")
}
```

### NormalizeColor

Normalizes a color input to standard #RRGGBB format.

```go
func NormalizeColor(color string) (string, error)
```

**Accepts:**

| Input | Output |
|-------|--------|
| `"ffffff"` | `"#ffffff"` |
| `"#ffffff"` | `"#ffffff"` |
| `"fff"` | `"#ffffff"` |
| `"#fff"` | `"#ffffff"` |
| `"white"` | `"#ffffff"` |
| `"BLACK"` | `"#000000"` |

**Named Colors:**

- `white`, `black`, `red`, `green`, `blue`
- `yellow`, `cyan`, `magenta`
- `gray`, `grey`
- `transparent` (returns "none")

**Example:**

```go
color, err := convert.NormalizeColor("fff")
if err != nil {
    log.Fatal(err)
}
fmt.Println(color) // Output: #ffffff
```

## Background Removal

When `RemoveBackground: true`, the following elements are removed:

| Element | Condition |
|---------|-----------|
| `<rect>` | Covers entire viewBox (within 1% tolerance) |
| `<circle>` | Centered and radius matches half viewBox |
| `<path>` | Draws rectangle covering entire viewBox |

## Mask Preservation

When `PreserveMasks: true`:

- Colors inside `<mask>` elements are not converted
- Colors inside `<clipPath>` elements are not converted
- This preserves visual appearance of masked content

## Common Patterns

### Convert to White

```go
convert.SVG("input.svg", "output.svg", convert.Options{
    Color:            "ffffff",
    IncludeStroke:    true,
    RemoveBackground: true,
    PreserveMasks:    true,
})
```

### Remove Background Only

```go
convert.SVG("input.svg", "output.svg", convert.Options{
    RemoveBackground: true,
})
```

### Convert to Brand Color

```go
convert.SVG("input.svg", "output.svg", convert.Options{
    Color: "ff5500", // Orange
})
```
