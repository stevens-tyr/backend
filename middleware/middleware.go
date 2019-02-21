package middleware

import (
	"backend/errors"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/stevens-tyr/tyr-gin"
)

func ObjectIDs() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Param("aid") != "" {
			val, err := primitive.ObjectIDFromHex(c.Param("aid"))
			if err != nil {
				c.AbortWithStatusJSON(
					errors.ErrorInvalidObjectID.StatusCode(),
					gin.H{
						"error": errors.ErrorInvalidObjectID.Error(),
					},
				)
			}

			c.Set("aid", val)
		}

		if c.Param("cid") != "" {
			val, err := primitive.ObjectIDFromHex(c.Param("cid"))
			if err != nil {
				c.AbortWithStatusJSON(
					errors.ErrorInvalidObjectID.StatusCode(),
					gin.H{
						"error": errors.ErrorInvalidObjectID.Error(),
					},
				)
			}

			c.Set("cids", val.Hex())
			c.Set("cid", val)
		}

		if c.Param("sid") != "" {
			val, err := primitive.ObjectIDFromHex(c.Param("sid"))
			if err != nil {
				c.AbortWithStatusJSON(
					errors.ErrorInvalidObjectID.StatusCode(),
					gin.H{
						"error": errors.ErrorInvalidObjectID.Error(),
					},
				)
			}

			c.Set("sid", val)
		}

		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		val, exists := c.Get("error")
		if !exists {
			return
		}

		apierr := val.(errors.APIError)
		if apierr.GetError() != nil {
			tyrgin.ErrorHandler(apierr.GetError(), c, apierr.StatusCode(), gin.H{
				"error": apierr.Error(),
			})
		}
	}
}
