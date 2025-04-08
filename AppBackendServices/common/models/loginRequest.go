package models

import "github.com/go-playground/validator/v10"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var loginVlaider = validator.New()

func (l *LoginRequest) Validate() error {
	return loginVlaider.Struct(l)
}
