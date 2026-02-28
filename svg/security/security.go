// Package security provides SVG security scanning to detect and remove malicious elements.
package security

import (
	"fmt"
	"os"
	"regexp"

	"github.com/grokify/brandkit/svg"
)

// ThreatType categorizes the type of security threat detected.
type ThreatType int

const (
	// ThreatScript indicates a script-based threat (script tags, javascript: URIs).
	ThreatScript ThreatType = iota
	// ThreatEventHandler indicates an event handler attribute (onclick, onload, etc).
	ThreatEventHandler
	// ThreatExternalRef indicates an external reference (http:// URLs, foreignObject).
	ThreatExternalRef
	// ThreatAnimation indicates animation elements that can trigger delayed attacks.
	ThreatAnimation
	// ThreatStyleBlock indicates a <style> element that could contain malicious CSS.
	ThreatStyleBlock
	// ThreatLink indicates an anchor element (unnecessary for static images).
	ThreatLink
	// ThreatXMLEntity indicates DOCTYPE or ENTITY declarations (XXE risk).
	ThreatXMLEntity
)

// String returns a human-readable name for the threat type.
func (t ThreatType) String() string {
	switch t {
	case ThreatScript:
		return "script"
	case ThreatEventHandler:
		return "event_handler"
	case ThreatExternalRef:
		return "external_ref"
	case ThreatAnimation:
		return "animation"
	case ThreatStyleBlock:
		return "style_block"
	case ThreatLink:
		return "link"
	case ThreatXMLEntity:
		return "xml_entity"
	default:
		return "unknown"
	}
}

// Severity returns the severity level for a threat type.
func (t ThreatType) Severity() string {
	switch t {
	case ThreatScript, ThreatEventHandler:
		return "critical"
	case ThreatExternalRef, ThreatXMLEntity:
		return "high"
	case ThreatAnimation, ThreatLink:
		return "medium"
	case ThreatStyleBlock:
		return "low"
	default:
		return "info"
	}
}

// Threat represents a detected security threat in an SVG file.
type Threat struct {
	Type        ThreatType
	Description string
	Match       string
}

// Result contains the result of scanning an SVG file for security threats.
type Result struct {
	FilePath     string
	IsSecure     bool
	Threats      []Threat
	ThreatCounts map[ThreatType]int
	Errors       []string
}

// IsSuccess returns true if the file is secure and has no errors.
func (r *Result) IsSuccess() bool {
	return r.IsSecure && len(r.Errors) == 0
}

// threatPattern defines a pattern to detect a specific security threat.
type threatPattern struct {
	pattern     *regexp.Regexp
	desc        string
	threatType  ThreatType
	matchLength int // max characters to include in match (0 = use pattern default)
}

// Script patterns detect script injection attacks.
var scriptPatterns = []threatPattern{
	{regexp.MustCompile(`(?i)<script\b[^>]*>.*?</script>`), "script element", ThreatScript, 100},
	{regexp.MustCompile(`(?i)<script\b[^>]*/>`), "self-closing script element", ThreatScript, 50},
	{regexp.MustCompile(`(?i)javascript\s*:`), "javascript: URI", ThreatScript, 30},
	{regexp.MustCompile(`(?i)vbscript\s*:`), "vbscript: URI", ThreatScript, 30},
	{regexp.MustCompile(`(?i)data\s*:\s*text/html`), "data:text/html URI", ThreatScript, 50},
}

// Event handler patterns detect inline event handlers.
var eventHandlerPatterns = []threatPattern{
	{regexp.MustCompile(`(?i)\s+on[a-z]+\s*=\s*"[^"]*"`), "event handler attribute", ThreatEventHandler, 80},
	{regexp.MustCompile(`(?i)\s+on[a-z]+\s*=\s*'[^']*'`), "event handler attribute", ThreatEventHandler, 80},
	{regexp.MustCompile(`(?i)\s+on[a-z]+\s*=\s*[^\s>"']+`), "unquoted event handler attribute", ThreatEventHandler, 60},
}

// External reference patterns detect external resource loading.
var externalRefPatterns = []threatPattern{
	{regexp.MustCompile(`(?i)href\s*=\s*["']https?://[^"']+["']`), "external href", ThreatExternalRef, 100},
	{regexp.MustCompile(`(?i)xlink:href\s*=\s*["']https?://[^"']+["']`), "external xlink:href", ThreatExternalRef, 100},
	{regexp.MustCompile(`(?i)<foreignObject\b`), "foreignObject element", ThreatExternalRef, 50},
	{regexp.MustCompile(`(?i)url\s*\(\s*["']?https?://[^)"']+`), "external URL in style", ThreatExternalRef, 100},
	// External use references (internal #id refs are OK)
	{regexp.MustCompile(`(?i)<use[^>]+xlink:href\s*=\s*["']https?://`), "external use reference", ThreatExternalRef, 100},
	{regexp.MustCompile(`(?i)<use[^>]+href\s*=\s*["']https?://`), "external use reference", ThreatExternalRef, 100},
}

// Animation patterns detect SVG animation elements.
var animationPatterns = []threatPattern{
	{regexp.MustCompile(`(?i)<animate\b`), "animate element", ThreatAnimation, 50},
	{regexp.MustCompile(`(?i)<animateTransform\b`), "animateTransform element", ThreatAnimation, 50},
	{regexp.MustCompile(`(?i)<animateMotion\b`), "animateMotion element", ThreatAnimation, 50},
	{regexp.MustCompile(`(?i)<animateColor\b`), "animateColor element", ThreatAnimation, 50},
	{regexp.MustCompile(`(?i)<set\b[^>]*\b(attributeName|to)\s*=`), "set element", ThreatAnimation, 50},
}

// Style block patterns detect <style> elements.
var styleBlockPatterns = []threatPattern{
	{regexp.MustCompile(`(?i)<style\b`), "style element", ThreatStyleBlock, 50},
}

// Link patterns detect anchor elements.
var linkPatterns = []threatPattern{
	{regexp.MustCompile(`(?i)<a\b[^>]*\bhref\s*=`), "anchor element with href", ThreatLink, 80},
}

// XML entity patterns detect DOCTYPE and ENTITY declarations.
var xmlEntityPatterns = []threatPattern{
	{regexp.MustCompile(`(?i)<!DOCTYPE\b`), "DOCTYPE declaration", ThreatXMLEntity, 50},
	{regexp.MustCompile(`(?i)<!ENTITY\b`), "ENTITY declaration", ThreatXMLEntity, 50},
}

// ScanLevel defines how strict the security scan should be.
type ScanLevel int

const (
	// ScanLevelStrict detects all threats including style blocks and animations.
	ScanLevelStrict ScanLevel = iota
	// ScanLevelStandard detects critical and high severity threats only.
	ScanLevelStandard
)

// patternsForLevel returns patterns based on scan level.
func patternsForLevel(level ScanLevel) []threatPattern {
	var all []threatPattern

	// Always include critical threats
	all = append(all, scriptPatterns...)
	all = append(all, eventHandlerPatterns...)

	// Always include high severity threats
	all = append(all, externalRefPatterns...)
	all = append(all, xmlEntityPatterns...)

	// Include medium/low severity threats only in strict mode
	if level == ScanLevelStrict {
		all = append(all, animationPatterns...)
		all = append(all, styleBlockPatterns...)
		all = append(all, linkPatterns...)
	}

	return all
}

// SVG scans a single SVG file for security threats using strict level.
func SVG(filePath string) (*Result, error) {
	return SVGWithLevel(filePath, ScanLevelStrict)
}

// SVGWithLevel scans a single SVG file with specified scan level.
func SVGWithLevel(filePath string, level ScanLevel) (*Result, error) {
	result := &Result{
		FilePath:     filePath,
		IsSecure:     true,
		Threats:      []Threat{},
		ThreatCounts: make(map[ThreatType]int),
		Errors:       []string{},
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ScanContentWithLevel(string(content), result, level), nil
}

// ScanContent scans SVG content for security threats using strict level.
func ScanContent(content string, result *Result) *Result {
	return ScanContentWithLevel(content, result, ScanLevelStrict)
}

// ScanContentWithLevel scans SVG content for security threats with specified level.
func ScanContentWithLevel(content string, result *Result, level ScanLevel) *Result {
	if result == nil {
		result = &Result{
			IsSecure:     true,
			Threats:      []Threat{},
			ThreatCounts: make(map[ThreatType]int),
			Errors:       []string{},
		}
	}

	for _, p := range patternsForLevel(level) {
		matches := p.pattern.FindAllString(content, -1)
		for _, match := range matches {
			// Truncate match for display
			displayMatch := match
			maxLen := p.matchLength
			if maxLen == 0 {
				maxLen = 50
			}
			if len(displayMatch) > maxLen {
				displayMatch = displayMatch[:maxLen] + "..."
			}

			result.Threats = append(result.Threats, Threat{
				Type:        p.threatType,
				Description: p.desc,
				Match:       displayMatch,
			})
			result.ThreatCounts[p.threatType]++
			result.IsSecure = false
		}
	}

	return result
}

// Directory scans all SVG files in a directory (non-recursive).
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
				IsSecure: false,
				Errors:   []string{err.Error()},
			})
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// DirectoryRecursive scans all SVG files in a directory tree.
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
				FilePath:     filePath,
				IsSecure:     false,
				ThreatCounts: make(map[ThreatType]int),
				Errors:       []string{err.Error()},
			})
			continue
		}
		results = append(results, result)
	}

	return results, nil
}
