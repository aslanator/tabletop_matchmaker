package link

import "tabletop_matchmaker/internal/types"

type CqData struct {
	types.CqData
	Confirmation bool    `json:"y"`
	Username     *string `json:"u"`
}

func newCqData(confirmation bool, username *string) CqData {
	return CqData{
		CqData: types.CqData{
			Command: Name(),
		},
		Confirmation: confirmation,
		Username:     username,
	}
}
