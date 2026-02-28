package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSVGSecure(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <path d="M 10 10 L 90 10 L 90 90 Z" fill="#ffffff"/>
  <circle cx="50" cy="50" r="20" fill="#000000"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if !result.IsSuccess() {
		t.Errorf("expected success, got threats: %v", result.Threats)
	}
	if !result.IsSecure {
		t.Error("expected IsSecure = true")
	}
	if len(result.Threats) != 0 {
		t.Errorf("expected no threats, got %d", len(result.Threats))
	}
}

func TestSVGScriptElement(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <script>alert('XSS')</script>
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
		t.Error("expected failure for script element")
	}
	if result.ThreatCounts[ThreatScript] == 0 {
		t.Error("expected ThreatScript count > 0")
	}
}

func TestSVGJavascriptURI(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <a href="javascript:alert('XSS')">
    <rect x="0" y="0" width="100" height="100"/>
  </a>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for javascript: URI")
	}
	if result.ThreatCounts[ThreatScript] == 0 {
		t.Error("expected ThreatScript count > 0")
	}
}

func TestSVGEventHandler(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg" onload="alert('XSS')">
  <rect x="0" y="0" width="100" height="100" onclick="doEvil()"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for event handlers")
	}
	if result.ThreatCounts[ThreatEventHandler] < 2 {
		t.Errorf("expected at least 2 event handler threats, got %d", result.ThreatCounts[ThreatEventHandler])
	}
}

func TestSVGExternalHref(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <image href="https://evil.com/tracker.png" width="100" height="100"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for external href")
	}
	if result.ThreatCounts[ThreatExternalRef] == 0 {
		t.Error("expected ThreatExternalRef count > 0")
	}
}

func TestSVGForeignObject(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <foreignObject x="0" y="0" width="100" height="100">
    <div xmlns="http://www.w3.org/1999/xhtml">
      <script>alert('XSS')</script>
    </div>
  </foreignObject>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for foreignObject")
	}
	if result.ThreatCounts[ThreatExternalRef] == 0 {
		t.Error("expected ThreatExternalRef count > 0 for foreignObject")
	}
}

func TestSVGExternalURLInStyle(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="100" height="100" style="fill: url(https://evil.com/tracker)"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for external URL in style")
	}
	if result.ThreatCounts[ThreatExternalRef] == 0 {
		t.Error("expected ThreatExternalRef count > 0")
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

	secure := `<?xml version="1.0"?><svg viewBox="0 0 10 10" xmlns="http://www.w3.org/2000/svg"><path d="M0 0L10 10"/></svg>`
	insecure := `<?xml version="1.0"?><svg viewBox="0 0 10 10" xmlns="http://www.w3.org/2000/svg" onclick="alert('XSS')"><path d="M0 0L10 10"/></svg>`

	if err := os.WriteFile(filepath.Join(dir, "secure.svg"), []byte(secure), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "insecure.svg"), []byte(insecure), 0600); err != nil {
		t.Fatal(err)
	}

	results, err := Directory(dir)
	if err != nil {
		t.Fatalf("Directory error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	secureCount := 0
	for _, r := range results {
		if r.IsSuccess() {
			secureCount++
		}
	}
	if secureCount != 1 {
		t.Errorf("got %d secure files, want 1", secureCount)
	}
}

func TestIsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		result   Result
		expected bool
	}{
		{"secure no errors", Result{IsSecure: true, Errors: nil}, true},
		{"secure with errors", Result{IsSecure: true, Errors: []string{"error"}}, false},
		{"insecure no errors", Result{IsSecure: false, Errors: nil}, false},
		{"insecure with errors", Result{IsSecure: false, Errors: []string{"error"}}, false},
	}

	for _, tt := range tests {
		if got := tt.result.IsSuccess(); got != tt.expected {
			t.Errorf("%s: IsSuccess() = %v, want %v", tt.name, got, tt.expected)
		}
	}
}

func TestThreatTypeString(t *testing.T) {
	tests := []struct {
		threatType ThreatType
		expected   string
	}{
		{ThreatScript, "script"},
		{ThreatEventHandler, "event_handler"},
		{ThreatExternalRef, "external_ref"},
		{ThreatAnimation, "animation"},
		{ThreatStyleBlock, "style_block"},
		{ThreatLink, "link"},
		{ThreatXMLEntity, "xml_entity"},
		{ThreatType(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.threatType.String(); got != tt.expected {
			t.Errorf("ThreatType(%d).String() = %q, want %q", tt.threatType, got, tt.expected)
		}
	}
}

func TestThreatTypeSeverity(t *testing.T) {
	tests := []struct {
		threatType ThreatType
		expected   string
	}{
		{ThreatScript, "critical"},
		{ThreatEventHandler, "critical"},
		{ThreatExternalRef, "high"},
		{ThreatXMLEntity, "high"},
		{ThreatAnimation, "medium"},
		{ThreatLink, "medium"},
		{ThreatStyleBlock, "low"},
	}

	for _, tt := range tests {
		if got := tt.threatType.Severity(); got != tt.expected {
			t.Errorf("ThreatType(%d).Severity() = %q, want %q", tt.threatType, got, tt.expected)
		}
	}
}

func TestSVGAnimation(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="100" height="100">
    <animate attributeName="width" from="0" to="100" dur="1s"/>
  </rect>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for animation element")
	}
	if result.ThreatCounts[ThreatAnimation] == 0 {
		t.Error("expected ThreatAnimation count > 0")
	}
}

func TestSVGStyleBlock(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <style>.cls-1 { fill: red; }</style>
  <rect class="cls-1" x="0" y="0" width="100" height="100"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for style block")
	}
	if result.ThreatCounts[ThreatStyleBlock] == 0 {
		t.Error("expected ThreatStyleBlock count > 0")
	}
}

func TestSVGXMLEntity(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	content := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="100" height="100"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := SVG(file)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure for DOCTYPE declaration")
	}
	if result.ThreatCounts[ThreatXMLEntity] == 0 {
		t.Error("expected ThreatXMLEntity count > 0")
	}
}

func TestScanLevelStandard(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.svg")

	// This has a style block which is LOW severity - not detected in standard mode
	content := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <style>.cls-1 { fill: red; }</style>
  <rect class="cls-1" x="0" y="0" width="100" height="100"/>
</svg>`

	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	// Standard level should not detect style blocks
	result, err := SVGWithLevel(file, ScanLevelStandard)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if !result.IsSuccess() {
		t.Error("expected success in standard mode for style block only")
	}

	// Strict level should detect style blocks
	result, err = SVGWithLevel(file, ScanLevelStrict)
	if err != nil {
		t.Fatalf("SVG error: %v", err)
	}

	if result.IsSuccess() {
		t.Error("expected failure in strict mode for style block")
	}
}

// Sanitization tests

func TestSanitizeScripts(t *testing.T) {
	content := `<?xml version="1.0"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <script>alert('XSS')</script>
  <a href="javascript:alert('XSS')">
    <rect x="0" y="0" width="100" height="100"/>
  </a>
  <path d="M 0 0 L 10 10"/>
</svg>`

	sanitized, threats := SanitizeContent(content, SanitizeOptions{RemoveScripts: true})

	if len(threats) != 2 {
		t.Errorf("expected 2 threats removed, got %d", len(threats))
	}

	if strings.Contains(sanitized, "<script>") {
		t.Error("sanitized content should not contain <script>")
	}
	if strings.Contains(sanitized, "javascript:") {
		t.Error("sanitized content should not contain javascript:")
	}
	if !strings.Contains(sanitized, "<path") {
		t.Error("sanitized content should still contain path element")
	}
}

func TestSanitizeEventHandlers(t *testing.T) {
	content := `<?xml version="1.0"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg" onload="alert('XSS')">
  <rect x="0" y="0" width="100" height="100" onclick="doEvil()" onmouseover="track()"/>
</svg>`

	sanitized, threats := SanitizeContent(content, SanitizeOptions{RemoveEventHandlers: true})

	if len(threats) != 3 {
		t.Errorf("expected 3 threats removed, got %d", len(threats))
	}

	if strings.Contains(sanitized, "onload=") {
		t.Error("sanitized content should not contain onload=")
	}
	if strings.Contains(sanitized, "onclick=") {
		t.Error("sanitized content should not contain onclick=")
	}
	if strings.Contains(sanitized, "onmouseover=") {
		t.Error("sanitized content should not contain onmouseover=")
	}
	if !strings.Contains(sanitized, "<rect") {
		t.Error("sanitized content should still contain rect element")
	}
}

func TestSanitizeExternalRefs(t *testing.T) {
	content := `<?xml version="1.0"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <image href="https://evil.com/tracker.png" width="100" height="100"/>
  <foreignObject x="0" y="0" width="100" height="100">
    <div>HTML content</div>
  </foreignObject>
  <rect style="fill: url(http://evil.com/track)"/>
</svg>`

	sanitized, threats := SanitizeContent(content, SanitizeOptions{RemoveExternalRefs: true})

	if len(threats) < 3 {
		t.Errorf("expected at least 3 threats removed, got %d", len(threats))
	}

	if strings.Contains(sanitized, "https://evil.com") {
		t.Error("sanitized content should not contain external URLs")
	}
	if strings.Contains(sanitized, "<foreignObject") {
		t.Error("sanitized content should not contain foreignObject")
	}
	if strings.Contains(sanitized, "http://evil.com") {
		t.Error("sanitized content should not contain external URLs in style")
	}
}

func TestSanitizeAll(t *testing.T) {
	content := `<?xml version="1.0"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg" onload="evil()">
  <script>alert('XSS')</script>
  <a href="javascript:void(0)">
    <rect x="0" y="0" width="100" height="100"/>
  </a>
  <image href="https://evil.com/tracker.png"/>
  <path d="M 0 0 L 10 10"/>
</svg>`

	sanitized, threats := SanitizeContent(content, DefaultSanitizeOptions())

	if len(threats) < 4 {
		t.Errorf("expected at least 4 threats removed, got %d", len(threats))
	}

	// Verify output is still valid SVG structure
	if !strings.Contains(sanitized, "<svg") {
		t.Error("sanitized content should still contain svg element")
	}
	if !strings.Contains(sanitized, "<path") {
		t.Error("sanitized content should still contain path element")
	}

	// Verify threats are removed
	if strings.Contains(sanitized, "<script>") {
		t.Error("sanitized content should not contain script")
	}
	if strings.Contains(sanitized, "javascript:") {
		t.Error("sanitized content should not contain javascript:")
	}
	if strings.Contains(sanitized, "onload=") {
		t.Error("sanitized content should not contain event handlers")
	}
	if strings.Contains(sanitized, "https://evil.com") {
		t.Error("sanitized content should not contain external URLs")
	}
}

func TestSanitizeFile(t *testing.T) {
	dir := t.TempDir()
	inputFile := filepath.Join(dir, "input.svg")
	outputFile := filepath.Join(dir, "output.svg")

	content := `<?xml version="1.0"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg" onclick="evil()">
  <script>alert('XSS')</script>
  <path d="M 0 0 L 10 10"/>
</svg>`

	if err := os.WriteFile(inputFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := Sanitize(inputFile, outputFile, DefaultSanitizeOptions())
	if err != nil {
		t.Fatalf("Sanitize error: %v", err)
	}

	if !result.Sanitized {
		t.Error("expected Sanitized = true")
	}
	if len(result.ThreatsRemoved) < 2 {
		t.Errorf("expected at least 2 threats removed, got %d", len(result.ThreatsRemoved))
	}

	// Verify output file exists and is valid
	outputContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if strings.Contains(string(outputContent), "<script>") {
		t.Error("output file should not contain script")
	}
	if !strings.Contains(string(outputContent), "<path") {
		t.Error("output file should still contain path element")
	}

	// Verify the sanitized file passes security scan
	scanResult, err := SVG(outputFile)
	if err != nil {
		t.Fatalf("failed to scan output file: %v", err)
	}
	if !scanResult.IsSuccess() {
		t.Errorf("sanitized file should pass security scan, got threats: %v", scanResult.Threats)
	}
}

func TestSanitizeNoChanges(t *testing.T) {
	dir := t.TempDir()
	inputFile := filepath.Join(dir, "input.svg")
	outputFile := filepath.Join(dir, "output.svg")

	content := `<?xml version="1.0"?>
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <path d="M 0 0 L 10 10"/>
</svg>`

	if err := os.WriteFile(inputFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := Sanitize(inputFile, outputFile, DefaultSanitizeOptions())
	if err != nil {
		t.Fatalf("Sanitize error: %v", err)
	}

	if result.Sanitized {
		t.Error("expected Sanitized = false for clean file")
	}
	if len(result.ThreatsRemoved) != 0 {
		t.Errorf("expected 0 threats removed, got %d", len(result.ThreatsRemoved))
	}
}
