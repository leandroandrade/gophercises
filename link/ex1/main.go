package main

import (
	"strings"
	"github.com/leandroandrade/gophercises/link"
	"fmt"
)

var exampleHTML = `<html>
<body>
<h1>Hello!</h1>
<a href="/other-page">A link to another page <span> using span </span></a>
</body>
</html>`

func main() {
	r := strings.NewReader(exampleHTML)
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", links)
}
