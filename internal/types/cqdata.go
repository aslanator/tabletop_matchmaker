package types

type CqData struct {
	Command string `json:"C"`
}

type PaginatorCqData interface {
	SetPage(page int)
}
