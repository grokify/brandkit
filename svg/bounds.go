// Package svg provides shared types and utilities for SVG processing.
package svg

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// BoundingBox represents a rectangular bounding box.
type BoundingBox struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

// NewBoundingBox creates an empty bounding box.
func NewBoundingBox() *BoundingBox {
	return &BoundingBox{
		MinX: math.MaxFloat64,
		MinY: math.MaxFloat64,
		MaxX: -math.MaxFloat64,
		MaxY: -math.MaxFloat64,
	}
}

// Width returns the width of the bounding box.
func (b *BoundingBox) Width() float64 {
	return b.MaxX - b.MinX
}

// Height returns the height of the bounding box.
func (b *BoundingBox) Height() float64 {
	return b.MaxY - b.MinY
}

// CenterX returns the X coordinate of the center.
func (b *BoundingBox) CenterX() float64 {
	return (b.MinX + b.MaxX) / 2
}

// CenterY returns the Y coordinate of the center.
func (b *BoundingBox) CenterY() float64 {
	return (b.MinY + b.MaxY) / 2
}

// IsValid returns true if the bounding box has been expanded with at least one point.
func (b *BoundingBox) IsValid() bool {
	return b.MinX != math.MaxFloat64 && b.MaxX != -math.MaxFloat64
}

// Expand expands the bounding box to include the given point.
func (b *BoundingBox) Expand(x, y float64) {
	if x < b.MinX {
		b.MinX = x
	}
	if x > b.MaxX {
		b.MaxX = x
	}
	if y < b.MinY {
		b.MinY = y
	}
	if y > b.MaxY {
		b.MaxY = y
	}
}

// Merge merges another bounding box into this one.
func (b *BoundingBox) Merge(other *BoundingBox) {
	if !other.IsValid() {
		return
	}
	b.Expand(other.MinX, other.MinY)
	b.Expand(other.MaxX, other.MaxY)
}

// ViewBox represents an SVG viewBox.
type ViewBox struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// CenterX returns the X coordinate of the viewBox center.
func (v *ViewBox) CenterX() float64 {
	return v.X + v.Width/2
}

// CenterY returns the Y coordinate of the viewBox center.
func (v *ViewBox) CenterY() float64 {
	return v.Y + v.Height/2
}

// String returns the viewBox as a string suitable for SVG attribute.
func (v *ViewBox) String() string {
	return fmt.Sprintf("%.1f %.1f %.1f %.1f", v.X, v.Y, v.Width, v.Height)
}

// ParseViewBox parses a viewBox string like "0 0 100 100".
func ParseViewBox(s string) (ViewBox, error) {
	parts := strings.Fields(s)
	if len(parts) != 4 {
		return ViewBox{}, fmt.Errorf("invalid viewBox format: %s", s)
	}

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return ViewBox{}, err
	}
	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return ViewBox{}, err
	}
	w, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return ViewBox{}, err
	}
	h, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return ViewBox{}, err
	}

	return ViewBox{X: x, Y: y, Width: w, Height: h}, nil
}

// ParseFloat parses a float with a default value on error.
func ParseFloat(s string, defaultVal float64) float64 {
	if s == "" {
		return defaultVal
	}
	// Remove "px" suffix if present
	s = strings.TrimSuffix(s, "px")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}
	return v
}
