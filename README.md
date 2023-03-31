# go-fingerprint

Go package for working with SVG files produced by the `aaronland/fingerprint` drawing tool

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-fingerprint.svg)](https://pkg.go.dev/github.com/aaronland/go-fingerprint)

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/convert cmd/convert/main.go
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

## See also

* https://github.com/aaronland/fingerprint
* https://github.com/fogleman/gg
* https://github.com/mandykoh/prism
* https://github.com/sfomuseum/go-exif-update