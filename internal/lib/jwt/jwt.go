package jwtUtils

import (
	"errors"
	"time"
	"todo/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

const (
	RefreshTokenType = "qazwsxedc"
	AccessTokenType  = "rtyhgfvnb"
)

type tokenClaims struct {
	jwt.Claims
}

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

func (j *Jwt) ValidateAccessToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.cfg.AccessSecret), nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}
	if !token.Valid {
		return jwt.MapClaims{}, errors.New("token not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return claims, errors.New("not get access token claims")
	}
	return claims, nil
}

func (j *Jwt) ValidateRefreshToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.cfg.RefreshSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("token not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("not get refresh token claims")
	}

	return claims["user_id"].(int), nil

}
