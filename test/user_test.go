package test

import (
	"github/LIOU2021/go-eloquent-mongodb/logger"
	"github/LIOU2021/go-eloquent-mongodb/test/services"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setup() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func cleanup() {
	logger.Close()
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

	userFind, ok := userService.Find("6380c2a141c9cfa264b345db")
	assert.True(t, ok, "find not ok")
	assert.True(t, userFind.ID != "", "id not find")
	logger.LogDebug.Infof("[userService@Find] - id : %s, name : %s, age : %d, created_time : %d, updated_time : %d\n", userFind.ID, userFind.Name, userFind.Age, userFind.CreatedAt, userFind.UpdatedAt)
}
