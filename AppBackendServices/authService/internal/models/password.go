package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type (
	ForgotPasswordRequest struct {
		Email       string `json:"email" validate:"required,email"`
		Code        string `json:"code" validate:"required,len=6"`
		NewPassword string `json:"new_password" validate:"required,len=8"`
	}

	ResetPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	ResetPasswordVerificationRequest struct {
		Email            string `json:"email" validate:"required,email"`
		Token            string `json:"token" validate:"required"`
		NewPassword      string `json:"new_password" validate:"required,min=8"`
		PreviousPassword string `json:"previous_password" validate:"required,min=8"`
	}

	ChangePasswordRequest struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=8"`
	}

	PasswordResetToken struct {
		Token     string    `json:"token"`
		Email     string    `json:"email"`
		ExpiresAt time.Time `json:"expires_at"`
	}
)

var passwordValidator = validator.New()

func (f *ForgotPasswordRequest) Validate() error {
	return passwordValidator.Struct(f)
}

func (r *ResetPasswordRequest) Validate() error {
	return passwordValidator.Struct(r)
}

func (c *ChangePasswordRequest) Validate() error {
	return passwordValidator.Struct(c)
}
