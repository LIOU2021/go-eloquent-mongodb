package repositories

import (
	"github.com/LIOU2021/go-eloquent-mongodb/orm"
	"github.com/LIOU2021/go-eloquent-mongodb/tests/models"
)

type UserRepository struct {
	Orm orm.IEloquent[models.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Orm: orm.NewEloquent[models.User]("users"),
	}
}
