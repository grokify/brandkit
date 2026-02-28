# Tools

BrandKit includes interactive tools for SVG development.

## Coordinate Picker

An interactive browser-based tool for creating precise SVG polygons from images.

<a href="coordinate-picker.html" class="md-button md-button--primary">Open Coordinate Picker</a>

### Features

- **Interactive Clicking** — Click on images to place coordinate points
- **Multi-Shape Support** — Create multiple shapes with color-coded markers
- **Pin Sharing** — Share coordinate pins across shapes
- **Zoom Controls** — 50% to 400% zoom for precise placement
- **SVG Export** — Generate SVG polygon code from coordinates
- **State Serialization** — Save and load coordinate sessions

### Use Cases

- Creating precise SVG icon shapes from reference images
- Tracing logos for vectorization
- Building complex multi-path SVG icons
- Generating mathematically precise geometry

### Getting Started

1. Open the [Coordinate Picker](coordinate-picker.html)
2. Load an image (drag & drop or file picker)
3. Click to place coordinate points
4. Adjust zoom for precision
5. Copy generated SVG code

### Technical Details

The coordinate picker is built with:

- TypeScript for type-safe coordinate handling
- Vanilla JavaScript for zero-dependency runtime
- HTML5 Canvas for image rendering
- 100% test coverage

Source code: [`tools/src/coordinate-picker.ts`](https://github.com/grokify/brandkit/tree/main/tools/src)
