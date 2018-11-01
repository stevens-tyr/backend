package utils

import (
	"context"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"

	"backend/models"
)

// GStorage to stor ingormation about google storage bucket info.
type GStorage struct {
	Client     *storage.Client
	BucketName string
	Bucket     *storage.BucketHandle
	Ctx        context.Context
	W          io.Writer
}

// UploadFile takes a file name and file content to upload to gcp. It also
// takes a meta type to help us mark files as teacher submissions or student.
func (g *GStorage) UploadFile(fileName, metaType string, fileContent []byte) (string, error) {
	wc := g.Bucket.Object(fileName).NewWriter(g.Ctx)
	wc.ContentType = "application/tar+gzip"
	wc.Metadata = map[string]string{
		// to specify teacher files or student submission as metadata
		"metaType": metaType,
	}

	if _, err := wc.Write(fileContent); err != nil {
		return "", models.ErrorUnableToWrite
	}

	if err := wc.Close(); err != nil {
		return "", models.ErrorUnableToClose
	}

	return "", nil
}

// ReadFile allows us to read a file from gc storage.
// Returns a string and error, for message purposes if an error occurs.
func (g *GStorage) ReadFile(fileName string) (content []byte, err error) {
	rc, err := g.Bucket.Object(fileName).NewReader(g.Ctx)
	if err != nil {
		return content, models.ErrorUnableToOpen
	}
	defer rc.Close()

	content, err = ioutil.ReadAll(rc)
	if err != nil {
		return content, models.ErrorUnableToRead
	}

	return content, nil
}
