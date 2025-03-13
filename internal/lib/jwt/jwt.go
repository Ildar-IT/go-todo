package jwtUtils

import (
	"errors"
	"time"
	"todo/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	UserId int       `json:"user_id"`
	Role   string    `json:"role"`
	Exp    time.Time `json:"exp"`
}
type RefreshClaims struct {
	UserId int       `json:"user_id"`
	Role   string    `json:"role"`
	Exp    time.Time `json:"exp"`
	jwt.RegisteredClaims
}

type Jwt struct {
	cfg *config.JwtConfig
}

func New(cfg *config.JwtConfig) *Jwt {
	return &Jwt{cfg: cfg}
}

func (j *Jwt) GenerateAccessToken(userId int, role string) (string, error) {
	// claims := jwt.MapClaims{
	// 	"user_id": userId,
	// 	"role":    role,
	// 	"exp":     time.Now().Add(time.Minute * time.Duration(j.cfg.AccessTTL)).Unix(),
	// }
	claims := AccessClaims{
		UserId: userId,
		Role:   role,
		Exp:    time.Now().Add(time.Minute * time.Duration(j.cfg.AccessTTL)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	return token.SignedString([]byte(j.cfg.AccessSecret))
}

func (j *Jwt) GenerateRefreshToken(userId int, role string) (string, error) {
	claims := RefreshClaims{
		UserId: userId,
		Role:   role,
		Exp:    time.Now().Add(time.Minute * time.Duration(j.cfg.RefreshTTL)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.cfg.RefreshSecret))
}

func (j *Jwt) ValidateAccessToken(tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.cfg.AccessSecret), nil
	})
	if err != nil {
		return &AccessClaims{}, err
	}
	if !token.Valid {
		return &AccessClaims{}, errors.New("token not valid")
	}

	claims, ok := token.Claims.(*AccessClaims)

	if !ok {
		return claims, errors.New("not get access token claims")
	}
	return claims, nil
}

func (j *Jwt) ValidateRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.cfg.RefreshSecret), nil
	})
	if err != nil {
		return &RefreshClaims{}, err
	}
	if !token.Valid {
		return &RefreshClaims{}, errors.New("token not valid")
	}

	claims, ok := token.Claims.(*RefreshClaims)

	if !ok {
		return claims, errors.New("not get access token claims")
	}
	return claims, nil

}
