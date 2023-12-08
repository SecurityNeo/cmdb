package middleware

import (
	"cmdb/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.RequestURI == "/api/v1/login" {
			c.Next()
			return
		}

		if c.Request.RequestURI == "/api/v1/webssh/ssh" {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Missing header: Authorization",
				"data":    nil,
			})
			c.Abort()
			return
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Invalid token or token has expired",
					"data":    nil,
				})
				c.Abort()
				return
			}
			c.Set("user", claims)
			c.Next()
		}
	}
}
