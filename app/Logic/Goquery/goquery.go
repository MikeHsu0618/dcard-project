package Goquery

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
)

type Meta struct {
	Title       string `json:"title" form:"title"`
	Image       string `json:"image" form:"image"`
	Description string `json:"description" form:"description"`
}

func GetHtmlMeta(resultBody io.ReadCloser) Meta {
	doc, err := goquery.NewDocumentFromReader(resultBody)
	if err != nil {
		log.Fatal(err)
	}

	meta := Meta{}
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		content, _ := s.Attr("content")
		if property, _ := s.Attr("property"); property == "og:title" {
			fmt.Printf("title field: %s\n", content)
			meta.Title = content
		}

		if property, _ := s.Attr("property"); property == "og:description" {
			fmt.Printf("description field: %s\n", content)
			meta.Description = content
		}

		if property, _ := s.Attr("property"); property == "og:image" {
			fmt.Printf("image field: %s\n", content)
			meta.Image = content
		}
	})
	return meta
}
