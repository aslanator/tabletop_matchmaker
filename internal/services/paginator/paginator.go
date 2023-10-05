package paginator

import (
	json2 "encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
	"tabletop_matchmaker/internal/helpers/keyboards"
	"tabletop_matchmaker/internal/types"
)

type Paginator[T any] struct {
	items      []T
	count      int
	pageSize   int
	pagesCount int
}

func NewPaginator[T any](items []T, pageSize int) Paginator[T] {
	count := len(items)
	pagesCount := int(math.Ceil(float64(count) / float64(pageSize)))

	return Paginator[T]{
		items:      items,
		count:      count,
		pageSize:   pageSize,
		pagesCount: pagesCount,
	}
}

func (p Paginator[T]) GetPage(page int) ([]T, error) {
	if page > p.pagesCount {
		return nil, errors.New("page does not exist")
	}

	offset := page - 1
	start := offset * p.pageSize
	end := start + p.pageSize
	if end > p.count {
		end = p.count
	}

	return p.items[start:end], nil
}

func (p Paginator[T]) GetPagesCount() int {
	return p.pagesCount
}

func (p Paginator[T]) GenKeyboardRow(page int, cqData types.PaginatorCqData) ([]tgbotapi.InlineKeyboardButton, error) {
	row := tgbotapi.NewInlineKeyboardRow()

	if p.pagesCount == 1 {
		return row, nil
	}

	if page-1 > 0 {
		cqData.SetPage(page - 1)
		json, err := json2.Marshal(cqData)

		if err != nil {
			return nil, fmt.Errorf("paginator error: %e", err)
		}

		row = append(row, tgbotapi.NewInlineKeyboardButtonData(keyboards.LEFT_ARROW, string(json)))
	} else {
		row = append(row, keyboards.GenBlankButton())
	}

	if page+1 < p.pagesCount {
		cqData.SetPage(page + 1)
		json, err := json2.Marshal(cqData)

		if err != nil {
			return nil, fmt.Errorf("paginator error: %e", err)
		}
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(keyboards.RIGHT_ARROW, string(json)))
	} else {
		row = append(row, keyboards.GenBlankButton())
	}

	return row, nil
}
