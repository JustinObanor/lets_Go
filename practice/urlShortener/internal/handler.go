package internal

import (
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
	yaml "gopkg.in/yaml.v2"
)

const DBBucket = "pathsToUrls"

type PathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

type Database struct {
	DB *bolt.DB
}

func (db *Database) DBHandler(fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		db.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(DBBucket))

			if url := b.Get([]byte(path)); url != nil {
				http.Redirect(w, r, string(url), http.StatusFound)
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

func parseYAML(y []byte) ([]PathURL, error) {
	var datas []PathURL

	if err := yaml.Unmarshal(y, &datas); err != nil {
		return nil, err
	}

	return datas, nil
}

func (db *Database) buildDB(parsedYAML []PathURL) {
	for _, data := range parsedYAML {
		db.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(DBBucket))

			return b.Put([]byte(data.Path), []byte(data.URL))
		})
	}
}
