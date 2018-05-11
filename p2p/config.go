package p2p

type Config struct {
	Bootnodes []Node
	Laddr     string
	Id        string
}

var bootnode = [...]string{
	"node://abe245d34da3434341@192.168.1.190:60101",
}
