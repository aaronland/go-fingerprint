package svg

import (
	"fmt"

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

// Type will return a string identifier for the type of path derived from the `D` property of 'p'.
func (p *Path) Type() string {

	return DerivePathType(p.D)
}

// Coodinates with return a `Coordinates` instance derived derived from the `D` property of 'p'.
func (p *Path) Coordinates() (Coordinates, error) {

	return DeriveCoordinates(p.D)
}

// Draw will render 'p' in to 'dc'.
func (p *Path) Draw(dc *gg.Context, scale float64) error {

	coordinates, err := p.Coordinates()

	if err != nil {
		return fmt.Errorf("Failed to determine type, %w", err)
	}

	return Draw(dc, coordinates, p.HexColor(), scale)
}
