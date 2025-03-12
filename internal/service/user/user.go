package userService

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

	"github.com/lib/pq"
)

type UserService struct {
	log      *slog.Logger
	repo     repository.User
	repoRole repository.Role
	jwt      *jwtUtils.Jwt
	salt     string
}

func NewUserService(log *slog.Logger, repo repository.User, repoRole repository.Role, jwt *jwtUtils.Jwt, salt string) *UserService {

	return &UserService{
		log:      log,
		repo:     repo,
		repoRole: repoRole,
		jwt:      jwt,
		salt:     salt,
	}
}

func (s *UserService) Login(user *entity.UserLoginReq) (entity.TokensRes, error, int) {
	const op = "services.user.login"
	log := s.log.With(
		slog.String("op", op),
	)

	res, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		log.Error("Get user by email", "error", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return entity.TokensRes{}, errors.New("user with email not exists"), http.StatusBadRequest
		}
		return entity.TokensRes{}, errors.New("Login error"), http.StatusInternalServerError
	}

	if res.Password_hash != generatePasswordHash(user.Password, s.salt) {
		return entity.TokensRes{}, errors.New("not valid login or password"), http.StatusBadRequest
	}

	role, err := s.repoRole.GetUserRole(res.Id)
	if err != nil {
		log.Error("Get User role", "error", err.Error())
		return entity.TokensRes{}, errors.New("Login error"), http.StatusInternalServerError
	}
	return s.GenerateTokens(res.Id, role.Name)
}

func (s *UserService) Register(user *entity.UserRegisterReq) (entity.TokensRes, error, int) {
	const op = "services.user.register"
	log := s.log.With(
		slog.String("op", op),
	)
	userEntity := entity.User{
		Username:      user.Username,
		Email:         user.Email,
		Password_hash: generatePasswordHash(user.Password, s.salt),
	}

	userId, err := s.repo.Create(&userEntity)
	if err != nil {
		log.Error("Repo user create error", "error", err.Error())
		pqErr := err.(*pq.Error).Code
		if pqErr == "23505" {
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

func (s *UserService) GenerateTokens(userId int, role string) (entity.TokensRes, error, int) {
	const op = "services.user.GenerateTokens"
	log := s.log.With(
		slog.String("op", op),
	)
	access, err := s.jwt.GenerateAccessToken(userId, role)
	if err != nil {
		log.Error("Generate access token error", "error", err.Error())
		return entity.TokensRes{}, errors.New("create tokens error"), http.StatusInternalServerError
	}
	refresh, err := s.jwt.GenerateRefreshToken(userId, role)
	if err != nil {
		log.Error("Generate refresh token error", "error", err.Error())
		return entity.TokensRes{}, errors.New("create tokens error"), http.StatusInternalServerError
	}
	return entity.TokensRes{
		Access:  access,
		Refresh: refresh,
	}, nil, 0
}

func generatePasswordHash(password string, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
