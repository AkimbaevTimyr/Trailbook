package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"tracking-backend/internal/delivery/http/requests"

	"golang.org/x/crypto/bcrypt"

	"tracking-backend/internal/domain"
	"tracking-backend/internal/repository"
)

// UserUsecase defines the interface for user-related business logic.
type UserUsecase interface {
	GetUsers(ctx context.Context) ([]domain.User, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	Create(ctx context.Context, req requests.CreateUserRequest) (*domain.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) GetUsers(ctx context.Context) ([]domain.User, error) {
	users, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	return users, nil
}

func (uc *userUsecase) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (uc *userUsecase) Create(ctx context.Context, req requests.CreateUserRequest) (*domain.User, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		RegionID:     req.RegionID,
		PasswordHash: string(hash),
	}

	created, err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	created.PasswordHash = ""
	return created, nil
}

func validateCreateUserRequest(req requests.CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}
