package analyze

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/brandkit/svg"
)

func TestSVGCentered(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "centered.svg")

	// Content perfectly centered in viewBox with ~10% padding
	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <path d="M 10 10 L 90 10 L 90 90 L 10 90 Z" fill="#000000"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.HasIssues {
		t.Errorf("expected no issues for centered content, got: %s", result.Assessment)
	}

	// Padding should be ~10% on all sides
	if result.PaddingLeft < 9 || result.PaddingLeft > 11 {
		t.Errorf("PaddingLeft = %.1f, want ~10", result.PaddingLeft)
	}
	if result.PaddingRight < 9 || result.PaddingRight > 11 {
		t.Errorf("PaddingRight = %.1f, want ~10", result.PaddingRight)
	}
}

func TestSVGOffCenter(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "offcenter.svg")

	// Content shifted right and down
	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="40" y="40" width="60" height="60" fill="#000"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if !result.HasIssues {
		t.Error("expected issues for off-center content")
	}
	if result.CenterOffsetX <= 0 {
		t.Errorf("CenterOffsetX = %.1f, expected positive (shifted right)", result.CenterOffsetX)
	}
	if result.CenterOffsetY <= 0 {
		t.Errorf("CenterOffsetY = %.1f, expected positive (shifted down)", result.CenterOffsetY)
	}
}

func TestSVGExcessivePadding(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "padded.svg")

	// Tiny content in large viewBox
	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="40" y="40" width="20" height="20" fill="#000"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if !result.HasIssues {
		t.Error("expected issues for excessive padding")
	}
	if result.PaddingLeft <= 20 {
		t.Errorf("PaddingLeft = %.1f, expected >20 for excessive padding", result.PaddingLeft)
	}
}

func TestSVGWithWidthHeight(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "wh.svg")

	// No viewBox, only width/height
	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
  <rect x="10" y="10" width="80" height="80" fill="#000"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.ViewBox.Width != 100 || result.ViewBox.Height != 100 {
		t.Errorf("ViewBox = %v, expected 100x100 from width/height", result.ViewBox)
	}
}

func TestSVGNoViewBoxOrDimensions(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "nodims.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg">
  <rect x="10" y="10" width="80" height="80"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := SVG(file)
	if err == nil {
		t.Error("expected error for SVG with no viewBox or dimensions")
	}
}

func TestSVGFileNotFound(t *testing.T) {
	_, err := SVG("/nonexistent/path.svg")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSuggestViewBox(t *testing.T) {
	box := svg.NewBoundingBox()
	box.Expand(10, 10)
	box.Expand(90, 90)

	suggested := SuggestViewBox(box)
	if suggested == "" {
		t.Error("expected non-empty suggested viewBox")
	}

	// Parse the suggestion and verify it has padding
	vb, err := svg.ParseViewBox(suggested)
	if err != nil {
		t.Fatalf("failed to parse suggested viewBox: %v", err)
	}

	// The viewBox should be larger than the content (80x80)
	if vb.Width <= 80 {
		t.Errorf("suggested width = %.1f, expected > 80", vb.Width)
	}
	if vb.Height <= 80 {
		t.Errorf("suggested height = %.1f, expected > 80", vb.Height)
	}
}

func TestSuggestViewBoxSquareAspect(t *testing.T) {
	box := svg.NewBoundingBox()
	box.Expand(0, 0)
	box.Expand(95, 100) // Nearly square

	suggested := SuggestViewBox(box)
	vb, err := svg.ParseViewBox(suggested)
	if err != nil {
		t.Fatalf("failed to parse suggested viewBox: %v", err)
	}

	// Should be made square since aspect ratio is close to 1:1
	if vb.Width != vb.Height {
		t.Errorf("expected square viewBox for near-square content, got %.1fx%.1f", vb.Width, vb.Height)
	}
}

func TestDirectory(t *testing.T) {
	dir := t.TempDir()

	content := `<?xml version="1.0"?><svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect x="10" y="10" width="80" height="80"/></svg>`

	if err := os.WriteFile(filepath.Join(dir, "a.svg"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.svg"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	results, err := Directory(dir)
	if err != nil {
		t.Fatalf("Directory error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("got %d results, want 2", len(results))
	}
}

func TestDirectoryWithBadFile(t *testing.T) {
	dir := t.TempDir()

	good := `<?xml version="1.0"?><svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect x="10" y="10" width="80" height="80"/></svg>`
	bad := `not valid svg at all`

	if err := os.WriteFile(filepath.Join(dir, "good.svg"), []byte(good), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bad.svg"), []byte(bad), 0600); err != nil {
		t.Fatal(err)
	}

	results, err := Directory(dir)
	if err != nil {
		t.Fatalf("Directory error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	// Bad file should have issues
	hasError := false
	for _, r := range results {
		if r.HasIssues {
			hasError = true
		}
	}
	if !hasError {
		t.Error("expected at least one result with issues")
	}
}
