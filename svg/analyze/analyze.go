// Package analyze provides SVG analysis for centering and padding.
package analyze

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/JoshVarga/svgparser"

	"github.com/grokify/brandkit/svg"
)

// Result contains the analysis of an SVG file.
type Result struct {
	FilePath         string
	ViewBox          svg.ViewBox
	ContentBox       svg.BoundingBox
	CenterOffsetX    float64
	CenterOffsetY    float64
	PaddingLeft      float64
	PaddingRight     float64
	PaddingTop       float64
	PaddingBottom    float64
	Assessment       string
	SuggestedViewBox string
	HasIssues        bool
}

// SVG analyzes an SVG file for centering and padding.
func SVG(filePath string) (*Result, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	svgDoc, err := svgparser.Parse(file, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SVG: %w", err)
	}

	// Get viewBox
	var viewBox svg.ViewBox
	if vb, ok := svgDoc.Attributes["viewBox"]; ok {
		viewBox, err = svg.ParseViewBox(vb)
		if err != nil {
			return nil, fmt.Errorf("failed to parse viewBox: %w", err)
		}
	} else {
		// Try to use width/height
		w := svg.ParseFloat(svgDoc.Attributes["width"], 0)
		h := svg.ParseFloat(svgDoc.Attributes["height"], 0)
		if w > 0 && h > 0 {
			viewBox = svg.ViewBox{X: 0, Y: 0, Width: w, Height: h}
		} else {
			return nil, fmt.Errorf("no viewBox or width/height found")
		}
	}

	// Calculate content bounds
	contentBox := svg.NewBoundingBox()
	for _, child := range svgDoc.Children {
		// Skip defs, mask, clipPath
		if child.Name == "defs" || child.Name == "mask" || child.Name == "clipPath" {
			continue
		}
		childBox := svg.GetElementBounds(child)
		contentBox.Merge(childBox)
	}

	if !contentBox.IsValid() {
		return nil, fmt.Errorf("no parseable content found")
	}

	// Calculate center offsets
	viewBoxCenterX := viewBox.CenterX()
	viewBoxCenterY := viewBox.CenterY()
	contentCenterX := contentBox.CenterX()
	contentCenterY := contentBox.CenterY()

	centerOffsetX := contentCenterX - viewBoxCenterX
	centerOffsetY := contentCenterY - viewBoxCenterY

	// Calculate padding percentages
	paddingLeft := ((contentBox.MinX - viewBox.X) / viewBox.Width) * 100
	paddingRight := ((viewBox.X + viewBox.Width - contentBox.MaxX) / viewBox.Width) * 100
	paddingTop := ((contentBox.MinY - viewBox.Y) / viewBox.Height) * 100
	paddingBottom := ((viewBox.Y + viewBox.Height - contentBox.MaxY) / viewBox.Height) * 100

	// Generate assessment
	var issues []string
	hasIssues := false

	// Check centering (threshold: 5% of viewBox dimension)
	centerThresholdX := viewBox.Width * 0.05
	centerThresholdY := viewBox.Height * 0.05

	if math.Abs(centerOffsetX) > centerThresholdX {
		if centerOffsetX > 0 {
			issues = append(issues, fmt.Sprintf("content shifted RIGHT by %.1f%%", (centerOffsetX/viewBox.Width)*100))
		} else {
			issues = append(issues, fmt.Sprintf("content shifted LEFT by %.1f%%", (-centerOffsetX/viewBox.Width)*100))
		}
		hasIssues = true
	}

	if math.Abs(centerOffsetY) > centerThresholdY {
		if centerOffsetY > 0 {
			issues = append(issues, fmt.Sprintf("content shifted DOWN by %.1f%%", (centerOffsetY/viewBox.Height)*100))
		} else {
			issues = append(issues, fmt.Sprintf("content shifted UP by %.1f%%", (-centerOffsetY/viewBox.Height)*100))
		}
		hasIssues = true
	}

	// Check for excessive padding (more than 20%)
	if paddingLeft > 20 || paddingRight > 20 || paddingTop > 20 || paddingBottom > 20 {
		maxPadding := math.Max(math.Max(paddingLeft, paddingRight), math.Max(paddingTop, paddingBottom))
		issues = append(issues, fmt.Sprintf("excessive padding (max %.1f%%)", maxPadding))
		hasIssues = true
	}

	// Check for uneven padding (difference > 10%)
	hPaddingDiff := math.Abs(paddingLeft - paddingRight)
	vPaddingDiff := math.Abs(paddingTop - paddingBottom)
	if hPaddingDiff > 10 {
		issues = append(issues, fmt.Sprintf("uneven horizontal padding (L:%.1f%% R:%.1f%%)", paddingLeft, paddingRight))
		hasIssues = true
	}
	if vPaddingDiff > 10 {
		issues = append(issues, fmt.Sprintf("uneven vertical padding (T:%.1f%% B:%.1f%%)", paddingTop, paddingBottom))
		hasIssues = true
	}

	assessment := "OK"
	if len(issues) > 0 {
		assessment = strings.Join(issues, "; ")
	}

	// Suggest fixed viewBox (5% padding on all sides)
	suggestedViewBox := SuggestViewBox(contentBox)

	return &Result{
		FilePath:         filePath,
		ViewBox:          viewBox,
		ContentBox:       *contentBox,
		CenterOffsetX:    centerOffsetX,
		CenterOffsetY:    centerOffsetY,
		PaddingLeft:      paddingLeft,
		PaddingRight:     paddingRight,
		PaddingTop:       paddingTop,
		PaddingBottom:    paddingBottom,
		Assessment:       assessment,
		SuggestedViewBox: suggestedViewBox,
		HasIssues:        hasIssues,
	}, nil
}

// SuggestViewBox suggests a viewBox with 5% padding that centers the content.
func SuggestViewBox(contentBox *svg.BoundingBox) string {
	targetPadding := 0.05 // 5%
	contentWidth := contentBox.Width()
	contentHeight := contentBox.Height()
	newWidth := contentWidth / (1 - 2*targetPadding)
	newHeight := contentHeight / (1 - 2*targetPadding)

	// Make it square if aspect ratio is close
	aspectRatio := newWidth / newHeight
	if aspectRatio > 0.9 && aspectRatio < 1.1 {
		size := math.Max(newWidth, newHeight)
		newWidth = size
		newHeight = size
	}

	newX := contentBox.MinX - (newWidth-contentWidth)/2
	newY := contentBox.MinY - (newHeight-contentHeight)/2

	return fmt.Sprintf("%.1f %.1f %.1f %.1f", newX, newY, newWidth, newHeight)
}

// Directory analyzes all SVG files in a directory.
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
				FilePath:   filePath,
				Assessment: fmt.Sprintf("Error: %v", err),
				HasIssues:  true,
			})
			continue
		}
		results = append(results, result)
	}

	return results, nil
}
