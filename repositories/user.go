package repositories

import "github/LIOU2021/go-eloquent-mongodb/orm"

type UserRepository struct {
	orm *orm.Eloquent
}

func NewUserModel() *UserRepository {
	return &UserRepository{
		orm: orm.NewEloquent("users"),
	}
}
