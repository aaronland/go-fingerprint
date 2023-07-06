package svg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

// Path is a struct representing an individual path element in a fingerprint SVG drawing.
type Path struct {
	// Fill is the colour assigned to the path.
	Fill string `xml:"fill,attr"`
	// FillOpacity is the opacity of the colour (assigned to the path).
	FillOpacity float64 `xml:"fill-opacity,attr"`
	// D is the SVG-encoded value of the path.
	D string `xml:"d,attr"`
}

// HexColor will return the 8-digit hexidecial color (fill + fill opacity) for the path
func (p *Path) HexColor() string {

	fill := p.Fill
	opacity := p.FillOpacity

	dec := int(opacity * 100.)

	color := dec * 255 / 100
	alpha := fmt.Sprintf("%02x", color)

	return fill + alpha
}

func (p *Path) Type() (string, error) {

	path_re, err := regexp.Compile(`^M\s{0,}\d+\,\d+\s{0,}(L|C)`)

	if err != nil {
		return "", fmt.Errorf("Failed to compile pattern, %w", err)
	}

	if !path_re.MatchString(p.D) {
		return "", fmt.Errorf("Unsupported path '%s'", p.D)
	}

	m := path_re.FindStringSubmatch(p.D)
	return m[1], nil
}

func (p *Path) Draw(dc *gg.Context, scale float64) error {

	t, err := p.Type()

	if err != nil {
		return nil
		// return fmt.Errorf("Failed to determine type, %w", err)
	}

	switch t {
	case "L":
		return p.drawLine(dc, scale)
	case "C":
		return p.drawCurve(dc, scale)
	default:
		return fmt.Errorf("Unsupported draw type, %s", t)
	}
}

func (p *Path) drawLine(dc *gg.Context, scale float64) error {

	coords, err := p.Coordinates()

	if err != nil {
		return fmt.Errorf("Failed to derive coordinates, %w", err)
	}

	count := len(coords)

	dc.NewSubPath()

	for i := 0; i < count; i++ {
		x := coords[i][0] * scale
		y := coords[i][1] * scale
		dc.LineTo(x, y)
	}

	dc.ClosePath()

	dc.SetHexColor(p.HexColor())
	dc.Fill()

	return nil
}

func (p *Path) drawCurve(dc *gg.Context, scale float64) error {

	coords, err := p.curves()

	if err != nil {
		return fmt.Errorf("Failed to derive curves, %w", err)
	}

	count := len(coords)

	dc.NewSubPath()

	for i := 0; i < count; i++ {
		mx := coords[i][0][0] * scale
		my := coords[i][0][1] * scale

		x1 := coords[i][1][0] * scale
		y1 := coords[i][1][1] * scale

		x2 := coords[i][2][0] * scale
		y2 := coords[i][2][1] * scale

		x := coords[i][3][0] * scale
		y := coords[i][3][1] * scale

		dc.MoveTo(mx, my)
		dc.CubicTo(x1, y1, x2, y2, x, y)
	}

	dc.ClosePath()

	dc.SetHexColor(p.HexColor())
	dc.Fill()

	return nil
}

func (p *Path) curves() ([][4][2]float64, error) {

	re, err := regexp.Compile(`^M\s{0,}(\d+\,\d+)\s{0,}C`)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile pattern, %w", err)
	}

	d := p.D

	if !re.MatchString(d) {
		return nil, fmt.Errorf("Invalid preamble, %w", err)
	}

	m := re.FindStringSubmatch(d)

	prefix := m[0]
	suffix := "Z"

	d = strings.Replace(d, prefix, "", 1)
	d = strings.TrimRight(d, suffix)

	xy := strings.Split(m[1], ",")

	mx, err := strconv.ParseFloat(xy[0], 64)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse mx '%s', %w", xy[0], err)
	}

	my, err := strconv.ParseFloat(xy[1], 64)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse my '%s', %w", xy[1], err)
	}

	curves := strings.Split(d, "C")
	coords := make([][4][2]float64, len(curves))

	for i, curve := range curves {

		xy := strings.Split(curve, ",")

		if len(xy) != 6 {
			return nil, fmt.Errorf("Invalid coordinate (%s) at offset %d", curve, i)
		}

		x1, err := strconv.ParseFloat(xy[0], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[0], i, err)
		}

		y1, err := strconv.ParseFloat(xy[1], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[1], i, err)
		}

		x2, err := strconv.ParseFloat(xy[2], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[2], i, err)
		}

		y2, err := strconv.ParseFloat(xy[3], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[3], i, err)
		}

		x, err := strconv.ParseFloat(xy[4], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[4], i, err)
		}

		y, err := strconv.ParseFloat(xy[5], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[5], i, err)
		}

		coords[i] = [4][2]float64{
			[2]float64{mx, my},
			[2]float64{x1, y1},
			[2]float64{x2, y2},
			[2]float64{x, y},
		}

		mx = x
		my = y
	}

	return coords, nil

}

// Coordinates will return a list of X and Y floating points for 'p' derived from its `D` property.
func (p *Path) Coordinates() ([][2]float64, error) {

	d := p.D
	d = strings.TrimLeft(d, "M")
	d = strings.TrimRight(d, "Z")

	pairs := strings.Split(d, "L")
	coords := make([][2]float64, len(pairs))

	for i, pair := range pairs {

		xy := strings.Split(pair, ",")

		if len(xy) != 2 {
			return nil, fmt.Errorf("Invalid coordinate (%s) at offset %d", pair, i)
		}

		x, err := strconv.ParseFloat(xy[0], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[0], i, err)
		}

		y, err := strconv.ParseFloat(xy[1], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[1], i, err)
		}

		coords[i] = [2]float64{x, y}
	}

	return coords, nil
}
