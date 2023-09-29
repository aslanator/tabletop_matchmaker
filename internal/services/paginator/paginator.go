package paginator

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const MAX_PAGE_SIZE int = 2000

type Paginator struct {}

func (paginator Paginator) genPageList(page int, pagesCount int) []int {
	pageList := []int{}

	if page > 1 {
		pageList = append(pageList, page-1)
	}

	if page < pagesCount {
		pageList = append(pageList, page+1)
	}

	return pageList
}

func (paginator Paginator) genPageButton(page int, command string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(page), 
	command + "|" + strconv.Itoa(page))
}

func (paginator Paginator) genPageButtons(pageList []int, command string) tgbotapi.InlineKeyboardMarkup {
	buttons := []tgbotapi.InlineKeyboardButton {}

	for _, page := range pageList {
		buttons = append(buttons, paginator.genPageButton(page, command))
	}
	row := tgbotapi.NewInlineKeyboardRow(buttons...)

	return tgbotapi.NewInlineKeyboardMarkup(row)
}

func (paginator Paginator) genPageKeyboardMarkup(page int, pagesCount int, command string) tgbotapi.InlineKeyboardMarkup {
	pageList := paginator.genPageList(page, pagesCount)
	return paginator.genPageButtons(pageList, command)
}

func (paginator Paginator) getPagesCount(text string) int {
	collectionSize := len(text)
	pagesCount := collectionSize / MAX_PAGE_SIZE
	if collectionSize % MAX_PAGE_SIZE != 0 {
		pagesCount++
	}
	return pagesCount
}

// make function that return page message ends on line before maximum size and keyboard markup
func (paginator Paginator) GenPageMessage(text string, page int, command string) (string, tgbotapi.InlineKeyboardMarkup) {
	pagesCount := paginator.getPagesCount(text)
	if page > pagesCount {
		page = pagesCount
	}
	startIndex := (page - 1) * MAX_PAGE_SIZE
	endIndex := page * MAX_PAGE_SIZE
	if endIndex > len(text) {
		endIndex = len(text)
	}
	pageText := text[startIndex:endIndex]
	lastIndex := strings.LastIndex(pageText, "\n")
	keyboardMarkup := paginator.genPageKeyboardMarkup(page, pagesCount, command)
	return pageText[0:lastIndex], keyboardMarkup
}