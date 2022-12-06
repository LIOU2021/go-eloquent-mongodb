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

func (service *UserService) All() (userAll []*models.User, err error) {
	userAll, err = service.repo.Orm.All()
	return
}

func (service *UserService) Find(id string) (user *models.User, err error) {
	user, err = service.repo.Orm.Find(id)
	return
}

func (service *UserService) FindMultiple(filter interface{}) (users []*models.User, err error) {
	users, err = service.repo.Orm.FindMultiple(filter)
	return
}

func (service *UserService) Insert(user *models.User) (insertId string, err error) {
	currentTime := uint64(time.Now().Unix())
	user.CreatedAt = &currentTime
	user.UpdatedAt = &currentTime

	insertId, err = service.repo.Orm.Insert(user)
	return
}

func (service *UserService) InsertMultiple(user []*models.User) (InsertedIDs []string, err error) {
	InsertedIDs, err = service.repo.Orm.InsertMultiple(user)
	return
}

func (service *UserService) Delete(id string) (deleteCount int, err error) {
	deleteCount, err = service.repo.Orm.Delete(id)
	return
}

func (service *UserService) DeleteMultiple(filter interface{}) (deleteCount int, err error) {
	deleteCount, err = service.repo.Orm.DeleteMultiple(filter)
	return
}

func (service *UserService) Update(id string, data *models.User) (updateCount int, err error) {
	currentTime := uint64(time.Now().Unix())
	data.UpdatedAt = &currentTime
	updateCount, err = service.repo.Orm.Update(id, data)
	return
}

func (service *UserService) UpdateMultiple(filter interface{}, data *models.User) (updateCount int, err error) {
	currentTime := uint64(time.Now().Unix())
	data.UpdatedAt = &currentTime
	updateCount, err = service.repo.Orm.UpdateMultiple(filter, data)
	return
}

func (e *UserService) Count(filter interface{}) (count int, ok bool) {
	count, ok = e.repo.Orm.Count(filter)
	return
}
