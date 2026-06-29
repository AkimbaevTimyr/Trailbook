package usecase

import (
	"context"
	"fmt"

	"tracking-backend/internal/domain"
	"tracking-backend/internal/repository"
)

// UserUsecase defines the interface for user-related business logic.
type UserUsecase interface {
	GetUsers(ctx context.Context) ([]domain.User, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

// NewUserUsecase creates a new UserUsecase.
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
