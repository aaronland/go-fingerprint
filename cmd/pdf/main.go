package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aaronland/go-fingerprint/pdf"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	for _, path := range flag.Args() {

		//

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %w", path, err)
		}

		defer r.Close()

		//

		pdf_doc, err := pdf.FromReader(ctx, r)

		if err != nil {
			log.Fatal(err)
		}

		//

		pdf_path := fmt.Sprintf("%s.pdf", path)

		err = pdf_doc.OutputFileAndClose(pdf_path)

		if err != nil {
			log.Fatal("WOMP", err)
		}

	}

}
