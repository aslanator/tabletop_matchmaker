package getcollection

import "tabletop_matchmaker/internal/types"

type CqData struct {
	types.CqData
	Username string `json:"u"`
	Page     int    `json:"p"`
}

func (cqData *CqData) SetPage(page int) {
	cqData.Page = page
}

func newCqData(username string, page int) CqData {
	return CqData{
		CqData: types.CqData{
			Command: Name(),
		},
		Username: username,
		Page:     page,
	}
}
