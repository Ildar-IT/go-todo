package authService

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"todo/internal/entity"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/repository"
)

type AuthService struct {
	log      *slog.Logger
	repoUser repository.User
	repoRole repository.Role
	jwt      *jwtUtils.Jwt
	salt     string
}

func NewAuthService(log *slog.Logger, repoUser repository.User, repoRole repository.Role, jwt *jwtUtils.Jwt, salt string) *AuthService {

	return &AuthService{
		log:      log,
		repoUser: repoUser,
		repoRole: repoRole,
		jwt:      jwt,
		salt:     salt,
	}
}

func (s *AuthService) Login(user *entity.UserLoginReq) (entity.TokensRes, error, int) {
	const op = "services.auth.login"
	log := s.log.With(
		slog.String("op", op),
	)

	res, err := s.repoUser.GetUserByEmail(user.Email)
	if err != nil {
		log.Error("Get user by email", "error", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return entity.TokensRes{}, errors.New("user with email not exists"), http.StatusBadRequest
		}
		return entity.TokensRes{}, errors.New("login error"), http.StatusInternalServerError
	}

	if res.Password_hash != generatePasswordHash(user.Password, s.salt) {
		return entity.TokensRes{}, errors.New("not valid login or password"), http.StatusBadRequest
	}

	role, err := s.repoRole.GetUserRole(res.Id)
	if err != nil {
		log.Error("Get User role", "error", err.Error())
		return entity.TokensRes{}, errors.New("login error"), http.StatusInternalServerError
	}
	return s.GenerateTokens(res.Id, role.Name)
}

func (s *AuthService) Register(user *entity.UserRegisterReq) (entity.TokensRes, error, int) {
	const op = "services.auth.register"
	log := s.log.With(
		slog.String("op", op),
	)
	userEntity := entity.User{
		Username:      user.Username,
		Email:         user.Email,
		Password_hash: generatePasswordHash(user.Password, s.salt),
	}

	userId, err := s.repoUser.Create(&userEntity)
	if err != nil {
		log.Error("Repo user create error", "error", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return entity.TokensRes{}, errors.New("user already exists"), http.StatusBadRequest
		}
		return entity.TokensRes{}, errors.New("user register error"), http.StatusInternalServerError
	}
	role, err := s.repoRole.SetUserRole(userId)
	if err != nil {
		log.Error("Repo user set roel error", "error", err.Error())
		return entity.TokensRes{}, errors.New("user register error"), http.StatusInternalServerError
	}

	return s.GenerateTokens(userId, role.Name)
}

func (s *AuthService) GenerateTokens(userId int, role string) (entity.TokensRes, error, int) {
	access, err, status := s.GenerateAccessToken(userId, role)
	if err != nil {
		return entity.TokensRes{}, err, status
	}
	refresh, err, status := s.GenerateRefreshToken(userId, role)
	if err != nil {
		return entity.TokensRes{}, err, status
	}
	return entity.TokensRes{
		Access:  access,
		Refresh: refresh,
	}, nil, http.StatusOK
}
func (s *AuthService) GenerateAccessToken(userId int, role string) (string, error, int) {
	const op = "services.auth.GenerateAccessToken"
	log := s.log.With(
		slog.String("op", op),
	)
	access, err := s.jwt.GenerateAccessToken(userId, role)
	if err != nil {
		log.Error("Generate access token error", "error", err.Error())
		return "", errors.New("create tokens error"), http.StatusInternalServerError
	}
	return access, nil, http.StatusOK
}
func (s *AuthService) GenerateRefreshToken(userId int, role string) (string, error, int) {
	const op = "services.auth.GenerateRefreshToken"
	log := s.log.With(
		slog.String("op", op),
	)
	refresh, err := s.jwt.GenerateRefreshToken(userId, role)
	if err != nil {
		log.Error("Generate refresh token error", "error", err.Error())
		return "", errors.New("create token error"), http.StatusInternalServerError
	}
	return refresh, nil, http.StatusOK
}

func generatePasswordHash(password string, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
