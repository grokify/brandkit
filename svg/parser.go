package svg

import (
	"regexp"
	"strconv"

	"github.com/JoshVarga/svgparser"
)

// PathCommand represents a single SVG path command.
type PathCommand struct {
	Command byte
	Params  []float64
}

// ParsePath parses an SVG path d attribute into commands.
func ParsePath(d string) []PathCommand {
	var commands []PathCommand

	// Match command letters followed by optional numbers
	cmdRe := regexp.MustCompile(`([MmLlHhVvCcSsQqTtAaZz])([^MmLlHhVvCcSsQqTtAaZz]*)`)
	// SVG path numbers can be:
	// - Optional sign, digits, optional decimal + digits: -123.45
	// - Optional sign, optional digits, decimal + digits: -.45 or .45
	// Numbers can be separated by whitespace, commas, or nothing when the next number starts with - or .
	numRe := regexp.MustCompile(`[+-]?(?:\d+\.?\d*|\.\d+)(?:[eE][+-]?\d+)?`)

	matches := cmdRe.FindAllStringSubmatch(d, -1)
	for _, match := range matches {
		cmd := match[1][0]
		params := numRe.FindAllString(match[2], -1)

		var floatParams []float64
		for _, p := range params {
			if v, err := strconv.ParseFloat(p, 64); err == nil {
				floatParams = append(floatParams, v)
			}
		}

		commands = append(commands, PathCommand{Command: cmd, Params: floatParams})
	}

	return commands
}

// CalculatePathBounds calculates the bounding box from path commands.
func CalculatePathBounds(d string) *BoundingBox {
	box := NewBoundingBox()
	commands := ParsePath(d)

	var curX, curY float64
	var startX, startY float64

	for _, cmd := range commands {
		switch cmd.Command {
		case 'M': // moveto absolute
			for i := 0; i+1 < len(cmd.Params); i += 2 {
				curX, curY = cmd.Params[i], cmd.Params[i+1]
				if i == 0 {
					startX, startY = curX, curY
				}
				box.Expand(curX, curY)
			}
		case 'm': // moveto relative
			for i := 0; i+1 < len(cmd.Params); i += 2 {
				curX += cmd.Params[i]
				curY += cmd.Params[i+1]
				if i == 0 {
					startX, startY = curX, curY
				}
				box.Expand(curX, curY)
			}
		case 'L': // lineto absolute
			for i := 0; i+1 < len(cmd.Params); i += 2 {
				curX, curY = cmd.Params[i], cmd.Params[i+1]
				box.Expand(curX, curY)
			}
		case 'l': // lineto relative
			for i := 0; i+1 < len(cmd.Params); i += 2 {
				curX += cmd.Params[i]
				curY += cmd.Params[i+1]
				box.Expand(curX, curY)
			}
		case 'H': // horizontal absolute
			for _, x := range cmd.Params {
				curX = x
				box.Expand(curX, curY)
			}
		case 'h': // horizontal relative
			for _, dx := range cmd.Params {
				curX += dx
				box.Expand(curX, curY)
			}
		case 'V': // vertical absolute
			for _, y := range cmd.Params {
				curY = y
				box.Expand(curX, curY)
			}
		case 'v': // vertical relative
			for _, dy := range cmd.Params {
				curY += dy
				box.Expand(curX, curY)
			}
		case 'C': // cubic bezier absolute
			for i := 0; i+5 < len(cmd.Params); i += 6 {
				box.Expand(cmd.Params[i], cmd.Params[i+1])
				box.Expand(cmd.Params[i+2], cmd.Params[i+3])
				curX, curY = cmd.Params[i+4], cmd.Params[i+5]
				box.Expand(curX, curY)
			}
		case 'c': // cubic bezier relative
			for i := 0; i+5 < len(cmd.Params); i += 6 {
				box.Expand(curX+cmd.Params[i], curY+cmd.Params[i+1])
				box.Expand(curX+cmd.Params[i+2], curY+cmd.Params[i+3])
				curX += cmd.Params[i+4]
				curY += cmd.Params[i+5]
				box.Expand(curX, curY)
			}
		case 'S': // smooth cubic absolute
			for i := 0; i+3 < len(cmd.Params); i += 4 {
				box.Expand(cmd.Params[i], cmd.Params[i+1])
				curX, curY = cmd.Params[i+2], cmd.Params[i+3]
				box.Expand(curX, curY)
			}
		case 's': // smooth cubic relative
			for i := 0; i+3 < len(cmd.Params); i += 4 {
				box.Expand(curX+cmd.Params[i], curY+cmd.Params[i+1])
				curX += cmd.Params[i+2]
				curY += cmd.Params[i+3]
				box.Expand(curX, curY)
			}
		case 'Q': // quadratic bezier absolute
			for i := 0; i+3 < len(cmd.Params); i += 4 {
				box.Expand(cmd.Params[i], cmd.Params[i+1])
				curX, curY = cmd.Params[i+2], cmd.Params[i+3]
				box.Expand(curX, curY)
			}
		case 'q': // quadratic bezier relative
			for i := 0; i+3 < len(cmd.Params); i += 4 {
				box.Expand(curX+cmd.Params[i], curY+cmd.Params[i+1])
				curX += cmd.Params[i+2]
				curY += cmd.Params[i+3]
				box.Expand(curX, curY)
			}
		case 'T': // smooth quadratic absolute
			for i := 0; i+1 < len(cmd.Params); i += 2 {
				curX, curY = cmd.Params[i], cmd.Params[i+1]
				box.Expand(curX, curY)
			}
		case 't': // smooth quadratic relative
			for i := 0; i+1 < len(cmd.Params); i += 2 {
				curX += cmd.Params[i]
				curY += cmd.Params[i+1]
				box.Expand(curX, curY)
			}
		case 'A': // arc absolute
			for i := 0; i+6 < len(cmd.Params); i += 7 {
				curX, curY = cmd.Params[i+5], cmd.Params[i+6]
				box.Expand(curX, curY)
			}
		case 'a': // arc relative
			for i := 0; i+6 < len(cmd.Params); i += 7 {
				curX += cmd.Params[i+5]
				curY += cmd.Params[i+6]
				box.Expand(curX, curY)
			}
		case 'Z', 'z': // closepath
			curX, curY = startX, startY
		}
	}

	return box
}

// GetElementBounds calculates bounds for an SVG element.
func GetElementBounds(elem *svgparser.Element) *BoundingBox {
	box := NewBoundingBox()

	switch elem.Name {
	case "path":
		if d, ok := elem.Attributes["d"]; ok {
			box.Merge(CalculatePathBounds(d))
		}
	case "circle":
		cx := ParseFloat(elem.Attributes["cx"], 0)
		cy := ParseFloat(elem.Attributes["cy"], 0)
		r := ParseFloat(elem.Attributes["r"], 0)
		box.Expand(cx-r, cy-r)
		box.Expand(cx+r, cy+r)
	case "ellipse":
		cx := ParseFloat(elem.Attributes["cx"], 0)
		cy := ParseFloat(elem.Attributes["cy"], 0)
		rx := ParseFloat(elem.Attributes["rx"], 0)
		ry := ParseFloat(elem.Attributes["ry"], 0)
		box.Expand(cx-rx, cy-ry)
		box.Expand(cx+rx, cy+ry)
	case "rect":
		x := ParseFloat(elem.Attributes["x"], 0)
		y := ParseFloat(elem.Attributes["y"], 0)
		w := ParseFloat(elem.Attributes["width"], 0)
		h := ParseFloat(elem.Attributes["height"], 0)
		box.Expand(x, y)
		box.Expand(x+w, y+h)
	case "line":
		x1 := ParseFloat(elem.Attributes["x1"], 0)
		y1 := ParseFloat(elem.Attributes["y1"], 0)
		x2 := ParseFloat(elem.Attributes["x2"], 0)
		y2 := ParseFloat(elem.Attributes["y2"], 0)
		box.Expand(x1, y1)
		box.Expand(x2, y2)
	case "polygon", "polyline":
		if points, ok := elem.Attributes["points"]; ok {
			box.Merge(parsePoints(points))
		}
	}

	// Recursively process children
	for _, child := range elem.Children {
		// Skip mask and clipPath elements - they define clipping regions, not visible content
		if child.Name == "mask" || child.Name == "clipPath" || child.Name == "defs" {
			continue
		}
		childBox := GetElementBounds(child)
		box.Merge(childBox)
	}

	return box
}

// parsePoints parses polygon/polyline points attribute.
func parsePoints(points string) *BoundingBox {
	box := NewBoundingBox()
	re := regexp.MustCompile(`-?[\d]+\.?[\d]*`)
	matches := re.FindAllString(points, -1)

	for i := 0; i+1 < len(matches); i += 2 {
		x, _ := strconv.ParseFloat(matches[i], 64)
		y, _ := strconv.ParseFloat(matches[i+1], 64)
		box.Expand(x, y)
	}

	return box
}
