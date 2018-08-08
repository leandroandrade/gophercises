package cyoa

import (
	"io"
	"encoding/json"
	"net/http"
	"html/template"
	"strings"
	"log"
)

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

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.tpl = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func NewHandler(s Story, opts ...HandlerOption) handler {
	h := handler{story: s, tpl: templateDefault, pathFn: defaultPath}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	story  Story
	tpl    *template.Template
	pathFn func(r *http.Request) string
}

func defaultPath(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	if chapter, ok := h.story[path]; ok {
		if err := h.tpl.Execute(w, chapter); err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found.", http.StatusNotFound)
}
