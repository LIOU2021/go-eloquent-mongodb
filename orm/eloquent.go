package orm

import (
	"context"
	"github/LIOU2021/go-eloquent-mongodb/logger"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Find(id string, model interface{}) bool
	Insert(data interface{}) (insertedID interface{}, ok bool)
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

func (e *Eloquent) connect() (client *mongo.Client) {
	uri := e.uri
	if uri == "" {
		logger.LogDebug.Error("You must set your 'mongodb_host' and 'mongodb_port' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		os.Exit(0)
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logger.LogDebug.Error(err)
		os.Exit(0)
	}

	return client
}

func (e *Eloquent) close(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		logger.LogDebug.Error(err)
	}
}

func (e *Eloquent) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(e.db).Collection(e.Collection)
}

func (e *Eloquent) All(models interface{}) bool {
	client := e.connect()

	defer e.close(client)

	coll := e.getCollection(client)

	cursor, err := coll.Find(context.TODO(), bson.M{})

	if err != nil {
		logger.LogDebug.Error(err)
		return false
	}

	if err = cursor.All(context.TODO(), models); err != nil {
		logger.LogDebug.Error(err)
		return false

	}

	return true
}

func (e *Eloquent) Find(id string, model interface{}) bool {
	idH, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogDebug.Errorf("collection - %s : _id Hex fail", e.Collection)
	}

	client := e.connect()

	defer e.close(client)

	coll := e.getCollection(client)

	err = coll.FindOne(context.TODO(), bson.M{"_id": idH}).Decode(model)

	if err == mongo.ErrNoDocuments {
		return true
	} else if err != nil {
		logger.LogDebug.Errorf("collection %s - FindOne with error %v =>", e.Collection, err)
		return false
	}

	return true
}

/**
 * @title insert a document
 * @param data interface{} your model struct
 * @return insertedID *primitive.ObjectID ObjectId of mongodb
 */
func (e *Eloquent) Insert(data interface{}) (insertedID interface{}, ok bool) {
	client := e.connect()

	defer e.close(client)

	coll := e.getCollection(client)

	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		ok = false
		panic(err)
	}

	ok = true
	insertedID = result.InsertedID
	return
}

func (e *Eloquent) Delete(filter interface{}) (ok bool) {
	return
}

func (e *Eloquent) Update(filter interface{}, data interface{}) (ok bool) {
	return
}
