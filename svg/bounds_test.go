package svg

import (
	"math"
	"testing"
)

func TestNewBoundingBox(t *testing.T) {
	box := NewBoundingBox()
	if box.MinX != math.MaxFloat64 {
		t.Errorf("MinX = %v, want MaxFloat64", box.MinX)
	}
	if box.MaxX != -math.MaxFloat64 {
		t.Errorf("MaxX = %v, want -MaxFloat64", box.MaxX)
	}
	if box.IsValid() {
		t.Error("new box should not be valid")
	}
}

func TestBoundingBoxExpand(t *testing.T) {
	box := NewBoundingBox()
	box.Expand(10, 20)
	box.Expand(30, 40)

	if box.MinX != 10 {
		t.Errorf("MinX = %v, want 10", box.MinX)
	}
	if box.MinY != 20 {
		t.Errorf("MinY = %v, want 20", box.MinY)
	}
	if box.MaxX != 30 {
		t.Errorf("MaxX = %v, want 30", box.MaxX)
	}
	if box.MaxY != 40 {
		t.Errorf("MaxY = %v, want 40", box.MaxY)
	}
	if !box.IsValid() {
		t.Error("expanded box should be valid")
	}
}

func TestBoundingBoxDimensions(t *testing.T) {
	box := NewBoundingBox()
	box.Expand(10, 20)
	box.Expand(50, 80)

	if w := box.Width(); w != 40 {
		t.Errorf("Width() = %v, want 40", w)
	}
	if h := box.Height(); h != 60 {
		t.Errorf("Height() = %v, want 60", h)
	}
	if cx := box.CenterX(); cx != 30 {
		t.Errorf("CenterX() = %v, want 30", cx)
	}
	if cy := box.CenterY(); cy != 50 {
		t.Errorf("CenterY() = %v, want 50", cy)
	}
}

func TestBoundingBoxMerge(t *testing.T) {
	box1 := NewBoundingBox()
	box1.Expand(10, 10)
	box1.Expand(50, 50)

	box2 := NewBoundingBox()
	box2.Expand(0, 0)
	box2.Expand(30, 60)

	box1.Merge(box2)

	if box1.MinX != 0 {
		t.Errorf("MinX = %v, want 0", box1.MinX)
	}
	if box1.MinY != 0 {
		t.Errorf("MinY = %v, want 0", box1.MinY)
	}
	if box1.MaxX != 50 {
		t.Errorf("MaxX = %v, want 50", box1.MaxX)
	}
	if box1.MaxY != 60 {
		t.Errorf("MaxY = %v, want 60", box1.MaxY)
	}
}

func TestBoundingBoxMergeInvalid(t *testing.T) {
	box1 := NewBoundingBox()
	box1.Expand(10, 20)

	invalid := NewBoundingBox()
	box1.Merge(invalid)

	if box1.MinX != 10 {
		t.Errorf("MinX = %v, want 10 (merge of invalid should be no-op)", box1.MinX)
	}
}

func TestViewBoxCenterAndString(t *testing.T) {
	vb := ViewBox{X: 0, Y: 0, Width: 100, Height: 200}

	if cx := vb.CenterX(); cx != 50 {
		t.Errorf("CenterX() = %v, want 50", cx)
	}
	if cy := vb.CenterY(); cy != 100 {
		t.Errorf("CenterY() = %v, want 100", cy)
	}
	if s := vb.String(); s != "0.0 0.0 100.0 200.0" {
		t.Errorf("String() = %q, want %q", s, "0.0 0.0 100.0 200.0")
	}
}

func TestParseViewBox(t *testing.T) {
	tests := []struct {
		input   string
		want    ViewBox
		wantErr bool
	}{
		{"0 0 100 100", ViewBox{0, 0, 100, 100}, false},
		{"-10 -20 200 300", ViewBox{-10, -20, 200, 300}, false},
		{"0.5 1.5 99.5 199.5", ViewBox{0.5, 1.5, 99.5, 199.5}, false},
		{"0 0 100", ViewBox{}, true},     // too few parts
		{"a b c d", ViewBox{}, true},     // non-numeric
		{"0 0 100 abc", ViewBox{}, true}, // partial non-numeric
	}

	for _, tt := range tests {
		got, err := ParseViewBox(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseViewBox(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("ParseViewBox(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		input      string
		defaultVal float64
		want       float64
	}{
		{"42", 0, 42},
		{"3.14", 0, 3.14},
		{"100px", 0, 100},
		{"", 99, 99},
		{"abc", 5, 5},
	}

	for _, tt := range tests {
		got := ParseFloat(tt.input, tt.defaultVal)
		if got != tt.want {
			t.Errorf("ParseFloat(%q, %v) = %v, want %v", tt.input, tt.defaultVal, got, tt.want)
		}
	}
}
