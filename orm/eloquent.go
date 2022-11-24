package orm

type Eloquent struct {
	DB         string
	Collection string
}

type IEloquent interface {
	All() []*Model
	Find(id string) *Model
	Insert(data interface{}) (ok bool)
	Delete(filter interface{}) (ok bool)
	Update(filter interface{}, data interface{}) (ok bool)
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

func (e *Eloquent) Find(id string) *Model {
	return &Model{}
}

func (e *Eloquent) Insert(data interface{}) (ok bool) {
	return
}

func (e *Eloquent) Delete(filter interface{}) (ok bool) {
	return
}

func (e *Eloquent) Update(filter interface{}, data interface{}) (ok bool) {
	return
}
