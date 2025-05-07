package model

import (
	"github.com/google/uuid"
	"time"
)

// Authentication Models
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	UserId       uuid.UUID `json:"user_id"`
}

type RegisterRequest struct {
	FullName      string `json:"full_name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8"`
	BioMetricHash string `json:"bio_metric_hash" validate:"required,min=8"`
}

type RegisterResponse struct {
	UserID            string `json:"user_id"`
	VerificationToken string `json:"verification_token,omitempty"`
}

type RegisterVerificationRequest struct {
	Token string `json:"token" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type LogoutRequest struct {
	AccessToken  string `json:"-"` // From header
	RefreshToken string `json:"-"` // From cookie
}

// Password Management Models
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetCodeVerificationRequest struct {
	Token       string `json:"token" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordCodeVerificationRequest struct {
	Token       string `json:"token" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// Email Models
type EmailVerification struct {
	To        string    `json:"to" validate:"required,email"`
	Code      string    `json:"code" validate:"required"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Token Models
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenBlacklist struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Response Wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
