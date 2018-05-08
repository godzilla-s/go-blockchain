package p2p

import (
	"net"
)

type udp struct {
}

func ListenUDP(cfg Config) {
	laddr, err := net.ResolveUDPAddr("udp", cfg.Laddr)
	if err != nil {
		return
	}

	net.ListenUDP("udp", laddr)
}
