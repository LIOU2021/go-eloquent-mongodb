package orm

import (
	"context"
	"os"

	"github.com/LIOU2021/go-eloquent-mongodb/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Eloquent[t interface{}] struct {
	db         string //db name
	Collection string
	uri        string
	logTitle   string
}

type IEloquent[T interface{}] interface {
	All() (models []*T, ok bool)
	Find(id string) (model *T, ok bool)
	Insert(data interface{}) (insertedID string, ok bool)
	Delete(id string) (deleteCount int, ok bool)
	Update(id string, data interface{}) (modifiedCount int, ok bool)
}

func NewEloquent[T interface{}](collection string) *Eloquent[T] {
	return &Eloquent[T]{
		db:         os.Getenv("mongodb_name"),
		Collection: collection,
		uri:        getUri(),
		logTitle:   getLogTitle(collection),
	}
}

/**
 * @title creates a new Client connect
 */
func (e *Eloquent[T]) Connect() (client *mongo.Client) {
	uri := e.uri
	if uri == "" {
		logger.LogDebug.Error(e.logTitle, "You must set your 'mongodb_host' and 'mongodb_port' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable", getCurrentFuncInfo())
		os.Exit(0)
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo())
		os.Exit(0)
	}

	return client
}

/**
 * @title Close connect
 */
func (e *Eloquent[T]) Close(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo())
	}
}

/**
 * @title get collection instance
 */
func (e *Eloquent[T]) GetCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(e.db).Collection(e.Collection)
}

/**
 * @title get all collection from document
 * @param models interface{}
 * @return bool query success or fail
 */
func (e *Eloquent[T]) All() (models []*T, ok bool) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	cursor, err := coll.Find(context.TODO(), bson.M{})

	if err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo())
		ok = false
		return
	}

	models = []*T{}

	if err = cursor.All(context.TODO(), &models); err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo())
		ok = false
		return

	}
	ok = true
	return
}

/**
 * @title find a document by _id
 * @param id string _id of document
 * @return model your model struct
 * @return ok bool query success or fail
 */
func (e *Eloquent[T]) Find(id string) (model *T, ok bool) {
	idH, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo())
		ok = false
		return
	}

	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)
	model = new(T)
	err = coll.FindOne(context.TODO(), bson.M{"_id": idH}).Decode(model)

	if err == mongo.ErrNoDocuments {
		ok = true
		return
	} else if err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo())
		ok = false
		return
	}
	ok = true
	return
}

/**
 * @title insert a document
 * @param data interface{} your model struct
 * @return insertedID *primitive.ObjectID ObjectId of mongodb
 */
func (e *Eloquent[T]) Insert(data interface{}) (insertedID string, ok bool) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		ok = false
		return
	}

	ok = true
	insertedID = result.InsertedID.(primitive.ObjectID).Hex()
	return
}

/**
 * @title delete a document
 * @param id string _id of document
 * @return deleteCount int delete document count
 * @return ok bool query success or fail
 */
func (e *Eloquent[T]) Delete(id string) (deleteCount int, ok bool) {
	idH, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo())
		ok = false
		return
	}

	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	filter := bson.M{"_id": idH}

	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		ok = false
		return
	}

	ok = true
	deleteCount = int(result.DeletedCount)
	return
}

func (e *Eloquent[T]) Update(id string, data interface{}) (modifiedCount int, ok bool) {
	idH, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo())
		ok = false
		return
	}

	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	filter := bson.M{"_id": idH}
	update := bson.M{"$set": data}

	result, err := coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo())
		ok = false
		return
	}

	ok = true
	modifiedCount = int(result.ModifiedCount)
	return
}
