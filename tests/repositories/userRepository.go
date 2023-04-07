package repositories

import (
	"context"

	"github.com/LIOU2021/go-eloquent-mongodb/orm"
	"github.com/LIOU2021/go-eloquent-mongodb/tests/models"
	"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
	orm.IEloquent[models.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		IEloquent: orm.NewEloquent[models.User]("users"),
	}
}

func (repo *UserRepository) GetUnderage(age int) (users []*models.User, err error) {
	filter := bson.M{
		"age": bson.M{"$lt": age},
	}
	users, err = repo.FindMultiple(context.Background(), filter)
	return
}

func (repo *UserRepository) GetOverage(age int) (users []*models.User, err error) {
	coll := repo.GetCollection()

	filter := bson.M{"age": bson.M{"$gt": age}}
	cursor, errF := coll.Find(context.TODO(), filter)

	if errF != nil {
		err = errF
		return
	}

	if errA := cursor.All(context.TODO(), &users); errA != nil {
		err = errA
		return
	}

	return
}
