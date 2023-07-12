package fpdf

import (
	"context"
	"flag"
)

// Options defines a struct containing configuration information for a given  instance.
type Options struct {
	// The orientation of the final . Valid options are "P" and "L" for portrait and landscape respectively.
	Orientation string
	// A string label corresponding to known size. Valid options are "a1", "a2", "a3", "a4", "a5", "a6", "a7", "letter", "legal" and "tabloid".
	Size string
	// The width of the final .
	Width float64
	// The height of the final .
	Height float64
	// The unit of measurement to use for the `Width` and `Height` options.
	Units string
	// The number dots per inch to use when calculating the size of the final . Valid options are "inches", "centimeters", "millimeters".
	DPI float64
	// The size of any border to apply to each image in the final .
	Border float64
	// The size of any additional bleed to apply to the final .
	Bleed float64
	// The size of any margin to add to the top of each page.
	MarginTop float64
	// The size of any margin to add to the bottom of each page.
	MarginBottom float64
	// The size of any margin to add to the left-hand side of each page.
	MarginLeft float64
	// The size of any margin to add to the right-hand side of each page.
	MarginRight float64
	// An optional `filter.Filter` instance used to determine whether or not an image should be included in the final .
	Verbose bool
	// A boolean value to enable to use of an OCRA font for writing captions.
	OCRAFont bool
	// A boolean value signaling that images should only be added on even-numbered pages.
	MaxPages int
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*Options, error) {

	opts := &Options{
		Orientation:  orientation,
		Size:         size,
		Width:        width,
		Height:       height,
		Units:        units,
		DPI:          dpi,
		Border:       border,
		Bleed:        bleed,
		MarginTop:    margin_top,
		MarginBottom: margin_bottom,
		MarginLeft:   margin_left,
		MarginRight:  margin_right,
		Verbose:      verbose,
	}

	return opts, nil
}

// NewDefaultOptions returns a `Options` with default settings.
func DefaultOptions(ctx context.Context) (*Options, error) {

	opts := &Options{
		Orientation:  "P",
		Size:         "letter",
		Width:        0.0,
		Height:       0.0,
		Units:        "inches",
		DPI:          150.0,
		Border:       0.01,
		Bleed:        0.0,
		MarginTop:    1.0,
		MarginBottom: 1.0,
		MarginLeft:   1.0,
		MarginRight:  1.0,
		Verbose:      false,
	}

	return opts, nil
}
