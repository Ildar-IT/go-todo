package jwtUtils_test

import (
	"testing"
	"todo/internal/config"
	jwtUtils "todo/internal/lib/jwt"
)

func TestGenerateAccessToken(t *testing.T) {
	cfg := config.JwtConfig{
		AccessSecret: "access_secret",
		AccessTTL:    15,
	}

	jwt := jwtUtils.New(&cfg)
	token, err := jwt.GenerateAccessToken(1, "admin")
	if err != nil {
		t.Fatalf("Generate token err: %s", err.Error())
	}

	if token == "" {
		t.Fatal("Token is empty")
	}
}
func TestGenerateRefreshToken(t *testing.T) {
	cfg := config.JwtConfig{
		AccessSecret: "refresh_secret",
		AccessTTL:    1000,
	}
	// mok := jwtUtils.AccessClaims{
	// 	UserId: 1,
	// 	Role:   "user",
	// }
	jwt := jwtUtils.New(&cfg)
	token, err := jwt.GenerateRefreshToken(1, "user")
	if err != nil {
		t.Fatalf("Generate token err: %s", err.Error())
	}

	if token == "" {
		t.Fatal("Token is empty")
	}
}

func TestValidateAccessToken(t *testing.T) {
	cfg := config.JwtConfig{
		AccessSecret: "access_secret",
		AccessTTL:    15,
	}
	mok := jwtUtils.AccessClaims{
		UserId: 1,
		Role:   "user",
	}
	jwt := jwtUtils.New(&cfg)
	token, _ := jwt.GenerateAccessToken(mok.UserId, mok.Role)

	claims, err := jwt.ValidateAccessToken(token)

	if err != nil {
		t.Fatalf("Validate token err: %s", err.Error())
	}
	if claims.UserId != mok.UserId || claims.Role != mok.Role {
		t.Fatalf("Validate token return not valid claims")
	}
}

func TestValidateRefreshToken(t *testing.T) {
	cfg := config.JwtConfig{
		RefreshSecret: "refresh_secret",
		RefreshTTL:    1000,
	}
	mok := jwtUtils.RefreshClaims{
		UserId: 1,
		Role:   "user",
	}
	jwt := jwtUtils.New(&cfg)
	token, _ := jwt.GenerateRefreshToken(mok.UserId, mok.Role)

	claims, err := jwt.ValidateRefreshToken(token)

	if err != nil {
		t.Fatalf("Validate token err: %s", err.Error())
	}
	if claims.UserId != mok.UserId || claims.Role != mok.Role {
		t.Fatalf("Validate token return not valid claims")
	}
}
