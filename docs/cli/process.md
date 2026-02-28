# brandkit process

Full SVG processing pipeline with all options.

## Synopsis

```bash
brandkit process <input> -o <output> [flags]
```

## Description

Process an SVG file through the complete pipeline:

1. **Remove background** — Remove full-bleed background elements (if `--remove-background`)
2. **Convert colors** — Convert to target color (if `--color` specified)
3. **Center content** — Analyze and fix viewBox for optimal centering (if `--center`)
4. **Verify vector** — Ensure output is pure vector, no embedded raster (if `--strict`)

## Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output file path (required) |
| `-c, --color` | Target color (hex or name) |
| `--remove-background` | Remove full-bleed background rect/circle |
| `--include-stroke` | Also convert stroke colors |
| `--center` | Auto-fix viewBox for centering |
| `--strict` | Fail on embedded binary (default: true) |
| `-h, --help` | Help for process |

## Examples

Full processing pipeline:

```bash
brandkit process icon_orig.svg -o icon_white.svg --color ffffff --center --strict
```

Remove background and convert:

```bash
brandkit process icon_orig.svg -o icon_white.svg --remove-background --color ffffff
```

Center without color conversion:

```bash
brandkit process input.svg -o output.svg --center --strict
```

## Pipeline Steps

### 1. Background Removal

When `--remove-background` is specified, the following elements are removed:

- Full-bleed `<rect>` elements (covering entire viewBox)
- Full-bleed `<circle>` elements
- First child elements that appear to be backgrounds

### 2. Color Conversion

When `--color` is specified, all fill and stroke colors are converted:

- Hex colors (`#ffffff`, `#fff`)
- RGB colors (`rgb(255,255,255)`)
- Named colors (`white`, `black`)

Mask and clipPath elements are preserved by default.

### 3. Centering

When `--center` is specified:

- Calculates content bounding box
- Adjusts viewBox for optimal centering
- Adds equal padding on all sides

### 4. Verification

When `--strict` is specified (default: true):

- Detects embedded base64 images
- Detects data URIs
- Detects binary references
- Fails if any embedded raster data found

## See Also

- [white](white.md) — Shortcut for white icon creation
- [color](color.md) — Shortcut for color icon creation
- [convert](convert.md) — Color conversion only
- [analyze](analyze.md) — Geometry analysis
- [verify](verify.md) — Vector verification
