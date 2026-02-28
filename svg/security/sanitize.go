package security

import (
	"fmt"
	"os"
	"regexp"
)

// SanitizeOptions specifies which threat types to remove during sanitization.
type SanitizeOptions struct {
	RemoveScripts       bool // Remove script elements and javascript: URIs
	RemoveEventHandlers bool // Remove on* event handler attributes
	RemoveExternalRefs  bool // Remove external URLs and foreignObject
	RemoveAll           bool // Remove all threat types (overrides individual flags)
}

// DefaultSanitizeOptions returns options that remove all threats.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		RemoveAll: true,
	}
}

// SanitizeResult contains the result of sanitizing an SVG file.
type SanitizeResult struct {
	InputPath      string
	OutputPath     string
	ThreatsRemoved []Threat
	Sanitized      bool
	Error          error
}

// sanitizePattern defines a pattern and its replacement for sanitization.
type sanitizePattern struct {
	pattern     *regexp.Regexp
	replacement string
	desc        string
	threatType  ThreatType
}

// Script removal patterns.
var scriptRemovalPatterns = []sanitizePattern{
	// Remove <script>...</script> elements
	{regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script>`), "", "script element", ThreatScript},
	// Remove self-closing <script/> elements
	{regexp.MustCompile(`(?i)<script\b[^>]*/>`), "", "self-closing script element", ThreatScript},
	// Remove javascript: URIs in href attributes - replace with empty href
	{regexp.MustCompile(`(?i)(href\s*=\s*["'])javascript:[^"']*["']`), `$1#"`, "javascript: URI in href", ThreatScript},
	// Remove vbscript: URIs in href attributes
	{regexp.MustCompile(`(?i)(href\s*=\s*["'])vbscript:[^"']*["']`), `$1#"`, "vbscript: URI in href", ThreatScript},
	// Remove data:text/html URIs in href attributes
	{regexp.MustCompile(`(?i)(href\s*=\s*["'])data:\s*text/html[^"']*["']`), `$1#"`, "data:text/html URI", ThreatScript},
}

// Event handler removal patterns.
var eventHandlerRemovalPatterns = []sanitizePattern{
	// Remove on* event handler attributes (double-quoted values)
	{regexp.MustCompile(`(?i)\s+on[a-z]+\s*=\s*"[^"]*"`), "", "event handler attribute", ThreatEventHandler},
	// Remove on* event handler attributes (single-quoted values)
	{regexp.MustCompile(`(?i)\s+on[a-z]+\s*=\s*'[^']*'`), "", "event handler attribute", ThreatEventHandler},
	// Remove on* event handler attributes (unquoted values)
	{regexp.MustCompile(`(?i)\s+on[a-z]+\s*=\s*[^\s>"']+`), "", "unquoted event handler attribute", ThreatEventHandler},
}

// XML entity removal patterns.
var xmlEntityRemovalPatterns = []sanitizePattern{
	// Remove DOCTYPE declarations (entire line)
	{regexp.MustCompile(`(?i)<!DOCTYPE[^>]*>`), "", "DOCTYPE declaration", ThreatXMLEntity},
	// Remove ENTITY declarations
	{regexp.MustCompile(`(?i)<!ENTITY[^>]*>`), "", "ENTITY declaration", ThreatXMLEntity},
}

// External reference removal patterns.
var externalRefRemovalPatterns = []sanitizePattern{
	// Replace external href with empty
	{regexp.MustCompile(`(?i)(href\s*=\s*["'])https?://[^"']*["']`), `$1#"`, "external href", ThreatExternalRef},
	// Replace external xlink:href with empty
	{regexp.MustCompile(`(?i)(xlink:href\s*=\s*["'])https?://[^"']*["']`), `$1#"`, "external xlink:href", ThreatExternalRef},
	// Remove foreignObject elements entirely
	{regexp.MustCompile(`(?is)<foreignObject\b[^>]*>.*?</foreignObject>`), "", "foreignObject element", ThreatExternalRef},
	// Remove self-closing foreignObject
	{regexp.MustCompile(`(?i)<foreignObject\b[^>]*/>`), "", "self-closing foreignObject", ThreatExternalRef},
	// Replace external URLs in style url() with none
	{regexp.MustCompile(`(?i)(url\s*\(\s*["']?)https?://[^)"']+([)"']?)`), "${1}none${2}", "external URL in style", ThreatExternalRef},
}

// Sanitize removes security threats from an SVG file and writes the result.
func Sanitize(inputPath, outputPath string, opts SanitizeOptions) (*SanitizeResult, error) {
	result := &SanitizeResult{
		InputPath:      inputPath,
		OutputPath:     outputPath,
		ThreatsRemoved: []Threat{},
		Sanitized:      false,
	}

	content, err := os.ReadFile(inputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to read input file: %w", err)
		return result, result.Error
	}

	sanitized, threats := SanitizeContent(string(content), opts)
	result.ThreatsRemoved = threats
	result.Sanitized = len(threats) > 0

	if err := os.WriteFile(outputPath, []byte(sanitized), 0600); err != nil {
		result.Error = fmt.Errorf("failed to write output file: %w", err)
		return result, result.Error
	}

	return result, nil
}

// SanitizeContent removes security threats from SVG content in memory.
func SanitizeContent(content string, opts SanitizeOptions) (string, []Threat) {
	var threats []Threat
	sanitized := content

	// Collect all patterns to apply based on options
	var patterns []sanitizePattern
	if opts.RemoveAll || opts.RemoveScripts {
		patterns = append(patterns, scriptRemovalPatterns...)
	}
	if opts.RemoveAll || opts.RemoveEventHandlers {
		patterns = append(patterns, eventHandlerRemovalPatterns...)
	}
	if opts.RemoveAll || opts.RemoveExternalRefs {
		patterns = append(patterns, externalRefRemovalPatterns...)
	}
	if opts.RemoveAll {
		patterns = append(patterns, xmlEntityRemovalPatterns...)
	}

	// Apply each pattern
	for _, p := range patterns {
		matches := p.pattern.FindAllString(sanitized, -1)
		for _, match := range matches {
			displayMatch := match
			if len(displayMatch) > 80 {
				displayMatch = displayMatch[:80] + "..."
			}
			threats = append(threats, Threat{
				Type:        p.threatType,
				Description: p.desc,
				Match:       displayMatch,
			})
		}
		sanitized = p.pattern.ReplaceAllString(sanitized, p.replacement)
	}

	return sanitized, threats
}
