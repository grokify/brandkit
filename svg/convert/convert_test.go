package convert

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeColor(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"ffffff", "#ffffff", false},
		{"#ffffff", "#ffffff", false},
		{"FFFFFF", "#ffffff", false},
		{"fff", "#ffffff", false},
		{"#fff", "#ffffff", false},
		{"000000", "#000000", false},
		{"000", "#000000", false},
		{"abc", "#aabbcc", false},
		{"#ABC", "#aabbcc", false},
		{"white", "#ffffff", false},
		{"black", "#000000", false},
		{"red", "#ff0000", false},
		{"transparent", "none", false},
		{"", "", false},
		{"gggggg", "", true},  // invalid hex chars
		{"12345", "", true},   // wrong length
		{"1234567", "", true}, // too long
	}

	for _, tt := range tests {
		got, err := NormalizeColor(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("NormalizeColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("NormalizeColor(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNormalizeColorTrimSpace(t *testing.T) {
	got, err := NormalizeColor("  ffffff  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "#ffffff" {
		t.Errorf("got %q, want %q", got, "#ffffff")
	}
}

func TestSVGColorConversion(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.svg")
	output := filepath.Join(dir, "output.svg")

	svgContent := `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <path fill="#ff0000" d="M 10 10 L 90 10 L 90 90 Z"/>
  <circle fill="#00ff00" cx="50" cy="50" r="20"/>
</svg>`

	if err := os.WriteFile(input, []byte(svgContent), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(input, output, Options{Color: "ffffff"})
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}
	if !result.Converted {
		t.Error("expected Converted = true")
	}
	if result.TargetColor != "#ffffff" {
		t.Errorf("TargetColor = %q, want %q", result.TargetColor, "#ffffff")
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if !contains(contentStr, `fill="#ffffff"`) {
		t.Error("output should contain fill=\"#ffffff\"")
	}
	if contains(contentStr, `fill="#ff0000"`) || contains(contentStr, `fill="#00ff00"`) {
		t.Error("output should not contain original colors")
	}
}

func TestSVGPreserveNone(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.svg")
	output := filepath.Join(dir, "output.svg")

	svgContent := `<svg viewBox="0 0 100 100">
  <path fill="none" d="M 0 0 L 10 10"/>
  <path fill="#ff0000" d="M 20 20 L 30 30"/>
</svg>`

	if err := os.WriteFile(input, []byte(svgContent), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := SVG(input, output, Options{Color: "ffffff"})
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if !contains(contentStr, `fill="none"`) {
		t.Error("fill=\"none\" should be preserved")
	}
}

func TestSVGIncludeStroke(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.svg")
	output := filepath.Join(dir, "output.svg")

	svgContent := `<svg viewBox="0 0 100 100">
  <path fill="#ff0000" stroke="#00ff00" d="M 0 0 L 10 10"/>
</svg>`

	if err := os.WriteFile(input, []byte(svgContent), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := SVG(input, output, Options{Color: "ffffff", IncludeStroke: true})
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if contains(contentStr, `stroke="#00ff00"`) {
		t.Error("stroke should have been converted")
	}
}

func TestSVGNoStrokeByDefault(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.svg")
	output := filepath.Join(dir, "output.svg")

	svgContent := `<svg viewBox="0 0 100 100">
  <path fill="#ff0000" stroke="#00ff00" d="M 0 0 L 10 10"/>
</svg>`

	if err := os.WriteFile(input, []byte(svgContent), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := SVG(input, output, Options{Color: "ffffff", IncludeStroke: false})
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if !contains(contentStr, `stroke="#00ff00"`) {
		t.Error("stroke should be preserved when IncludeStroke is false")
	}
}

func TestSVGRemoveBackground(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.svg")
	output := filepath.Join(dir, "output.svg")

	svgContent := `<svg viewBox="0 0 100 100">
  <rect x="0" y="0" width="100" height="100" fill="#000000"/>
  <path fill="#ff0000" d="M 10 10 L 90 90"/>
</svg>`

	if err := os.WriteFile(input, []byte(svgContent), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(input, output, Options{RemoveBackground: true})
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}
	if !result.BackgroundRemoved {
		t.Error("expected BackgroundRemoved = true")
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if contains(contentStr, `width="100" height="100"`) {
		t.Error("background rect should have been removed")
	}
	if !contains(contentStr, `<path`) {
		t.Error("non-background path should be preserved")
	}
}

func TestSVGPreserveMasks(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.svg")
	output := filepath.Join(dir, "output.svg")

	svgContent := `<svg viewBox="0 0 100 100">
  <mask id="m1"><rect fill="#ffffff" width="100" height="100"/></mask>
  <path fill="#ff0000" d="M 10 10 L 90 90"/>
</svg>`

	if err := os.WriteFile(input, []byte(svgContent), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := SVG(input, output, Options{Color: "000000", PreserveMasks: true})
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	// The mask content should be preserved (still #ffffff)
	if !contains(contentStr, `fill="#ffffff"`) {
		t.Error("mask fill should be preserved")
	}
	// The path fill should be converted
	if contains(contentStr, `fill="#ff0000"`) {
		t.Error("path fill should have been converted")
	}
}

func TestSVGNoColor(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.svg")
	output := filepath.Join(dir, "output.svg")

	svgContent := `<svg viewBox="0 0 100 100"><path fill="#ff0000" d="M 0 0 L 10 10"/></svg>`

	if err := os.WriteFile(input, []byte(svgContent), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(input, output, Options{})
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}
	if !result.Converted {
		t.Error("expected Converted = true (copy)")
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != svgContent {
		t.Error("without color option, output should be a copy of input")
	}
}

func TestSVGStyleAttribute(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.svg")
	output := filepath.Join(dir, "output.svg")

	svgContent := `<svg viewBox="0 0 100 100">
  <path style="fill:#ff0000;opacity:1" d="M 0 0 L 10 10"/>
</svg>`

	if err := os.WriteFile(input, []byte(svgContent), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := SVG(input, output, Options{Color: "ffffff"})
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if contains(contentStr, "#ff0000") {
		t.Error("style fill color should have been converted")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
