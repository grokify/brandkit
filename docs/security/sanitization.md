# Sanitization

BrandKit can remove security threats from SVG files while preserving valid content.

## CLI Usage

### Basic Sanitization

Remove all threats:

```bash
brandkit sanitize malicious.svg -o clean.svg
```

### In-Place Sanitization

Overwrite the original file:

```bash
brandkit sanitize icon.svg -o icon.svg
```

### Selective Removal

Remove only specific threat types:

```bash
# Scripts only
brandkit sanitize icon.svg -o clean.svg --remove-scripts

# Event handlers only
brandkit sanitize icon.svg -o clean.svg --remove-event-handlers

# External references only
brandkit sanitize icon.svg -o clean.svg --remove-external-refs
```

## Command Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output file path (required) |
| `--remove-all` | Remove all threat types (default: true) |
| `--remove-scripts` | Remove script elements only |
| `--remove-event-handlers` | Remove event handler attributes only |
| `--remove-external-refs` | Remove external URLs only |

## What Gets Removed

### Scripts (--remove-scripts)

| Element/Attribute | Removed |
|-------------------|---------|
| `<script>...</script>` | Entire element |
| `<script/>` | Self-closing element |
| `href="javascript:..."` | Attribute value |
| `href="vbscript:..."` | Attribute value |
| `href="data:text/html..."` | Attribute value |

### Event Handlers (--remove-event-handlers)

| Attribute | Removed |
|-----------|---------|
| `onclick="..."` | Entire attribute |
| `onload="..."` | Entire attribute |
| `onerror="..."` | Entire attribute |
| `onmouseover="..."` | Entire attribute |
| All `on*="..."` | Entire attribute |

### External References (--remove-external-refs)

| Element/Attribute | Removed |
|-------------------|---------|
| `href="http://..."` | Attribute |
| `href="https://..."` | Attribute |
| `xlink:href="http://..."` | Attribute |
| `<foreignObject>` | Entire element |
| `url(http://...)` | In style attributes |

## What Gets Preserved

- Internal ID references (`href="#myid"`)
- Local file references (`href="other.svg"`)
- Inline styles (unless containing external URLs)
- All visual elements (paths, shapes, text)
- Valid SVG structure

## Batch Processing

### Using Shell

```bash
for svg in brands/*/*.svg; do
  brandkit sanitize "$svg" -o "$svg"
done
```

### Using Make

```makefile
sanitize-all: build
	@for svg in $$(find brands -name "*.svg"); do \
		$(BUILD_DIR)/$(BINARY_NAME) sanitize $$svg -o $$svg.tmp && mv $$svg.tmp $$svg; \
	done
```

## Library Usage

### Sanitize File

```go
import "github.com/grokify/brandkit/svg/security"

result, err := security.Sanitize("input.svg", "output.svg", security.SanitizeOptions{
    RemoveAll: true,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Removed %d threats\n", len(result.ThreatsRemoved))
```

### Sanitize Content

```go
content := `<svg onclick="alert(1)"><rect/></svg>`

cleaned, removed := security.SanitizeContent(content, security.SanitizeOptions{
    RemoveEventHandlers: true,
})

fmt.Println(cleaned)
// Output: <svg><rect/></svg>

fmt.Printf("Removed: %v\n", removed)
```

### Selective Sanitization

```go
// Remove only critical threats
result, _ := security.Sanitize("input.svg", "output.svg", security.SanitizeOptions{
    RemoveScripts:       true,
    RemoveEventHandlers: true,
})

// Remove external refs only
result, _ := security.Sanitize("input.svg", "output.svg", security.SanitizeOptions{
    RemoveExternalRefs: true,
})
```

## Verification After Sanitization

Always verify the sanitized file:

```bash
brandkit sanitize input.svg -o output.svg
brandkit security-scan output.svg
```

Or programmatically:

```go
// Sanitize
security.Sanitize("input.svg", "output.svg", security.SanitizeOptions{
    RemoveAll: true,
})

// Verify
result, _ := security.SVG("output.svg")
if !result.IsSecure {
    log.Fatal("Sanitization incomplete")
}
```

## Limitations

Sanitization removes threats but cannot guarantee:

- Visual appearance is unchanged (removed elements may affect rendering)
- All possible attack vectors are covered (SVG is complex)
- Future vulnerabilities are addressed

For maximum security:

1. Sanitize incoming SVGs
2. Scan after sanitization
3. Use Content Security Policy (CSP)
4. Serve SVGs with proper MIME type
