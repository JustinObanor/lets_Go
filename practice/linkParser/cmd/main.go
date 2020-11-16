package main

import (
	"fmt"
	"strings"

	"github.com/lets_Go/practice/linkParser/link"
)

var exampleHTML = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">
    A link to another page
    <span> some span  </span>
  </a>
  <a href="/page-two">A link to a second page</a>
</body>
</html>
`

func main() {
	r := strings.NewReader(exampleHTML)

	link, err := link.Parse(r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", link)
}
