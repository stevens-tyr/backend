package cms

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	//jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/stevens-tyr/tyr-gin"
)

func DownloadSubmission(c *gin.Context) {
	fdb, err := tyrgin.GetMongoDB(os.Getenv("GRIDFS_DB_NAME"))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get Mongo Session/DB.",
		"error":       err,
	}) {
		return
	}

	bucketSize, err := strconv.Atoi(os.Getenv("UPLOAD_SIZE"))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"staus_code": 500,
		"message":    "Failed to get gridfs bucket chunk size.",
		"error":      err,
	}) {
		return
	}

	bucket, err := tyrgin.GetGridFSBucket(fdb, fmt.Sprintf("%s%s", c.Param("cid"), c.Param("aid")), int32(bucketSize))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"staus_code": 500,
		"message":    "Failed to get assignments bucket.",
		"error":      err,
	}) {
		return
	}

	//claims := jwt.ExtractClaims(c)
	//submittedFileName := fmt.Sprintf("%s%s%s%s.tar.gz", c.Param("cid"), c.Param("aid"), claims["uid"].(string), c.Param("sid"))
	fileID, err := objectid.FromHex(c.Param("sid"))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"staus_code": 500,
		"message":    "Inavlid course submission id.",
		"error":      err,
	}) {
		return
	}

	file, err := bucket.GridFSDownloadFile(fileID)
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"staus_code": 500,
		"message":    "Failed to get submission.",
		"error":      err,
	}) {
		return
	}

	additonalHeaders := map[string]string {
		"Content-Disposition": fmt.Sprintf(`attachment; filename=" $%s-%s.tar.gz"`, c.Param("sid"), c.Param("num")) ,
	}

	c.DataFromReader(200, int64(file.Len()), "application/tar+gzi", bytes.NewReader(file.Bytes()), additonalHeaders)
}