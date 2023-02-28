package db

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

var defaultBucket = []byte("defaultBucket")

// Database is an open bolt database
type Database struct {
	db *bolt.DB
}

// NewDatabase return an  instance  of  a database that we can work with
func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {
	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db = &Database{
		db: boltDb,
	}
	closeFunc = boltDb.Close
	if err := db.CreateDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("create default bucket: %w", err)
	}
	return db, closeFunc, nil
}

func (d *Database) CreateDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

// SetKey sets the key to the requested value into the default database or returns an error
func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})
}

// GetKey gets the key to the requested value into the default database or returns an error
func (d *Database) GetKey(key string) ([]byte, error) {
	var result []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		return nil
	})

	if err == nil {
		return result, nil
	}
	return nil, err
}
