package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func main() {
	var (
		bucket   string
		keyFile  string
		object   string
		fileName string
	)

	flag.StringVar(&bucket, "b", "", "Bucket name")
	flag.StringVar(&keyFile, "k", "", "GCP key.json file")
	flag.StringVar(&object, "o", "", "object to download")
	flag.StringVar(&fileName, "f", "", "filename to call object")
	flag.Parse()

	if bucket == "" {
		fmt.Println("Bucket Required")
		os.Exit(1)
	}

	if keyFile == "" {
		fmt.Println("GCP key.json required")
		os.Exit(1)
	}

	if fileName == "" {
		fmt.Println("Filename required")
		os.Exit(1)
	}

	err := downloadFile(os.Stdout, bucket, object, fileName, keyFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

// downloadFile downloads an object to a file.
func downloadFile(w io.Writer, bucket, object string, destFileName, jsonPath string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	f, err := os.Create(destFileName)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %v", err)
	}

	fmt.Fprintf(w, "Blob %v downloaded to local file %v\n", object, destFileName)

	return nil

}
