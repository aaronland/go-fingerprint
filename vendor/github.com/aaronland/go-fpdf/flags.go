package fpdf

import (
	"flag"
)

// String label defining the orientation of  PDF files. Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively.
var orientation string

// A common paper size to use for the size of your . Valid sizes are: "a3", "a4", "a5", "letter", "legal", or "tabloid".
var size string

// A width height to use as the size for a  PDF file.
var width float64

// A custom height to use as the size for a  PDF file.
var height float64

// The unit of measurement to apply to the height and width of a  PDF file.
var units string

// The "dots per inch" (DPI) resolution for a  PDF file.
var dpi float64

// The size of the border to apply to each image in a  PDF file.
var border float64

// The size of the margin to be applied to all sides of a .
var margin float64

// The size of the top margin for a .
var margin_top float64

// The size of the bottom margin for a .
var margin_bottom float64

// The size of the left margin for a .
var margin_left float64

// The size of the right margin for a .
var margin_right float64

// The size of an exterior "bleed" margin for a .
var bleed float64

var ocra_font bool

// Boolean flag to signal verbose logging during the creation of a .
var verbose bool

func AppendFlags(fs *flag.FlagSet) error {

	fs.StringVar(&orientation, "orientation", "P", "The orientation of your . Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively.")
	fs.StringVar(&size, "size", "letter", `A common paper size to use for the size of your . Valid sizes are: "a3", "a4", "a5", "letter", "legal", or "tabloid".`)

	fs.Float64Var(&width, "width", 0.0, "A custom height to use as the size of your . Units are defined in inches by default. This flag overrides the -size flag when used in combination with the -height flag.")

	fs.Float64Var(&height, "height", 0.0, "A custom width to use as the size of your . Units are defined in inches by default. This flag overrides the -size flag when used in combination with the -width flag.")

	fs.StringVar(&units, "units", "inches", "The unit of measurement to apply to the -height and -width flags. Valid options are inches, millimeters, centimeters")

	fs.Float64Var(&dpi, "dpi", 150, "The DPI (dots per inch) resolution for your .")

	fs.Float64Var(&border, "border", 0.01, "The size of the border around images.")

	fs.Float64Var(&margin_top, "margin-top", 1.0, "The margin around the top of each page.")
	fs.Float64Var(&margin_bottom, "margin-bottom", 1.0, "The margin around the bottom of each page.")
	fs.Float64Var(&margin_left, "margin-left", 1.0, "The margin around the left-hand side of each page.")
	fs.Float64Var(&margin_right, "margin-right", 1.0, "The margin around the right-hand side of each page.")
	fs.Float64Var(&margin, "margin", 0.0, "The margin around all sides of a page. If non-zero this value will be used to populate all the other -margin-(N) flags.")

	fs.Float64Var(&bleed, "bleed", 0.0, "An additional bleed area to add (on all four sides) to the size of your .")

	fs.BoolVar(&verbose, "verbose", false, "Display verbose output as the  is created.")

	fs.BoolVar(&ocra_font, "ocra-font", false, "Use an OCR-compatible font for captions.")

	return nil
}
