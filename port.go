package tsb

func (s *Server) PortInit(jack byte) (err error) {
	CheckJack(jack)
	s.jack[jack].ReadChan[TypPort] = make(chan byte, 1024)
	return nil
}

func (s *Server) PortPutc(jack byte, c byte) (err error) {
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypPort}, Payload: []byte{c}}
	s.tdPutCh <- td
	return nil
}
