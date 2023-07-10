package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aaronland/go-fingerprint/fpdf"
	"github.com/aaronland/go-fingerprint/pdf"
	"github.com/sfomuseum/go-flags/flagset"
)

func main() {

	fs := flagset.NewFlagSet("pdf")

	fpdf.AppendFlags(fs)

	flagset.Parse(fs)

	ctx := context.Background()

	opts, err := fpdf.DefaultOptions(ctx)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range fs.Args() {

		//

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %w", path, err)
		}

		defer r.Close()

		//

		pdf_doc, err := pdf.FromReader(ctx, r, opts)

		if err != nil {
			log.Fatal(err)
		}

		//

		pdf_path := fmt.Sprintf("%s.pdf", path)

		err = pdf_doc.Save(pdf_path)

		if err != nil {
			log.Fatal("WOMP", err)
		}

	}

}
