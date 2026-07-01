package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tracking-backend/internal/delivery/http/requests"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"tracking-backend/internal/domain"
	"tracking-backend/internal/repository"
)

const (
	accessTokenExpiry = 24 * time.Hour
)

// AuthUsecase defines the interface for authentication business logic.
type AuthUsecase interface {
	Login(ctx context.Context, req requests.LoginRequest) (*requests.TokenResponse, error)
	VerifyToken(tokenString string) (*requests.JWTClaims, error)
}

type authUsecase struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthUsecase creates a new AuthUsecase.
func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (uc *authUsecase) Login(ctx context.Context, req requests.LoginRequest) (*requests.TokenResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := uc.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &requests.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(accessTokenExpiry.Seconds()),
	}, nil
}

func (uc *authUsecase) VerifyToken(tokenString string) (*requests.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &requests.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*requests.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (uc *authUsecase) generateToken(user *domain.User) (string, error) {
	now := time.Now()
	claims := requests.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
