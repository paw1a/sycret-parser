package storage

import (
	"errors"
)

const (
	FolderName = "sycret"
)

var (
	ErrInitCloud  = errors.New("failed to initialize cloudinary instance")
	ErrUploadFile = errors.New("failed to upload file")
)

func UploadDocument(docData []byte, filename string) (string, error) {
	//var cld, err = cloudinary.NewFromParams(viper.GetString("cloud_name"),
	//	viper.GetString("cloud_key"), viper.GetString("cloud_secret"))
	//if err != nil {
	//	return "", ErrInitCloud
	//}
	//
	//var ctx = context.Background()
	//
	//uploadResult, err := cld.Upload.Upload(
	//	ctx,
	//	bytes.NewBuffer(docData),
	//	uploader.UploadParams{
	//		PublicID: filename,
	//		Folder:   FolderName,
	//	})
	//if err != nil {
	//	return "", ErrUploadFile
	//}

	return "url.doc", nil
}
