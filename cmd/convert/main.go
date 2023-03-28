package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-fingerprint"
)

func main() {

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

		err = fingerprint.Convert(r, wr, 4096)

		if err != nil {
			log.Fatalf("Failed to derive info for %s, %v", out_path, err)
		}

		err = wr.Close()

		if err != nil {
			log.Fatalf("Failed to close %s, %v", out_path, err)
		}
	}
}
