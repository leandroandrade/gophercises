package database

import (
	"github.com/coreos/bbolt"
	"log"
	"time"
)

const Bucket = "paths"
const DBName = "urlshortener.db"

type BoltDB struct {
	boltdb *bolt.DB
}

func Connect() (*BoltDB, error) {
	db, err := bolt.Open(DBName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	if err := createBucket(db); err != nil {
		return nil, err
	}

	return &BoltDB{boltdb: db}, nil
}

func createBucket(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(Bucket))
		if err != nil {
			log.Printf("bucket %v already exist!\n", Bucket)
		}
		return nil
	})

	return err
}

func (b BoltDB) BatchSavePathURL(pathToUrls map[string]string) error {
	err := b.boltdb.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))
		for key, value := range pathToUrls {
			err := bucket.Put([]byte(key), []byte(value))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (b BoltDB) FindUrl(path string) (string, error) {
	var url string
	err := b.boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))
		url = string(bucket.Get([]byte(path)))
		return nil
	})

	return url, err
}
