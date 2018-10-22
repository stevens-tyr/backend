package api

import (
	"errors"
	"io"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

// gStorage types/structs

type gStorage struct {
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle
	ctx        context.Context
	w          io.Writer
}

var (
	ErrorUnableToWrite = errors.New("UNABLE TO WRITE TO BUCKET")
	ErrorUnableToClose = errors.New("UNABLE TO CLOSE BUCKET")
	ErrorUnableToOpen  = errors.New("UNABLE TO OPEN BUCKET FILE")
	ErrorUnableToRead  = errors.New("UNABLE TO READ FILE")
)

// Assignment Creation/Submission types/structs
type (
	createAssignment struct {
		lang  string `json:"language" binding:"requried"`
		ver   string `json:"version"`
		files string `json:"files" binding:"requried"`
	}

	submitAssignment struct {
		files `json:"" binding:"requried"`
	}
)
