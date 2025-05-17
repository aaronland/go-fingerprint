package image

import (
	"image"
	"image/color"

	"github.com/aaronland/go-image/v2/background"
)

// AddBackground draws 'im' on to a new `image.Image` instance of the same dimensions but with
// a white background.
func AddBackground(im image.Image) image.Image {

	bg_colour := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	return background.AddBackground(im, bg_colour)
}
