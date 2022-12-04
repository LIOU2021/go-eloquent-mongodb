package models

type User struct {
	ID        *string `bson:"_id,omitempty"`
	Name      *string `bson:"name,omitempty"`
	Age       *uint16 `bson:"age,omitempty"`
	CreatedAt *uint64 `bson:"created_at,omitempty"`
	UpdatedAt *uint64 `bson:"updated_at,omitempty"`
}
