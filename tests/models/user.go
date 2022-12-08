package models

type User struct {
	ID        *string `bson:"_id,omitempty" json:"id"`
	Name      *string `bson:"name,omitempty" json:"name"`
	Age       *int    `bson:"age,omitempty" json:"age"`
	CreatedAt *int64  `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt *int64  `bson:"updated_at,omitempty" json:"updated_at"`
}
