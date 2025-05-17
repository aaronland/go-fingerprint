package image

import (
	"fmt"
	"image"
	"io"
	"time"

	"github.com/aaronland/go-image/v2/exif"
)

// AppendTime appends the RFC3339-encoded value or 't' to the JPEG encoding of 'im'
// in its `DateTime`, `DateTimeDigitized` and `DateTimeOriginal` EXIF headers. The
// final JPEG encoding of 'im' is written to 'wr'.
func AppendTime(im image.Image, wr io.Writer, t time.Time) error {

	// https://github.com/rwcarlsen/goexif/blob/go1/exif/exif.go#L385
	jpeg_dt := t.Format("2006:01:02 15:04:05") //  time.RFC3339)

	exif_props := map[string]interface{}{
		"DateTime":          jpeg_dt,
		"DateTimeDigitized": jpeg_dt,
		"DateTimeOriginal":  jpeg_dt,
		"Software":          "fingerprint",
	}

	err := exif.UpdateExif(im, wr, exif_props)

	if err != nil {
		return fmt.Errorf("Failed to update EXIF data, %w", err)
	}

	return nil
}
