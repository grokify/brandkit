# svg/analyze Package

```go
import "github.com/grokify/brandkit/svg/analyze"
```

Provides SVG geometry analysis for centering and padding assessment.

## Types

### Result

Contains the analysis results for an SVG file.

```go
type Result struct {
    FilePath         string
    ViewBox          svg.ViewBox
    ContentBox       svg.BoundingBox
    CenterOffsetX    float64
    CenterOffsetY    float64
    PaddingLeft      float64
    PaddingRight     float64
    PaddingTop       float64
    PaddingBottom    float64
    Assessment       string
    SuggestedViewBox string
    HasIssues        bool
}
```

| Field | Description |
|-------|-------------|
| `FilePath` | Path to the analyzed file |
| `ViewBox` | Current viewBox of the SVG |
| `ContentBox` | Calculated bounding box of visual content |
| `CenterOffsetX` | Horizontal offset from center (positive = right) |
| `CenterOffsetY` | Vertical offset from center (positive = down) |
| `PaddingLeft/Right/Top/Bottom` | Padding percentages on each side |
| `Assessment` | Human-readable assessment ("OK" or list of issues) |
| `SuggestedViewBox` | Optimized viewBox with 5% padding |
| `HasIssues` | True if any centering/padding issues detected |

## Functions

### SVG

Analyzes a single SVG file.

```go
func SVG(filePath string) (*Result, error)
```

**Example:**

```go
result, err := analyze.SVG("icon.svg")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Padding: L:%.1f%% R:%.1f%% T:%.1f%% B:%.1f%%\n",
    result.PaddingLeft, result.PaddingRight,
    result.PaddingTop, result.PaddingBottom)

if result.HasIssues {
    fmt.Printf("Issues: %s\n", result.Assessment)
    fmt.Printf("Suggested: %s\n", result.SuggestedViewBox)
}
```

### Directory

Analyzes all SVG files in a directory (non-recursive).

```go
func Directory(dirPath string) ([]*Result, error)
```

**Example:**

```go
results, err := analyze.Directory("brands/anthropic/")
if err != nil {
    log.Fatal(err)
}

for _, r := range results {
    if r.HasIssues {
        fmt.Printf("%s: %s\n", r.FilePath, r.Assessment)
    }
}
```

### SuggestViewBox

Suggests an optimized viewBox with 5% padding.

```go
func SuggestViewBox(contentBox *svg.BoundingBox) string
```

**Example:**

```go
box := svg.NewBoundingBox()
box.Expand(10, 10)
box.Expand(90, 90)

suggested := analyze.SuggestViewBox(box)
fmt.Println(suggested)
// Output: "5.3 5.3 89.5 89.5" (approximate)
```

## Issue Detection

The analyzer detects these issues:

| Issue | Threshold |
|-------|-----------|
| Content shifted left/right | > 5% of viewBox width |
| Content shifted up/down | > 5% of viewBox height |
| Excessive padding | > 20% on any side |
| Uneven horizontal padding | > 10% difference left vs right |
| Uneven vertical padding | > 10% difference top vs bottom |

## Use Cases

### Finding Uncentered Icons

```go
results, _ := analyze.Directory("brands/")
for _, r := range results {
    if r.HasIssues {
        fmt.Printf("%s needs fixing: %s\n", r.FilePath, r.Assessment)
    }
}
```

### Generating Fix Commands

```go
results, _ := analyze.Directory("brands/")
for _, r := range results {
    if r.HasIssues {
        fmt.Printf("// Fix %s: viewBox=\"%s\"\n", r.FilePath, r.SuggestedViewBox)
    }
}
```
