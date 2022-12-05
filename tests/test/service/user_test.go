package service

import (
	"strconv"
	"testing"
	"time"

	"github.com/LIOU2021/go-eloquent-mongodb/core"
	"github.com/LIOU2021/go-eloquent-mongodb/logger"
	"github.com/LIOU2021/go-eloquent-mongodb/tests/models"
	"github.com/LIOU2021/go-eloquent-mongodb/tests/services"
	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
)

func setup() {
	core.Setup()
}

func cleanup() {
	core.Cleanup()
}

var testId string

func Test_User_Insert_A_Document(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	name := "c8"
	age := uint16(110)
	data := &models.User{
		Name: &name,
		Age:  &age,
	}

	insertId, ok := userService.Insert(data)
	logger.LogDebug.Info("insertId : ", insertId)

	testId = insertId

	assert.True(t, ok, "insert not ok")
	assert.True(t, insertId != "", "id was null")
}

func Test_User_InsertMultiple(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	var data []*models.User
	currentTime := uint64(time.Now().Unix())
	count := 10
	for i := 0; i < count; i++ {
		age := uint16(1 + i*10)
		name := "serviceT_" + strconv.FormatInt(int64(i), 10)
		data = append(data, &models.User{
			Name:      &name,
			Age:       &age,
			CreatedAt: &currentTime,
			UpdatedAt: &currentTime,
		})
	}

	InsertedIDs, ok := userService.InsertMultiple(data)
	logger.LogDebug.Info("InsertedIDs  : ", InsertedIDs)

	assert.True(t, ok, "insertMultiple not ok")
	assert.Equal(t, count, len(InsertedIDs), "insertMultiple count miss match")
}

func Test_User_Find_A_Document(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	userFind, ok := userService.Find(testId)
	assert.True(t, ok, "find not ok")
	assert.True(t, *userFind.ID != "", "id not find")
	logger.LogDebug.Infof("[userService@Find] - id : %s, name : %s, age : %d, created_time : %d, updated_time : %d\n", *userFind.ID, *userFind.Name, *userFind.Age, *userFind.CreatedAt, *userFind.UpdatedAt)
}

func Test_User_Find__Multiple_Document(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	ageCondition := 30

	userFindMultiple, ok := userService.FindMultiple(bson.M{"age": bson.M{"$lte": ageCondition}})
	assert.True(t, ok, "findMultiple not ok")
	for _, value := range userFindMultiple {
		assert.True(t, *value.ID != "", "id not find")
		assert.LessOrEqual(t, *value.Age, uint16(ageCondition), "not less than 30")
		logger.LogDebug.Infof("[user@FindMultiple] - id : %s, name : %s, age : %d, created_time : %d, updated_time : %d\n", *value.ID, *value.Name, *value.Age, *value.CreatedAt, *value.UpdatedAt)
	}
}

func Test_User_All(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	userAll, ok := userService.All()
	assert.True(t, ok, "all not ok")
	assert.GreaterOrEqual(t, len(userAll), 1, "no data")
	for i, v := range userAll {
		logger.LogDebug.Infof("index : %d, id : %s, name : %s, age : %d, created_time : %d, updated_time : %d\n", i, *v.ID, *v.Name, *v.Age, *v.CreatedAt, *v.UpdatedAt)
		assert.True(t, *v.ID != "", "_id is empty")
	}
}

func Test_User_Update_a_document_by_full(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	name := "LaLa"
	age := uint16(30)

	data := &models.User{
		Name: &name,
		Age:  &age,
	}
	updateCount, ok := userService.Update(testId, data)
	assert.True(t, ok, "updateCount not ok")
	assert.Equal(t, 1, updateCount, "find not ok")

	user, userOk := userService.Find(testId)
	assert.True(t, userOk, "updateCount for find user not ok")
	assert.Equal(t, name, *user.Name, "update name not working")
	assert.Equal(t, age, *user.Age, "update age not working")

}

func Test_User_Update_A_Document_By_Part(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	age := uint16(time.Now().Second())

	data := &models.User{
		Age: &age,
	}

	updateCount, ok := userService.Update(testId, data)
	assert.True(t, ok, "updateCount not ok")
	assert.Equal(t, 1, updateCount, "find not ok")

	user, userOk := userService.Find(testId)
	assert.True(t, userOk, "updateCount for find user not ok")
	assert.True(t, "" != *user.Name, "name was change")
	assert.True(t, 0 != *user.CreatedAt, "CreatedAt was change")
	assert.True(t, 0 != *user.UpdatedAt, "UpdatedAt was change")
	assert.Equal(t, age, *user.Age, "update age not working")
}

func Test_User_Count_All(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()
	count, ok := userService.Count(nil)
	assert.True(t, ok, "count not ok")
	assert.GreaterOrEqual(t, count, 6, "count not working")
	t.Logf("total count : %d", count)
}

func Test_User_Count_Filter(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()
	filter := bson.M{"age": bson.M{"$lte": 30}} //less than or equal 30
	count, ok := userService.Count(filter)
	assert.True(t, ok, "count not ok")
	assert.GreaterOrEqual(t, count, 3, "count not working")
	t.Logf("filter count : %d", count)
}

func Test_User_Delete(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	deleteCount, ok := userService.Delete(testId)
	assert.True(t, ok, "delete not ok")
	assert.Equal(t, 1, deleteCount, "find not ok")
}
