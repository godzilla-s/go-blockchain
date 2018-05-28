package storage

import (
	"errors"
	"fmt"
	"go-blockchain/storage/boltdb"
	"go-blockchain/storage/leveldb"
)

// 定义一个存储接口: 增，删，改
type Database interface {
	Get(key []byte) (val []byte, err error)
	Put(key, val []byte) (err error)
	Delete(key []byte) (err error)
	Close() (err error)
	//NewBatch()
}

// 数据
type DataBaseType int

const (
	TypeLevelDB DataBaseType = 1 << iota
	TypeBoltDB
)

// NewDatabase 创建一个database
func NewDatabase(typ DataBaseType) (Database, error) {
	switch typ {
	case TypeLevelDB:
		return leveldb.NewDB("")
	case TypeBoltDB:
		return boltdb.NewDB("")
	}
	return nil, errors.New("")
}

func (dtyp DataBaseType) String() string {
	switch dtyp {
	case TypeLevelDB:
		return "leveldb"
	case TypeBoltDB:
		return "boltdb"
	}
	return fmt.Sprintf("unknow:%d", dtyp)
}
