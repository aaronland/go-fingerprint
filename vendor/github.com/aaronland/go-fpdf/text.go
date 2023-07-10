package fpdf

// type Text defines a struct for storing information for how text should be displayed in a .
type Text struct {
	// The name of the font to use for text strings.
	Font string
	// The style of the font to use for text strings.
	Style string
	// The size of the font to use for text strings.
	Size float64
	// The margin to apply to text strings.
	Margin float64
	// The colour of the font to use for text strings.
	Colour []int
}
