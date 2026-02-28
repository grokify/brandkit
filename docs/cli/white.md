# brandkit white

Create a white icon on transparent background.

## Synopsis

```bash
brandkit white <input> -o <output> [flags]
```

## Description

Shortcut for the common workflow of creating a white icon on transparent background. This is equivalent to:

```bash
brandkit process <input> -o <output> --color ffffff --remove-background --include-stroke --center --strict
```

The command performs these steps:

1. Removes any solid background
2. Converts all colors to white (`#ffffff`)
3. Includes stroke colors in conversion
4. Centers the content in the viewBox
5. Verifies the result is pure vector
6. Scans for security threats (fails by default if threats found)

## Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output file path (required) |
| `--insecure` | Warn on security threats instead of failing |
| `-h, --help` | Help for white |

## Examples

Create a white icon:

```bash
brandkit white icon_orig.svg -o icon_white.svg
```

Process a brand icon:

```bash
brandkit white brands/anthropic/icon_orig.svg -o brands/anthropic/icon_white.svg
```

Allow files with security warnings (style blocks, etc.):

```bash
brandkit white icon_orig.svg -o icon_white.svg --insecure
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (invalid input, processing failure, security threats detected) |

## See Also

- [color](color.md) — Create color icon preserving original colors
- [process](process.md) — Full processing pipeline with all options
- [security-scan](security-scan.md) — Scan for security threats
