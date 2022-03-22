package goquery

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
)

type Meta struct {
	Title       string `json:"title" form:"title" example:"Meta title"`
	Image       string `json:"image" form:"image" example:"Meta Image Url"`
	Description string `json:"description" form:"description" example:"Meta Description"`
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
			meta.Title = content
		}

		if property, _ := s.Attr("property"); property == "og:description" {
			meta.Description = content
		}

		if property, _ := s.Attr("property"); property == "og:image" {
			meta.Image = content
		}
	})
	return meta
}
