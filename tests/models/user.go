package models

type User struct {
	ID        *string `bson:"_id,omitempty"`
	Name      *string `bson:"name,omitempty"`
	Age       *int    `bson:"age,omitempty"`
	CreatedAt *int64  `bson:"created_at,omitempty"`
	UpdatedAt *int64  `bson:"updated_at,omitempty"`
}
