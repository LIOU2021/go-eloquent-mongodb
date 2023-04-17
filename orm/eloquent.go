package orm

import (
	"context"

	"github.com/LIOU2021/go-eloquent-mongodb/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	logger.Init()
}

var conf *config
var conn *mongo.Client

type config struct {
	// db name ex:go-eloquent-mongodb
	DB string
	// mongodb host ex:127.0.0.1
	Host string
	// mongodb port ex:27017
	Port string
	// mongodb user
	User string
	// mongodb password
	Password string
}

// setup mongodb connect config
func Setup(db, host, port, password string) {
	if conf != nil {
		return
	}
	conf = &config{
		DB:       db,
		Host:     host,
		Port:     port,
		Password: password,
	}
}

type Eloquent[t any] struct {
	db         string //db name
	Collection string
	uri        string
	logTitle   string
}

type IEloquent[T any] interface {
	GetCollection() *mongo.Collection
	All(ctx context.Context, opts ...*options.FindOptions) (models []*T, err error)
	Find(ctx context.Context, id string) (model *T, err error)
	FindMultiple(ctx context.Context, filter any, opts ...*options.FindOptions) (models []*T, err error)
	Insert(ctx context.Context, data *T) (insertedID string, err error)
	InsertMultiple(ctx context.Context, data []*T) (InsertedIDs []string, err error)
	Delete(ctx context.Context, id string) (deleteCount int, err error)
	DeleteMultiple(ctx context.Context, filter any) (deleteCount int, err error)
	Update(ctx context.Context, id string, data *T) (modifiedCount int, err error)
	UpdateMultiple(ctx context.Context, filter any, data *T) (modifiedCount int, err error)
	Count(ctx context.Context, filter any) (count int, err error)
	Paginate(ctx context.Context, limit int, page int, filter any) (paginated *Pagination[T], err error)
}

func NewEloquent[T any](collection string) *Eloquent[T] {
	return &Eloquent[T]{
		db:         conf.DB,
		Collection: collection,
		uri:        getUri(),
		logTitle:   getLogTitle(collection),
	}
}

/**
 * @title connect mongodb server
 */
func Connect(ctx context.Context) *mongo.Client {
	if conn != nil {
		return conn
	}
	uri := getUri()
	if uri == "" {
		logger.LogDebug.Fatal(`[connect fail]: `, "You must set your 'mongodb_host' and 'mongodb_port' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable", getCurrentFuncInfo(1))
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.LogDebug.Fatal(`[connect fail]: `, err, getCurrentFuncInfo(1))
	}
	conn = client
	return conn
}

/**
 * @title disconnect mongodb server
 */
func Disconnect(ctx context.Context) {
	if conn == nil {
		return
	}

	if err := conn.Disconnect(ctx); err != nil {
		logger.LogDebug.Error(`[disconnect fail]: `, err, getCurrentFuncInfo(1))
	}
}

/**
 * @title get collection instance
 */
func (e *Eloquent[T]) GetCollection() *mongo.Collection {
	if conn == nil {
		return nil
	}
	return conn.Database(e.db).Collection(e.Collection)
}

/**
 * @title get all document from collection
 * @return models []*T your model slice
 * @return err error fail message from query
 */
func (e *Eloquent[T]) All(ctx context.Context, opts ...*options.FindOptions) (models []*T, err error) {
	coll := e.GetCollection()
	cursor, errF := coll.Find(ctx, bson.M{}, opts...)

	if errF != nil {
		logger.LogDebug.Error(e.logTitle, errF, getCurrentFuncInfo(1))
		err = e.errMsg(errF)
		return
	}
	defer cursor.Close(ctx)

	models = []*T{}

	for cursor.Next(ctx) {
		model := new(T)
		if errNext := cursor.Decode(&model); errNext != nil {
			logger.LogDebug.Error(e.logTitle, errNext, getCurrentFuncInfo(1))
			err = e.errMsg(errNext)
			return
		}
		models = append(models, model)
	}

	if errC := cursor.Err(); errC != nil {
		logger.LogDebug.Error(e.logTitle, errC, getCurrentFuncInfo(1))
		err = e.errMsg(errC)
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
func (e *Eloquent[T]) Find(ctx context.Context, id string) (model *T, err error) {
	idH, errP := primitive.ObjectIDFromHex(id)
	if errP != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo(1))
		err = e.errMsg(errP)
		return
	}

	coll := e.GetCollection()
	model = new(T)
	errF := coll.FindOne(ctx, bson.M{"_id": idH}).Decode(model)

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
func (e *Eloquent[T]) FindMultiple(ctx context.Context, filter any, opts ...*options.FindOptions) (models []*T, err error) {
	coll := e.GetCollection()
	cursor, errF := coll.Find(ctx, filter, opts...)

	if errF != nil {
		logger.LogDebug.Error(e.logTitle, errF, getCurrentFuncInfo(1))
		err = e.errMsg(errF)
		return
	}
	defer cursor.Close(ctx)
	models = []*T{}

	if errA := cursor.All(ctx, &models); errA != nil {
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
func (e *Eloquent[T]) Insert(ctx context.Context, data *T) (insertedID string, err error) {
	coll := e.GetCollection()

	result, errI := coll.InsertOne(ctx, data)
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
func (e *Eloquent[T]) InsertMultiple(ctx context.Context, data []*T) (InsertedIDs []string, err error) {
	coll := e.GetCollection()
	var slice []any
	for _, value := range data {
		slice = append(slice, value)
	}

	result, errI := coll.InsertMany(ctx, slice)
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
func (e *Eloquent[T]) Delete(ctx context.Context, id string) (deleteCount int, err error) {
	idH, errP := primitive.ObjectIDFromHex(id)
	if errP != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo(1))
		err = e.errMsg(errP)
		return
	}

	coll := e.GetCollection()

	filter := bson.M{"_id": idH}

	result, errD := coll.DeleteOne(ctx, filter)
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
func (e *Eloquent[T]) DeleteMultiple(ctx context.Context, filter any) (deleteCount int, err error) {
	coll := e.GetCollection()

	results, errD := coll.DeleteMany(ctx, filter)
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
func (e *Eloquent[T]) Update(ctx context.Context, id string, data *T) (modifiedCount int, err error) {
	idH, errP := primitive.ObjectIDFromHex(id)
	if errP != nil {
		logger.LogDebug.Error(e.logTitle, "_id Hex fail", getCurrentFuncInfo(1))
		err = e.errMsg(errP)
		return
	}

	coll := e.GetCollection()

	filter := bson.M{"_id": idH}
	update := bson.M{"$set": data}

	result, errU := coll.UpdateOne(ctx, filter, update)

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
func (e *Eloquent[T]) UpdateMultiple(ctx context.Context, filter any, data *T) (modifiedCount int, err error) {
	coll := e.GetCollection()
	update := bson.M{"$set": data}

	result, errU := coll.UpdateMany(ctx, filter, update)
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
func (e *Eloquent[T]) Count(ctx context.Context, filter any) (count int, err error) {
	coll := e.GetCollection()

	if filter == nil {
		estCount, estCountErr := coll.EstimatedDocumentCount(context.TODO())
		if estCountErr != nil {
			err = e.errMsg(estCountErr)
			logger.LogDebug.Error(e.logTitle, estCountErr, getCurrentFuncInfo(1))
			return
		}

		count = int(estCount)
	} else {
		countD, errD := coll.CountDocuments(ctx, filter)
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
func (e *Eloquent[T]) Paginate(ctx context.Context, limit int, page int, filter any) (paginated *Pagination[T], err error) {
	coll := e.GetCollection()

	total, totalErr := e.Count(ctx, filter)
	if totalErr != nil {
		err = e.errMsg(totalErr)
		logger.LogDebug.Error(e.logTitle, totalErr, getCurrentFuncInfo(1))
		return
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
		if page == lastPage {
			to = total
		} else {
			to = limit * page
		}
	} else {
		from = 0
		to = 0
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

	cursor, errF := coll.Find(ctx, filter, findOptions)
	if errF != nil {
		logger.LogDebug.Error(e.logTitle, errF, getCurrentFuncInfo(1))
		err = e.errMsg(errF)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		model := new(T)
		if errNext := cursor.Decode(&model); err != nil {
			logger.LogDebug.Error(e.logTitle, errNext, getCurrentFuncInfo(1))
			err = e.errMsg(errNext)
			return
		}
		data = append(data, model)
	}

	if errC := cursor.Err(); err != nil {
		logger.LogDebug.Error(e.logTitle, errC, getCurrentFuncInfo(1))
		err = e.errMsg(errC)
		return
	}
	return
}
