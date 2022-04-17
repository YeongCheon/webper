package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/yeongcheon/webper/internal/image"
	"github.com/yeongcheon/webper/internal/storage"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
	imageName := strings.TrimPrefix(req.URL.Path, "/")
	log.Println(imageName)

	var b bytes.Buffer
	r, err := storage.ReadImage(context.Background(), imageName)
	if err != nil {
		panic(err)
	}

	tee := io.TeeReader(r, &b)
	_, imageType, err := image.GetImageType(tee)
	if err != nil {
		panic(err)
	}
	io.ReadAll(tee)

	ctx := context.Background()

	switch imageType {
	case image.BMP:
		fallthrough
	case image.PNG:
		fallthrough
	case image.WEBP:
		fallthrough
	case image.JPG:
		err = image.RunCwebp(ctx, &b, w)
	case image.GIF:
		err = image.RunGif2Webp(ctx, &b, w)
	default:
		err = errors.New("invalid image format")
	}

	// if err != nil {
	// 	panic(err)
	// }
}
