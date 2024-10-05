package middlewares

import (
	"net/http"
	"strings"

	tokens "working/super_task/package/token"

	"github.com/gin-gonic/gin"
)

// method for working with jwt based authentication
func JwtAuthMiddleWare(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Bearer token is missing"})
			c.Abort()
			return
		}
		_, err := tokens.VerifyToken(tokenString, secret)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		claims, err := tokens.GetUserClaims(tokenString, secret)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		if username, ok := claims["username"]; ok {
			c.Set("username", username)
		} else {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Username not found in token claims"})
			c.Abort()
			return
		}

		if id, ok := claims["id"]; ok {
			c.Set("id", id)
		} else {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "ID not found in token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}
