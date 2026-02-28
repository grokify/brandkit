# brandkit verify

Verify SVG files are pure vector images.

## Synopsis

```bash
brandkit verify [path] [flags]
brandkit verify-all [path] [flags]
```

## Description

Verify SVG files are pure vector images without:

- Embedded binary data (base64 images)
- Data URIs (`data:image/png;base64,...`)
- External binary image references (`<image href="photo.png">`)

This ensures SVG files are scalable and don't contain hidden raster content.

## Commands

### verify

Verify a single file or all SVGs in a directory (non-recursive):

```bash
brandkit verify icon.svg
brandkit verify brands/anthropic/
```

### verify-all

Recursively verify all SVG files in a directory tree. Designed for CI pipelines:

```bash
brandkit verify-all brands/
```

## Flags

| Flag | Description |
|------|-------------|
| `-h, --help` | Help for verify |

## Examples

Verify a single file:

```bash
brandkit verify icon.svg
```

Verify all files in a directory:

```bash
brandkit verify brands/
```

Recursive verification (for CI):

```bash
brandkit verify-all brands/
```

## Output

### Success

```
✓ icon.svg: Pure vector SVG
```

### Failure

```
✗ icon.svg: Embedded data detected
  - Base64 image at line 15
  - Data URI at line 23
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All files are pure vector |
| 1 | One or more files contain embedded raster data |

## Detection Patterns

The verification detects:

| Pattern | Description |
|---------|-------------|
| `data:image/` | Data URI images |
| `base64,` | Base64 encoded content |
| `<image>` with external `href` | External image references |
| Binary signatures | PNG, JPEG, GIF headers |

## CI Integration

Add to your CI pipeline:

```yaml
- name: Verify SVG icons
  run: brandkit verify-all brands/
```

Or using Make:

```bash
make verify
```

## See Also

- [security-scan](security-scan.md) — Scan for security threats
- [process](process.md) — Processing pipeline with verification
