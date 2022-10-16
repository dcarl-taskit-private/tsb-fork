package tsb

func (s *Server) PortInit(jack byte) (err error) {
	CheckJack(jack)
	s.jack[jack].ReadChan[TypPort] = make(chan byte, 1024)
	return nil
}
