package test

import (
	"github/LIOU2021/go-eloquent-mongodb/core"
	"github/LIOU2021/go-eloquent-mongodb/logger"
	"github/LIOU2021/go-eloquent-mongodb/test/models"
	"github/LIOU2021/go-eloquent-mongodb/test/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	core.Setup()
}

func cleanup() {
	core.Cleanup()
}

var testId string

func Test_Insert(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	data := &models.UserCreateData{
		Name: "c8",
		Age:  110,
	}

	insertId, ok := userService.Insert(data)
	logger.LogDebug.Info("insertId : ", insertId)

	testId = insertId

	assert.True(t, ok, "insert not ok")
	assert.True(t, insertId != "", "id was null")
}

func Test_Find(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	userFind, ok := userService.Find(testId)
	assert.True(t, ok, "find not ok")
	assert.True(t, userFind.ID != "", "id not find")
	logger.LogDebug.Infof("[userService@Find] - id : %s, name : %s, age : %d, created_time : %d, updated_time : %d\n", userFind.ID, userFind.Name, userFind.Age, userFind.CreatedAt, userFind.UpdatedAt)
}

func Test_All(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	userAll, ok := userService.All()
	assert.True(t, ok, "all not ok")
	assert.GreaterOrEqual(t, len(userAll), 1, "no data")
	for i, v := range userAll {
		logger.LogDebug.Infof("index : %d, id : %s, name : %s, age : %d, created_time : %d, updated_time : %d\n", i, v.ID, v.Name, v.Age, v.CreatedAt, v.UpdatedAt)
		assert.True(t, v.ID != "", "_id is empty")
	}
}

func Test_Update(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	name := "LaLa"
	age := uint16(30)

	data := &models.UserUpdateData{
		Name: name,
		Age:  age,
	}
	updateCount, ok := userService.Update(testId, data)
	assert.True(t, ok, "updateCount not ok")
	assert.Equal(t, 1, updateCount, "find not ok")

	user, userOk := userService.Find(testId)
	assert.True(t, userOk, "updateCount for find user not ok")
	assert.Equal(t, name, user.Name, "update name not working")
	assert.Equal(t, age, user.Age, "update age not working")

}

func Test_Delete(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	deleteCount, ok := userService.Delete(testId)
	assert.True(t, ok, "delete not ok")
	assert.Equal(t, 1, deleteCount, "find not ok")
}
