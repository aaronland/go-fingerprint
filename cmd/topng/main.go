package main

import (
	"flag"
	"image/png"
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


	for _, path := range flag.Args(){

		root := filepath.Dir(path)
		fname := filepath.Base(path)
		
		r, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer r.Close()
		
		im, err := fingerprint.ToImage(r, w, h)
		
		if err != nil {
			panic(err)
		}

		im = fingerprint.ToAdobeRGB(im)
			
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
