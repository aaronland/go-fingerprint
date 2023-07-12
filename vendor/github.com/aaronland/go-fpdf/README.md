# go-fpdf

An opionated Go package for working with the `jung-kurt/gofpdf` package.

## Documentation

Documentation is incomplete at this time.

## Motivation

Most of this code is derived from the `aaronland/go-picturebook` package and is concerned with the details of creating a new `jung-kurt/gofpdf.Fpdf` document and assigning margins, borders, bleeds and fonts. It does implement or abstract away any of the `jung-kurt/gofpdf.Fpdf` functionality which is expected to be handled on a per-use basis accessing a public `gofpdf.Fpdf` instance.

The `jung-kurt/gofpdf.Fpdf` package stopped being updated a while ago. Since then the package has beeen forked and appears to be actively maintained by the `go-pdf/fpdf` project. Eventually this package will probably be updated to use that.

## See also

* https://github.com/jung-kurt/gofpdf
* https://github.com/go-pdf/fpdf