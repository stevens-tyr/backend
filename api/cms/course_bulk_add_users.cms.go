package cms

import (
  "fmt"
  "github.com/gin-gonic/gin"

  "backend/errors"
  "backend/forms"
)

func CourseAddUsers(c *gin.Context) {
  cid, _ := c.Get("cid")

  var addUsers []forms.CourseAddUserForm
  if err := c.ShouldBindJSON(&addUsers); err != nil {
    c.Set("error", errors.ErrorInvlaidJSON)
    return
  }

  for _, addUser := range addUsers {
    fmt.Println(addUser.Email)

    user, err := um.FindOne(addUser.Email)
    if err != nil {
      c.Set("error", err)
      return
    }
    fmt.Println("plz")

    err = um.AddCourse(addUser.Level, cid, user.ID)
    if err != nil {
      c.Set("error", err)
      return
    }

    err = cm.AddUser(addUser.Level, user.ID, cid)
    if err != nil {
      c.Set("error", err)

      return
    }
  }

  c.JSON(200, gin.H{
    "message": "User added.",
  })
}
