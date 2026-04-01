package postgres

import (
	"context"
	"errors"

	"github.com/SKIND0A/online-shop/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(
	ctx context.Context,
	email string,
	passwordHash string,
	role string,
) (*domain.User, error) {
	const q = `
		INSERT INTO users (email, password_hash, role, is_active)
		VALUES ($1, $2, $3, true)
		RETURNING id, email, password_hash, role, is_active, created_at, updated_at
	`

	var u domain.User
	err := r.pool.QueryRow(ctx, q, email, passwordHash, role).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var u domain.User
	err := r.pool.QueryRow(ctx, q, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}
