package fingerprint

import (
	"encoding/xml"
	"fmt"
	"io"
)

// https://pkg.go.dev/github.com/rustyoz/svg#Info

type Svg struct {
	Title   string `xml:"title"`
	Width   string `xml:"width,attr"`
	Height  string `xml:"height,attr"`
	ViewBox string `xml:"viewBox,attr"`
	Name    string
	Date    string `xml:"x-fingerprint-date,attr"`
}

func Info(r io.Reader) (*Svg, error) {

	var info *Svg

	dec := xml.NewDecoder(r)
	err := dec.Decode(&info)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode record, %w", err)
	}

	return info, nil
}
