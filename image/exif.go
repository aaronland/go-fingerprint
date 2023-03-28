package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"time"

	"github.com/sfomuseum/go-exif-update"
)

func AppendTime(im image.Image, wr io.Writer, t time.Time) error {

	temp_wr, err := os.CreateTemp("", "fingerprint.*.jpg")

	if err != nil {
		return fmt.Errorf("Failed to create temp file, %w", err)
	}

	defer os.Remove(temp_wr.Name())

	jpeg_opts := &jpeg.Options{
		Quality: 100,
	}

	err = jpeg.Encode(temp_wr, im, jpeg_opts)

	if err != nil {
		return fmt.Errorf("Failed to write JPEG, %w", err)
	}

	err = temp_wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close, %w", err)
	}

	jpeg_r, err := os.Open(temp_wr.Name())

	if err != nil {
		return fmt.Errorf("Failed to open %s, %v", temp_wr.Name(), err)
	}

	defer jpeg_r.Close()

	jpeg_dt := t.Format(time.RFC3339)

	exif_props := map[string]interface{}{
		"DateTime":          jpeg_dt,
		"DateTimeDigitized": jpeg_dt,
		"DateTimeOriginal":  jpeg_dt,
	}

	err = update.UpdateExif(jpeg_r, wr, exif_props)

	if err != nil {
		return fmt.Errorf("Failed to update EXIF data, %w", err)
	}

	return nil
}
