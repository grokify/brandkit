# Release Notes v0.3.0

Release date: 2026-02-14

## Overview

This release adds an interactive coordinate picker tool for creating precise SVG polygons from images, a TypeScript library with 100% test coverage, and new Saviynt brand icons with mathematically precise geometry.

## New Features

### Coordinate Picker Tool

Interactive browser-based tool for picking coordinates from images to create SVG polygons. Useful for tracing brand icons and creating precise vector graphics.

**Access via GitHub Pages:** `https://grokify.github.io/brandkit/`

**Features:**

- Load images (JPG, PNG, SVG) via drag-and-drop or file picker
- Zoom controls (50% to 400%) for precise coordinate selection
- Multi-shape support with color-coded markers
- Pin sharing across shapes for shared vertices
- Context menu for pin management (add to shape, remove from shape, delete)
- Live SVG preview with real-time updates
- Click-to-copy for individual polygons or full SVG code

### TypeScript Library

Reusable TypeScript library (`tools/src/coordinate-picker.ts`) with:

- State management for pins and shapes
- Coordinate calculations with zoom support
- SVG polygon generation
- Serialization for save/load functionality

**Test Coverage:** 63 unit tests with 100% coverage on statements, branches, functions, and lines.

### Saviynt Brand Icons

New brand icons for Saviynt identity platform with mathematically precise geometry:

- `icon_orig.svg` — Original green (#00FF00) on black background
- `icon_white.svg` — White variant for dark backgrounds
- `icon_color.svg` — Color variant without background

**Geometry:**

- 4 parallelograms forming an "S" shape
- X-coordinates: 70, 135, 200, 265, 330 (65px spacing)
- Diagonal slope: ±37/65 for parallel edges
- Vertical edges at x=70 and x=330

## CI/CD

- New `test-tools.yaml` workflow for TypeScript tests
- Matrix: Node.js 20.x/22.x across ubuntu, macos, windows
- Path filtering: triggers only on `tools/` changes

## File Structure

```
docs/
  index.html              # Coordinate picker (GitHub Pages)

tools/
  src/coordinate-picker.ts   # TypeScript library
  tests/coordinate-picker.test.ts  # Unit tests (63 tests, 100% coverage)
  package.json            # npm config
  tsconfig.json           # TypeScript config
  vitest.config.ts        # Test config

brands/saviynt/
  icon_orig.svg           # Original (green on black)
  icon_white.svg          # White variant
  icon_color.svg          # Color variant
```

## Requirements

- Go 1.24+ (for brandkit CLI)
- Node.js 20+ (for tools development)
