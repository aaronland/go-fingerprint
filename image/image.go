package image

import (
	"image"
	"image/color"
	"image/draw"
	"runtime"

	"github.com/mandykoh/prism"
	"github.com/mandykoh/prism/adobergb"
	"github.com/mandykoh/prism/srgb"
)

// ToAdobeRGB converts all the coloura in 'im' to match the Adobe RGB colour profile.
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

// AddBackground draws 'im' on to a new `image.Image` instance of the same dimensions but with
// a white background.
func AddBackground(im image.Image) image.Image {

	backgroundColor := color.NRGBA{0xff, 0xff, 0xff, 0xff}

	new_im := image.NewNRGBA(im.Bounds())

	draw.Draw(new_im, new_im.Bounds(), image.NewUniform(backgroundColor), image.Point{}, draw.Src)
	draw.Draw(new_im, new_im.Bounds(), im, im.Bounds().Min, draw.Over)

	return new_im
}
