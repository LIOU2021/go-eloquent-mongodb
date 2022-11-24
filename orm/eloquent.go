package orm

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Eloquent struct {
	db         string //db name
	Collection string
	uri        string
}

type IEloquent interface {
	All(models interface{}) bool
	Find(id string) *Model
	Insert(data interface{}) (ok bool)
	Delete(filter interface{}) (ok bool)
	Update(filter interface{}, data interface{}) (ok bool)
}

func NewEloquent(collection string) *Eloquent {
	return &Eloquent{
		db:         os.Getenv("mongodb_name"),
		Collection: collection,
		uri:        getUri(),
	}
}

func (e *Eloquent) All(models interface{}) bool {
	uri := e.uri
	if uri == "" {
		log.Fatal("You must set your 'mongodb_host' and 'mongodb_port' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			fmt.Println(err)
		}
	}()

	coll := client.Database(e.db).Collection(e.Collection)

	cursor, err := coll.Find(context.TODO(), bson.M{})

	if err != nil {
		fmt.Println(err)
		return false
	}

	if err = cursor.All(context.TODO(), models); err != nil {
		fmt.Println(err)
		return false

	}

	return true
}

func (e *Eloquent) Find(id string) *Model {
	return &Model{}
}

func (e *Eloquent) Insert(data interface{}) (ok bool) {
	return
}

func (e *Eloquent) Delete(filter interface{}) (ok bool) {
	return
}

func (e *Eloquent) Update(filter interface{}, data interface{}) (ok bool) {
	return
}
