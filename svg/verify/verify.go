// Package verify validates SVG files are pure vector images.
package verify

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/brandkit/svg"
)

// Result contains the result of validating an SVG file.
type Result struct {
	FilePath        string
	IsValid         bool
	IsPureVector    bool
	HasEmbeddedData bool
	VectorElements  []string
	Errors          []string
}

// embeddedPattern defines a pattern to detect embedded binary data.
type embeddedPattern struct {
	pattern *regexp.Regexp
	desc    string
}

var embeddedPatterns = []embeddedPattern{
	{regexp.MustCompile(`data:image/(png|jpeg|jpg|gif|webp|bmp)`), "base64 embedded image"},
	{regexp.MustCompile(`xlink:href\s*=\s*["']data:`), "xlink:href with data URI"},
	{regexp.MustCompile(`href\s*=\s*["']data:image`), "href with embedded image data"},
	{regexp.MustCompile(`<image[^>]+xlink:href\s*=\s*["'][^"']*\.(png|jpg|jpeg|gif|webp|bmp)`), "image element referencing binary file"},
}

var vectorPatterns = map[string]*regexp.Regexp{
	"path":     regexp.MustCompile(`<path\b`),
	"rect":     regexp.MustCompile(`<rect\b`),
	"circle":   regexp.MustCompile(`<circle\b`),
	"ellipse":  regexp.MustCompile(`<ellipse\b`),
	"line":     regexp.MustCompile(`<line\b`),
	"polyline": regexp.MustCompile(`<polyline\b`),
	"polygon":  regexp.MustCompile(`<polygon\b`),
	"text":     regexp.MustCompile(`<text\b`),
}

// SVG checks if an SVG file is a pure vector image without embedded binary data.
func SVG(filePath string) (*Result, error) {
	result := &Result{
		FilePath:       filePath,
		IsValid:        true,
		IsPureVector:   true,
		VectorElements: []string{},
		Errors:         []string{},
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)

	// Check for valid XML/SVG structure
	if !strings.Contains(contentStr, "<svg") {
		result.IsValid = false
		result.Errors = append(result.Errors, "missing <svg> element")
	}

	// Check for embedded binary patterns
	for _, p := range embeddedPatterns {
		if p.pattern.MatchString(contentStr) {
			result.IsPureVector = false
			result.HasEmbeddedData = true
			result.Errors = append(result.Errors, fmt.Sprintf("contains %s", p.desc))
		}
	}

	// Count vector elements
	for name, pattern := range vectorPatterns {
		matches := pattern.FindAllString(contentStr, -1)
		if len(matches) > 0 {
			result.VectorElements = append(result.VectorElements, fmt.Sprintf("%s:%d", name, len(matches)))
		}
	}

	// Verify it's valid XML
	var svgDoc any
	if err := xml.Unmarshal(content, &svgDoc); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("invalid XML: %v", err))
	}

	return result, nil
}

// Directory validates all SVG files in a directory.
func Directory(dirPath string) ([]*Result, error) {
	files, err := svg.ListSVGFiles(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var results []*Result
	for _, filePath := range files {
		result, err := SVG(filePath)
		if err != nil {
			results = append(results, &Result{
				FilePath: filePath,
				IsValid:  false,
				Errors:   []string{err.Error()},
			})
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// IsSuccess returns true if the result indicates a valid pure vector SVG.
func (r *Result) IsSuccess() bool {
	return r.IsValid && r.IsPureVector
}

// DirectoryRecursive validates all SVG files in a directory tree.
func DirectoryRecursive(dirPath string) ([]*Result, error) {
	files, err := svg.ListSVGFilesRecursive(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var results []*Result
	for _, filePath := range files {
		result, err := SVG(filePath)
		if err != nil {
			results = append(results, &Result{
				FilePath: filePath,
				IsValid:  false,
				Errors:   []string{err.Error()},
			})
			continue
		}
		results = append(results, result)
	}

	return results, nil
}
