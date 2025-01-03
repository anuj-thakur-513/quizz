package middlewares

import (
	"context"
	"strings"

	"github.com/anuj-thakur-513/quizz/internal/models"
	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/anuj-thakur-513/quizz/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authToken string
		cookie, err := c.Cookie("token")
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
			authToken = cookie
		}

		verified, err := utils.VerifyToken(authToken)
		if err != nil {
			c.JSON(401, core.NewAppError(401, "Unauthorized"))
			c.Abort()
			return
		}
		claims, ok := verified.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, core.NewAppError(401, "Invalid Token"))
		}

		users := models.GetUsersCollection()
		var user *models.User

		if err := users.FindOne(context.Background(), bson.M{"email": claims["email"]}).Decode(&user); err != nil {
			c.JSON(401, core.NewAppError(401, "Unauthorized"))
			c.Abort()
			return
		}

		c.Set("email", claims["email"])
		c.Set("user", user)
		c.Next()
	}
}
