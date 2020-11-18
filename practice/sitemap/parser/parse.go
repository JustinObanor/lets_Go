package parser

import (
	"io"

	"golang.org/x/net/html"
)

func ParsePage(r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return getLinks(doc), nil
}

func getLinks(n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				return []string{attr.Val}
			}
		}
	}

	var ret []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, getLinks(c)...)
	}

	return ret
}
