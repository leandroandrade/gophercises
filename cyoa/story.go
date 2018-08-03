package cyoa

import (
	"io"
	"encoding/json"
	"net/http"
	"html/template"
	"strings"
	"log"
)

var defaultHandlerTemplate = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Choose Your Own Adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
        <p>{{.}}</p>
    {{end}}

    <ul>
    {{range .Options}}
        <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
    {{end}}
    </ul>
</body>`

var templateDefault *template.Template

func init() {
	templateDefault = template.Must(template.ParseFiles("page.html"))
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func JsonStory(r io.Reader) (Story, error) {
	var story Story
	if err := json.NewDecoder(r).Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

func NewHandler(s Story) handler {
	return handler{story: s}
}

type handler struct {
	story Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]

	if chapter, ok := h.story[path]; ok {
		if err := templateDefault.Execute(w, chapter); err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found.", http.StatusNotFound)
}
