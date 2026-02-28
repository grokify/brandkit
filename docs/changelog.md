# Changelog

All notable changes to BrandKit are documented here.

## [v0.4.0] - 2026-02-26

### Highlights

- SVG security scanning to detect XSS and session hijacking threats
- SVG sanitization to remove malicious elements while preserving valid content
- Security scanning integrated into processing pipelines with CI support

### Added

- Public `svg/security` package for programmatic security scanning with 7 `ThreatType` values: `ThreatScript`, `ThreatEventHandler`, `ThreatExternalRef`, `ThreatAnimation`, `ThreatStyleBlock`, `ThreatLink`, `ThreatXMLEntity`
- Security scanning functions: `security.SVG()`, `security.SVGWithLevel()`, `security.Directory()`, `security.DirectoryRecursive()`, `security.ScanContent()`, `security.ScanContentWithLevel()`
- Scan levels: `ScanLevelStrict` (all threats) and `ScanLevelStandard` (critical/high only)
- Sanitization functions: `security.Sanitize()`, `security.SanitizeContent()` with configurable `SanitizeOptions`
- Team report generation: `security.GenerateReport()` outputs multi-agent-spec team-report JSON format
- CLI command `brandkit security-scan` with `--report`, `--strict`, `--project`, `--version` flags
- CLI command `brandkit security-scan-all` for recursive security scanning with JSON report output
- CLI command `brandkit sanitize` for removing threats from SVG files with selective removal options
- Makefile targets `security-scan-all` and `sanitize-all` for batch operations

### Changed

- CLI commands `brandkit white` and `brandkit color` now perform security scanning by default
- Added `--insecure` flag to `white` and `color` commands to warn instead of fail on threats
- Library functions `ProcessWhite()` and `ProcessColor()` now include security scanning in pipeline
- `ProcessResult` struct extended with `SecurityScanned` and `SecurityThreats` fields

### Security

- Detects script elements (`<script>`) and self-closing script tags (critical)
- Detects dangerous URI schemes: `javascript:`, `vbscript:`, `data:text/html` (critical)
- Detects event handler attributes (`onclick`, `onload`, `onerror`, `onmouseover`, etc.) (critical)
- Detects external references: `href="http://..."`, `xlink:href`, `foreignObject`, `url()` in styles, external `<use>` refs (high)
- Detects XML entities: `<!DOCTYPE>`, `<!ENTITY>` declarations for XXE prevention (high)
- Detects animation elements: `<animate>`, `<animateTransform>`, `<animateMotion>`, `<set>` (medium)
- Detects `<style>` blocks that may contain malicious CSS (low)
- Detects `<a>` anchor/link elements unnecessary for static images (medium)

### Infrastructure

- GitHub Actions workflow `verify.yaml` updated to include security scanning step

### Tests

- 24 unit tests for security scanning covering all 7 threat types and scan levels
- Tests verify sanitized output remains valid SVG and passes security scan
- Tests cover ScanLevelStrict vs ScanLevelStandard behavior differences

---

## [v0.3.0] - 2026-02-14

### Highlights

- Interactive coordinate picker tool for creating precise SVG polygons from images
- TypeScript library with 100% test coverage for coordinate picker functionality
- Saviynt brand icons with mathematically precise geometry

### Added

- Saviynt brand icons (`icon_orig.svg`, `icon_white.svg`, `icon_color.svg`) with precise parallelogram geometry
- Interactive coordinate picker tool (`docs/coordinate-picker.html`) for SVG polygon creation via GitHub Pages
- TypeScript library (`tools/src/coordinate-picker.ts`) with state management, SVG generation, and serialization
- Multi-shape support with color-coded markers and pin sharing across shapes
- Zoom controls (50%-400%) for precise coordinate selection
- GitHub Actions workflow (`test-tools.yaml`) for TypeScript tests on Node.js 20.x/22.x across platforms

### Tests

- 63 unit tests for coordinate picker with 100% coverage on statements, branches, functions, and lines

### Infrastructure

- Move coordinate picker to `docs/` for GitHub Pages static hosting

---

## [v0.2.0] - 2026-01-25

### Highlights

- Go library APIs for programmatic icon retrieval and processing
- New CLI commands `brandkit color` and `brandkit verify-all`
- 17 new brand icons and 42 new `icon_color.svg` files (total: 52 brands)

### Added

- Go library API for embedded icon retrieval: `GetIcon`, `GetIconWhite`, `GetIconColor`, `GetIconOrig`, `ListIcons`, `IconExists`, `NormalizeIconName`
- Go library API for icon processing: `ProcessWhite`, `ProcessColor` for programmatic SVG processing
- Go library API for recursive verification: `svg.ListSVGFilesRecursive`, `verify.DirectoryRecursive`
- CLI command `brandkit color` for creating centered color icons preserving original colors
- CLI command `brandkit verify-all` for recursive pure vector verification (CI-friendly)
- 17 new brand icons: bolt, bootstrap, dart, flutter, go, javascript, kotlin, lovable, openapi, postgresql, postman, python, react, replit, spring, v0, windsurf (total: 52 brands)
- 42 new `icon_color.svg` files for all brands with transparent backgrounds
- GitHub Actions workflow `verify.yaml` for automated icon validation on PR/push

### Changed

- Renamed `gcp` brand directory to `google-gcp` for clarity and consistency
- Refactored `brandkit white` CLI to use `ProcessWhite` library function

---

## [v0.1.0] - 2026-01-24

### Highlights

- CLI toolkit for one-command SVG icon processing workflows
- Go library packages for SVG analysis, conversion, and verification
- Brand asset library with 36 brand icon directories

### Added

- CLI command `brandkit white` for one-step white icon generation with background removal, color conversion, centering, and verification
- CLI command `brandkit convert` for SVG color conversion with support for hex, shorthand, and named colors
- CLI command `brandkit process` for full SVG processing pipeline (convert, center, verify)
- CLI command `brandkit analyze` for SVG geometry analysis (centering, padding, viewBox optimization)
- CLI command `brandkit verify` for pure vector validation (detects embedded base64, data URIs, binary references)
- Public `svg` package with `BoundingBox`, `ViewBox`, `ParsePath`, `CalculatePathBounds`, and file utilities
- Public `svg/convert` package for programmatic SVG color conversion with background removal and mask preservation
- Public `svg/analyze` package for programmatic SVG geometry analysis
- Public `svg/verify` package for programmatic pure vector validation
- Brand asset library with 36 brand icon directories and standardized SVG variants
