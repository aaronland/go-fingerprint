package pdf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aaronland/go-fingerprint"
	"github.com/aaronland/go-fingerprint/fpdf"
	"github.com/aaronland/go-fingerprint/svg"
	"github.com/jung-kurt/gofpdf"
)

func FromReader(ctx context.Context, r io.ReadSeeker, opts *fpdf.Options) (*fpdf.Document, error) {

	pdf_doc, err := fpdf.NewDocument(ctx, opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create PDF document, %w", err)
	}

	pdf := pdf_doc.PDF

	cell_h := .15
	max_d := 11.0 * pdf_doc.Options.DPI

	doc, err := svg.Unmarshal(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal document, %w", err)
	}

	_, err = r.Seek(0, 0)

	if err != nil {
		return nil, fmt.Errorf("Failed to rewind reader, %w", err)
	}

	wr, err := os.CreateTemp("", "fingerprint.*.jpg")

	if err != nil {
		return nil, fmt.Errorf("Failed to create temporary image file, %w", err)
	}

	defer os.Remove(wr.Name())

	im, err := fingerprint.Convert(r, wr, max_d)

	if err != nil {
		return nil, fmt.Errorf("Failed to create image, %w", err)
	}

	err = wr.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close image writer, %w", err)
	}

	// Draw the image to the PDF

	image_opts := gofpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: "jpeg",
	}

	dims := im.Bounds()

	w := float64(dims.Max.X)
	h := float64(dims.Max.Y)

	max_w := pdf_doc.Canvas.Width
	max_h := pdf_doc.Canvas.Height

	im_r, _ := os.Open(wr.Name())
	defer im_r.Close()

	info := pdf.RegisterImageOptionsReader(wr.Name(), image_opts, im_r)

	if info == nil {
		return nil, fmt.Errorf("SAD 1")
	}

	info.SetDpi(pdf_doc.Options.DPI)

	if w == 0.0 || h == 0.0 {
		return nil, fmt.Errorf("SAD 2")
	}

	// Remember: margins have been calculated inclusive of page bleeds

	margins := pdf_doc.Margins

	x := margins.Left
	y := margins.Top

	for {

		if w >= max_w || h >= max_h {

			if w > max_w {

				ratio := max_w / w
				w = max_w
				h = h * ratio

			}

			if h > max_h {

				ratio := max_h / h
				w = w * ratio
				h = max_h

			}

		}

		// TO DO: ENSURE ! h < max_h && ! w < max_w

		if w <= max_w && h <= max_h {
			break

			if h < max_h {
				h = max_h
			}

		}
	}

	if w < max_w {
		padding := max_w - w
		x = x + (padding / 2.0)
	}

	if h < max_h {
		padding := max_h - h
		y = y + (padding / 2.0)
	}

	image_x := x / pdf_doc.Options.DPI
	image_y := y / pdf_doc.Options.DPI
	image_w := w / pdf_doc.Options.DPI
	image_h := h / pdf_doc.Options.DPI

	//

	pdf.AddPage()
	// pdf.ImageOptions(wr.Name(), 0, 0, -1, -1, false, opt, 0, "")

	// var opt gofpdf.ImageOptions
	// opt.ImageType = "jpg"

	pdf.ImageOptions(wr.Name(), image_x, image_y, image_w, image_h, false, image_opts, 0, "")

	// Write the data to the PDF

	enc_doc, err := json.Marshal(doc)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal document, %w", err)
	}

	pdf.AddPage()
	pdf.MultiCell(0, cell_h, string(enc_doc), "", "L", false)

	//

	return pdf_doc, nil
}
