package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/leandroandrade/gophercises/cyoa"
	"log"
	"net/http"
	"strings"
	"html/template"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA Story")
	flag.Parse()

	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	//tpl := template.Must(template.New("").Parse("Hello!"))
	//handler := cyoa.NewHandler(story, cyoa.WithTemplate(tpl))

	//handler := cyoa.NewHandler(story)
	//fmt.Printf("Starting the server on port: %d\n", *port)
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))

	customTemplate := template.Must(template.New("").Parse(customTemplateHTML))
	h := cyoa.NewHandler(story,
		cyoa.WithTemplate(customTemplate),
		cyoa.WithPathFunc(customPrefixPath),
	)

	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))

}

func customPrefixPath(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var customTemplateHTML = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Choose Your Own Adventure **CUSTOM**</title>
</head>
<body>
    <section class="page">
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
            <p>{{.}}</p>
        {{end}}

        <ul>
        {{range .Options}}
            <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
    </section>
    <style>
        body {
            font-family: helvetica, arial;
        }
        h1 {
            text-align:center;
            position:relative;
        }
        .page {
            width: 80%;
            max-width: 500px;
            margin: 40px auto;
            padding: 80px;
            background: #FCF6FC;
            border: 1px solid #eee;
            box-shadow: 0 10px 6px -6px #797;
        }
        ul {
            border-top: 1px dotted #ccc;
            padding: 10px 0 0 0;
            -webkit-padding-start: 0;
        }
        li {
            padding-top: 10px;
        }
        a,
        a:visited {
            text-decoration: underline;
            color: #555;
        }
        a:active,
        a:hover {
            color: #222;
        }
        p {
            text-indent: 1em;
        }
    </style>
</body>
</html>`
