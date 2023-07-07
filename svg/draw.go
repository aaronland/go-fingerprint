package svg

import (
	"fmt"

	"github.com/fogleman/gg"
)

// Draw renders the path defined by 'coords' on to 'gg'.
func Draw(dc *gg.Context, coords Coordinates, color string, scale float64) error {

	t := coords.Type()

	switch t {
	case LINE:
		return DrawLine(dc, coords.(Line), color, scale)
	case CUBIC:
		return DrawCubic(dc, coords.(Cubic), color, scale)
	default:
		return fmt.Errorf("Unsupported coordinate type '%s'", t)
	}
}

// Draw renders the path defined by 'coords' as a closed line string on to 'gg'.
func DrawLine(dc *gg.Context, points Line, color string, scale float64) error {

	count := len(points)

	dc.NewSubPath()

	for i := 0; i < count; i++ {
		x := points[i][0] * scale
		y := points[i][1] * scale
		dc.LineTo(x, y)
	}

	dc.ClosePath()

	dc.SetHexColor(color)
	dc.Fill()

	return nil
}

// Draw renders the path defined by 'coords' as a closed cubic curve path on to 'gg'.
func DrawCubic(dc *gg.Context, curves Cubic, color string, scale float64) error {

	count := len(curves)

	dc.NewSubPath()

	for i := 0; i < count; i++ {

		x1 := curves[i][0][0] * scale
		y1 := curves[i][0][1] * scale

		x2 := curves[i][1][0] * scale
		y2 := curves[i][1][1] * scale

		x := curves[i][2][0] * scale
		y := curves[i][2][1] * scale

		dc.CubicTo(x1, y1, x2, y2, x, y)
	}

	dc.ClosePath()

	dc.SetHexColor(color)
	dc.Fill()

	return nil
}
