package image

import (
	"bytes"
	"context"
	"fmt"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"
)

type ImageType int

const (
	PNG ImageType = iota
	JPG
	GIF
	WEBP
	BMP
	TIFF
	// ICO
	ERR
)

// ex: cat color.jpg | cwebp -o - -- - > result.webp
func RunCwebp(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
) error {
	var stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, ".bin/webp/cwebp", "-o", "-", "--", "-")
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, stderr.String())
		return err
	}

	return nil
}

// ex: cat animation.gif | gif2webp -lossy -o - -- - > result.webp
func RunGif2Webp(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
) error {
	var stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, ".bin/webp/gif2webp", "-lossy", "-o", "-", "--", "-")
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, stderr.String())
		return err
	}

	return nil
}

func GetImageType(r io.Reader) (image.Config, ImageType, error) {
	config, f, err := image.DecodeConfig(r)

	var format ImageType

	switch f {
	case "jpeg":
		format = JPG
	case "png":
		format = PNG
	case "gif":
		format = GIF
	case "bmp":
		format = BMP
	case "webp":
		format = WEBP
	case "TIFF":
		format = TIFF
	default:
		format = ERR
	}

	return config, format, err
}
