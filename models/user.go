package models

type User struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
	Age  uint16 `bson:"age"`
}
