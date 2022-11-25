package services

import (
	"fmt"
	"github/LIOU2021/go-eloquent-mongodb/tests/models"
	"github/LIOU2021/go-eloquent-mongodb/tests/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (service *UserService) Insert(user *models.UserCreateData) (_id string, ok bool) {
	user.CreatedAt = uint64(time.Now().Unix())
	user.UpdatedAt = uint64(time.Now().Unix())

	insertId, ok := service.repo.Orm.Insert(user)
	_id = insertId.(primitive.ObjectID).Hex()
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
