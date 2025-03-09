package userService

import (
	"crypto/sha1"
	"fmt"
	"log/slog"
	"todo/internal/entity"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/repository"
)

const (
	secretRefresh = "refresh"
	secretAccess  = "access"
	salt          = "passhashsalt"
)

type UserService struct {
	log  *slog.Logger
	repo repository.User
	jwt  *jwtUtils.Jwt
}

func NewUserService(log *slog.Logger, repo repository.User, jwt *jwtUtils.Jwt, salt string) *UserService {

	return &UserService{
		log:  log,
		repo: repo,
		jwt:  jwt,
	}
}

func (s *UserService) Login(user *entity.UserLoginReq) (int, error) {

	return 0, nil
}

func (s *UserService) Register(user *entity.UserRegisterReq) (entity.UserRegisterRes, error) {
	const op = "service.user"
	log := s.log.With(
		slog.String("op", op),
	)
	userEntity := entity.User{
		Username:      user.Username,
		Email:         user.Email,
		Password_hash: generatePasswordHash(user.Password),
	}

	// res, err := s.repo.GetUserByEmail(userEntity.Email)
	// if err != nil {
	// 	return entity.UserRegisterRes{}, err
	// }
	// if res.Id == 0 {
	// 	return entity.UserRegisterRes{}, errors.New("User with email exists")
	// }
	userId, err := s.repo.Create(&userEntity)

	if err != nil {
		log.Error("Repo user create error", "error", err.Error())
		return entity.UserRegisterRes{}, err
	}
	access, err := s.jwt.GenerateAccessToken(userId, "user")
	if err != nil {
		s.log.Error("Create access token error", "error", err.Error())
		return entity.UserRegisterRes{}, err
	}
	refresh, err := s.jwt.GenerateRefreshToken(userId, "user")
	if err != nil {
		s.log.Error("Create refresh token error", "error", err.Error())
		return entity.UserRegisterRes{}, err
	}
	return entity.UserRegisterRes{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
