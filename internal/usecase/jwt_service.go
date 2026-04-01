package usecase

import (
	"errors"
	"strconv"
	"time"

	"github.com/SKIND0A/online-shop/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTService(secret string, ttl time.Duration) *JWTService {
	return &JWTService{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (s *JWTService) GenerateAccessToken(userID int64, role domain.UserRole) (string, int64, error) {
	now := time.Now()
	exp := now.Add(s.ttl)

	claims := jwt.MapClaims{
		"sub":  strconv.FormatInt(userID, 10),
		"role": string(role),
		"iat":  now.Unix(),
		"exp":  exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", 0, err
	}

	return signed, int64(s.ttl.Seconds()), nil
}

func (s *JWTService) Parse(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
