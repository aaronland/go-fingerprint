package svg

import (
	"fmt"

	"github.com/fogleman/gg"
)

// DrawOutline renders the path defined by 'coords' on to 'gg'.
func DrawOutline(dc *gg.Context, coords Coordinates, scale float64) error {

	t := coords.Type()

	switch t {
	case LINE:
		return DrawOutlineLine(dc, coords.(Line), scale)
	case CUBIC:
		return DrawOutlineCubic(dc, coords.(Cubic), scale)
	default:
		return fmt.Errorf("Unsupported coordinate type '%s'", t)
	}
}

// DrawOutlineLine renders the path defined by 'coords' as a closed line string on to 'gg'.
func DrawOutlineLine(dc *gg.Context, points Line, scale float64) error {

	count := len(points)

	dc.NewSubPath()

	for i := 0; i < count; i++ {
		x := points[i][0] * scale
		y := points[i][1] * scale
		dc.LineTo(x, y)
	}

	dc.ClosePath()

	dc.SetHexColor("#000000")
	dc.Stroke()

	return nil
}

// DrawOutlineCubic renders the path defined by 'coords' as a closed cubic curve path on to 'gg'.
func DrawOutlineCubic(dc *gg.Context, curves Cubic, scale float64) error {

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

	dc.SetHexColor("#000000")
	dc.Stroke()

	return nil
}
