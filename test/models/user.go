package models

type User struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	Age       uint16 `bson:"age"`
	CreatedAt uint64 `bson:"created_at"`
	UpdatedAt uint64 `bson:"updated_at"`
}

type UserCreateData struct {
	Name      string `bson:"name"`
	Age       uint16 `bson:"age"`
	CreatedAt uint64 `bson:"created_at"`
	UpdatedAt uint64 `bson:"updated_at"`
}
