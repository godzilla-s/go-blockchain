package p2p

import (
	"flag"
	"fmt"
	"net"
)

type Config struct {
	Bootnodes []*Node
	Laddr     string
	Id        string
	args      map[string]interface{}
}

var bootnode = [...]string{
	"node://abe245d34da3434341@192.168.1.100:60101",
}

func DefaultNodes() []*Node {
	var nodes []*Node
	for _, s := range bootnode {
		nodes = append(nodes, MustParseNode(s))
	}

	return nodes
}

var (
	host    string
	port    int
	id      string
	node    string
	defmode bool
)

func parseArgs(args []string) {
	flag.StringVar(&host, "host", "", "p2p host addr")
	flag.IntVar(&port, "port", 0, "p2p port")
	flag.StringVar(&id, "id", "", "identify of node")
	flag.BoolVar(&defmode, "d", false, "use default node")
}

func NewConfig(args []string) Config {
	parseArgs(args)
	flag.CommandLine.Parse(args)

	var cfg Config
	if host == "" {
		host = GetLocalAddr()
	}

	defnodes := DefaultNodes()
	if defmode {
		host = defnodes[0].IP.String()
		port = defnodes[0].UDP
		id = defnodes[0].ID.String()
	}

	cfg.Laddr = fmt.Sprintf("%s:%d", host, port)
	cfg.Id = id
	cfg.Bootnodes = DefaultNodes()
	return cfg
}

func GetLocalAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
