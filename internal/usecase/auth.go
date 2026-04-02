package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/SKIND0A/online-shop/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("user is inactive")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid token")
)

type UserRepo interface {
	Create(ctx context.Context, email, passwordHash, role, displayName string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
}

type TokenService interface {
	GenerateAccessToken(userID int64, role domain.UserRole) (string, int64, error)
	ParseAccessToken(token string) (userID int64, role domain.UserRole, err error)
}

type AuthUsecase struct {
	users  UserRepo
	tokens TokenService
}

func NewAuthUsecase(users UserRepo, tokens TokenService) *AuthUsecase {
	return &AuthUsecase{users: users, tokens: tokens}
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"` // пока не используем, но оставляем под API контракт
}

type RegisterResult struct {
	UserID int64           `json:"user_id"`
	Email  string          `json:"email"`
	Name   string          `json:"name"`
	Role   domain.UserRole `json:"role"`
}

func (u *AuthUsecase) Register(ctx context.Context, in RegisterInput) (*RegisterResult, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if in.Email == "" || len(in.Password) < 8 {
		return nil, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	displayName := strings.TrimSpace(in.Name)

	user, err := u.users.Create(ctx, in.Email, string(hash), string(domain.RoleCustomer), displayName)
	if err != nil {
		// на уровне handler можно маппить repo-ошибку duplicate -> 409
		return nil, err
	}

	return &RegisterResult{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.DisplayName,
		Role:   user.Role,
	}, nil
}

type MeResult struct {
	Email string          `json:"email"`
	Name  string          `json:"name"`
	Role  domain.UserRole `json:"role"`
}

func (u *AuthUsecase) Me(ctx context.Context, accessToken string) (*MeResult, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, ErrInvalidToken
	}

	userID, _, err := u.tokens.ParseAccessToken(accessToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	return &MeResult{
		Email: user.Email,
		Name:  user.DisplayName,
		Role:  user.Role,
	}, nil
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (u *AuthUsecase) Login(ctx context.Context, in LoginInput) (*LoginResult, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if in.Email == "" || in.Password == "" {
		return nil, ErrInvalidInput
	}

	user, err := u.users.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrInactiveUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, expiresIn, err := u.tokens.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}, nil
}
