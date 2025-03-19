package userService

import (
	"errors"
	"log/slog"
	"net/http"
	"todo/internal/entity"
	"todo/internal/repository"
)

type UserService struct {
	log      *slog.Logger
	repoUser repository.User
}

func NewUserService(log *slog.Logger, repoUser repository.User) *UserService {

	return &UserService{
		log:      log,
		repoUser: repoUser,
	}
}

func (s *UserService) GetAllUsers() ([]entity.User, error, int) {
	const op = "services.user.login"
	log := s.log.With(
		slog.String("op", op),
	)

	users, err := s.repoUser.GetUsers()
	if err != nil {
		log.Error("Repo user create error", "error", err.Error())
		return nil, errors.New("user register error"), http.StatusInternalServerError
	}

	return users, nil, http.StatusOK

}
