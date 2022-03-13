package storage

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"log"
	"os"
)

var client *storage.Client
var bucket *storage.BucketHandle

func init() {
	var err error
	client, err = storage.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	bucketName := os.Getenv("BUCKET_NAME")
	bucket = client.Bucket(bucketName)
}

func ReadImage(ctx context.Context, path string) (io.Reader, error) {
	originalImage := bucket.Object(path)
	return originalImage.NewReader(ctx)
}
