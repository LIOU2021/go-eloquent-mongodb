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

func Test_Find(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()

	userFind, ok := userService.Find("6380c8a9185309a5944a3171")
	assert.True(t, ok, "find not ok")
	assert.True(t, userFind.ID != "", "id not find")
	logger.LogDebug.Infof("[userService@Find] - id : %s, name : %s, age : %d, created_time : %d, updated_time : %d\n", userFind.ID, userFind.Name, userFind.Age, userFind.CreatedAt, userFind.UpdatedAt)
}

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
	assert.True(t, ok, "insert not ok")
	assert.True(t, insertId != "", "id was null")
}