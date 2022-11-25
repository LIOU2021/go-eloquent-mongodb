package repositories

import "github/LIOU2021/go-eloquent-mongodb/orm"

type UserRepository struct {
	Orm orm.IEloquent
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Orm: orm.NewEloquent("users"),
	}
}