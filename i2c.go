package tsb

import "fmt"

func (s *Server) I2cInit(jack byte) (get chan []byte, put chan []byte, err error) {
	checkJack(jack)
	s.jack[jack].ReadChan[TypI2c] = make(chan []byte, 10)
	get = s.jack[jack].ReadChan[TypI2c]
	put = make(chan []byte, 10)
	go func(jack uint8) {
		for {
			select {
			case msg := <-put:
				{
					td := Packet{Ch: []byte{jack}, Typ: []byte{TypI2c}, Payload: msg}
					s.tdPutCh <- td
					//s.redirect((td))
				}
			case <-s.done:
				{
					fmt.Printf("I2C %d closed!\n", jack)
					return
				}
			}
			td := Packet{Ch: []byte{jack}, Typ: []byte{TypI2c}, Payload: <-put}
			s.tdPutCh <- td
		}
	}(jack)
	return get, put, nil
}

func (s *Server) I2cRead(jack byte, adr byte, count byte) (read []byte, err error) {
	return nil, nil
}

func (s *Server) I2cWrite(jack byte, adr byte, data []byte) (written byte, err error) {
	return 0, nil
}
