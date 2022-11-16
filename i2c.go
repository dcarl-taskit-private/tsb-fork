package tsb

import (
	"fmt"
	"time"
)

func (s *Server) I2cInit(jack byte) (err error) {
	CheckJack(jack)
	s.jack[jack].ReadChan[TypI2c] = make(chan byte, 1024)
	return nil
}

func (s *Server) I2cRead(jack byte, adr byte, b []byte) (n int, err error) {
	var w [3]byte
	n = len(b)
	if n > 127 {
		return 0, fmt.Errorf("only 127 bytes to write are allowed")
	}
	w[0] = 128
	w[1] = adr
	w[2] = byte(n)
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypI2c}, Payload: w[:3]}
	s.tdPutCh <- td
	for i := 0; i < n; i++ {
		select {
		case b[i] = <-s.jack[jack].ReadChan[TypI2c]:
			fmt.Printf("Read from i2c: %x\n", b[i])
		case <-time.After(1 * time.Second):
			fmt.Printf("I2cRead Timeout\n")
		}
	}
	return n, nil
}

func (s *Server) I2cWrite(jack byte, adr byte, b []byte) (n int, err error) {
	var w [256]byte
	if len(b) > 127 {
		return 0, fmt.Errorf("only 127 bytes to write are allowed")
	}
	w[0] = 128
	w[1] = adr
	w[2] = byte(len(b) + 128)
	var i int
	for i = range b {
		w[i+3] = b[i]
	}
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypI2c}, Payload: w[:i+3]}
	s.tdPutCh <- td
	return i + 3, nil
}
