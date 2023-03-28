package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-fingerprint"
)

func main() {

	max_dimension := flag.Float64("max-dimension", 4096, "The maximum dimension to scale an image to.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Convert one or more fingerprint SVG documents in to JPEG images.\n\n")
		fmt.Fprintf(os.Stderr, "The final JPEG image is scaled to ensure that its maximum dimension is\n")
		fmt.Fprintf(os.Stderr, "'max_dimension'. Date information defined in the SVG document's `x-fingerprint-date`\n")
		fmt.Fprintf(os.Stderr, "attribute is written to the final JPEG image's `DateTime`, `DateTimeDigitized` and\n")
		fmt.Fprintf(os.Stderr, "`DateTimeOriginal` EXIF headers. The final JPEG representation is updated to ensure that\n")
		fmt.Fprintf(os.Stderr, "all pixel values match the Adobe RGB colour profile. JPEG images are written to the same\n")
		fmt.Fprintf(os.Stderr, "location as the source SVG document with a '.jpg' extension.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	for _, path := range flag.Args() {

		root := filepath.Dir(path)
		fname := filepath.Base(path)

		out_fname := strings.Replace(fname, ".svg", ".jpg", 1)
		out_path := filepath.Join(root, out_fname)

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", path, err)
		}

		defer r.Close()

		wr, err := os.OpenFile(out_path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatalf("Failed to open %s for writing, %v", out_path, err)
		}

		err = fingerprint.Convert(r, wr, *max_dimension)

		if err != nil {
			log.Fatalf("Failed to derive info for %s, %v", out_path, err)
		}

		err = wr.Close()

		if err != nil {
			log.Fatalf("Failed to close %s, %v", out_path, err)
		}
	}
}
