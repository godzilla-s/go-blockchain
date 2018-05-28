## 存储 
作为区块链重要部分之一，存储选择的数据库与一般数据库不一样，区块链的存储基本都是本地存储，像关系型这样笨重数据库显然不太符合要求。能够做到读写快速，且使用便捷，是选取存储的几个要点之一。K/V型的数据便是首选之一，下面介绍比较常用的K/V存储数据库。

### leveldb  
原生是google用C++写的，go也有相应的版本，可以通过下面方式获取：
``` 
go get -u github.com/syndtr/goleveldb/leveldb
``` 
缺点：levelDB不支持transaction； 

### boltdb   
原生go写的数据库，支持事务性，通过下面获取:
```
go get -u github.com/boltdb/bolt
```

### leveldb API 

```go
// 创建leveldb
db, err := leveldb.OpenFile(file, &opt.Options{
	OpenFilesCacheCapacity: x,  // 文件缓存的大小
	BlockCacheCapacity:     x,  // 每个块缓存的大小
	WriteBuffer:            x,  // 每次写入大小
	Filter:                 filter.NewBloomFilter(10),
})
// 如果发生错误，则恢复文件
if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
	db, err = leveldb.RecoverFile(file, nil)
}
```