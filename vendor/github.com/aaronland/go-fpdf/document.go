package fpdf

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/jung-kurt/gofpdf"
	"github.com/sfomuseum/go-font-ocra"
)

// MM2INCH defines the number if millimeters in an inch.
const MM2INCH float64 = 25.4

// type Document provides a struct for creating a PDF file from a folder of images (a picturebook).
type Document struct {
	// A `gofpdf.Fpdf` instance used to produce the picturebook PDF file.
	PDF *gofpdf.Fpdf
	// A `sync.Mutex` instance used to add images in an orderly fashion.
	Mutex *sync.Mutex
	// The `PictureBookBorders` definition to use for this picturebook
	Borders *Borders
	// The `PictureBookMargins` definition to use for this picturebook
	Margins *Margins
	// The `PictureBookCanvas` definition to use for this picturebook
	Canvas Canvas
	// The `PictureBookText` definition to use for this picturebook
	Text Text
	// The `PictureBookOptions` used to create this picturebook
	Options *Options
	// The number of pages in this picturebook
	pages int
	// A list of temporary files used in the creation of a picturebook and to be removed when the picturebook is saved
	tmpfiles []string
}

// New returns a new “ instances configured according to the settings in 'opts'.
func NewDocument(ctx context.Context, opts *Options) (*Document, error) {

	var pdf *gofpdf.Fpdf

	// opts_w := opts.Width
	// opts_h := opts.Height
	// opts_b := opts.Bleed

	// Start by convert everything to inches - not because it's better but
	// just because it's expedient right now (20210218/straup)

	if opts.Width == 0.0 && opts.Height == 0.0 {

		switch strings.ToLower(opts.Size) {
		case "a1":
			opts.Width = 584.0 / MM2INCH
			opts.Height = 841.0 / MM2INCH
		case "a2":
			opts.Width = 420 / MM2INCH
			opts.Height = 594 / MM2INCH
		case "a3":
			opts.Width = 297 / MM2INCH
			opts.Height = 420 / MM2INCH
		case "a4":
			opts.Width = 210.0 / MM2INCH
			opts.Height = 297.0 / MM2INCH
		case "a5":
			opts.Width = 148 / MM2INCH
			opts.Height = 210 / MM2INCH
		case "a6":
			opts.Width = 105 / MM2INCH
			opts.Height = 148 / MM2INCH
		case "a7":
			opts.Width = 74 / MM2INCH
			opts.Height = 105 / MM2INCH
		case "letter":
			opts.Width = 8.5
			opts.Height = 11.0
		case "legal":
			opts.Width = 11.0
			opts.Height = 17.0
		case "tabloid":
			opts.Width = 11.0
			opts.Height = 17.0
		default:
			return nil, fmt.Errorf("Unrecognized page size '%s'", opts.Size)
		}
	} else {

		switch opts.Units {
		case "inches":
			// pass
		case "millimeters":
			opts.Width = opts.Width / MM2INCH
			opts.Height = opts.Height / MM2INCH
		case "centimeters":
			opts.Width = (opts.Width * 10.0) / MM2INCH
			opts.Height = (opts.Height * 10.0) / MM2INCH
		default:
			return nil, fmt.Errorf("Invalid or unsupported unit '%s'", opts.Units)
		}
	}

	// log.Printf("%0.2f x %0.2f (%s)\n", opts.Width, opts.Height, opts.Size)

	sz := gofpdf.SizeType{
		Wd: opts.Width + (opts.Bleed * 2.0),
		Ht: opts.Height + (opts.Bleed * 2.0),
	}

	init := gofpdf.InitType{
		OrientationStr: opts.Orientation,
		UnitStr:        "in",
		SizeStr:        "",
		Size:           sz,
		FontDirStr:     "",
	}

	pdf = gofpdf.NewCustom(&init)

	t := Text{
		Font:   "Helvetica",
		Style:  "",
		Size:   8.0,
		Margin: 0.1,
		Colour: []int{128, 128, 128},
	}

	if opts.OCRAFont {

		font, err := ocra.LoadFPDFFont()

		if err != nil {
			return nil, fmt.Errorf("Failed to load OCRA font, %w", err)
		}

		pdf.AddFontFromBytes(font.Family, font.Style, font.JSON, font.Z)
		pdf.SetFont(font.Family, "", 8.0)

		pdf.SetTextColor(t.Colour[0], t.Colour[1], t.Colour[2])

	} else {

		pdf.SetFont(t.Font, t.Style, t.Size)
	}

	w, h, _ := pdf.PageSize(1)

	page_w := w * opts.DPI
	page_h := h * opts.DPI

	// https://github.com/aaronland/go-picturebook/issues/22

	// margin around each page (inclusive of page bleed)

	margin_top := (opts.MarginTop + (opts.Bleed * 2.0)) * opts.DPI
	margin_bottom := (opts.MarginBottom + (opts.Bleed * 2.0)) * opts.DPI
	margin_left := (opts.MarginLeft + (opts.Bleed * 2.0)) * opts.DPI
	margin_right := (opts.MarginRight + (opts.Bleed * 2.0)) * opts.DPI

	margins := &Margins{
		Top:    margin_top,
		Bottom: margin_bottom,
		Left:   margin_left,
		Right:  margin_right,
	}

	// border around each image

	border_top := opts.Border * opts.DPI
	border_bottom := opts.Border * opts.DPI
	border_left := opts.Border * opts.DPI
	border_right := opts.Border * opts.DPI

	borders := &Borders{
		Top:    border_top,
		Bottom: border_bottom,
		Left:   border_left,
		Right:  border_right,
	}

	// Remember: margins have been calculated inclusive of page bleeds

	canvas_w := page_w - (margin_left + margin_right + border_left + border_right)
	canvas_h := page_h - (margin_top + margin_bottom + border_top + border_bottom)

	// pdf.SetAutoPageBreak(false, border_bottom)

	canvas := Canvas{
		Width:  canvas_w,
		Height: canvas_h,
	}

	left := margins.Left / opts.DPI
	right := margins.Right / opts.DPI
	top := margins.Top / opts.DPI
	bottom := margins.Bottom / opts.DPI

	pdf.SetMargins(left, top, right)
	pdf.SetAutoPageBreak(true, bottom)
	
	tmpfiles := make([]string, 0)
	mu := new(sync.Mutex)

	pb := Document{
		PDF:      pdf,
		Mutex:    mu,
		Borders:  borders,
		Margins:  margins,
		Canvas:   canvas,
		Text:     t,
		Options:  opts,
		pages:    0,
		tmpfiles: tmpfiles,
	}

	return &pb, nil
}

func (d *Document) Save(path string) error {
	return d.PDF.OutputFileAndClose(path)
}
