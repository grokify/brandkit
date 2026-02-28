# brandkit analyze

Analyze SVG geometry including centering, padding, and viewBox optimization.

## Synopsis

```bash
brandkit analyze [path] [flags]
```

## Description

Analyze SVG files to check:

- ViewBox dimensions
- Content centering
- Padding percentages (left, right, top, bottom)
- Suggested viewBox fixes for optimal centering

Can analyze a single file or all SVG files in a directory.

## Flags

| Flag | Description |
|------|-------------|
| `--fix` | Show suggested viewBox fixes |
| `-h, --help` | Help for analyze |

## Examples

Analyze a single file:

```bash
brandkit analyze icon.svg
```

Analyze all SVGs in a directory:

```bash
brandkit analyze brands/
```

Show suggested fixes:

```bash
brandkit analyze brands/ --fix
```

## Output

The analysis output includes:

```
File: brands/react/icon_orig.svg
  ViewBox: 0 0 100 100
  Content Bounds: x=10 y=15 w=80 h=70
  Padding: L:10.0% R:10.0% T:15.0% B:15.0%
  Centered: No (horizontal: yes, vertical: no)
  Suggested ViewBox: 10 15 80 70
```

### Output Fields

| Field | Description |
|-------|-------------|
| ViewBox | Current viewBox attribute |
| Content Bounds | Calculated bounding box of all visual content |
| Padding | Percentage of empty space on each side |
| Centered | Whether content is centered horizontally and vertically |
| Suggested ViewBox | Optimized viewBox for zero padding |

## Use Cases

### Checking Brand Consistency

Ensure all brand icons have consistent padding:

```bash
brandkit analyze brands/ | grep "Padding:"
```

### Finding Uncentered Icons

Identify icons that need centering:

```bash
brandkit analyze brands/ | grep -A1 "Centered: No"
```

### Generating Fix Commands

Get suggested viewBox values for optimization:

```bash
brandkit analyze brands/ --fix
```

## See Also

- [process](process.md) — Full processing pipeline with `--center` option
- [white](white.md) — Creates centered white icons
- [color](color.md) — Creates centered color icons
