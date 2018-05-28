package boltdb

import (
	"github.com/boltdb/bolt"
)

type blDatabase struct {
	db *bolt.DB
}

func NewDB(file string) (bdb *blDatabase, err error) {
	db, err := bolt.Open(file, 0640, nil)
	if err != nil {
		return
	}
	bdb = &blDatabase{
		db: db,
	}
	return
}

func (bl *blDatabase) Put(key, val []byte) (err error) {
	err = bl.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(""))
		if err != nil {
			return err
		}
		if err = b.Put(key, val); err != nil {
			return err
		}
		return nil
	})
	return
}

func (bl *blDatabase) Get(key []byte) (val []byte, err error) {
	err = bl.db.View(func(tx *bolt.Tx) error {
		val = tx.Bucket([]byte("")).Get(key)
		return nil
	})
	return
}

func (bl *blDatabase) Delete(key []byte) (err error) {
	err = bl.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(""))
		if err != nil {
			return err
		}
		if err = b.Delete(key); err != nil {
			return err
		}
		return nil
	})

	return
}
func (bl *blDatabase) Close() (err error) {
	err = bl.db.Close()
	return
}
