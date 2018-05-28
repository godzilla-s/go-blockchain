package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type lvDatabase struct {
	db *leveldb.DB
}

// NewDB 创建一个leveldb
func NewDB(filepath string) (ldb *lvDatabase, err error) {
	db, err := leveldb.OpenFile(filepath, nil)
	ldb = &lvDatabase{
		db: db,
	}
	return
}

func (ldb *lvDatabase) Get(key []byte) (val []byte, err error) {
	val, err = ldb.db.Get(key, nil)
	return
}

func (ldb *lvDatabase) Put(key, val []byte) (err error) {
	err = ldb.db.Put(key, val, nil)
	return
}

func (ldb *lvDatabase) Delete(key []byte) (err error) {
	var exist bool
	exist, err = ldb.db.Has(key, nil)
	if err == nil {
		if exist {
			err = ldb.db.Delete(key, nil)
		}
	}
	return
}

func (ldb *lvDatabase) Close() (err error) {
	err = ldb.db.Close()
	return
}
