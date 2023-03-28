package fingerprint

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/aaronland/go-fingerprint/image"
	"github.com/aaronland/go-fingerprint/svg"
)

// Convert writes a fingerprint SVG document defined by 'r' to a JPEG image defined by 'wr'. The final
// JPEG image is scaled to ensure that its maximum dimension is 'max_dimension'. Date information defined
// in the SVG document's `x-fingerprint-date` attribute is written to the final JPEG image's `DateTime`,
// `DateTimeDigitized` and `DateTimeOriginal` EXIF headers. The final JPEG representation is updated to
// ensure that all pixel values match the Adobe RGB colour profile.
func Convert(r io.ReadSeeker, wr io.Writer, max_dimension float64) error {

	doc, err := svg.Parse(r)

	if err != nil {
		return fmt.Errorf("Failed to derive doc, %w", err)
	}

	// 2023-03-19T06:50:28.965Z
	layout := "2006-01-02T15:04:05.000Z"

	t, err := time.Parse(layout, doc.Date)

	if err != nil {
		return fmt.Errorf("Failed to parse date (%s), %w", doc.Date, err)
	}

	w := doc.Width
	h := doc.Height

	if max_dimension > 0 {

		max := math.Max(float64(w), float64(h))
		scale := 1.0

		if max_dimension > max {
			scale = max_dimension / max
		}

		w = int(float64(w) * scale)
		h = int(float64(h) * scale)
	}

	_, err = r.Seek(0, 0)

	if err != nil {
		return fmt.Errorf("Failed to rewind reader, %w", err)
	}

	im, err := svg.ToImage(r, w, h)

	if err != nil {
		return fmt.Errorf("Failed to render image, %w", err)
	}

	im = image.ToAdobeRGB(im)
	im = image.AddBackground(im)

	err = image.AppendTime(im, wr, t)

	if err != nil {
		return fmt.Errorf("Failed to add time, %w", err)
	}

	return nil
}
