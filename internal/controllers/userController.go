package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/anuj-thakur-513/quizz/internal/models"
	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/anuj-thakur-513/quizz/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

const DB_NAME = utils.DB_NAME

var ctx context.Context = context.Background()
var validate = validator.New()

func Signup(c *gin.Context) {
	var newUser *models.User
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

	c.JSON(200, core.ApiResponse(200, "User created successfully", map[string]string{
		"email": newUser.Email,
		"name":  newUser.Name,
	}))
}

func Login(c *gin.Context) {
	ginBody := c.Request.Body
	// convert to json
	body, err := io.ReadAll(ginBody)
	if err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid JSON body", err.Error()))
		return
	}
	// convert body to struct
	jsonBody := map[string]interface{}{}
	if err = json.Unmarshal(body, &jsonBody); err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid JSON body", err.Error()))
		return
	}

	email := jsonBody["email"].(string)
	password := jsonBody["password"].(string)

	if email == "" || password == "" {
		c.JSON(400, core.NewAppError(400, "Invalid JSON body", "email and password are required"))
		return
	}
	if !utils.ValidateEmail(email) {
		c.JSON(400, core.NewAppError(400, "Invalid JSON body", "email is invalid"))
		return
	}

	users := models.GetUsersCollection()
	var user models.User
	dbErr := users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if dbErr != nil {
		log.Printf("MongoDB query error: %v", dbErr)
		c.JSON(500, core.NewAppError(500, "Database error", dbErr.Error()))
		return
	}

	if !user.ComparePassword(password) {
		c.JSON(401, core.NewAppError(401, "Unauthorized", "Password and Email don't match"))
		return
	}
	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to generate token", err.Error()))
		return
	}
	fmt.Println("is live reload working")
	c.SetCookie("token", token, 3600*24*30, "/", "localhost", true, true)

	c.JSON(200, core.ApiResponse(200, "User logged in successfully", nil))
}
