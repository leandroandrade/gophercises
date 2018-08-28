package main

import (
	"os"
	"golang.org/x/net/html"
	"github.com/leandroandrade/gophercises/link"
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	file, err := os.Open("ex2.html")
	if err != nil {
		panic(err)
	}

	document, err := html.Parse(file)
	if err != nil {
		panic(err)
	}

	jsonValue, err := process(document)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonValue))

}

func process(document *html.Node) ([]byte, error) {
	links := make([]link.Link, 0)

	parserHTML(document, &links)

	jsonValue, err := json.MarshalIndent(links, "", " ")
	if err != nil {
		return nil, err
	}
	return jsonValue, nil
}

func parserHTML(tag *html.Node, links *[]link.Link) {
	if tag.Type == html.ElementNode && tag.Data == "a" {
		for _, attr := range tag.Attr {
			if attr.Key == "href" {
				*links = append(*links, link.Link{
					Href: attr.Val,
					Text: strings.TrimSpace(tag.FirstChild.Data),
				})
			}
		}
	}

	for c := tag.FirstChild; c != nil; c = c.NextSibling {
		parserHTML(c, links)
	}
}
