package urlshort

import (
	"net/http"
	"encoding/json"
	"github.com/leandroandrade/gophercises/url-shortener/second/database"
	"log"
	"github.com/leandroandrade/gophercises/url-shortener/second/file"
	"gopkg.in/yaml.v2"
)

func Handler(pathsToUrls map[string]string, db *database.BoltDB, fallback http.Handler) http.HandlerFunc {
	if err := db.BatchSavePathURL(pathsToUrls); err != nil {
		log.Fatal(err)
	}

	if err := YAMLHandler(db); err != nil {
		log.Fatal(err)
	}

	if err := JSONHandler(db); err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if destination, err := db.FindUrl(path); err == nil && destination != "" {
			http.Redirect(w, r, destination, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(db *database.BoltDB) error {
	yml, err := file.Read("./sources/urls.yml")
	if err != nil {
		panic(err)
	}

	var pathUrls []pathUrl
	if err := yaml.Unmarshal(yml, &pathUrls); err != nil {
		return err
	}
	return processBatchURLs(db, pathUrls)
}

func JSONHandler(db *database.BoltDB) error {
	jsonfile, err := file.Read("./sources/urls.json")
	if err != nil {
		panic(err)
	}

	var pathUrls []pathUrl
	if err := json.Unmarshal(jsonfile, &pathUrls); err != nil {
		return err
	}
	return processBatchURLs(db, pathUrls)
}

func processBatchURLs(db *database.BoltDB, pathUrls []pathUrl) error {
	if err := db.BatchSavePathURL(parserToMap(pathUrls)); err != nil {
		return err
	}
	return nil
}

func parserToMap(pathUrls []pathUrl) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.URL
	}
	return pathToUrls
}

type pathUrl struct {
	Path string `yaml:"path" json:"path,omitempty"`
	URL  string `yaml:"url" json:"url,omitempty"`
}
