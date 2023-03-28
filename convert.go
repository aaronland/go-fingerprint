package fingerprint

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/aaronland/go-fingerprint/image"
	"github.com/aaronland/go-fingerprint/svg"
)

func Convert(r io.ReadSeeker, wr io.Writer, max_dimension float64) error {

	info, err := svg.Info(r)

	if err != nil {
		return fmt.Errorf("Failed to derive info, %w", err)
	}

	// 2023-03-19T06:50:28.965Z
	layout := "2006-01-02T15:04:05.000Z"

	t, err := time.Parse(layout, info.Date)

	if err != nil {
		return fmt.Errorf("Failed to parse date (%s), %w", info.Date, err)
	}

	w := info.Width
	h := info.Height

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
