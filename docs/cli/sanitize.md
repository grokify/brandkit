# brandkit sanitize

Remove security threats from SVG files.

## Synopsis

```bash
brandkit sanitize <input> -o <output> [flags]
```

## Description

Remove security threats from an SVG file while preserving valid content. By default, all threat types are removed. Use flags to selectively remove specific threat types.

## Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output file path (required) |
| `--remove-all` | Remove all threat types (default: true) |
| `--remove-scripts` | Remove script elements only |
| `--remove-event-handlers` | Remove event handler attributes only |
| `--remove-external-refs` | Remove external URLs only |
| `-h, --help` | Help for sanitize |

## Examples

Remove all threats:

```bash
brandkit sanitize malicious.svg -o clean.svg
```

In-place sanitization (overwrites original):

```bash
brandkit sanitize icon.svg -o icon.svg
```

Remove only scripts:

```bash
brandkit sanitize icon.svg -o clean.svg --remove-scripts
```

Remove only event handlers:

```bash
brandkit sanitize icon.svg -o clean.svg --remove-event-handlers
```

## What Gets Removed

### --remove-all (default)

Removes all detected threats:

| Category | Elements Removed |
|----------|------------------|
| Scripts | `<script>` tags, `javascript:` URIs, `vbscript:` URIs |
| Event Handlers | `onclick`, `onload`, `onerror`, `onmouseover`, etc. |
| External Refs | `href="http://..."`, `xlink:href`, `foreignObject` |

### --remove-scripts

- `<script>` elements (inline and external)
- Self-closing `<script/>` tags
- `javascript:` URI schemes
- `vbscript:` URI schemes
- `data:text/html` URIs

### --remove-event-handlers

- All `on*` attributes: `onclick`, `onload`, `onerror`, `onmouseover`, etc.
- Both quoted and unquoted attribute values
- Handles nested quotes correctly

### --remove-external-refs

- `href="http://..."` and `href="https://..."`
- `xlink:href` with external URLs
- `<foreignObject>` elements
- `url(http://...)` in style attributes
- External `<use>` references

## Output

The sanitized SVG maintains:

- Valid SVG structure
- Internal ID references (`#id`)
- Local file references
- Inline styles (unless containing threats)
- All visual elements

## Verification

After sanitization, verify the result:

```bash
brandkit sanitize icon.svg -o clean.svg
brandkit security-scan clean.svg
```

## Batch Processing

Sanitize all SVG files in a directory:

```bash
for svg in brands/*/*.svg; do
  brandkit sanitize "$svg" -o "$svg"
done
```

Or using Make:

```bash
make sanitize-all
```

## See Also

- [security-scan](security-scan.md) — Scan for security threats
- [Security Guide](../security/index.md) — Full security documentation
- [Sanitization](../security/sanitization.md) — Detailed sanitization docs
