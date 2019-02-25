package cms

import (
  "fmt"

  "github.com/gin-gonic/gin"
)

func AssignmentAsFile(c *gin.Context) {
  aid, _ := c.Get("aid")

  file, filename, numBytes, err := am.AsFile(aid)
  if err != nil {
    c.Set("error", err)
    return
  }

  additonalHeaders := map[string]string{
    "Content-Disposition": fmt.Sprintf(`attachment; filename="%s.json"`, filename),
  }

  c.DataFromReader(200, numBytes, "application/tar+gzip", file, additonalHeaders)
}
