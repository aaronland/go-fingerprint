package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aaronland/go-fingerprint/pdf"
	"github.com/aaronland/go-fpdf"
	"github.com/sfomuseum/go-flags/flagset"
)

func main() {

	fs := flagset.NewFlagSet("pdf")
	fpdf.AppendFlags(fs)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Convert one or more fingerprint SVG documents into PDF documents.\n\n")
		fmt.Fprintf(os.Stderr, "Each document contains a full-page rasterized rendering of the SVG document\n")
		fmt.Fprintf(os.Stderr, "followed by one or more pages containing the body of the aaronland/go-fingerprint/svg.Document\n")
		fmt.Fprintf(os.Stderr, "representation of the SVG document as JSON-encoded text.\n\n")
		fmt.Fprintf(os.Stderr, "The final PDF document will be saved as (svg-path).pdf.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] path(N) path(N)\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	flagset.Parse(fs)

	ctx := context.Background()

	pdf_opts, err := fpdf.OptionsFromFlagSet(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to derive options from flagset, %w", err)
	}

	for _, path := range fs.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %w", path, err)
		}

		defer r.Close()

		title := filepath.Base(path)

		pdf_doc, err := pdf.FromReader(ctx, r, title, pdf_opts)

		if err != nil {
			log.Fatalf("Failed to create PDF from reader for %s, %w", path, err)
		}

		pdf_path := fmt.Sprintf("%s.pdf", path)
		// pdf_path = "test.pdf"

		err = pdf_doc.Save(pdf_path)

		if err != nil {
			log.Fatalf("Failed to save %s, %w", pdf_path, err)
		}

		log.Println(pdf_path)
	}

}
