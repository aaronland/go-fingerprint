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

	h := .15
	max_d := 11.0 * 72.0
	// font_size := 8.0

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

	err = fingerprint.Convert(r, wr, max_d)

	if err != nil {
		return nil, fmt.Errorf("Failed to create image, %w", err)
	}

	err = wr.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close image writer, %w", err)
	}

	// Draw the image to the PDF

	// pdf = gofpdf.New("P", "in", "letter", "")
	// pdf.SetFont("Helvetica", "", font_size)

	pdf.AddPage()

	var opt gofpdf.ImageOptions
	opt.ImageType = "jpg"

	pdf.ImageOptions(wr.Name(), 0, 0, -1, -1, false, opt, 0, "")

	// Write the data to the PDF

	enc_doc, err := json.Marshal(doc)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal document, %w", err)
	}

	pdf.AddPage()
	pdf.MultiCell(0, h, string(enc_doc), "", "L", false)

	//

	return pdf_doc, nil
}
