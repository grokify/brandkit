# Release Notes — v0.1.0

*January 24, 2026*

Initial release of brandkit, a Go toolkit for processing brand SVG icons.

## Features

### CLI Commands

- **`brandkit white`** — One-command workflow to create white icons on transparent backgrounds (remove background, convert colors, center, verify)
- **`brandkit convert`** — Convert SVG fill/stroke colors to any target color (hex or named), with background removal and mask preservation
- **`brandkit process`** — Full pipeline: color conversion, viewBox centering, and pure vector verification
- **`brandkit analyze`** — Analyze SVG geometry for centering offsets, padding percentages, and viewBox optimization
- **`brandkit verify`** — Validate SVG files contain only vector elements (no embedded base64, data URIs, or binary references)

### Library Packages

- **`svg`** — Shared types (`BoundingBox`, `ViewBox`), SVG path parsing with full command support (M, L, H, V, C, S, Q, T, A, Z + relative variants), and file utilities
- **`svg/analyze`** — Programmatic SVG geometry analysis with configurable thresholds for centering (5%), excessive padding (20%), and uneven padding (10%)
- **`svg/convert`** — Color normalization (hex, shorthand, named colors), fill/stroke conversion, background element detection and removal, mask/clipPath preservation
- **`svg/verify`** — Pure vector validation with pattern-based detection of embedded images, data URIs, and external binary references

### Brand Assets

- 36 brand icon directories with standardized SVG variants
- Includes: Anthropic, AWS, Azure, GCP, Docker, Kubernetes, GitHub, OpenAI, and more

## Requirements

- Go 1.25.6+
