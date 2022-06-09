package storage

import (
	"bytes"
	"context"
	"errors"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

const (
	CloudName   = "miragost"
	CloudKey    = "535661528615272"
	CloudSecret = "aHyZLIYT5rpMmkUljmewhfZlpuk"
	FolderName  = "sycret"
)

var (
	ErrInitCloud  = errors.New("failed to initialize cloudinary instance")
	ErrUploadFile = errors.New("failed to upload file")
)

func UploadDocument(docData []byte, filename string) (string, error) {
	var cld, err = cloudinary.NewFromParams(CloudName, CloudKey, CloudSecret)
	if err != nil {
		return "", ErrInitCloud
	}

	var ctx = context.Background()

	uploadResult, err := cld.Upload.Upload(
		ctx,
		bytes.NewBuffer(docData),
		uploader.UploadParams{
			PublicID: filename,
			Folder:   FolderName,
		})
	if err != nil {
		return "", ErrUploadFile
	}

	return uploadResult.URL, nil
}
