package orm

type pagination[T any] struct {
	Total       int  `json:"total"`
	PerPage     int  `json:"per_page"`
	CurrentPage int  `json:"current_page"`
	LastPage    int  `json:"last_page"`
	From        int  `json:"from"`
	To          int  `json:"to"`
	Data        []*T `json:"data"`
}

func newPagination[T any](total int, limit int, page int, lastPage int, from int, to int, data []*T) *pagination[T] {
	return &pagination[T]{
		Total:       total,
		PerPage:     limit,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        from,
		To:          to,
		Data:        data,
	}
}
