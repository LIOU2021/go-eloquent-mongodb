package orm

type pagination[T any] struct {
	Total       uint `json:"total"`
	PerPage     uint `json:"per_page"`
	CurrentPage uint `json:"current_page"`
	LastPage    uint `json:"last_page"`
	From        uint `json:"from"`
	To          uint `json:"to"`
	Data        []*T `json:"data"`
}

func newPagination[T any]() *pagination[T] {
	return &pagination[T]{}
}
