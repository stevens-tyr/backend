package cms

import (
  "github.com/gin-gonic/gin"

  "backend/errors"
  "backend/forms"
)

func CourseAddUsers(c *gin.Context) {
  cid, _ := c.Get("cid")

  var addUsers forms.CourseBulkAddUserForm
  if err := c.ShouldBindJSON(&addUsers); err != nil {
    c.Set("error", errors.ErrorInvlaidJSON)
    return
  }

  for _, email := range addUsers.Emails {
    user, err := um.FindOne(email)
    if err != nil {
      c.Set("error", err)
      return
    }

		err = um.AddCourse(addUsers.Level, cid, user.ID)
    if err != nil {
      c.Set("error", err)
      return
    }

    err = cm.AddUser(addUsers.Level, user.ID, cid)
    if err != nil {
      c.Set("error", err)
      return
    }
  }

  c.JSON(200, gin.H{
    "message": "User added.",
  })
}
