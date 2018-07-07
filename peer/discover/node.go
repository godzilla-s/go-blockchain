package discover

import "net"

type Endpoint struct {
	ID   string
	IP   net.IP
	Port int
}
