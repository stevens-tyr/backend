package gridfsmodels

import (
	"bytes"
	"io"
	"os"
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	"backend/errors"

	"github.com/stevens-tyr/tyr-gin"
)

type (
	GridFSInterface struct {
		bucket *tyrgin.Bucket
		db     *mongo.Database
	}
)

func New() *GridFSInterface {
	db, _ := tyrgin.GetMongoDB(os.Getenv("GRIDFS_DB_NAME"))
	bucketSize, _ := strconv.Atoi(os.Getenv("UPLOAD_SIZE"))
	bucket, _ := tyrgin.GetGridFSBucket(db, "assignments", int32(bucketSize))

	return &GridFSInterface{
		bucket,
		db,
	}
}

func (g *GridFSInterface) Upload(id *primitive.ObjectID, filename string, file io.Reader) errors.APIError {
	var nid primitive.ObjectID
	if id == nil {
		nid = primitive.NewObjectID()
	} else {
		nid = *id
	}

	err := g.bucket.GridFSUploadFile(nid, filename, file)
	if err != nil {
		return errors.ErrorGridFSUploadFailure
	}

	return nil
}

func (g *GridFSInterface) Delete(fileID interface{}) errors.APIError {
	err := g.bucket.GridFSDeleteFile(fileID.(primitive.ObjectID))
	if err != nil {
		return errors.ErrorGridFSDeleteFailure
	}

	return nil
}

func (g *GridFSInterface) Download(fileID interface{}) (*bytes.Reader, int64, errors.APIError) {
	file, err := g.bucket.GridFSDownloadFile(fileID.(primitive.ObjectID))
	if err != nil {
		return nil, 0, errors.ErrorGridFSDownloadFailure
	}

	return bytes.NewReader(file.Bytes()), int64(file.Len()), nil
}
