package main

import (
	"fmt"
	"net/http"

	"github.com/leandroandrade/gophercises/url-shortener/first/urlshort"
	"github.com/leandroandrade/gophercises/url-shortener/first/file"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yaml, err := file.Read("./sources/urls.yml")
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	jsonfile, err := file.Read("./sources/urls.json")
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler([]byte(jsonfile), yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
