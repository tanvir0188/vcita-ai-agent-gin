package medicationreminder

import (
	"log"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

func ConvertHtmlToMarkDown(html string) string {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(html)
	if err != nil {
		log.Printf("failed to convert html to markdown: %v", err)
		return html
	}

	return markdown
}
