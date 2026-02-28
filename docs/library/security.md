# svg/security Package

```go
import "github.com/grokify/brandkit/svg/security"
```

Provides SVG security scanning and sanitization to detect and remove malicious elements.

## Types

### ThreatType

Categorizes the type of security threat detected.

```go
type ThreatType int

const (
    ThreatScript       ThreatType = iota // Script elements, javascript: URIs
    ThreatEventHandler                   // Event handler attributes (onclick, etc.)
    ThreatExternalRef                    // External references (http:// URLs)
    ThreatAnimation                      // Animation elements
    ThreatStyleBlock                     // Style elements
    ThreatLink                           // Anchor elements
    ThreatXMLEntity                      // DOCTYPE/ENTITY declarations
)
```

#### Methods

| Method | Description |
|--------|-------------|
| `String() string` | Returns threat type name (e.g., "script") |
| `Severity() string` | Returns severity level (critical/high/medium/low) |

### Threat

Represents a detected security threat.

```go
type Threat struct {
    Type        ThreatType
    Description string
    Match       string
}
```

### Result

Contains the result of scanning an SVG file.

```go
type Result struct {
    FilePath     string
    IsSecure     bool
    Threats      []Threat
    ThreatCounts map[ThreatType]int
    Errors       []string
}
```

#### Methods

| Method | Description |
|--------|-------------|
| `IsSuccess() bool` | Returns `r.IsSecure && len(r.Errors) == 0` |

### ScanLevel

Defines how strict the security scan should be.

```go
type ScanLevel int

const (
    ScanLevelStrict   ScanLevel = iota // All threats
    ScanLevelStandard                  // Critical/high only
)
```

## Scanning Functions

### SVG

Scans a single SVG file using strict level.

```go
func SVG(filePath string) (*Result, error)
```

### SVGWithLevel

Scans a single SVG file with specified scan level.

```go
func SVGWithLevel(filePath string, level ScanLevel) (*Result, error)
```

**Example:**

```go
// Strict scan (all threats)
result, err := security.SVG("icon.svg")

// Standard scan (critical/high only)
result, err := security.SVGWithLevel("icon.svg", security.ScanLevelStandard)

if !result.IsSecure {
    for _, threat := range result.Threats {
        fmt.Printf("[%s] %s: %s\n",
            threat.Type.Severity(),
            threat.Description,
            threat.Match)
    }
}
```

### ScanContent

Scans SVG content string for threats.

```go
func ScanContent(content string, result *Result) *Result
func ScanContentWithLevel(content string, result *Result, level ScanLevel) *Result
```

### Directory

Scans all SVG files in a directory (non-recursive).

```go
func Directory(dirPath string) ([]*Result, error)
```

### DirectoryRecursive

Scans all SVG files in a directory tree.

```go
func DirectoryRecursive(dirPath string) ([]*Result, error)
```

## Sanitization

### SanitizeOptions

Configures sanitization behavior.

```go
type SanitizeOptions struct {
    RemoveScripts       bool
    RemoveEventHandlers bool
    RemoveExternalRefs  bool
    RemoveAll           bool
}
```

### Sanitize

Sanitizes an SVG file and writes to output.

```go
func Sanitize(inputPath, outputPath string, opts SanitizeOptions) (*SanitizeResult, error)
```

### SanitizeContent

Sanitizes SVG content in memory.

```go
func SanitizeContent(content string, opts SanitizeOptions) (string, []Threat)
```

**Example:**

```go
// Remove all threats
result, err := security.Sanitize("input.svg", "output.svg", security.SanitizeOptions{
    RemoveAll: true,
})

// Remove only scripts
result, err := security.Sanitize("input.svg", "output.svg", security.SanitizeOptions{
    RemoveScripts: true,
})

fmt.Printf("Removed %d threats\n", len(result.ThreatsRemoved))
```

## Report Generation

### GenerateReport

Creates a TeamReport from scan results.

```go
func GenerateReport(results []*Result, project, version string) *TeamReport
```

**Example:**

```go
results, _ := security.DirectoryRecursive("brands/")
report := security.GenerateReport(results, "myproject", "1.0.0")

jsonBytes, _ := report.ToJSON()
os.WriteFile("report.json", jsonBytes, 0644)
```

### TeamReport

JSON report following multi-agent-spec format.

```go
type TeamReport struct {
    Schema        string            `json:"$schema,omitempty"`
    Title         string            `json:"title,omitempty"`
    Project       string            `json:"project"`
    Version       string            `json:"version"`
    Phase         string            `json:"phase"`
    Status        Status            `json:"status"` // GO, NO-GO, WARN, SKIP
    Teams         []TeamSection     `json:"teams"`
    GeneratedAt   string            `json:"generated_at"`
    // ...
}
```

## Threat Severity

| ThreatType | Severity | Detected Patterns |
|------------|----------|-------------------|
| `ThreatScript` | critical | `<script>`, `javascript:`, `vbscript:` |
| `ThreatEventHandler` | critical | `onclick`, `onload`, `onerror`, etc. |
| `ThreatExternalRef` | high | `href="http://..."`, `foreignObject` |
| `ThreatXMLEntity` | high | `<!DOCTYPE>`, `<!ENTITY>` |
| `ThreatAnimation` | medium | `<animate>`, `<animateTransform>` |
| `ThreatLink` | medium | `<a>` elements |
| `ThreatStyleBlock` | low | `<style>` elements |
