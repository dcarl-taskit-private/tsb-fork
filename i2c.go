package tsb

func (s *Server) I2cInit(jack byte) (err error) {
	CheckJack(jack)
	s.jack[jack].ReadChan[TypI2c] = make(chan byte, 1024)
	return nil
}

func (s *Server) I2cRead(jack byte, adr byte, count byte) (read []byte, err error) {
	return nil, nil
}

func (s *Server) I2cWrite(jack byte, adr byte, data []byte) (written byte, err error) {
	return 0, nil
}
