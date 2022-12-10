package repositories

import (
	"github.com/LIOU2021/go-eloquent-mongodb/orm"
	"github.com/LIOU2021/go-eloquent-mongodb/tests/models"
	"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
	Orm orm.IEloquent[models.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Orm: orm.NewEloquent[models.User]("users"),
	}
}

func (repo *UserRepository) GetUnderage(age int) (users []*models.User, err error) {
	filter := bson.M{
		"age": bson.M{"$lt": age},
	}
	users, err = repo.Orm.FindMultiple(filter)
	return
}
