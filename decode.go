package fingerprint

import (
	"context"
	"fmt"
	go_image "image"
	"io"

	"github.com/aaronland/go-fingerprint/image"
	"github.com/aaronland/go-fingerprint/svg"
	"github.com/aaronland/go-image/decode"
)

type FingerprintDecoder struct {
	decode.Decoder
	max_dimension float64
}

func init() {
	ctx := context.Background()
	decode.RegisterDecoder(ctx, NewFingerprintDecoder, "svg")
}

func NewFingerprintDecoder(ctx context.Context, uri string) (decode.Decoder, error) {

	d := &FingerprintDecoder{
		max_dimension: 4096,
	}

	return d, nil
}

func (d *FingerprintDecoder) Decode(ctx context.Context, r io.ReadSeeker) (go_image.Image, string, error) {

	doc, err := svg.Unmarshal(r)

	if err != nil {
		return nil, "", fmt.Errorf("Failed to derive doc, %w", err)
	}

	// 2023-03-19T06:50:28.965Z

	/*
		layout := "2006-01-02T15:04:05.000Z"

		t, err := time.Parse(layout, doc.Date)

		if err != nil {
			return nil, "", fmt.Errorf("Failed to parse date (%s), %w", doc.Date, err)
		}
	*/

	im, err := doc.ToImage(d.max_dimension)

	if err != nil {
		return nil, "", fmt.Errorf("Failed to create image, %w", err)
	}

	im = image.ToAdobeRGB(im)
	im = image.AddBackground(im)

	return im, "", nil
}
