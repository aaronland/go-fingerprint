package fingerprint

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"runtime"

	"github.com/mandykoh/prism"
	"github.com/mandykoh/prism/adobergb"
	"github.com/mandykoh/prism/srgb"
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

func ToAdobeRGB(im image.Image) image.Image {

	input_im := prism.ConvertImageToNRGBA(im, runtime.NumCPU())
	new_im := image.NewNRGBA(input_im.Rect)

	for i := input_im.Rect.Min.Y; i < input_im.Rect.Max.Y; i++ {

		for j := input_im.Rect.Min.X; j < input_im.Rect.Max.X; j++ {

			inCol, alpha := adobergb.ColorFromNRGBA(input_im.NRGBAAt(j, i))
			outCol := srgb.ColorFromXYZ(inCol.ToXYZ())
			new_im.SetNRGBA(j, i, outCol.ToNRGBA(alpha))
		}
	}

	return new_im
}

func AddBackground(im image.Image) image.Image {

	backgroundColor := color.NRGBA{0xff, 0xff, 0xff, 0xff}

	new_im := image.NewNRGBA(im.Bounds())

	draw.Draw(new_im, new_im.Bounds(), image.NewUniform(backgroundColor), image.Point{}, draw.Src)
	draw.Draw(new_im, new_im.Bounds(), im, im.Bounds().Min, draw.Over)

	return new_im
}
