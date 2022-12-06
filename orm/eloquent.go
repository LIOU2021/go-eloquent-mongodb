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
	All() (models []*T, err error)
	Find(id string) (model *T, err error)
	FindMultiple(filter interface{}) (models []*T, err error)
	Insert(data *T) (insertedID string, ok bool)
	InsertMultiple(data []*T) (InsertedIDs []string, ok bool)
	Delete(id string) (deleteCount int, ok bool)
	DeleteMultiple(filter interface{}) (deleteCount int, ok bool)
	Update(id string, data *T) (modifiedCount int, ok bool)
	UpdateMultiple(filter interface{}, data *T) (modifiedCount int, ok bool)
	Count(filter interface{}) (count int, ok bool)
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
		logger.LogDebug.Error(e.logTitle, "You must set your 'mongodb_host' and 'mongodb_port' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable", getCurrentFuncInfo(1))
		os.Exit(0)
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		os.Exit(0)
	}

	return client
}

/**
 * @title Close connect
 */
func (e *Eloquent[T]) Close(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
	}
}

/**
 * @title get collection instance
 */
func (e *Eloquent[T]) GetCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(e.db).Collection(e.Collection)
}

/**
 * @title get all document from collection
 * @return models []*T your model slice
 * @return err error fail message from query
 */
func (e *Eloquent[T]) All() (models []*T, err error) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	cursor, err := coll.Find(context.TODO(), bson.M{})

	if err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		err = e.errMsg(err)
		return
	}

	models = []*T{}

	if err = cursor.All(context.TODO(), &models); err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		err = e.errMsg(err)
		return

	}
	return
}

/**
 * @title find a document by _id
 * @param id string _id of document
 * @return model struct your model struct
 * @return err error fail message from query
 */
func (e *Eloquent[T]) Find(id string) (model *T, err error) {
	idH, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo(1))
		err = e.errMsg(err)
		return
	}

	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)
	model = new(T)
	err = coll.FindOne(context.TODO(), bson.M{"_id": idH}).Decode(model)

	if err == mongo.ErrNoDocuments {
		err = e.errMsg(err)
		return
	} else if err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		err = e.errMsg(err)
		return
	}

	return
}

/**
 * @title find multiple document
 * @param filter interface{}
 * @return models []*T your model slice
 * @return err error fail message from query
 */
func (e *Eloquent[T]) FindMultiple(filter interface{}) (models []*T, err error) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	cursor, err := coll.Find(context.TODO(), filter)

	if err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		err = e.errMsg(err)
		return
	}

	models = []*T{}

	if err = cursor.All(context.TODO(), &models); err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		err = e.errMsg(err)
		return
	}

	return
}

/**
 * @title insert a document
 * @param data interface{} your model struct
 * @return insertedID *primitive.ObjectID ObjectId of mongodb
 */
func (e *Eloquent[T]) Insert(data *T) (insertedID string, ok bool) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		ok = false
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		return
	}

	ok = true
	insertedID = result.InsertedID.(primitive.ObjectID).Hex()
	return
}

/**
 * @title insert multiple document
 * @param data []*T{} your model slice
 * @return InsertedIDs []string ObjectId of insert
 */
func (e *Eloquent[T]) InsertMultiple(data []*T) (InsertedIDs []string, ok bool) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)
	var slice []interface{}
	for _, value := range data {
		slice = append(slice, value)
	}
	InsertedIDs = []string{}

	result, err := coll.InsertMany(context.TODO(), slice)
	if err != nil {
		ok = false
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		return
	}

	for _, id := range result.InsertedIDs {
		idString := id.(primitive.ObjectID).Hex()
		InsertedIDs = append(InsertedIDs, idString)
	}
	ok = true
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
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo(1))
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
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		return
	}

	ok = true
	deleteCount = int(result.DeletedCount)
	return
}

/**
 * @title delete Multiple document
 * @param filter interface{} ex:bson.M{}, struct, bson.D{}
 * @return deleteCount int delete document count
 * @return ok bool query success or fail
 */
func (e *Eloquent[T]) DeleteMultiple(filter interface{}) (deleteCount int, ok bool) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	results, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		ok = false
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		return
	}

	ok = true
	deleteCount = int(results.DeletedCount)
	return
}

func (e *Eloquent[T]) Update(id string, data *T) (modifiedCount int, ok bool) {
	idH, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo(1))
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
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		ok = false
		return
	}

	ok = true
	modifiedCount = int(result.ModifiedCount)
	return
}

func (e *Eloquent[T]) UpdateMultiple(filter interface{}, data *T) (modifiedCount int, ok bool) {

	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)
	update := bson.M{"$set": data}

	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
		ok = false
		return
	}

	ok = true
	modifiedCount = int(result.ModifiedCount)
	return
}

func (e *Eloquent[T]) Count(filter interface{}) (count int, ok bool) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	if filter == nil {
		estCount, estCountErr := coll.EstimatedDocumentCount(context.TODO())
		if estCountErr != nil {
			ok = false
			logger.LogDebug.Error(e.logTitle, estCountErr, getCurrentFuncInfo(1))
			return
		}

		ok = true
		count = int(estCount)
	} else {
		countD, err := coll.CountDocuments(context.TODO(), filter)
		if err != nil {
			ok = false
			logger.LogDebug.Error(e.logTitle, err, getCurrentFuncInfo(1))
			return
		}

		ok = true
		count = int(countD)

	}

	return
}
