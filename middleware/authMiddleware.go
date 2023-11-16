package middleware

import (
	helper "golangsidang/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Hanya memeriksa token jika metode permintaan adalah GET
		if c.Request.Method == "" {
			clientToken := c.Request.Header.Get("Authorization")
			if clientToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Token not found"})
				c.Abort()
				return
			}
			claims, err := helper.ValidateToken(clientToken)
			if err != "" {
				c.JSON(http.StatusUnauthorized, gin.H{"message": err})
				c.Abort()
				return
			}
			// Set data pengguna dalam konteks jika token valid
			c.Set("email", claims.Email)
			c.Set("first_name", claims.First_name)
			c.Set("last_name", claims.Last_name)
			c.Set("uid", claims.Uid)
			c.Set("user_type", claims.User_type)
		}
		c.Next()
	}
}
