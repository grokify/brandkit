# brandkit convert

Convert SVG colors with fine-grained control.

## Synopsis

```bash
brandkit convert <input> -o <output> [flags]
```

## Description

Convert colors in an SVG file. Supports hex colors, shorthand hex, and named colors. Can optionally remove backgrounds and convert stroke colors.

## Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output file path (required) |
| `-c, --color` | Target color (hex or name, e.g., `ffffff`, `white`, `ff5500`) |
| `--remove-background` | Remove full-bleed background rect/circle |
| `--include-stroke` | Also convert stroke colors |
| `--preserve-masks` | Don't modify colors in mask/clipPath (default: true) |
| `-h, --help` | Help for convert |

## Color Formats

The `--color` flag accepts:

| Format | Example | Result |
|--------|---------|--------|
| 6-digit hex | `ffffff` | White |
| 3-digit hex | `fff` | White |
| Named color | `white` | White |
| Mixed case | `FF5500` | Orange |

## Examples

Convert to white:

```bash
brandkit convert icon_orig.svg -o icon_white.svg --color ffffff
```

Convert using named color:

```bash
brandkit convert icon.svg -o output.svg --color black
```

Remove background and convert:

```bash
brandkit convert icon.svg -o output.svg --color ffffff --remove-background --include-stroke
```

Copy without color change (just reformats):

```bash
brandkit convert icon.svg -o output.svg
```

## Background Removal

The `--remove-background` flag removes full-bleed background elements:

- `<rect>` elements covering the entire viewBox
- `<circle>` elements covering the entire viewBox
- Elements with fill colors that appear to be backgrounds

## Mask Preservation

By default, colors inside `<mask>` and `<clipPath>` elements are not converted. This preserves the visual appearance of masked content. Use `--preserve-masks=false` to convert all colors.

## See Also

- [white](white.md) — Shortcut for white icon creation
- [color](color.md) — Shortcut for color icon creation
- [process](process.md) — Full processing pipeline
