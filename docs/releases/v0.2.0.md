# Release Notes v0.2.0

Release date: 2026-01-25

## Overview

This release adds Go library APIs for programmatic icon processing, new CLI commands for color icon generation and recursive verification, 42 new color icon assets, and CI workflow for automated icon validation.

## New Features

### Go Library API

**Icon Retrieval** (added in this release)

- `GetIcon(brand, variant)` - retrieve icon by name and variant
- `GetIconWhite(brand)` - white variant for dark backgrounds
- `GetIconColor(brand)` - full color variant
- `GetIconOrig(brand)` - original unmodified variant
- `ListIcons()` - list all available brand names
- `IconExists(brand)` - check if brand exists
- `NormalizeIconName(name)` - convert aliases (golang→go, k8s→kubernetes)

**Icon Processing**

- `ProcessWhite(input, output)` - create white icon on transparent background
- `ProcessColor(input, output)` - create centered color icon preserving original colors

**Recursive Verification**

- `svg.ListSVGFilesRecursive(dir)` - walk directory tree for SVG files
- `verify.DirectoryRecursive(dir)` - verify all SVGs in directory tree

### CLI Commands

**`brandkit color`** - Create centered color icon on transparent background

```bash
brandkit color icon_orig.svg -o icon_color.svg
```

Removes background, centers content, verifies pure vector while preserving original colors.

**`brandkit verify-all`** - Recursively verify all SVGs are pure vector

```bash
brandkit verify-all brands/
```

Designed for CI pipelines to ensure all brand icons remain valid.

### Brand Assets

- **17 new brands**: bolt, bootstrap, dart, flutter, go, javascript, kotlin, lovable, openapi, postgresql, postman, python, react, replit, spring, v0, windsurf
- **42 new `icon_color.svg` files** for all brands with transparent backgrounds
- **Total**: 52 brands, 157 SVG files

### CI/CD

- New `verify.yaml` workflow for automated icon validation
- Triggers only when `brands/` directory changes
- Runs `make verify-all` to ensure all icons are pure vector

## Breaking Changes

- Renamed `gcp` brand directory to `google-gcp` for clarity

## Upgrade Notes

Update your import to use the new brand name:

```go
// Before
svg, _ := brandkit.GetIconWhite("gcp")

// After
svg, _ := brandkit.GetIconWhite("google-gcp")
```

Or use `NormalizeIconName` which maps `gcloud` → `gcp` (but not `gcp` → `google-gcp`).
