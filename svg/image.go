package svg

import (
	"fmt"
	"image"
	"io"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func ToImage(r io.Reader, w int, h int) (image.Image, error) {

	icon, err := oksvg.ReadIconStream(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to read stream, %w", err)
	}

	icon.SetTarget(0, 0, float64(w), float64(h))

	im := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, im, im.Bounds())), 1)

	// background stuff

	return im, nil
}
