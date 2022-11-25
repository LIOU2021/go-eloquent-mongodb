package services

import (
	"fmt"
	"github/LIOU2021/go-eloquent-mongodb/test/models"
	"github/LIOU2021/go-eloquent-mongodb/test/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repositories.NewUserRepository(),
	}
}

func (service *UserService) All() (userAll []*models.User, ok bool) {
	userAll = []*models.User{}
	ok = service.repo.Orm.All(&userAll)

	if !ok {
		fmt.Println("user all query fail !")
	}

	return
}

func (service *UserService) Find(id string) (user *models.User, ok bool) {
	user = &models.User{}
	ok = service.repo.Orm.Find(id, user)
	return
}
