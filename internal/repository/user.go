package repository

import (
	"context"
	"database/sql"
	"fmt"

	"tracking-backend/internal/domain"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func scanUser(scanner interface {
	Scan(dest ...interface{}) error
}) (domain.User, error) {
	var user domain.User
	var avatarURL sql.NullString
	var updatedAt sql.NullTime

	err := scanner.Scan(
		&user.ID,
		&user.Name,
		&user.RegionID,
		&user.Email,
		&user.PasswordHash,
		&avatarURL,
		&user.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return user, err
	}

	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}

	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, region_id, email, password_hash, avatar_url, created_at, updated_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := scanUser(r.db.QueryRowContext(
		ctx,
		"SELECT id, name, region_id, email, password_hash, avatar_url, created_at, updated_at FROM users WHERE id = $1",
		id,
	))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := scanUser(r.db.QueryRowContext(
		ctx,
		"SELECT id, name, region_id, email, password_hash, avatar_url, created_at, updated_at FROM users WHERE email = $1",
		email,
	))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (name, region_id, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, region_id, email, password_hash, avatar_url, created_at, updated_at
	`

	created, err := scanUser(r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.RegionID,
		user.Email,
		user.PasswordHash,
	))
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &created, nil
}
