package main

import (
	"flag"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-fingerprint"
)

func main() {

	flag.Parse()

	w, h := 375, 511

	factor := 8

	w = w * factor
	h = h * factor

	for _, path := range flag.Args() {

		root := filepath.Dir(path)
		fname := filepath.Base(path)

		r, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer r.Close()

		info, err := fingerprint.Info(r)

		if err != nil {
			log.Fatalf("Failed to derive info for %s, %v", path, err)
		}

		log.Println(info.Date)

		// Do max dim stuff here...

		_, err = r.Seek(0, 0)

		if err != nil {
			log.Fatalf("Failed to rewind reader for %s, %v", path, err)
		}

		im, err := fingerprint.ToImage(r, w, h)

		if err != nil {
			panic(err)
		}

		im = fingerprint.ToAdobeRGB(im)
		im = fingerprint.AddBackground(im)

		out_fname := strings.Replace(fname, ".svg", ".png", 1)
		out_path := filepath.Join(root, out_fname)

		out, err := os.Create(out_path)

		if err != nil {
			panic(err)
		}
		defer out.Close()

		err = png.Encode(out, im)

		if err != nil {
			panic(err)
		}
	}
}
