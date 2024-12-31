package controllers

import (
	"context"

	"github.com/anuj-thakur-513/quizz/internal/models"
	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/anuj-thakur-513/quizz/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const DB_NAME = utils.DB_NAME

var ctx context.Context = context.Background()
var validate = validator.New()

func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid JSON body", err.Error()))
		return
	}
	newUser.PreSave()

	if err := validate.Struct(newUser); err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid JSON body", err.Error()))
		return
	}

	users := models.GetUsersCollection()
	if _, err := users.InsertOne(ctx, newUser); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to create user", err.Error()))
		return
	}

	token, err := utils.GenerateToken(newUser.Email)
	if err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to create user", err.Error()))
	}
	c.SetCookie("token", token, 3600*24*30, "/", "localhost", true, true)

	c.JSON(200, core.ApiResponse(200, "User created successfully", nil))
}
