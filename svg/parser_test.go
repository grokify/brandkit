package svg

import (
	"testing"
)

func TestParsePathSimpleMoveTo(t *testing.T) {
	commands := ParsePath("M 10 20 L 30 40")
	if len(commands) != 2 {
		t.Fatalf("got %d commands, want 2", len(commands))
	}
	if commands[0].Command != 'M' {
		t.Errorf("commands[0].Command = %c, want M", commands[0].Command)
	}
	if len(commands[0].Params) != 2 || commands[0].Params[0] != 10 || commands[0].Params[1] != 20 {
		t.Errorf("commands[0].Params = %v, want [10 20]", commands[0].Params)
	}
	if commands[1].Command != 'L' {
		t.Errorf("commands[1].Command = %c, want L", commands[1].Command)
	}
	if len(commands[1].Params) != 2 || commands[1].Params[0] != 30 || commands[1].Params[1] != 40 {
		t.Errorf("commands[1].Params = %v, want [30 40]", commands[1].Params)
	}
}

func TestParsePathNegativeAndDecimal(t *testing.T) {
	commands := ParsePath("M-10.5-20.3L30.7 40.9")
	if len(commands) != 2 {
		t.Fatalf("got %d commands, want 2", len(commands))
	}
	if commands[0].Params[0] != -10.5 || commands[0].Params[1] != -20.3 {
		t.Errorf("M params = %v, want [-10.5 -20.3]", commands[0].Params)
	}
}

func TestParsePathMultipleCoordinatePairs(t *testing.T) {
	commands := ParsePath("M 0 0 10 20 30 40")
	if len(commands) != 1 {
		t.Fatalf("got %d commands, want 1", len(commands))
	}
	if len(commands[0].Params) != 6 {
		t.Errorf("got %d params, want 6", len(commands[0].Params))
	}
}

func TestParsePathClosePath(t *testing.T) {
	commands := ParsePath("M 0 0 L 10 0 L 10 10 Z")
	if len(commands) != 4 {
		t.Fatalf("got %d commands, want 4", len(commands))
	}
	if commands[3].Command != 'Z' {
		t.Errorf("last command = %c, want Z", commands[3].Command)
	}
}

func TestParsePathCubicBezier(t *testing.T) {
	commands := ParsePath("M 0 0 C 10 20 30 40 50 60")
	if len(commands) != 2 {
		t.Fatalf("got %d commands, want 2", len(commands))
	}
	if commands[1].Command != 'C' {
		t.Errorf("commands[1].Command = %c, want C", commands[1].Command)
	}
	if len(commands[1].Params) != 6 {
		t.Errorf("C params count = %d, want 6", len(commands[1].Params))
	}
}

func TestParsePathArc(t *testing.T) {
	commands := ParsePath("M 0 0 A 25 25 0 1 1 50 50")
	if len(commands) != 2 {
		t.Fatalf("got %d commands, want 2", len(commands))
	}
	if commands[1].Command != 'A' {
		t.Errorf("commands[1].Command = %c, want A", commands[1].Command)
	}
	if len(commands[1].Params) != 7 {
		t.Errorf("A params count = %d, want 7", len(commands[1].Params))
	}
}

func TestCalculatePathBoundsSquare(t *testing.T) {
	// A simple square: 0,0 to 100,100
	box := CalculatePathBounds("M 0 0 L 100 0 L 100 100 L 0 100 Z")
	if !box.IsValid() {
		t.Fatal("box should be valid")
	}
	if box.MinX != 0 || box.MinY != 0 {
		t.Errorf("min = (%v, %v), want (0, 0)", box.MinX, box.MinY)
	}
	if box.MaxX != 100 || box.MaxY != 100 {
		t.Errorf("max = (%v, %v), want (100, 100)", box.MaxX, box.MaxY)
	}
}

func TestCalculatePathBoundsRelative(t *testing.T) {
	// Start at 10,10, relative line to +20,+30
	box := CalculatePathBounds("M 10 10 l 20 30")
	if box.MinX != 10 || box.MinY != 10 {
		t.Errorf("min = (%v, %v), want (10, 10)", box.MinX, box.MinY)
	}
	if box.MaxX != 30 || box.MaxY != 40 {
		t.Errorf("max = (%v, %v), want (30, 40)", box.MaxX, box.MaxY)
	}
}

func TestCalculatePathBoundsHorizontalVertical(t *testing.T) {
	box := CalculatePathBounds("M 0 0 H 50 V 80")
	if box.MaxX != 50 {
		t.Errorf("MaxX = %v, want 50", box.MaxX)
	}
	if box.MaxY != 80 {
		t.Errorf("MaxY = %v, want 80", box.MaxY)
	}
}

func TestCalculatePathBoundsRelativeHV(t *testing.T) {
	box := CalculatePathBounds("M 10 20 h 40 v 60")
	if box.MaxX != 50 {
		t.Errorf("MaxX = %v, want 50", box.MaxX)
	}
	if box.MaxY != 80 {
		t.Errorf("MaxY = %v, want 80", box.MaxY)
	}
}

func TestCalculatePathBoundsRelativeMoveTo(t *testing.T) {
	box := CalculatePathBounds("M 10 10 m 5 5 l 10 10")
	// After M 10 10, m 5 5 -> 15, 15, then l 10 10 -> 25, 25
	if box.MinX != 10 || box.MinY != 10 {
		t.Errorf("min = (%v, %v), want (10, 10)", box.MinX, box.MinY)
	}
	if box.MaxX != 25 || box.MaxY != 25 {
		t.Errorf("max = (%v, %v), want (25, 25)", box.MaxX, box.MaxY)
	}
}

func TestCalculatePathBoundsCubicBezier(t *testing.T) {
	// Control points extend the bounding box
	box := CalculatePathBounds("M 0 0 C 50 -20 80 120 100 100")
	if box.MinY >= 0 {
		t.Errorf("MinY = %v, should be negative (control point at y=-20)", box.MinY)
	}
	if box.MaxY <= 100 {
		t.Errorf("MaxY = %v, should be >100 (control point at y=120)", box.MaxY)
	}
}

func TestCalculatePathBoundsArc(t *testing.T) {
	box := CalculatePathBounds("M 0 0 A 25 25 0 1 1 50 50")
	if box.MaxX != 50 || box.MaxY != 50 {
		t.Errorf("max = (%v, %v), want (50, 50)", box.MaxX, box.MaxY)
	}
}

func TestCalculatePathBoundsClosePath(t *testing.T) {
	// After Z, cursor returns to start
	box := CalculatePathBounds("M 10 20 L 50 60 Z L 100 100")
	// After Z, curX=10, curY=20, then L 100 100
	if box.MaxX != 100 || box.MaxY != 100 {
		t.Errorf("max = (%v, %v), want (100, 100)", box.MaxX, box.MaxY)
	}
}

func TestParsePointsPolygon(t *testing.T) {
	box := parsePoints("10,20 30,40 50,60")
	if box.MinX != 10 || box.MinY != 20 {
		t.Errorf("min = (%v, %v), want (10, 20)", box.MinX, box.MinY)
	}
	if box.MaxX != 50 || box.MaxY != 60 {
		t.Errorf("max = (%v, %v), want (50, 60)", box.MaxX, box.MaxY)
	}
}
