package urlshort

import (
	"net/http"
	"gopkg.in/yaml.v2"
	"encoding/json"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if destination, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, destination, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []pathUrl
	if err := yaml.Unmarshal(yml, &pathUrls); err != nil {
		return nil, err
	}
	return MapHandler(parserToMap(pathUrls), fallback), nil
}

func JSONHandler(jsonbytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []pathUrl
	if err := json.Unmarshal(jsonbytes, &pathUrls); err != nil {
		return nil, err
	}
	return MapHandler(parserToMap(pathUrls), fallback), nil
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
