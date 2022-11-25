package models

type User struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	Age       uint16 `bson:"age"`
	CreatedAt uint16 `bson:"created_at"`
	UpdatedAt uint16 `bson:"updated_at"`
}
