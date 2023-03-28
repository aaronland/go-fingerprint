package svg

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Path is a struct representing an individual path element in a fingerprint SVG drawing.
type Path struct {
	// Fill is the colour assigned to the path.
	Fill string `xml:"fill,attr"`
	// FillOpacity is the opacity of the colour (assigned to the path).
	FillOpacity float64 `xml:fill-opacity,attr"`
	// D is the SVG-encoded value of the path.
	D string `xml:d,attr"`
}

// Coordinates will return a list of X and Y floating points for 'p' derived from its `D` property.
func (p *Path) Coordinates() ([][2]float64, error) {

	d := p.D
	d = strings.TrimLeft(d, "L")
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

// Parse will parse the fingerprint SVG document defined by 'r' in a new `Document` instance.
func Parse(r io.Reader) (*Document, error) {

	var doc *Document

	dec := xml.NewDecoder(r)
	err := dec.Decode(&doc)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode record, %w", err)
	}

	return doc, nil
}
