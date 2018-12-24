package gcs

// https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/storage/example_test.go
// https://github.com/GoogleCloudPlatform/golang-samples/blob/master/storage/objects/main.go

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
)

var (
	client      *storage.Client
	clientMutex = new(sync.Mutex)
)

func init() {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		fmt.Println("Missing GCS environment variable")
	}
}

func getClient() (c *storage.Client, ctx context.Context, err error) {

	clientMutex.Lock()

	ctx = context.Background()

	if client == nil {

		client, err = storage.NewClient(ctx)
		if err != nil {
			return client, ctx, nil
		}
	}

	clientMutex.Unlock()

	return client, ctx, nil
}

type UploadPayload struct {
	Bucket      string
	Path        string
	Transformer transformer
	Data        []byte
	Public      bool
}

func Upload(payload UploadPayload) (err error) {

	// Encode
	if payload.Transformer != nil {
		payload.Data, err = payload.Transformer(payload.Data)
		if err != nil {
			return err
		}
	}

	// Get client
	client, ctx, err := getClient()
	if err != nil {
		return err
	}

	object := client.Bucket(payload.Bucket).Object(strings.TrimLeft(payload.Path, "/"))

	// Upload bytes
	wc := object.NewWriter(ctx)
	_, err = io.Copy(wc, bytes.NewReader(payload.Data))
	if err != nil {
		return err
	}

	// Close writer
	err = wc.Close()
	if err != nil {
		return err
	}

	// Make public
	if payload.Public {
		return object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
	}

	return nil
}

type DownloadPayload struct {
	Bucket      string
	Path        string
	Transformer transformer
}

func Download(payload DownloadPayload) (bytes []byte, err error) {

	// Get client
	client, ctx, err := getClient()
	if err != nil {
		return bytes, err
	}

	// Download
	reader, err := client.Bucket(payload.Bucket).Object(strings.TrimLeft(payload.Path, "/")).NewReader(ctx)
	if err != nil {
		return bytes, err
	}

	bytes, err = ioutil.ReadAll(reader)
	if err != nil {
		return bytes, err
	}

	// Decode
	if payload.Transformer != nil {
		bytes, err = payload.Transformer(bytes)
		if err != nil {
			return bytes, err
		}
	}

	// Close reader
	err = reader.Close()

	return bytes, err
}
