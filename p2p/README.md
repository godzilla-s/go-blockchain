## P2P 

TODO... 


测试用例: 
```go
import (
	"go-blockchain/p2p"
	"os"
)

func main() {
	cfg := p2p.NewConfig(os.Args[1:])
	server := p2p.NewP2PServer(cfg)
	if server == nil {
		return
	}

	server.Start()
	server.Stop()
}
```   

运行:  
```
启动中心节点 
go run main.go -d
启动其他节点
go run main.go -port xxx -id xxxx   
``` 

