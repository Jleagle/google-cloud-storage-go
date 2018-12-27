package main

import (
	"testing"

	"github.com/Jleagle/google-cloud-storage-go/gcs"
)

// Tests require you to manually create the bucket first

func TestUploadDownload(t *testing.T) {

	//
	uploadPayload := gcs.UploadPayload{}
	uploadPayload.Bucket = "testingx"
	uploadPayload.Path = "/dir/file.txt"
	uploadPayload.Data = []byte("testing")

	err := gcs.Upload(uploadPayload)

	if err != nil {
		t.Errorf(err.Error())
	}

	//
	downloadPayload := gcs.DownloadPayload{}
	downloadPayload.Bucket = uploadPayload.Bucket
	downloadPayload.Path = uploadPayload.Path

	b, err := gcs.Download(downloadPayload)

	if err != nil {
		t.Errorf(err.Error())
	}

	//
	if string(b) != string(uploadPayload.Data) {
		t.Errorf(err.Error())
	}
}

func TestUploadDownloadSnappy(t *testing.T) {

	//
	uploadPayload := gcs.UploadPayload{}
	uploadPayload.Bucket = "testingx"
	uploadPayload.Path = "/dir/file.txt"
	uploadPayload.Data = []byte("testing")
	uploadPayload.Transformer = gcs.TransformerSnappyEncode

	err := gcs.Upload(uploadPayload)

	if err != nil {
		t.Errorf(err.Error())
	}

	//
	downloadPayload := gcs.DownloadPayload{}
	downloadPayload.Bucket = uploadPayload.Bucket
	downloadPayload.Path = uploadPayload.Path
	downloadPayload.Transformer = gcs.TransformerSnappyDecode

	b, err := gcs.Download(downloadPayload)

	if err != nil {
		t.Errorf(err.Error())
	}

	//
	if string(b) != string(uploadPayload.Data) {
		t.Errorf(err.Error())
	}
}
