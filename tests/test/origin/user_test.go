package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/LIOU2021/go-eloquent-mongodb/logger"
	"github.com/LIOU2021/go-eloquent-mongodb/orm"
	"github.com/LIOU2021/go-eloquent-mongodb/tests/models"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
)

var testId string

func TestMain(m *testing.M) {
	orm.Setup("go-eloquent-mongo", "127.0.0.1", "27017", "")
	ctx := context.Background()
	orm.Connect(ctx)
	exitCode := m.Run()
	defer func() {
		orm.Disconnect(ctx)
		os.Exit(exitCode)
	}()
}

func Test_User_Insert_A_Document(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")
	currentTime := time.Now().Unix()

	name := "c6"
	age := 54
	data := &models.User{
		Name:      &name,
		Age:       &age,
		CreatedAt: &currentTime,
		UpdatedAt: &currentTime,
	}
	insertId, err := userOrm.Insert(context.Background(), data)
	logger.LogDebug.Info("insertId : ", insertId)

	testId = insertId

	assert.NoError(t, err, "insert not ok")
	assert.True(t, insertId != "", "id was null")
}

func Test_User_Insert_Multiple(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")

	var data []*models.User

	count := 10

	for i := 0; i < count; i++ {
		currentTime := time.Now().Unix()
		currentTime = currentTime + int64(i)
		name := "c8_" + strconv.FormatInt(int64(i), 10)
		age := 1 + i*10
		data = append(data, &models.User{
			Name:      &name,
			Age:       &age,
			CreatedAt: &currentTime,
			UpdatedAt: &currentTime,
		})
	}

	InsertedIDs, err := userOrm.InsertMultiple(context.Background(), data)
	logger.LogDebug.Info("InsertedIDs  : ", InsertedIDs)

	assert.NoError(t, err, "insertMultiple not ok")
	assert.Equal(t, count, len(InsertedIDs), "insertMultiple count miss match")
}

func Test_User_Find_A_Document(t *testing.T) {
	userOrm := orm.NewEloquent[models.User]("users")
	userFind, err := userOrm.Find(context.Background(), testId)
	assert.NotNil(t, userFind, "userFind was nil")
	assert.NoError(t, err, "find not ok")
	if userFind != nil {
		assert.True(t, *userFind.ID != "", "id not find")
		logger.LogDebug.Infof("[user@Find] - id : %s, name : %s, age : %d, created_at : %d, updated_at : %d\n", *userFind.ID, *userFind.Name, *userFind.Age, *userFind.CreatedAt, *userFind.UpdatedAt)
	}
}

func Test_User_Find_Multiple_Document(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")

	ageCondition := 30

	userFindMultiple, err := userOrm.FindMultiple(context.Background(), bson.M{"age": bson.M{"$lte": ageCondition}})
	assert.NoError(t, err, "findMultiple not ok")
	for _, value := range userFindMultiple {
		assert.True(t, *value.ID != "", "id not find")
		assert.LessOrEqual(t, *value.Age, ageCondition, fmt.Sprintf("not less than %d", ageCondition))
		logger.LogDebug.Infof("[user@FindMultiple] - id : %s, name : %s, age : %d, created_at : %d, updated_at : %d\n", *value.ID, *value.Name, *value.Age, *value.CreatedAt, *value.UpdatedAt)
	}
}

func Test_User_Find_Multiple_Document_With_Option(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")

	ageCondition := 30
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	filter := bson.M{"age": bson.M{"$lte": ageCondition}}
	userFindMultiple, err := userOrm.FindMultiple(context.Background(), filter, opts)
	assert.NoError(t, err, "findMultiple not ok")
	preCreatedAt := int64(0)
	for index, value := range userFindMultiple {
		if index != 0 {
			assert.Less(t, *value.CreatedAt, preCreatedAt)
		}
		preCreatedAt = *value.CreatedAt

		assert.True(t, *value.ID != "", "id not find")
		assert.LessOrEqual(t, *value.Age, ageCondition, fmt.Sprintf("not less than %d", ageCondition))
		logger.LogDebug.Infof("[user@FindMultiple] - id : %s, name : %s, age : %d, created_at : %d, updated_at : %d\n", *value.ID, *value.Name, *value.Age, *value.CreatedAt, *value.UpdatedAt)
	}
}

func Test_User_All(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")

	userAll, err := userOrm.All(context.Background())
	assert.NoError(t, err, "all not ok")
	assert.GreaterOrEqual(t, len(userAll), 1, "no data")
	for i, v := range userAll {
		logger.LogDebug.Infof("index : %d, id : %s, name : %s, age : %d, created_at : %d, updated_at : %d\n", i, *v.ID, *v.Name, *v.Age, *v.CreatedAt, *v.UpdatedAt)
		assert.True(t, *v.ID != "", "_id is empty")
	}
}

func Test_User_All_With_Option(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	userAll, err := userOrm.All(context.Background(), opts)
	assert.NoError(t, err, "all not ok")
	assert.GreaterOrEqual(t, len(userAll), 1, "no data")
	preCreatedAt := int64(0)
	for i, v := range userAll {
		if i != 0 {
			assert.Less(t, *v.CreatedAt, preCreatedAt)
		}
		preCreatedAt = *v.CreatedAt
		logger.LogDebug.Infof("index : %d, id : %s, name : %s, age : %d, created_at : %d, updated_at : %d\n", i, *v.ID, *v.Name, *v.Age, *v.CreatedAt, *v.UpdatedAt)
		assert.True(t, *v.ID != "", "_id is empty")
	}
}

func Test_User_Update_A_Document_By_Full(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")

	name := "LaLa"
	age := 30
	currentTime := time.Now().Unix()

	data := &models.User{
		Name:      &name,
		Age:       &age,
		UpdatedAt: &currentTime,
	}
	updateCount, err := userOrm.Update(context.Background(), testId, data)
	assert.NoError(t, err, "updateCount not ok")
	assert.Equal(t, 1, updateCount, "update not ok")

	user, userErr := userOrm.Find(context.Background(), testId)
	assert.NoError(t, userErr, "updateCount for find user not ok")
	assert.Equal(t, name, *user.Name, "update name not working")
	assert.Equal(t, age, *user.Age, "update age not working")
	assert.Equal(t, currentTime, *user.UpdatedAt, "update time not working")
}

func Test_User_Update_A_Document_By_Part(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")

	age := time.Now().Second()

	data := &models.User{
		Age: &age,
	}

	updateCount, err := userOrm.Update(context.Background(), testId, data)
	assert.NoError(t, err, "updateCount not ok")
	assert.Equal(t, 1, updateCount, "update not ok")

	user, userErr := userOrm.Find(context.Background(), testId)
	assert.NoError(t, userErr, "updateCount for find user not ok")
	assert.True(t, "" != *user.Name, "name was change")
	assert.True(t, 0 != *user.CreatedAt, "CreatedAt was change")
	assert.True(t, 0 != *user.UpdatedAt, "UpdatedAt was change")
	assert.Equal(t, age, *user.Age, "update age not working")
}

func Test_User_Count_All(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")
	count, err := userOrm.Count(context.Background(), nil)
	assert.NoError(t, err, "count not ok")
	assert.GreaterOrEqual(t, count, 6, "count not working")
	t.Logf("total count : %d", count)
}

func Test_User_Count_Filter(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")
	filter := bson.M{"age": bson.M{"$lte": 30}} //less than or equal 30
	count, err := userOrm.Count(context.Background(), filter)
	assert.NoError(t, err, "count not ok")
	assert.GreaterOrEqual(t, count, 3, "count not working")
	t.Logf("filter count : %d", count)
}

func Test_User_Paginate_Full(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")
	currentPage := 2
	limit := 3
	pagination, err := userOrm.Paginate(context.Background(), limit, currentPage, bson.M{})
	assert.NoError(t, err, "paginate not ok")
	t.Logf("pagination response : %+v", pagination)
	assert.Equal(t, 4, pagination.LastPage, "lastPage err")
	assert.Equal(t, 11, pagination.Total, "total err")
	assert.Equal(t, currentPage, pagination.CurrentPage, "currentPage err")
	assert.Equal(t, limit, pagination.PerPage, "PerPage err")
	assert.Equal(t, limit*(currentPage-1)+1, pagination.From, "From err")
	assert.Equal(t, limit*currentPage, pagination.To, "To err")

	var preCreatedAt int64
	for index, value := range pagination.Data {
		if index > 0 {
			assert.Less(t, *value.CreatedAt, preCreatedAt, "order by created_at desc fail")
		}
		preCreatedAt = *value.CreatedAt
		t.Logf("pagination data - id : %s, name : %s. age : %d, created_at : %d, updated_at : %d", *value.ID, *value.Name, *value.Age, *value.CreatedAt, *value.UpdatedAt)
	}

	jsonResponse, err := json.Marshal(pagination)
	assert.NoError(t, err)
	t.Log(string(jsonResponse))
}

func Test_User_Paginate_Filter(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")
	currentPage := 2
	limit := 4
	filter := bson.M{"age": bson.M{"$gte": 30}}
	pagination, err := userOrm.Paginate(context.Background(), limit, currentPage, filter)
	assert.NoError(t, err, "paginate not ok")
	t.Logf("pagination response : %+v", pagination)
	assert.Equal(t, 2, pagination.LastPage, "lastPage err")
	assert.Equal(t, 8, pagination.Total, "total err")
	assert.Equal(t, currentPage, pagination.CurrentPage, "currentPage err")
	assert.Equal(t, limit, pagination.PerPage, "PerPage err")
	assert.Equal(t, limit*(currentPage-1)+1, pagination.From, "From err")
	expectTo := 0
	if pagination.Total%limit == 0 {
		expectTo = limit * currentPage
	} else {
		expectTo = limit*currentPage - (limit*pagination.LastPage - pagination.Total)
	}
	assert.Equal(t, expectTo, pagination.To, "To err")

	var preCreatedAt int64
	for index, value := range pagination.Data {
		if index > 0 {
			assert.Less(t, *value.CreatedAt, preCreatedAt, "order by created_at desc fail")
		}
		preCreatedAt = *value.CreatedAt
		t.Logf("pagination data - id : %s, name : %s. age : %d, created_at : %d, updated_at : %d", *value.ID, *value.Name, *value.Age, *value.CreatedAt, *value.UpdatedAt)
	}

	jsonResponse, err := json.Marshal(pagination)
	assert.NoError(t, err)
	t.Log(string(jsonResponse))
}

func Test_User_Update_Multiple_Document_By_Full(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")

	age := 134
	currentTime := time.Now().Unix()
	name := "LaLa-UpdateMultiple_" + strconv.FormatInt(int64(currentTime), 10)
	data := &models.User{
		Name:      &name,
		Age:       &age,
		UpdatedAt: &currentTime,
	}

	ageCondition := 30
	filter := bson.M{"age": bson.M{"$lte": ageCondition}}
	updateCount, err := userOrm.UpdateMultiple(context.Background(), filter, data)
	assert.NoError(t, err, "updateCount not ok")
	assert.GreaterOrEqual(t, updateCount, 1, "update multiple not ok")
	t.Log("UpdateMultiple Count : ", updateCount)
	userFindMultiple, err := userOrm.FindMultiple(context.Background(), bson.M{"name": name})
	assert.NoError(t, err, "findMultiple not ok")

	assert.Equal(t, updateCount, len(userFindMultiple), "update count not match")
	for _, value := range userFindMultiple {
		assert.True(t, *value.ID != "", "id not find")
		assert.Equal(t, age, *value.Age, "update age not working")
		assert.Equal(t, currentTime, *value.UpdatedAt, "update time not working")
	}
}

func Test_User_Update_Multiple_Document_By_Part(t *testing.T) {
	userOrm := orm.NewEloquent[models.User]("users")

	currentTime := time.Now().Unix()
	name := "LaLa-UpdateMultiple_test2_" + strconv.FormatInt(int64(currentTime), 10)
	data := &models.User{
		Name: &name,
	}

	ageCondition := 30
	filter := bson.M{"age": bson.M{"$gte": ageCondition}}
	updateCount, err := userOrm.UpdateMultiple(context.Background(), filter, data)
	assert.NoError(t, err, "updateCount not ok")
	assert.GreaterOrEqual(t, updateCount, 1, "update multiple not ok")
	t.Log("UpdateMultiple Count : ", updateCount)
	userFindMultiple, err := userOrm.FindMultiple(context.Background(), bson.M{"name": name})
	assert.NoError(t, err, "findMultiple not ok")

	assert.Equal(t, updateCount, len(userFindMultiple), "update count not match")
	for _, value := range userFindMultiple {
		assert.True(t, 0 != *value.CreatedAt, "CreatedAt was change")
		assert.True(t, 0 != *value.UpdatedAt, "UpdatedAt was change")
		assert.True(t, 0 != *value.Age, "age was change")
	}
}

func Test_User_Delete_A_Document(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")

	deleteCount, err := userOrm.Delete(context.Background(), testId)
	assert.NoError(t, err, "delete not ok")
	assert.Equal(t, 1, deleteCount, "delete not working")
}

func Test_User_Delete_Multiple_Document(t *testing.T) {

	userOrm := orm.NewEloquent[models.User]("users")

	ageCondition := 99999
	deleteCount, err := userOrm.DeleteMultiple(context.Background(), bson.M{"age": bson.M{"$lte": ageCondition}})
	assert.NoError(t, err, "delete not ok")
	assert.Greater(t, deleteCount, 1, "deleteMultiple not working")
	t.Logf("DeleteMultiple : %d", deleteCount)
}
