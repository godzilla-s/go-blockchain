package rpc

import (
	"reflect"
	"sync"
)

type registerService map[string]*service
type callbacks map[string]*callback

// 服务
type service struct {
	name      string
	typ       reflect.Type
	callbacks callbacks
}

// 回调处理函数
type callback struct {
	rval     reflect.Value  // 参数
	method   reflect.Method // 调用方法
	argTypes []reflect.Type // 调用参数
}

type RPCServer struct {
	services registerService
	run      int32
	mux      sync.Mutex
}

type RPCService struct {
	server *RPCServer
}

func NewServer() *RPCServer {
	server := &RPCServer{
		services: make(registerService),
		run:      1,
	}

	rpcService := &RPCService{server}
	server.RegisterName("", rpcService)
	return server
}

func (s *RPCServer) RegisterName(name string, svc interface{}) error {
	if s.services == nil {
		s.services = make(registerService)
	}

	service := new(service)
	service.typ = reflect.TypeOf(svc)
	service.name = name
	service.callbacks = nil

	s.services[name] = service
	return nil
}
