package userservice

import (
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/repository/userrepo"
)

type UserService interface {
	AddUser(user models.User) error
	GetUserByTelegramId(telegramId int64) (models.User, error)
	GetUserByToken(token string) (models.User, error)
}

type userService struct {
	Repo userrepo.Repo
}

func NewUserService(usrRepo userrepo.Repo) UserService {
	return &userService{Repo: usrRepo}
}

func (u userService) AddUser(user models.User) error {
	return u.Repo.AddUser(user)
}

func (u userService) GetUserByTelegramId(telegramId int64) (models.User, error) {
	return u.Repo.GetUserByTelegramId(telegramId)
}

func (u userService) GetUserByToken(token string) (models.User, error) {
	return u.Repo.GetUserByToken(token)
}
