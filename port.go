package tsb

import (
	"fmt"
)

func (s *Server) PortInit(jack byte) (get chan []byte, put chan []byte, err error) {
	checkJack(jack)
	s.jack[jack].GetChan[TypPort] = make(chan []byte, 10)
	get = s.jack[jack].GetChan[TypPort]
	put = make(chan []byte, 10)
	go func(jack uint8) {
		for {
			select {
			case msg := <-put:
				{
					td := Packet{Ch: []byte{jack}, Typ: []byte{TypPort}, Payload: msg}
					s.tdPutCh <- td
					//s.redirect((td))
				}
			case <-s.done:
				{
					fmt.Printf("Port %d closed!\n", jack)
					return
				}
			}
			td := Packet{Ch: []byte{jack}, Typ: []byte{TypPort}, Payload: <-put}
			s.tdPutCh <- td
		}
	}(jack)
	return get, put, nil
}
