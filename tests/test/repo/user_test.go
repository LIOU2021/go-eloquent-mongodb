package repo

import (
	"os"
	"testing"

	"github.com/LIOU2021/go-eloquent-mongodb/core"
	"github.com/LIOU2021/go-eloquent-mongodb/tests/repositories"

	"github.com/stretchr/testify/assert"
)

func setup() {
	core.Setup()
}

func cleanup() {
	core.Cleanup()
}

var testId string

func TestMain(m *testing.M) {
	setup()

	exitCode := m.Run()

	defer func() {
		cleanup()

		os.Exit(exitCode)
	}()
}

func Test_User_GetUnderage(t *testing.T) {

	userRep := repositories.NewUserRepository()

	ageCondition := 30

	users, err := userRep.GetUnderage(ageCondition)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 1)

	for _, user := range users {
		assert.Less(t, *user.Age, ageCondition)
	}
}

func Test_User_GetOverage(t *testing.T) {

	userRep := repositories.NewUserRepository()

	ageCondition := 30

	users, err := userRep.GetOverage(ageCondition)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 1)

	for _, user := range users {
		assert.Greater(t, *user.Age, ageCondition)
	}
}
