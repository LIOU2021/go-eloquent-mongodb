package main

import (
	"fmt"
	"log"

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
	orm.Setup("go-eloquent-mongo", "127.0.0.1", "27017", "")
	orm.Connect()
	defer orm.Disconnect()

	userOrm := orm.NewEloquent[User]("users")
	id := "642d5b2298ba2bb73c55e5c4"
	user, err := userOrm.Find(id)

	if err != nil {
		log.Fatal("user id not found !")
	}

	fmt.Printf("id : %s, name : %s, age : %d, created_at : %d, updated_at : %d\n", *user.ID, *user.Name, *user.Age, *user.CreatedAt, *user.UpdatedAt)
}
