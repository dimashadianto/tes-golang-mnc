package middleware

import (
	"fmt"
	"net/http"
	"os"
	"tes-mnc/initializers"
	"tes-mnc/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication token not provided",
		})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid token",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if expirationTime, ok := claims["exp"].(float64); ok {
			if float64(time.Now().Unix()) > expirationTime {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":   "Unauthorized",
					"message": "Token has expired",
				})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid expiration time in the token",
			})
			return
		}

		var user models.User
		if result := initializers.DB.First(&user, claims["sub"]); result.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not found",
			})
			return
		}

		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid token",
		})
	}
}
