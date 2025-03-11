package userService

import (
	"crypto/sha1"
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

func (s *UserService) Login(user *entity.UserLoginReq) (entity.UserLoginRes, error, int) {
	const op = "services.user.login"
	log := s.log.With(
		slog.String("op", op),
	)

	res, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		log.Error("Get user by email", "error", err.Error())
		return entity.UserLoginRes{}, errors.New("Login error"), http.StatusInternalServerError
	}
	if res.Id == 0 {
		return entity.UserLoginRes{}, errors.New("user with email not exists"), http.StatusBadRequest
	}

	if res.Password_hash != generatePasswordHash(user.Password, s.salt) {
		return entity.UserLoginRes{}, errors.New("not valid login or password"), http.StatusBadRequest
	}

	role, err := s.repoRole.GetUserRole(res.Id)
	log.Info("LOGIN USER ROLE", "roel", role.Name)
	if err != nil {
		log.Error("Get User role", "error", err.Error())
		return entity.UserLoginRes{}, errors.New("Login error"), http.StatusInternalServerError
	}

	access, err := s.jwt.GenerateAccessToken(res.Id, role.Name)
	if err != nil {
		log.Error("Generate access token error", "error", err.Error())
		return entity.UserLoginRes{}, errors.New("create tokens error"), http.StatusInternalServerError
	}
	refresh, err := s.jwt.GenerateRefreshToken(res.Id, role.Name)
	if err != nil {
		log.Error("Generate refresh token error", "error", err.Error())
		return entity.UserLoginRes{}, errors.New("create tokens error"), http.StatusInternalServerError
	}
	return entity.UserLoginRes{
		Access:  access,
		Refresh: refresh,
	}, nil, 0
}

func (s *UserService) Register(user *entity.UserRegisterReq) (entity.UserRegisterRes, error, int) {
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
			return entity.UserRegisterRes{}, errors.New("user already exists"), http.StatusBadRequest
		}
		return entity.UserRegisterRes{}, errors.New("user register error"), http.StatusInternalServerError
	}
	role, err := s.repoRole.SetUserRole(userId)
	log.Info("register USER ROLE", role.Name)
	if err != nil {
		log.Error("Repo user set roel error", "error", err.Error())
		return entity.UserRegisterRes{}, errors.New("user register error"), http.StatusInternalServerError
	}
	access, err := s.jwt.GenerateAccessToken(userId, role.Name)
	if err != nil {
		log.Error("Generate access token error", "error", err.Error())
		return entity.UserRegisterRes{}, errors.New("create tokens error"), http.StatusInternalServerError
	}
	refresh, err := s.jwt.GenerateRefreshToken(userId, role.Name)
	if err != nil {
		log.Error("Generate refresh token error", "error", err.Error())
		return entity.UserRegisterRes{}, errors.New("create tokens error"), http.StatusInternalServerError
	}
	return entity.UserRegisterRes{
		Access:  access,
		Refresh: refresh,
	}, nil, 0
}

func generatePasswordHash(password string, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
