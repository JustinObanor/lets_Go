package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	href string
	text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, nil
	}

	var links []Link
	nodes := buildNode(doc)
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func buildText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += buildText(c)
	}

	return strings.Join(strings.Fields(ret), " ")
}

func buildLink(n *html.Node) Link {
	link := Link{}
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link.href = attr.Val
		}
	}

	link.text = buildText(n)

	return link
}

func buildNode(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, buildNode(c)...)
	}

	return ret
}
