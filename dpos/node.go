package dpos

import "net"

type RoleType uint8

const (
	Follow RoleType = iota + 1
	Candidate
	Leader
)

type endpoing struct {
	IP   net.IP
	Port int
}

type Node struct {
	role RoleType
	ID   string
}
