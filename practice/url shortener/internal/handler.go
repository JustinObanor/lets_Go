package internal

import (
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
	yaml "gopkg.in/yaml.v2"
)

const DBBucket = "pathsToUrls"

type YAMLData struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

type Database struct {
	DB *bolt.DB
}

func (db *Database) DBHandler(fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(DBBucket))

			if url := b.Get([]byte(r.URL.Path)); url != nil {
				http.RedirectHandler(string(url), http.StatusPermanentRedirect).ServeHTTP(w, r)
			} else {
				fallback.ServeHTTP(w, r)
			}
			return nil
		})
	}
}

func (db *Database) YAMLHandler(y []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(y)
	if err != nil {
		return nil, err
	}

	db.buildDB(parsedYaml)

	return db.DBHandler(fallback), nil
}

func ReadYAML(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func parseYAML(y []byte) ([]YAMLData, error) {
	var data []YAMLData

	if err := yaml.Unmarshal(y, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) buildDB(parsedYAML []YAMLData) {
	for _, data := range parsedYAML {
		db.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(DBBucket))

			return b.Put([]byte(data.Path), []byte(data.URL))
		})
	}
}
