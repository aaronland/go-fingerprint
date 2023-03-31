package svg

import (
	"fmt"
	"image"
	"math"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

// Document defines a struct representing a fingerprint SVG document. At this time it is not
type Document struct {
	// Width is the width of the fingerprint SVG drawing.
	Width int `xml:"width,attr"`
	// Height is the width of the fingerprint SVG drawing.
	Height int `xml:"height,attr"`
	// ViewBox is the viewbox boundary of the fingerprint SVG drawing.
	ViewBox string `xml:"viewBox,attr"`
	// Date is the date that a fingerprint SVG drawing was produced.
	Date string `xml:"x-fingerprint-date,attr"`
	// Paths is the list of SVG path that define a fingerprint SVG drawing.
	Paths []*Path `xml:"path"`
}

func (doc *Document) ToImage(max_dimension float64) (image.Image, error) {

	w := doc.Width
	h := doc.Height

	scale := 1.0

	max := math.Max(float64(w), float64(h))

	if max_dimension > max {
		scale = max_dimension / max
	}

	w = int(float64(w) * scale)
	h = int(float64(h) * scale)

	dc := gg.NewContext(w, h)

	for _, p := range doc.Paths {

		coords, err := p.Coordinates()

		if err != nil {
			return nil, fmt.Errorf("Failed to derive coordinates, %w", err)
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
	}

	return dc.Image(), nil
}

// Path is a struct representing an individual path element in a fingerprint SVG drawing.
type Path struct {
	// Fill is the colour assigned to the path.
	Fill string `xml:"fill,attr"`
	// FillOpacity is the opacity of the colour (assigned to the path).
	FillOpacity float64 `xml:"fill-opacity,attr"`
	// D is the SVG-encoded value of the path.
	D string `xml:"d,attr"`
}

func (p *Path) HexColor() string {

	fill := p.Fill
	opacity := p.FillOpacity

	dec := int(opacity * 100.)

	color := dec * 255 / 100
	alpha := fmt.Sprintf("%02x", color)

	return fill + alpha
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
