package fingerprint

import (
	"fmt"
	go_image "image"
	_ "image/png"
	"io"
	"time"

	"github.com/aaronland/go-fingerprint/image"
	"github.com/aaronland/go-fingerprint/svg"
)

func Outline(r io.ReadSeeker, wr io.Writer, max_dimension float64) (go_image.Image, error) {

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

	im, err := doc.ToOutline(max_dimension)

	if err != nil {
		return nil, fmt.Errorf("Failed to create image, %w", err)
	}

	im = image.AddBackground(im)

	err = image.AppendTime(im, wr, t)

	if err != nil {
		return nil, fmt.Errorf("Failed to add time, %w", err)
	}

	return im, nil
}
