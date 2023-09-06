package markdown

import "fmt"

const parse_mode = "MarkdownV2"

func GenLink(text string, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

func GetParseMode() string {
	return parse_mode
}
