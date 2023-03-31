package svg

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Unmarshal will parse the fingerprint SVG document defined by 'r' in a new `Document` instance.
func Unmarshal(r io.Reader) (*Document, error) {

	var doc *Document

	dec := xml.NewDecoder(r)
	err := dec.Decode(&doc)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode record, %w", err)
	}

	return doc, nil
}
