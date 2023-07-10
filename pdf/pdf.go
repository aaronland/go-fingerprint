package pdf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	_ "log"
	"math"
	"os"

	"github.com/aaronland/go-fingerprint"
	"github.com/aaronland/go-fingerprint/fpdf"
	"github.com/aaronland/go-fingerprint/svg"
	"github.com/jung-kurt/gofpdf"
)

func FromReader(ctx context.Context, r io.ReadSeeker, title string, opts *fpdf.Options) (*fpdf.Document, error) {

	cell_h := .15 // This should derived from page and font dimensions...

	pdf_doc, err := fpdf.NewDocument(ctx, opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create PDF document, %w", err)
	}

	pdf := pdf_doc.PDF

	doc, err := svg.Unmarshal(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal document, %w", err)
	}

	_, err = r.Seek(0, 0)

	if err != nil {
		return nil, fmt.Errorf("Failed to rewind reader, %w", err)
	}

	// START OF move this in to NewDocument

	margins := pdf_doc.Margins

	left := margins.Left / pdf_doc.Options.DPI
	right := margins.Right / pdf_doc.Options.DPI
	top := margins.Top / pdf_doc.Options.DPI
	bottom := margins.Bottom / pdf_doc.Options.DPI

	pdf.SetMargins(left, top, right)
	pdf.SetAutoPageBreak(true, bottom)

	// END OF move this in to NewDocument

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "", 6)

	pdf.SetFooterFunc(func() {

		x := pdf_doc.Margins.Left / pdf_doc.Options.DPI
		y := (pdf_doc.Canvas.Height + (pdf_doc.Margins.Top * 1.35)) / pdf_doc.Options.DPI

		pdf.SetXY(x, y)
		pdf.SetFont("Courier", "", 6)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(0, cell_h, fmt.Sprintf("%s/%d", title, pdf.PageNo()), "", 0, "C", false, 0, "")
	})

	// Render SVG

	// Note, we could also use fpdf.SVGBasicWrite` which might save the toruble
	// of all the scaling and positioning code below. Still untested...
	// https://pkg.go.dev/github.com/jung-kurt/gofpdf#example-Fpdf.SVGBasicWrite

	wr, err := os.CreateTemp("", "fingerprint.*.jpg")

	if err != nil {
		return nil, fmt.Errorf("Failed to create temporary image file, %w", err)
	}

	defer os.Remove(wr.Name())

	max_d := math.Max(pdf_doc.Canvas.Width, pdf_doc.Canvas.Height)

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
		return nil, fmt.Errorf("Failed to register image options")
	}

	info.SetDpi(pdf_doc.Options.DPI)

	// Remember: margins have been calculated inclusive of page bleeds

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

	pdf.AddPage()
	pdf.ImageOptions(wr.Name(), image_x, image_y, image_w, image_h, false, image_opts, 0, "")

	pdf.AddPage()

	// Write the data to the PDF

	enc_doc, err := json.Marshal(doc)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal document, %w", err)
	}

	pdf.AddPage()

	// cell_w := pdf_doc.Canvas.Width / pdf_doc.Options.DPI

	pdf.MultiCell(0, cell_h, string(enc_doc), "", "L", false)

	//

	return pdf_doc, nil
}
