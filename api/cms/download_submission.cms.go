package cms

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func DownloadSubmission(c *gin.Context) {
	sid, _ := c.Get("sid")
	file, numBytes, err := gfs.Download(sid)
	if err != nil {
		c.Set("error", err)
		return
	}

	additonalHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s-%s.tar.gz"`, c.Param("sid"), c.Param("num")),
	}

	c.DataFromReader(200, numBytes, "application/tar+gzip", file, additonalHeaders)
}
