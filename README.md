# go-fingerprint

Go package for working with SVG files produced by the `aaronland/fingerprint` drawing tool

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-fingerprint.svg)](https://pkg.go.dev/github.com/aaronland/go-fingerprint)

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/convert cmd/convert/main.go
go build -mod vendor -ldflags="-s -w" -o bin/pdf cmd/pdf/main.go
```

### convert

```
$> ./bin/convert -h
Convert one or more fingerprint SVG documents in to JPEG images.

The final JPEG image is scaled to ensure that its maximum dimension is
'max_dimension'. Date information defined in the SVG document's `x-fingerprint-date`
attribute is written to the final JPEG image's `DateTime`, `DateTimeDigitized` and
`DateTimeOriginal` EXIF headers. The final JPEG representation is updated to ensure that
all pixel values match the Adobe RGB colour profile. JPEG images are written to the same
location as the source SVG document with a '.jpg' extension.

Usage:
	 ./bin/convert path(N) path(N)

Valid options are:
  -max-dimension float
    	The maximum dimension to scale an image to. (default 4096)
```

### pdf

```
$> ./bin/pdf -h
Convert one or more fingerprint SVG documents into PDF documents.

Each document contains a full-page rasterized rendering of the SVG document
followed by one or more pages containing the body of the aaronland/go-fingerprint/svg.Document
representation of the SVG document as JSON-encoded text.

The final PDF document will be saved as (svg-path).pdf.

Usage:
	 ./bin/pdf [options] path(N) path(N)

Valid options are:
  -bleed float
    	An additional bleed area to add (on all four sides) to the size of your .
  -border float
    	The size of the border around images. (default 0.01)
  -dpi float
    	The DPI (dots per inch) resolution for your . (default 150)
  -height float
    	A custom width to use as the size of your . Units are defined in inches by default. This flag overrides the -size flag when used in combination with the -width flag.
  -margin float
    	The margin around all sides of a page. If non-zero this value will be used to populate all the other -margin-(N) flags.
  -margin-bottom float
    	The margin around the bottom of each page. (default 1)
  -margin-left float
    	The margin around the left-hand side of each page. (default 1)
  -margin-right float
    	The margin around the right-hand side of each page. (default 1)
  -margin-top float
    	The margin around the top of each page. (default 1)
  -ocra-font
    	Use an OCR-compatible font for captions.
  -orientation string
    	The orientation of your . Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively. (default "P")
  -size string
    	A common paper size to use for the size of your . Valid sizes are: "a3", "a4", "a5", "letter", "legal", or "tabloid". (default "letter")
  -units string
    	The unit of measurement to apply to the -height and -width flags. Valid options are inches, millimeters, centimeters (default "inches")
  -verbose
    	Display verbose output as the  is created.
  -width float
    	A custom height to use as the size of your . Units are defined in inches by default. This flag overrides the -size flag when used in combination with the -height flag.
```

## See also

* https://github.com/aaronland/fingerprint
* https://github.com/fogleman/gg
* https://github.com/mandykoh/prism
* https://github.com/sfomuseum/go-exif-update
* https://github.com/aaronland/go-fpdf