package main

import (
	"fmt"
	"log"

	"github.com/LIOU2021/go-eloquent-mongodb/core"
	"github.com/LIOU2021/go-eloquent-mongodb/orm"
)

type User struct {
	ID        *string `bson:"_id,omitempty"`
	Name      *string `bson:"name,omitempty"`
	Age       *int    `bson:"age,omitempty"`
	CreatedAt *int64  `bson:"created_at,omitempty"`
	UpdatedAt *int64  `bson:"updated_at,omitempty"`
}

func main() {
	core.Setup()
	defer core.Cleanup()

	userOrm := orm.NewEloquent[User]("users")
	id := "638f8be88d53c89a1c7d2e39"
	user, err := userOrm.Find(id)

	if err != nil {
		log.Fatal("user id not found !")
	}

	fmt.Printf("id : %s, name : %s, age : %d, created_time : %d, updated_time : %d\n", *user.ID, *user.Name, *user.Age, *user.CreatedAt, *user.UpdatedAt)
}
