package discover

import "log"

type P2PServer struct {
	udp  *udp
	cfg  Config
	exit chan struct{}
}

func NewP2PServer(cfg Config) *P2PServer {
	if cfg.Laddr == "" {
		log.Println("no address implement to udp")
		return nil
	}

	if cfg.Id == "" {
		log.Println("no id implement to udp")
		return nil
	}

	server := &P2PServer{
		cfg:  cfg,
		exit: make(chan struct{}, 1),
	}

	return server
}

func (s *P2PServer) Start() {
	t := ListenUDP(s.cfg)
	if t == nil {
		s.exit <- struct{}{}
		return
	}
	s.udp = t

	s.run()
	log.Println("server start up...")
}

func (s *P2PServer) run() {

}

func (s *P2PServer) Stop() {
	for {
		select {
		case <-s.exit:
			return
		}
	}
}
