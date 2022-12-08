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

type Eloquent[t any] struct {
	db         string //db name
	Collection string
	uri        string
	logTitle   string
}

type IEloquent[T any] interface {
	All() (models []*T, err error)
	Find(id string) (model *T, err error)
	FindMultiple(filter any) (models []*T, err error)
	Insert(data *T) (insertedID string, err error)
	InsertMultiple(data []*T) (InsertedIDs []string, err error)
	Delete(id string) (deleteCount int, err error)
	DeleteMultiple(filter any) (deleteCount int, err error)
	Update(id string, data *T) (modifiedCount int, err error)
	UpdateMultiple(filter any, data *T) (modifiedCount int, err error)
	Count(filter any) (count int, err error)
	Paginate(limit int, page int, filter any) (paginated *Pagination[T], err error)
}

func NewEloquent[T any](collection string) *Eloquent[T] {
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

	cursor, errF := coll.Find(context.TODO(), bson.M{})

	if errF != nil {
		logger.LogDebug.Error(e.logTitle, errF, getCurrentFuncInfo(1))
		err = e.errMsg(errF)
		return
	}

	models = []*T{}

	if errA := cursor.All(context.TODO(), &models); errA != nil {
		logger.LogDebug.Error(e.logTitle, errA, getCurrentFuncInfo(1))
		err = e.errMsg(errA)
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
	idH, errP := primitive.ObjectIDFromHex(id)
	if errP != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo(1))
		err = e.errMsg(errP)
		return
	}

	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)
	model = new(T)
	errF := coll.FindOne(context.TODO(), bson.M{"_id": idH}).Decode(model)

	if errF == mongo.ErrNoDocuments {
		err = e.errMsg(errF)
		return
	} else if errF != nil {
		logger.LogDebug.Error(e.logTitle, errF, getCurrentFuncInfo(1))
		err = e.errMsg(errF)
		return
	}

	return
}

/**
 * @title find multiple document
 * @param filter any
 * @return models []*T your model slice
 * @return err error fail message from query
 */
func (e *Eloquent[T]) FindMultiple(filter any) (models []*T, err error) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	cursor, errF := coll.Find(context.TODO(), filter)

	if errF != nil {
		logger.LogDebug.Error(e.logTitle, errF, getCurrentFuncInfo(1))
		err = e.errMsg(errF)
		return
	}

	models = []*T{}

	if errA := cursor.All(context.TODO(), &models); errA != nil {
		logger.LogDebug.Error(e.logTitle, errA, getCurrentFuncInfo(1))
		err = e.errMsg(errA)
		return
	}

	return
}

/**
 * @title insert a document
 * @param data *T your model struct
 * @return insertedID string _id of mongodb
 * @return err error fail message from query
 */
func (e *Eloquent[T]) Insert(data *T) (insertedID string, err error) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	result, errI := coll.InsertOne(context.TODO(), data)
	if errI != nil {
		err = e.errMsg(errI)
		logger.LogDebug.Error(e.logTitle, errI, getCurrentFuncInfo(1))
		return
	}
	insertedID = result.InsertedID.(primitive.ObjectID).Hex()
	return
}

/**
 * @title insert multiple document
 * @param data []*T{} your model slice
 * @return InsertedIDs []string _id of mongodb
 * @return err error fail message from query
 */
func (e *Eloquent[T]) InsertMultiple(data []*T) (InsertedIDs []string, err error) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)
	var slice []any
	for _, value := range data {
		slice = append(slice, value)
	}

	result, errI := coll.InsertMany(context.TODO(), slice)
	if errI != nil {
		err = e.errMsg(errI)
		logger.LogDebug.Error(e.logTitle, errI, getCurrentFuncInfo(1))
		return
	}

	InsertedIDs = []string{}

	for _, id := range result.InsertedIDs {
		idString := id.(primitive.ObjectID).Hex()
		InsertedIDs = append(InsertedIDs, idString)
	}
	return
}

/**
 * @title delete a document
 * @param id string _id of document
 * @return deleteCount int delete document count
 * @return err error fail message from query
 */
func (e *Eloquent[T]) Delete(id string) (deleteCount int, err error) {
	idH, errP := primitive.ObjectIDFromHex(id)
	if errP != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo(1))
		err = e.errMsg(errP)
		return
	}

	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	filter := bson.M{"_id": idH}

	result, errD := coll.DeleteOne(context.TODO(), filter)
	if errD != nil {
		err = e.errMsg(errD)
		logger.LogDebug.Error(e.logTitle, errD, getCurrentFuncInfo(1))
		return
	}

	deleteCount = int(result.DeletedCount)
	return
}

/**
 * @title delete Multiple document
 * @param filter any ex:bson.M{}, struct, bson.D{}
 * @return deleteCount int delete document count
 * @return err error fail message from query
 */
func (e *Eloquent[T]) DeleteMultiple(filter any) (deleteCount int, err error) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	results, errD := coll.DeleteMany(context.TODO(), filter)
	if errD != nil {
		err = e.errMsg(errD)
		logger.LogDebug.Error(e.logTitle, errD, getCurrentFuncInfo(1))
		return
	}

	deleteCount = int(results.DeletedCount)
	return
}

/**
 * @title update a document
 * @param id string _id of mongodb
 * @return modifiedCount int modified document count
 * @return err error fail message from query
 */
func (e *Eloquent[T]) Update(id string, data *T) (modifiedCount int, err error) {
	idH, errP := primitive.ObjectIDFromHex(id)
	if errP != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo(1))
		err = e.errMsg(errP)
		return
	}

	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	filter := bson.M{"_id": idH}
	update := bson.M{"$set": data}

	result, errU := coll.UpdateOne(context.TODO(), filter, update)

	if errU != nil {
		logger.LogDebug.Error(e.logTitle, errU, getCurrentFuncInfo(1))
		err = e.errMsg(errU)
		return
	}

	modifiedCount = int(result.ModifiedCount)
	return
}

/**
 * @title update multiple document
 * @param filter any ex:struct, bson
 * @return modifiedCount int modified document count
 * @return err error fail message from query
 */
func (e *Eloquent[T]) UpdateMultiple(filter any, data *T) (modifiedCount int, err error) {

	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)
	update := bson.M{"$set": data}

	result, errU := coll.UpdateMany(context.TODO(), filter, update)
	if errU != nil {
		logger.LogDebug.Error(e.logTitle, errU, getCurrentFuncInfo(1))
		err = e.errMsg(errU)
		return
	}

	modifiedCount = int(result.ModifiedCount)
	return
}

/**
 * @title counter document
 * @param filter any you can use struct, bson,etc .., or nil
 * @return count int filter document count
 * @return err error fail message from query
 */
func (e *Eloquent[T]) Count(filter any) (count int, err error) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	if filter == nil {
		estCount, estCountErr := coll.EstimatedDocumentCount(context.TODO())
		if estCountErr != nil {
			err = e.errMsg(estCountErr)
			logger.LogDebug.Error(e.logTitle, estCountErr, getCurrentFuncInfo(1))
			return
		}

		count = int(estCount)
	} else {
		countD, errD := coll.CountDocuments(context.TODO(), filter)
		if errD != nil {
			err = e.errMsg(errD)
			logger.LogDebug.Error(e.logTitle, errD, getCurrentFuncInfo(1))
			return
		}

		count = int(countD)

	}

	return
}

/**
 * @title create pagination for your model data
 * @param limit int how many data display in each page. default=10
 * @param page int choose page for pagination. default=1
 * @param filter any you can use struct, bson,etc ... . but reject pass nil to filter
 * @return pagination *pagination[T]
 * @return err error fail message from query
 */
func (e *Eloquent[T]) Paginate(limit int, page int, filter any) (paginated *Pagination[T], err error) {
	client := e.Connect()

	defer e.Close(client)

	coll := e.GetCollection(client)

	total, totalErr := e.Count(filter)
	if totalErr != nil {
		err = e.errMsg(totalErr)
		logger.LogDebug.Error(e.logTitle, totalErr, getCurrentFuncInfo(1))
	}

	if limit < 1 {
		limit = 10
	}

	if page < 1 {
		page = 1
	}

	var to int
	var from int
	data := []*T{}
	lastPage := total / limit

	if total%limit >= 1 {
		lastPage++
	}

	if total > 0 {
		from = limit*(page-1) + 1
	} else {
		from = 0
	}

	if page == lastPage {
		to = total
	} else {
		to = limit * page
	}

	defer func() {
		paginated = newPagination(total, limit, page, lastPage, from, to, data)
	}()

	if total == 0 {
		return
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(limit * (page - 1)))

	cursor, errF := coll.Find(context.TODO(), filter, findOptions)
	if errF != nil {
		logger.LogDebug.Error(e.logTitle, errF, getCurrentFuncInfo(1))
		err = e.errMsg(errF)
		return
	}

	if errA := cursor.All(context.TODO(), &data); errA != nil {
		logger.LogDebug.Error(e.logTitle, errA, getCurrentFuncInfo(1))
		err = e.errMsg(errA)
		return
	}

	return
}
