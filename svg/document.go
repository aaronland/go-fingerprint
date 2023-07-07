package svg

import (
	"fmt"
	"image"
	_ "log"
	"math"

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

// ToImage will rasterize 'doc' and return it as an `image.Image` instance whose maximum dimension is scaled to 'max_dimension'
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

	for idx, p := range doc.Paths {

		err := p.Draw(dc, scale)

		if err != nil {
			return nil, fmt.Errorf("Failed to draw path at offset %d, %w", idx, err)
		}
	}

	return dc.Image(), nil
}
