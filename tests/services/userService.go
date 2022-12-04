package services

import (
	"time"

	"github.com/LIOU2021/go-eloquent-mongodb/tests/models"
	"github.com/LIOU2021/go-eloquent-mongodb/tests/repositories"
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
	userAll, ok = service.repo.Orm.All()
	return
}

func (service *UserService) Find(id string) (user *models.User, ok bool) {
	user, ok = service.repo.Orm.Find(id)
	return
}

func (service *UserService) Insert(user *models.User) (insertId string, ok bool) {
	user.CreatedAt = uint64(time.Now().Unix())
	user.UpdatedAt = uint64(time.Now().Unix())

	insertId, ok = service.repo.Orm.Insert(user)
	return
}

func (service *UserService) Delete(id string) (deleteCount int, ok bool) {
	deleteCount, ok = service.repo.Orm.Delete(id)
	return
}

func (service *UserService) Update(id string, data *models.UserUpdateData) (updateCount int, ok bool) {
	data.UpdatedAt = uint64(time.Now().Unix())
	updateCount, ok = service.repo.Orm.Update(id, data)
	return
}
