package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/aaronland/go-fingerprint"
	"github.com/aaronland/go-fingerprint/svg"
	"github.com/jung-kurt/gofpdf"
)

func main() {

	flag.Parse()

	pdf := gofpdf.New("P", "in", "letter", "")
	pdf.SetFont("Courier", "", 6)

	h := .15

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %w", path, err)
		}

		defer r.Close()

		//

		wr, err := os.CreateTemp("", "example.*.jpg")

		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(wr.Name()) // clean up

		err = fingerprint.Convert(r, wr, 800.00)

		if err != nil {
			log.Fatal(err)
		}

		err = wr.Close()

		if err != nil {
			log.Fatal(err)
		}

		pdf.AddPage()

		var opt gofpdf.ImageOptions
		opt.ImageType = "jpg"

		pdf.ImageOptions(wr.Name(), 0, 0, -1, -1, false, opt, 0, "")

		//

		r.Seek(0, 0)

		doc, err := svg.Unmarshal(r)

		if err != nil {
			log.Fatalf("Failed to read %s, %w", err)
		}

		enc_doc, err := json.Marshal(doc)

		if err != nil {
			log.Fatalf("Failed to encode, %w", err)
		}

		pdf.AddPage()
		pdf.MultiCell(0, h, string(enc_doc), "", "L", false)
	}

	err := pdf.OutputFileAndClose("hello.pdf")

	if err != nil {
		log.Fatal("WOMP", err)
	}
}
