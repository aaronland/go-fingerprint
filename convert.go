package fingerprint

import (
	"fmt"
	go_image "image"
	_ "image/png"
	"io"
	"time"

	"github.com/aaronland/go-fingerprint/image"
	"github.com/aaronland/go-fingerprint/svg"
	"github.com/aaronland/go-image/colour"
)

// Convert writes a fingerprint SVG document defined by 'r' to a JPEG image defined by 'wr'. The final
// JPEG image is scaled to ensure that its maximum dimension is 'max_dimension'. Date information defined
// in the SVG document's `x-fingerprint-date` attribute is written to the final JPEG image's `DateTime`,
// `DateTimeDigitized` and `DateTimeOriginal` EXIF headers. The final JPEG representation is updated to
// ensure that all pixel values match the Adobe RGB colour profile.
func Convert(r io.ReadSeeker, wr io.Writer, max_dimension float64) (go_image.Image, error) {

	doc, err := svg.Unmarshal(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive doc, %w", err)
	}

	// 2023-03-19T06:50:28.965Z
	layout := "2006-01-02T15:04:05.000Z"

	t, err := time.Parse(layout, doc.Date)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse date (%s), %w", doc.Date, err)
	}

	im, err := doc.ToImage(max_dimension)

	if err != nil {
		return nil, fmt.Errorf("Failed to create image, %w", err)
	}

	// im = colour.ToAdobeRGB(im)
	im = colour.ToDisplayP3(im)
	im = image.AddBackground(im)

	err = image.AppendTime(im, wr, t)

	if err != nil {
		return nil, fmt.Errorf("Failed to add time, %w", err)
	}

	return im, nil
}
