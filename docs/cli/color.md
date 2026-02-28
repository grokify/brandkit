# brandkit color

Create a centered color icon preserving original colors.

## Synopsis

```bash
brandkit color <input> -o <output> [flags]
```

## Description

Shortcut for creating a well-sized, centered icon that preserves original colors. This is equivalent to:

```bash
brandkit process <input> -o <output> --remove-background --center --strict
```

The command performs these steps:

1. Removes any solid background
2. Preserves original colors (no color conversion)
3. Centers the content in the viewBox
4. Verifies the result is pure vector
5. Scans for security threats (fails by default if threats found)

## Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output file path (required) |
| `--insecure` | Warn on security threats instead of failing |
| `-h, --help` | Help for color |

## Examples

Create a color icon:

```bash
brandkit color icon_orig.svg -o icon_color.svg
```

Process a brand icon:

```bash
brandkit color brands/react/icon_orig.svg -o brands/react/icon_color.svg
```

Allow files with security warnings:

```bash
brandkit color icon_orig.svg -o icon_color.svg --insecure
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (invalid input, processing failure, security threats detected) |

## See Also

- [white](white.md) — Create white icon on transparent background
- [process](process.md) — Full processing pipeline with all options
- [convert](convert.md) — Convert to specific colors
