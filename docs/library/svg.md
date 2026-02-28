# svg Package

```go
import "github.com/grokify/brandkit/svg"
```

Core SVG types and utilities shared across all packages.

## Types

### BoundingBox

Represents a rectangular bounding box.

```go
type BoundingBox struct {
    MinX float64
    MinY float64
    MaxX float64
    MaxY float64
}
```

#### Methods

| Method | Description |
|--------|-------------|
| `Width() float64` | Returns the width of the bounding box |
| `Height() float64` | Returns the height of the bounding box |
| `CenterX() float64` | Returns the X coordinate of the center |
| `CenterY() float64` | Returns the Y coordinate of the center |
| `IsValid() bool` | Returns true if the box has been expanded with at least one point |
| `Expand(x, y float64)` | Expands the box to include the given point |
| `Merge(other *BoundingBox)` | Merges another bounding box into this one |

#### Example

```go
box := svg.NewBoundingBox()
box.Expand(10, 10)
box.Expand(90, 90)

fmt.Printf("Size: %.1f x %.1f\n", box.Width(), box.Height())
// Output: Size: 80.0 x 80.0

fmt.Printf("Center: (%.1f, %.1f)\n", box.CenterX(), box.CenterY())
// Output: Center: (50.0, 50.0)
```

### ViewBox

Represents an SVG viewBox attribute.

```go
type ViewBox struct {
    X      float64
    Y      float64
    Width  float64
    Height float64
}
```

#### Methods

| Method | Description |
|--------|-------------|
| `CenterX() float64` | Returns the X coordinate of the viewBox center |
| `CenterY() float64` | Returns the Y coordinate of the viewBox center |
| `String() string` | Returns the viewBox as a string for SVG attribute |

#### Example

```go
vb, err := svg.ParseViewBox("0 0 100 100")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Center: (%.1f, %.1f)\n", vb.CenterX(), vb.CenterY())
// Output: Center: (50.0, 50.0)

fmt.Println(vb.String())
// Output: 0.0 0.0 100.0 100.0
```

## Functions

### ParseViewBox

Parses a viewBox string.

```go
func ParseViewBox(s string) (ViewBox, error)
```

**Parameters:**

- `s` — ViewBox string (e.g., "0 0 100 100")

**Returns:**

- `ViewBox` — Parsed viewBox
- `error` — Error if format is invalid

### ParseFloat

Parses a float with a default value on error.

```go
func ParseFloat(s string, defaultVal float64) float64
```

**Parameters:**

- `s` — String to parse (handles "px" suffix)
- `defaultVal` — Default value if parsing fails

### NewBoundingBox

Creates an empty bounding box.

```go
func NewBoundingBox() *BoundingBox
```

### GetElementBounds

Calculates bounds for an SVG element.

```go
func GetElementBounds(element *svgparser.Element) *BoundingBox
```

## File Utilities

### ListSVGFiles

Lists all SVG files in a directory (non-recursive).

```go
func ListSVGFiles(dirPath string) ([]string, error)
```

### ListSVGFilesRecursive

Lists all SVG files in a directory tree.

```go
func ListSVGFilesRecursive(dirPath string) ([]string, error)
```

**Example:**

```go
files, err := svg.ListSVGFilesRecursive("brands/")
if err != nil {
    log.Fatal(err)
}

for _, file := range files {
    fmt.Println(file)
}
```
