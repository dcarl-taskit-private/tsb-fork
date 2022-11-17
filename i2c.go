package tsb

import (
	"fmt"
	"time"
)

func (s *Server) I2cInit(jack byte) (err error) {
	CheckJack(jack)
	/*
		for i := 0; i <= int(MaxJacks); i++ {
			s.jack[i].ReadChan[TypI2c] = make(chan byte, 1024)
			s.jack[i].ReadChan[TypPort] = make(chan byte, 1024)
			s.jack[i].ReadChan[TypRaw] = make(chan byte, 1024)
		}
	*/
	//s.jack[jack].ReadChan[TypI2c] = make(chan byte, 1024)
	return nil
}

func (s *Server) I2cSetAdr(jack byte, adr byte) (err error) {
	w := make([]byte, 2)
	w[0] = 0x80
	w[1] = adr
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypI2c}, Payload: w}
	s.tdPutCh <- td
	select {
	case a := <-s.jack[jack].ReadChan[TypI2c]:
		if a == 1 {
			return nil
		} else {
			return fmt.Errorf("wrong response: %x", a)
		}
	case <-time.After(1 * time.Second):
		return fmt.Errorf("timeout")
	}
}

func (s *Server) I2cRead(jack byte, b []byte) (n int, err error) {
	w := make([]byte, 1)
	n = len(b)
	if n > 127 {
		return 0, fmt.Errorf("only 127 bytes to write are allowed")
	}
	w[0] = byte(n)
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypI2c}, Payload: w}
	s.tdPutCh <- td
	for i := 0; i < n; i++ {
		select {
		case b[i] = <-s.jack[jack].ReadChan[TypI2c]:
		case <-time.After(1 * time.Second):
			return 0, fmt.Errorf("timeout")
		}
	}
	return n, nil
}

func (s *Server) I2cWrite(jack byte, b []byte) (n int, err error) {
	n = len(b)
	if n > 127 {
		return 0, fmt.Errorf("only 127 bytes to write are allowed")
	}
	w := make([]byte, n+1)
	w[0] = byte(n + 128)
	for i := 0; i < n; i++ {
		w[i+1] = b[i]
	}
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypI2c}, Payload: w[:n+1]}
	s.tdPutCh <- td
	select {
	case a := <-s.jack[jack].ReadChan[TypI2c]:
		if a == byte(n) {
			return 0, nil
		} else {
			return 0, fmt.Errorf("wrong response: %x", a)
		}
	case <-time.After(1 * time.Second):
		return 0, fmt.Errorf("timeout")
	}
}
