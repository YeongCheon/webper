package main

import (
	"bytes"
	"context"
	"github.com/disintegration/imaging"
	imgWebp "github.com/yeongcheon/webper/internal/image"
	"github.com/yeongcheon/webper/internal/storage"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"errors"
	"github.com/chai2010/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
)

func main() {
	log.Print("starting server...")
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	width, widthErr := strconv.Atoi(query.Get("width"))

	w.Header().Set("Cache-Control", "max-age=31536000, public")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if widthErr != nil {
		width = 0
	}

	imageName := strings.TrimPrefix(req.URL.Path, "/")

	var b bytes.Buffer
	r, err := storage.ReadImage(context.Background(), imageName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tee := io.TeeReader(r, &b)
	_, imageType, err := imgWebp.GetImageType(tee)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	io.ReadAll(tee)

	var resizedBuf bytes.Buffer

	if width > 0 {
		resize(&b, &resizedBuf, imageType, width)
	} else {
		resizedBuf = b
	}

	ctx := context.Background()

	switch imageType {
	case imgWebp.BMP:
		fallthrough
	case imgWebp.PNG:
		fallthrough
	case imgWebp.WEBP:
		fallthrough
	case imgWebp.JPG:
		err = imgWebp.RunCwebp(ctx, &resizedBuf, w)
	case imgWebp.GIF:
		err = imgWebp.RunGif2Webp(ctx, &resizedBuf, w)
	default:
		err = errors.New("invalid image format")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func resize(
	r io.Reader,
	w io.Writer,
	imageType imgWebp.ImageType,
	width int,
) error {
	img, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	var resizeImg image.Image

	if width <= 0 {
		resizeImg = img
	} else {
		resizeImg = imaging.Resize(img, width, 0, imaging.Lanczos)
	}

	switch imageType {
	case imgWebp.JPG:
		err = jpeg.Encode(w, resizeImg, nil)
	case imgWebp.PNG:
		err = png.Encode(w, resizeImg)
	case imgWebp.WEBP:
		err = webp.Encode(w, resizeImg, nil)
	case imgWebp.GIF:
		err = gif.Encode(w, resizeImg, nil)
	case imgWebp.BMP:
		// err = bmp.Encode(&tmp, resizeImg)
		fallthrough
	case imgWebp.TIFF:
		// err = tiff.Encode(&tmp, resizeImg, nil)
		fallthrough
	default:
		return errors.New("unknown file type")
	}

	return nil
}
