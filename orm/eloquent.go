package orm

type Eloquent struct {
	DB         string
	Collection string
}

func NewEloquent(collection string) *Eloquent {
	return &Eloquent{
		DB:         "go-eloquent-mongodb",
		Collection: collection,
	}
}

func (e *Eloquent) All() []*Model {
	return []*Model{}
}

func (e *Eloquent) find(id string) *Model {
	return &Model{}
}
