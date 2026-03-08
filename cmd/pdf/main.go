package main

/*

> ./bin/pdf -root /usr/local/data/fingerprint-pdf -group-by-year /usr/local/data/fingerprint/fingerprint*.svg
2026/03/08 10:35:08 /usr/local/data/fingerprint-pdf/2023/2023-03-13-fingerprint-1678772041.svg.pdf
2026/03/08 10:35:08 /usr/local/data/fingerprint-pdf/2023/2023-03-17-fingerprint-1679068267.svg.pdf
2026/03/08 10:35:08 /usr/local/data/fingerprint-pdf/2023/2023-03-19-fingerprint-1679211124.svg.pdf
2026/03/08 10:35:08 /usr/local/data/fingerprint-pdf/2023/2023-03-19-fingerprint-1679239781.svg.pdf
...

*/

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

	var root string
	var group_by_year bool

	fs := flagset.NewFlagSet("pdf")
	fpdf.AppendFlags(fs)

	fs.StringVar(&root, "root", "", "The root directory where files will be saved. If empty then PDF files will be saved in the same directory as the source SVG documents.")
	fs.BoolVar(&group_by_year, "group-by-year", false, "Save PDF files in YYYY subdirectory (of root)")

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
		log.Fatalf("Failed to derive options from flagset, %v", err)
	}

	if root != "" {

		abs_root, err := filepath.Abs(root)

		if err != nil {
			log.Fatalf("Failed to derive absolute root, %v", err)
		}

		info, err := os.Stat(root)

		if err != nil {
			log.Fatalf("Failed to stat root, %v", err)
		}

		if !info.IsDir() {
			log.Fatalf("Root is not a directory")
		}

		root = abs_root
	}

	for _, path := range fs.Args() {

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatalf("Failed to derive abs path for %s, %v", path, err)
		}

		r, err := os.Open(abs_path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", abs_path, err)
		}

		defer r.Close()

		fname := filepath.Base(abs_path)

		pdf_doc, err := pdf.FromReader(ctx, r, fname, pdf_opts)

		if err != nil {
			log.Fatalf("Failed to create PDF from reader for %s, %v", abs_path, err)
		}

		pdf_t := pdf_doc.PDF.GetCreationDate()

		froot := root

		if froot == "" {
			froot = filepath.Dir(path)
		}

		if group_by_year {
			froot = filepath.Join(froot, pdf_t.Format("2006"))
		}

		err = os.MkdirAll(froot, 0755)

		if err != nil {
			log.Fatalf("Failed to create %s, %v", froot, err)
		}

		fname = fmt.Sprintf("%s-%s.pdf", pdf_t.Format("2006-01-02"), fname)
		pdf_path := filepath.Join(froot, fname)

		// pdf_path = "test.pdf"

		err = pdf_doc.Save(pdf_path)

		if err != nil {
			log.Fatalf("Failed to save %s, %v", pdf_path, err)
		}

		log.Println(pdf_path)
	}

}
