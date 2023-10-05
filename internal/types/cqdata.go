package types

type CqData struct {
	Command string `json:"C"`
}

type PaginatorCqData interface {
	GetPage() int
	SetPage(page int)
}
