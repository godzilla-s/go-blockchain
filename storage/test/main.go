// 测试用例
package main

import (
	"go-blockchain/storage"
	"log"
)

func lvlDatabase() {
	db, err := storage.NewDatabase(storage.TypeLevelDB)
	if err != nil {
		log.Fatalf("fail to new database: %v", err)
	}
	log.Println("write data ======")
	db.Put([]byte("key001"), []byte("value001"))
	log.Println("get data ======")
	db.Get([]byte("key001"))
}

func main() {
	lvlDatabase()
}
