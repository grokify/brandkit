// Package brandkit provides embedded brand icons and SVG processing utilities.
package brandkit

import (
	"fmt"
	"os"
	"regexp"

	"github.com/grokify/brandkit/svg/analyze"
	"github.com/grokify/brandkit/svg/convert"
	"github.com/grokify/brandkit/svg/verify"
)

// ProcessResult contains the result of a processing operation.
type ProcessResult struct {
	InputPath         string
	OutputPath        string
	BackgroundRemoved bool
	ColorConverted    bool
	TargetColor       string
	Centered          bool
	SuggestedViewBox  string
	Verified          bool
	VectorElements    []string
}

// ProcessWhite creates a white icon on transparent background.
// It removes background elements, converts all colors to white,
// centers the content, and verifies the result is pure vector.
//
// Equivalent to CLI: brandkit white <input> -o <output>
func ProcessWhite(inputPath, outputPath string) (*ProcessResult, error) {
	return process(inputPath, outputPath, processOptions{
		color:            "ffffff",
		removeBackground: true,
		includeStroke:    true,
		center:           true,
		strict:           true,
	})
}

// ProcessColor creates a centered color icon on transparent background.
// It removes background elements, centers the content, and verifies
// the result is pure vector while preserving original colors.
//
// Equivalent to CLI: brandkit color <input> -o <output>
func ProcessColor(inputPath, outputPath string) (*ProcessResult, error) {
	return process(inputPath, outputPath, processOptions{
		color:            "", // No color conversion - keep originals
		removeBackground: true,
		includeStroke:    false, // Irrelevant since color is empty (no conversion happens)
		center:           true,
		strict:           true,
	})
}

type processOptions struct {
	color            string
	removeBackground bool
	includeStroke    bool
	center           bool
	strict           bool
}

func process(inputPath, outputPath string, opts processOptions) (*ProcessResult, error) {
	result := &ProcessResult{
		InputPath:  inputPath,
		OutputPath: outputPath,
	}

	// Step 1: Convert colors (to a temp file if we need to modify viewBox)
	tempOutput := outputPath
	if opts.center {
		tempOutput = outputPath + ".tmp"
	}

	convertOpts := convert.Options{
		Color:            opts.color,
		IncludeStroke:    opts.includeStroke,
		PreserveMasks:    true,
		RemoveBackground: opts.removeBackground,
	}

	convertResult, err := convert.SVG(inputPath, tempOutput, convertOpts)
	if err != nil {
		return result, fmt.Errorf("conversion failed: %w", err)
	}

	result.BackgroundRemoved = convertResult.BackgroundRemoved
	if convertResult.TargetColor != "" {
		result.ColorConverted = true
		result.TargetColor = convertResult.TargetColor
	}

	// Step 2: Analyze (and optionally fix centering)
	analysisResult, err := analyze.SVG(tempOutput)
	if err != nil {
		if opts.center {
			_ = os.Remove(tempOutput)
		}
		return result, fmt.Errorf("analysis failed: %w", err)
	}

	if opts.center && analysisResult.HasIssues {
		// Apply the suggested viewBox fix
		content, err := os.ReadFile(tempOutput)
		if err != nil {
			_ = os.Remove(tempOutput)
			return result, fmt.Errorf("failed to read for centering: %w", err)
		}

		contentStr := string(content)

		// Replace viewBox with suggested value
		viewBoxRe := regexp.MustCompile(`viewBox\s*=\s*["'][^"']*["']`)
		newViewBox := fmt.Sprintf(`viewBox="%s"`, analysisResult.SuggestedViewBox)

		if viewBoxRe.MatchString(contentStr) {
			contentStr = viewBoxRe.ReplaceAllString(contentStr, newViewBox)
		}

		if err := os.WriteFile(outputPath, []byte(contentStr), 0600); err != nil {
			_ = os.Remove(tempOutput)
			return result, fmt.Errorf("failed to write centered file: %w", err)
		}

		if tempOutput != outputPath {
			_ = os.Remove(tempOutput)
		}

		result.Centered = true
		result.SuggestedViewBox = analysisResult.SuggestedViewBox
	} else if opts.center {
		// No issues, just rename temp to final
		if tempOutput != outputPath {
			if err := os.Rename(tempOutput, outputPath); err != nil {
				return result, fmt.Errorf("failed to finalize output: %w", err)
			}
		}
	}

	// Step 3: Verify (if strict mode)
	if opts.strict {
		verifyResult, err := verify.SVG(outputPath)
		if err != nil {
			return result, fmt.Errorf("verification failed: %w", err)
		}

		if !verifyResult.IsSuccess() {
			return result, fmt.Errorf("SVG contains embedded binary data: %v", verifyResult.Errors)
		}

		result.Verified = true
		result.VectorElements = verifyResult.VectorElements
	}

	return result, nil
}
