package cms

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	//jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/stevens-tyr/tyr-gin"
)

func DownloadSubmission(c *gin.Context) {
	fdb, err := tyrgin.GetMongoDB(os.Getenv("GRIDFS_DB_NAME"))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"message":     "Failed to get Mongo Session/DB.",
			"error":       err.Error(),
		})
		return
	}

	bucketSize, err := strconv.Atoi(os.Getenv("UPLOAD_SIZE"))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"staus_code": 500,
			"message":    "Failed to get gridfs bucket chunk size.",
			"error":      err,
		})
		return
	}

	bucket, err := tyrgin.GetGridFSBucket(fdb, "fs", int32(bucketSize))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"staus_code": 500,
			"message":    "Failed to get assignments bucket.",
			"error":      err,
		})
		return
	}

	//claims := jwt.ExtractClaims(c)
	//submittedFileName := fmt.Sprintf("%s%s%s%s.tar.gz", c.Param("cid"), c.Param("aid"), claims["uid"].(string), c.Param("sid"))
	fileID, err := primitive.ObjectIDFromHex(c.Param("sid"))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"staus_code": 500,
			"message":    "Inavlid submission id.",
			"error":      err,
		})
		return
	}

	file, err := bucket.GridFSDownloadFile(fileID)
	fmt.Println("after", fileID, err)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"staus_code": 500,
			"message":    "Failed to get submission.",
			"error":      err,
		})
		return
	}

	additonalHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename=" $%s-%s.tar.gz"`, c.Param("sid"), c.Param("num")),
	}

	fmt.Println(file.Bytes())

	c.DataFromReader(200, int64(len(file.Bytes())), "application/tar+gzi", bytes.NewReader(file.Bytes()), additonalHeaders)
}
