package fpdf

// type Margins defines a struct for storing margins to be applied to a PDF document.
type Margins struct {
	// The size of any margin to add to the top of each page.
	Top float64
	// The size of any margin to add to the bottom of each page.
	Bottom float64
	// The size of any margin to add to the left-hand side of each page.
	Left float64
	// The size of any margin to add to the right-hand side of each page.
	Right float64
}
