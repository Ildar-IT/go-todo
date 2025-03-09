package jwtUtils

import (
	"errors"
	"time"
	"todo/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	cfg *config.JwtConfig
}

func New(cfg *config.JwtConfig) *Jwt {

	return &Jwt{cfg: cfg}
}

func (j *Jwt) GenerateAccessToken(userId int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     time.Now().Add(time.Minute * time.Duration(j.cfg.AccessTTL)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.cfg.AccessSecret))
}

func (j *Jwt) GenerateRefreshToken(userId int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * time.Duration(j.cfg.RefreshTTL)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.cfg.RefreshSecret))
}

func ValidateAccessToken(secret string, tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("Token not valid")
	}
	return nil
}

func ValidateRefreshToken(secret string, tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("Token not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Not get refresh token claims")
	}

	return claims["user_id"].(int), nil

}
