package verify

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSVGPureVector(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <path d="M 10 10 L 90 10 L 90 90 Z" fill="#ffffff"/>
  <circle cx="50" cy="50" r="20" fill="#000000"/>
  <rect x="10" y="10" width="30" height="30"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if !result.IsSuccess() {
		t.Errorf("expected success, got errors: %v", result.Errors)
	}
	if !result.IsValid {
		t.Error("expected IsValid = true")
	}
	if !result.IsPureVector {
		t.Error("expected IsPureVector = true")
	}
	if len(result.VectorElements) == 0 {
		t.Error("expected vector elements to be detected")
	}
}

func TestSVGEmbeddedBase64(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <image href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUg..." width="100" height="100"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for embedded base64 image")
	}
	if result.IsPureVector {
		t.Error("expected IsPureVector = false")
	}
	if !result.HasEmbeddedData {
		t.Error("expected HasEmbeddedData = true")
	}
}

func TestSVGXlinkHrefData(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <image xlink:href="data:image/jpeg;base64,/9j/4AAQ..." width="50" height="50"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for xlink:href data URI")
	}
	if !result.HasEmbeddedData {
		t.Error("expected HasEmbeddedData = true")
	}
}

func TestSVGExternalBinaryRef(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <image xlink:href="background.png" width="100" height="100"/>
  <path d="M 0 0 L 10 10"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for external binary reference")
	}
}

func TestSVGMissingSVGElement(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<div><p>Not an SVG</p></div>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsValid {
		t.Error("expected IsValid = false for missing svg element")
	}
}

func TestSVGFileNotFound(t *testing.T) {
	_, err := SVG("/nonexistent/path.svg")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDirectory(t *testing.T) {
	dir := t.TempDir()

	good := `<?xml version="1.0"?><svg viewBox="0 0 10 10" xmlns="http://www.w3.org/2000/svg"><path d="M0 0L10 10"/></svg>`
	bad := `<?xml version="1.0"?><svg viewBox="0 0 10 10" xmlns="http://www.w3.org/2000/svg"><image href="data:image/png;base64,abc"/></svg>`

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

	successCount := 0
	for _, r := range results {
		if r.IsSuccess() {
			successCount++
		}
	}
	if successCount != 1 {
		t.Errorf("got %d successes, want 1", successCount)
	}
}

func TestIsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		result   Result
		expected bool
	}{
		{"valid pure vector", Result{IsValid: true, IsPureVector: true}, true},
		{"invalid XML", Result{IsValid: false, IsPureVector: true}, false},
		{"embedded data", Result{IsValid: true, IsPureVector: false}, false},
		{"both bad", Result{IsValid: false, IsPureVector: false}, false},
	}

	for _, tt := range tests {
		if got := tt.result.IsSuccess(); got != tt.expected {
			t.Errorf("%s: IsSuccess() = %v, want %v", tt.name, got, tt.expected)
		}
	}
}
