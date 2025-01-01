package controllers

import (
	"context"

	"github.com/go-playground/validator/v10"
)

var ctx context.Context = context.Background()
var validate = validator.New()
