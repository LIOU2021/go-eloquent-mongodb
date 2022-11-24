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
	All() (results []*bson.M)
	Find(id string) *Model
	Insert(data interface{}) (ok bool)
	Delete(filter interface{}) (ok bool)
	Update(filter interface{}, data interface{}) (ok bool)
}

func NewEloquent(collection string) *Eloquent {
	return &Eloquent{
		db:         os.Getenv("mongodb_name"),
		Collection: collection,
		uri:        fmt.Sprintf("mongodb://%s:%s@%s:%s/?retryWrites=true&w=majority", os.Getenv("mongodb_user"), os.Getenv("mongodb_password"), os.Getenv("mongodb_host"), os.Getenv("mongodb_port")),
	}
}

func (e *Eloquent) All() (results []*bson.M) {
	uri := e.uri
	if uri == "" {
		log.Fatal("You must set your 'mongodb_host' and 'mongodb_port' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database(e.db).Collection(e.Collection)

	cursor, err := coll.Find(context.TODO(), bson.M{})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found for collection %s\n", e.Collection)
			return
		}
		panic(err)
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return
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
