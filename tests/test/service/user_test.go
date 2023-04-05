package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/LIOU2021/go-eloquent-mongodb/logger"
	"github.com/LIOU2021/go-eloquent-mongodb/orm"
	"github.com/LIOU2021/go-eloquent-mongodb/tests/models"
	"github.com/LIOU2021/go-eloquent-mongodb/tests/services"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
)

var testId string

func TestMain(m *testing.M) {
	orm.Setup("go-eloquent-mongo", "127.0.0.1", "27017", "")

	exitCode := m.Run()

	defer func() {
		os.Exit(exitCode)
	}()

}

func Test_User_Insert_A_Document(t *testing.T) {

	userService := services.NewUserService()

	name := "c8"
	age := 110
	data := &models.User{
		Name: &name,
		Age:  &age,
	}

	insertId, err := userService.Insert(data)
	logger.LogDebug.Info("insertId : ", insertId)

	testId = insertId

	assert.NoError(t, err, "insert not ok")
	assert.True(t, insertId != "", "id was null")
}

func Test_User_InsertMultiple(t *testing.T) {

	userService := services.NewUserService()

	var data []*models.User

	count := 10
	for i := 0; i < count; i++ {
		currentTime := time.Now().Unix()
		currentTime = currentTime + int64(i)
		age := 1 + i*10
		name := "serviceT_" + strconv.FormatInt(int64(i), 10)
		data = append(data, &models.User{
			Name:      &name,
			Age:       &age,
			CreatedAt: &currentTime,
			UpdatedAt: &currentTime,
		})
	}

	InsertedIDs, err := userService.InsertMultiple(data)
	logger.LogDebug.Info("InsertedIDs  : ", InsertedIDs)

	assert.NoError(t, err, "insertMultiple not ok")
	assert.Equal(t, count, len(InsertedIDs), "insertMultiple count miss match")
}

func Test_User_Find_A_Document(t *testing.T) {

	userService := services.NewUserService()

	userFind, userFindErr := userService.Find(testId)
	assert.NotNil(t, userFind, "userFind was nil")
	assert.NoError(t, userFindErr, "find not ok")
	if userFind != nil {
		assert.True(t, *userFind.ID != "", "id not find")
		logger.LogDebug.Infof("[userService@Find] - id : %s, name : %s, age : %d, created_at : %d, updated_at : %d\n", *userFind.ID, *userFind.Name, *userFind.Age, *userFind.CreatedAt, *userFind.UpdatedAt)
	}
}

func Test_User_Find_Multiple_Document(t *testing.T) {

	userService := services.NewUserService()

	ageCondition := 30

	userFindMultiple, err := userService.FindMultiple(bson.M{"age": bson.M{"$lte": ageCondition}})
	assert.NoError(t, err, "findMultiple not ok")
	for _, value := range userFindMultiple {
		assert.True(t, *value.ID != "", "id not find")
		assert.LessOrEqual(t, *value.Age, ageCondition, "not less than 30")
		logger.LogDebug.Infof("[user@FindMultiple] - id : %s, name : %s, age : %d, created_at : %d, updated_at : %d\n", *value.ID, *value.Name, *value.Age, *value.CreatedAt, *value.UpdatedAt)
	}
}

func Test_User_Find_Multiple_Document_With_Option(t *testing.T) {

	userService := services.NewUserService()

	ageCondition := 30
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	filter := bson.M{"age": bson.M{"$lte": ageCondition}}
	userFindMultiple, err := userService.FindMultiple(filter, opts)
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

	userService := services.NewUserService()

	userAll, err := userService.All()
	assert.NoError(t, err, "all not ok")
	assert.GreaterOrEqual(t, len(userAll), 1, "no data")
	for i, v := range userAll {
		logger.LogDebug.Infof("index : %d, id : %s, name : %s, age : %d, created_at : %d, updated_at : %d\n", i, *v.ID, *v.Name, *v.Age, *v.CreatedAt, *v.UpdatedAt)
		assert.True(t, *v.ID != "", "_id is empty")
	}
}

func Test_User_All_With_Option(t *testing.T) {

	userService := services.NewUserService()
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	userAll, err := userService.All(opts)
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

	userService := services.NewUserService()

	user, err := userService.Find(testId)
	assert.NoError(t, err)

	*user.Age = 30
	*user.Name = "user service update by full"
	updateCount, err := userService.Update(user)
	assert.NoError(t, err, "updateCount not ok")
	assert.Equal(t, 1, updateCount, "update not ok")

	user, userErr := userService.Find(testId)
	assert.NoError(t, userErr, "updateCount for find user not ok")
	assert.Equal(t, "user service update by full", *user.Name, "update name not working")
	assert.Equal(t, 30, *user.Age, "update age not working")
}

func Test_User_Update_A_Document_By_Part(t *testing.T) {

	userService := services.NewUserService()

	age := time.Now().Second()

	user, err := userService.Find(testId)
	assert.NoError(t, err)
	user.Age = &age
	updateCount, err := userService.Update(user)
	assert.NoError(t, err, "updateCount not ok")
	assert.Equal(t, 1, updateCount, "update not ok")

	user, userErr := userService.Find(testId)
	assert.NoError(t, userErr, "updateCount for find user not ok")
	assert.True(t, "" != *user.Name, "name was change")
	assert.True(t, 0 != *user.CreatedAt, "CreatedAt was change")
	assert.True(t, 0 != *user.UpdatedAt, "UpdatedAt was change")
	assert.Equal(t, age, *user.Age, "update age not working")
}

func Test_User_Count_All(t *testing.T) {

	userService := services.NewUserService()
	count, err := userService.Count(nil)
	assert.NoError(t, err, "count not ok")
	assert.GreaterOrEqual(t, count, 6, "count not working")
	t.Logf("total count : %d", count)
}

func Test_User_Count_Filter(t *testing.T) {

	userService := services.NewUserService()
	filter := bson.M{"age": bson.M{"$lte": 30}} //less than or equal 30
	count, err := userService.Count(filter)
	assert.NoError(t, err, "count not ok")
	assert.GreaterOrEqual(t, count, 3, "count not working")
	t.Logf("filter count : %d", count)
}

func Test_User_Paginate_Full(t *testing.T) {

	userService := services.NewUserService()
	currentPage := 2
	limit := 3
	pagination, err := userService.Paginate(limit, currentPage, bson.M{})
	assert.NoError(t, err, "paginate not ok")
	logger.LogDebug.Infof("pagination response : %+v", pagination)
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

	userService := services.NewUserService()
	currentPage := 2
	limit := 4
	filter := bson.M{"age": bson.M{"$gte": 30}}
	pagination, err := userService.Paginate(limit, currentPage, filter)
	assert.NoError(t, err, "paginate not ok")
	logger.LogDebug.Infof("pagination response : %+v", pagination)
	assert.Equal(t, 2, pagination.LastPage, "lastPage err")
	assert.Equal(t, 8, pagination.Total, "total err")
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

func Test_User_Update_Multiple_Document_By_Full(t *testing.T) {

	userService := services.NewUserService()

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
	updateCount, err := userService.UpdateMultiple(filter, data)
	assert.NoError(t, err, "updateCount not ok")
	assert.GreaterOrEqual(t, updateCount, 1, "update multiple not ok")
	t.Log("UpdateMultiple Count : ", updateCount)
	userFindMultiple, err := userService.FindMultiple(bson.M{"name": name})
	assert.NoError(t, err, "findMultiple not ok")

	assert.Equal(t, updateCount, len(userFindMultiple), "update count not match")
	for _, value := range userFindMultiple {
		assert.True(t, *value.ID != "", "id not find")
		assert.Equal(t, age, *value.Age, "update age not working")
		assert.Equal(t, currentTime, *value.UpdatedAt, "update time not working")
	}
}

func Test_User_Update_Multiple_Document_By_Part(t *testing.T) {

	userService := services.NewUserService()

	currentTime := time.Now().Unix()
	name := "LaLa-UpdateMultiple_test2_" + strconv.FormatInt(int64(currentTime), 10)
	data := &models.User{
		Name: &name,
	}

	ageCondition := 30
	filter := bson.M{"age": bson.M{"$gte": ageCondition}}
	updateCount, err := userService.UpdateMultiple(filter, data)
	assert.NoError(t, err, "updateCount not ok")
	assert.GreaterOrEqual(t, updateCount, 1, "update multiple not ok")
	t.Log("UpdateMultiple Count : ", updateCount)
	userFindMultiple, err := userService.FindMultiple(bson.M{"name": name})
	assert.NoError(t, err, "findMultiple not ok")

	assert.Equal(t, updateCount, len(userFindMultiple), "update count not match")
	for _, value := range userFindMultiple {
		assert.True(t, 0 != *value.CreatedAt, "CreatedAt was change")
		assert.True(t, 0 != *value.UpdatedAt, "UpdatedAt was change")
		assert.True(t, 0 != *value.Age, "age was change")
	}
}

func Test_User_GetUnderage(t *testing.T) {

	userService := services.NewUserService()

	ageCondition := 30

	users, err := userService.GetUnderage(ageCondition)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 1)

	for _, user := range users {
		assert.Less(t, *user.Age, ageCondition)
	}
}

func Test_User_GetOverage(t *testing.T) {

	userService := services.NewUserService()

	ageCondition := 30

	users, err := userService.GetOverage(ageCondition)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 1)
	for _, user := range users {
		assert.Greater(t, *user.Age, ageCondition)
	}
}

func Test_User_Delete_A_Document(t *testing.T) {

	userService := services.NewUserService()

	deleteCount, err := userService.Delete(testId)
	assert.NoError(t, err, "delete not ok")
	assert.Equal(t, 1, deleteCount, "delete not working")
}

func Test_User_Delete_Multiple_Document(t *testing.T) {

	userService := services.NewUserService()

	ageCondition := 99999
	deleteCount, err := userService.DeleteMultiple(bson.M{"age": bson.M{"$lte": ageCondition}})
	assert.NoError(t, err, "delete not ok")
	assert.Greater(t, deleteCount, 1, "deleteMultiple not working")
	t.Logf("DeleteMultiple : %d", deleteCount)
}
