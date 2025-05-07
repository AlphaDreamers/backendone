package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/model"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/repo"
	"github.com/SwanHtetAungPhyo/service-one/auth/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"sync"
	"time"
)

var ServiceModule = fx.Provide(
	NewAuthService,
	NewUserService,
)

const (
	refreshTokenPrefix         = "refresh_token_"
	emailVerificationPrefix    = "email_verification_"
	userEmailVerificationTopic = "user.email_verification"
	blacklistTokenPrefix       = "blacklist_token_"
	tokenExpiration            = 24 * time.Hour * 7
	emailTokenExpiration       = 10 * time.Minute
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailNotVerified      = errors.New("email not verified")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTokenGenerationFailed = errors.New("token generation failed")
	ErrTokenBlacklisted      = errors.New("token is blacklisted")
)

type AuthServiceBehaviour interface {
	Login(ctx context.Context, req *model.LoginRequest, deviceId string) (*model.LoginResponse, error)
	Register(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error)
	VerifyEmail(ctx context.Context, email, token string) error
	ResetPassword(ctx context.Context, email string) error
	VerifyResetPasswordToken(ctx context.Context, token, email string) error
	ForgotPassword(ctx context.Context, email string) error
	VerifyForgotPasswordToken(ctx context.Context, token, email, newPassword string) error
	RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
}

type AuthService struct {
	log          *logrus.Logger
	repo         *repo.AuthRepo
	jwtGenerator *utils.JwtTokenGenerator
	redisClient  *redis.Client
	natsConn     *nats.Conn
	wg           sync.WaitGroup
}

func NewAuthService(
	log *logrus.Logger,
	repo *repo.AuthRepo,
	jwtGenerator *utils.JwtTokenGenerator,
	redisClient *redis.Client,
	natsConn *nats.Conn,
) *AuthService {
	return &AuthService{
		log:          log,
		repo:         repo,
		jwtGenerator: jwtGenerator,
		redisClient:  redisClient,
		natsConn:     natsConn,
	}
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest, deviceId string) (*model.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.log.WithError(err).Error("failed to get user by email")
		return nil, ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, ErrEmailNotVerified
	}

	if !s.VerifyPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	refreshToken := s.jwtGenerator.GenerateJwtToken("refresh_token", user.ID.String(), deviceId)
	accessToken := s.jwtGenerator.GenerateJwtToken("access_token", user.ID.String(), deviceId)

	if refreshToken == "" || accessToken == "" {
		return nil, ErrTokenGenerationFailed
	}

	err = s.storeRefreshToken(ctx, user.ID.String(), deviceId, refreshToken)
	if err != nil {
		s.log.WithError(err).Error("failed to store refresh token")
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(time.Hour.Seconds() * 72),
		UserId:       user.ID,
	}, nil
}

func (s *AuthService) storeRefreshToken(ctx context.Context, userId, deviceId, token string) error {
	key := fmt.Sprintf("%s:%s:%s", refreshTokenPrefix, userId, deviceId)
	return s.redisClient.Set(ctx, key, token, tokenExpiration).Err()
}

func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*string, error) {
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return nil, errors.New("missing required fields")
	}

	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	user := &model.User{
		FullName:   req.FullName,
		Email:      req.Email,
		Password:   s.hashPassword(req.Password),
		IsVerified: false,
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil || !created {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	token := s.GenerateToken(6)
	key := fmt.Sprintf("%s:%s", emailVerificationPrefix, user.ID.String())
	err = s.redisClient.Set(ctx, key, token, emailTokenExpiration).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store verification token: %v", err)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		js, err := s.natsConn.JetStream()
		if err != nil {
			s.log.WithError(err).Error("failed to connect to nats")
			return
		}
		emailData := s.emailInfoTemplate(user.Email, token)
		_, err = js.Publish(userEmailVerificationTopic, emailData)
		if err != nil {
			s.log.WithError(err).Error("failed to publish email verification event")
		}
	}()

	return &token, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	token, err := s.jwtGenerator.ValidateToken(refreshToken)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	deviceId, ok := claims["device_id"].(string)
	if !ok {
		return nil, errors.New("invalid device ID in token")
	}

	if s.isTokenBlacklisted(ctx, refreshToken) {
		return nil, ErrTokenBlacklisted
	}

	newAccessToken := s.jwtGenerator.GenerateJwtToken("access_token", userId, deviceId)
	newRefreshToken := s.jwtGenerator.GenerateJwtToken("refresh_token", userId, deviceId)

	if newAccessToken == "" || newRefreshToken == "" {
		return nil, ErrTokenGenerationFailed
	}

	err = s.storeRefreshToken(ctx, userId, deviceId, newRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %v", err)
	}

	s.blacklistToken(ctx, refreshToken)

	return &model.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(time.Hour.Seconds() * 72),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	s.blacklistToken(ctx, accessToken)
	s.blacklistToken(ctx, refreshToken)
	return nil
}

func (s *AuthService) blacklistToken(ctx context.Context, token string) {
	key := fmt.Sprintf("%s%s", blacklistTokenPrefix, token)
	s.redisClient.Set(ctx, key, "1", tokenExpiration)
}

func (s *AuthService) isTokenBlacklisted(ctx context.Context, token string) bool {
	key := fmt.Sprintf("%s%s", blacklistTokenPrefix, token)
	val, err := s.redisClient.Get(ctx, key).Result()
	return err == nil && val == "1"
}

func (s *AuthService) VerifyEmail(ctx context.Context, email, token string) error {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	key := fmt.Sprintf("%s:%s", emailVerificationPrefix, user.ID.String())
	storedToken, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("verification token not found or expired: %v", err)
	}

	if storedToken != token {
		return errors.New("invalid verification token")
	}

	success, err := s.repo.PartialUpdateAfterEmailCode(ctx, email)
	if err != nil || !success {
		return fmt.Errorf("failed to verify email: %v", err)
	}

	s.redisClient.Del(ctx, key)

	return nil
}

func (s *AuthService) GenerateToken(n int) string {
	const charSet = "0123456789"
	rand.Seed(time.Now().UnixNano())
	token := make([]byte, n)
	for i := range token {
		token[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(token)
}

func (s *AuthService) hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func (s *AuthService) VerifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *AuthService) emailInfoTemplate(to, code string) []byte {
	md := &model.EmailVerification{
		To:      to,
		Code:    code,
		Message: "This email verification code will be expired in 10 min",
	}
	jsonBytes, err := json.Marshal(md)
	if err != nil {
		s.log.WithError(err).Panicf("failed to marshal email verification")
		return nil
	}
	return jsonBytes
}
func (s *AuthService) ResetPassword(ctx context.Context, email string) error {
	// 1. Check if user exists
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	// 2. Generate reset token
	resetToken := s.GenerateToken(32)
	key := fmt.Sprintf("password_reset:%s", user.ID.String())

	// 3. Store in Redis with expiration
	err = s.redisClient.Set(ctx, key, resetToken, emailTokenExpiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store reset token: %v", err)
	}

	// 4. Publish reset email event
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		emailData := s.passwordResetEmailTemplate(user.Email, resetToken)
		if err := s.natsConn.Publish("user.password_reset", emailData); err != nil {
			s.log.WithError(err).Error("failed to publish password reset email")
		}
	}()

	return nil
}

func (s *AuthService) VerifyResetPasswordToken(ctx context.Context, token, email string) error {
	// 1. Get user by email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	// 2. Get stored token from Redis
	key := fmt.Sprintf("password_reset:%s", user.ID.String())
	storedToken, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("reset token not found or expired: %v", err)
	}

	// 3. Compare tokens
	if storedToken != token {
		return errors.New("invalid reset token")
	}

	// Token is valid - deletion happens when actually resetting password
	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	// Similar to ResetPassword but with different messaging
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	resetToken := s.GenerateToken(32)
	key := fmt.Sprintf("forgot_password:%s", user.ID.String())

	err = s.redisClient.Set(ctx, key, resetToken, emailTokenExpiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store forgot password token: %v", err)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		emailData := s.forgotPasswordEmailTemplate(user.Email, resetToken)
		if err := s.natsConn.Publish("user.forgot_password", emailData); err != nil {
			s.log.WithError(err).Error("failed to publish forgot password email")
		}
	}()

	return nil
}

func (s *AuthService) VerifyForgotPasswordToken(ctx context.Context, token, email, newPassword string) error {
	// 1. Get user and validate token
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	key := fmt.Sprintf("forgot_password:%s", user.ID.String())
	storedToken, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("forgot password token not found or expired: %v", err)
	}

	if storedToken != token {
		return errors.New("invalid forgot password token")
	}

	// 2. Update password
	hashedPassword := s.hashPassword(newPassword)
	success, err := s.repo.UpdatePassword(ctx, email, hashedPassword)
	if err != nil || !success {
		return fmt.Errorf("failed to update password: %v", err)
	}

	// 3. Delete the used token
	s.redisClient.Del(ctx, key)

	return nil
}

// Helper methods for email templates
func (s *AuthService) passwordResetEmailTemplate(to, code string) []byte {
	md := &model.EmailVerification{
		To:      to,
		Code:    code,
		Message: "Your password reset code (expires in 10 minutes)",
	}
	jsonBytes, _ := json.Marshal(md)
	return jsonBytes
}

func (s *AuthService) forgotPasswordEmailTemplate(to, code string) []byte {
	md := &model.EmailVerification{
		To:      to,
		Code:    code,
		Message: "Your password recovery code (expires in 10 minutes)",
	}
	jsonBytes, _ := json.Marshal(md)
	return jsonBytes
}
