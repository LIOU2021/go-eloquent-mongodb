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

func (service *UserService) Find(id string) (user *models.User, err error) {
	user, err = service.repo.Orm.Find(id)
	return
}

func (service *UserService) FindMultiple(filter interface{}) (users []*models.User, ok bool) {
	users, ok = service.repo.Orm.FindMultiple(filter)
	return
}

func (service *UserService) Insert(user *models.User) (insertId string, ok bool) {
	currentTime := uint64(time.Now().Unix())
	user.CreatedAt = &currentTime
	user.UpdatedAt = &currentTime

	insertId, ok = service.repo.Orm.Insert(user)
	return
}

func (service *UserService) InsertMultiple(user []*models.User) (InsertedIDs []string, ok bool) {
	InsertedIDs, ok = service.repo.Orm.InsertMultiple(user)
	return
}

func (service *UserService) Delete(id string) (deleteCount int, ok bool) {
	deleteCount, ok = service.repo.Orm.Delete(id)
	return
}

func (service *UserService) DeleteMultiple(filter interface{}) (deleteCount int, ok bool) {
	deleteCount, ok = service.repo.Orm.DeleteMultiple(filter)
	return
}

func (service *UserService) Update(id string, data *models.User) (updateCount int, ok bool) {
	currentTime := uint64(time.Now().Unix())
	data.UpdatedAt = &currentTime
	updateCount, ok = service.repo.Orm.Update(id, data)
	return
}

func (service *UserService) UpdateMultiple(filter interface{}, data *models.User) (updateCount int, ok bool) {
	currentTime := uint64(time.Now().Unix())
	data.UpdatedAt = &currentTime
	updateCount, ok = service.repo.Orm.UpdateMultiple(filter, data)
	return
}

func (e *UserService) Count(filter interface{}) (count int, ok bool) {
	count, ok = e.repo.Orm.Count(filter)
	return
}
