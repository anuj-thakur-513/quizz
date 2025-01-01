package middlewares

import (
	"strings"

	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/anuj-thakur-513/quizz/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var authToken string

func AdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			header := c.GetHeader("Authorization")
			if header == "" {
				c.JSON(401, core.NewAppError(401, "Unauthorized"))
				c.Abort()
				return
			}
			parts := strings.Split(header, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				c.JSON(401, core.NewAppError(401, "Invalid Authorization header format"))
				c.Abort()
				return
			}
			authToken = parts[1]
		} else {
			authToken = token
		}

		verified, err := utils.VerifyToken(authToken)

		if err != nil {
			c.JSON(401, core.NewAppError(401, "Unauthorized"))
			c.Abort()
			return
		}
		if !verified.Valid {
			c.JSON(401, core.NewAppError(401, "Token expired"))
			c.Abort()
			return
		}

		claims, ok := verified.Claims.(jwt.MapClaims)

		if !ok || claims["role"] != "admin" {
			c.JSON(401, core.NewAppError(401, "Invalid Token"))
		}

		c.Set("email", claims["email"])
		c.Next()
	}
}
