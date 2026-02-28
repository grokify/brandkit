# svg/verify Package

```go
import "github.com/grokify/brandkit/svg/verify"
```

Validates SVG files are pure vector images without embedded binary data.

## Types

### Result

Contains the result of validating an SVG file.

```go
type Result struct {
    FilePath        string
    IsValid         bool
    IsPureVector    bool
    HasEmbeddedData bool
    VectorElements  []string
    Errors          []string
}
```

| Field | Description |
|-------|-------------|
| `FilePath` | Path to the validated file |
| `IsValid` | True if file is valid XML/SVG |
| `IsPureVector` | True if no embedded binary data detected |
| `HasEmbeddedData` | True if base64 or data URIs found |
| `VectorElements` | List of vector element counts (e.g., "path:5") |
| `Errors` | List of validation errors |

### Methods

#### IsSuccess

Returns true if the result indicates a valid pure vector SVG.

```go
func (r *Result) IsSuccess() bool
```

Equivalent to `r.IsValid && r.IsPureVector`.

## Functions

### SVG

Validates a single SVG file.

```go
func SVG(filePath string) (*Result, error)
```

**Example:**

```go
result, err := verify.SVG("icon.svg")
if err != nil {
    log.Fatal(err)
}

if result.IsSuccess() {
    fmt.Println("Pure vector SVG")
    fmt.Printf("Elements: %v\n", result.VectorElements)
} else {
    for _, e := range result.Errors {
        fmt.Printf("Error: %s\n", e)
    }
}
```

### Directory

Validates all SVG files in a directory (non-recursive).

```go
func Directory(dirPath string) ([]*Result, error)
```

**Example:**

```go
results, err := verify.Directory("brands/anthropic/")
if err != nil {
    log.Fatal(err)
}

for _, r := range results {
    if !r.IsSuccess() {
        fmt.Printf("%s: %v\n", r.FilePath, r.Errors)
    }
}
```

### DirectoryRecursive

Validates all SVG files in a directory tree.

```go
func DirectoryRecursive(dirPath string) ([]*Result, error)
```

**Example:**

```go
results, err := verify.DirectoryRecursive("brands/")
if err != nil {
    log.Fatal(err)
}

passed := 0
failed := 0
for _, r := range results {
    if r.IsSuccess() {
        passed++
    } else {
        failed++
        fmt.Printf("FAIL: %s\n", r.FilePath)
    }
}
fmt.Printf("Results: %d passed, %d failed\n", passed, failed)
```

## Detection Patterns

The verifier detects these embedded binary patterns:

| Pattern | Description |
|---------|-------------|
| `data:image/(png\|jpeg\|...)` | Base64 embedded image data |
| `xlink:href="data:..."` | Data URI in xlink:href |
| `href="data:image..."` | Data URI in href |
| `<image>` with binary href | Image element referencing .png, .jpg, etc. |

## Vector Elements

The verifier counts these SVG vector elements:

- `<path>` — Path elements
- `<rect>` — Rectangle elements
- `<circle>` — Circle elements
- `<ellipse>` — Ellipse elements
- `<line>` — Line elements
- `<polyline>` — Polyline elements
- `<polygon>` — Polygon elements
- `<text>` — Text elements

## CI Integration

```go
func main() {
    results, err := verify.DirectoryRecursive("brands/")
    if err != nil {
        log.Fatal(err)
    }

    hasFailures := false
    for _, r := range results {
        if !r.IsSuccess() {
            fmt.Printf("FAIL: %s - %v\n", r.FilePath, r.Errors)
            hasFailures = true
        }
    }

    if hasFailures {
        os.Exit(1)
    }
    fmt.Println("All SVG files are pure vector")
}
```
