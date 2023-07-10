package fpdf

// type Borders defines a struct for storing borders to be applied to a images in a PDF document.
type Borders struct {
	// The size of any border to add to the top of each image.
	Top float64
	// The size of any border to add to the bottom of each image.
	Bottom float64
	// The size of any border to add to the left-hand side of each image.
	Left float64
	// The size of any border to add to the right-hand side of each image.
	Right float64
}
