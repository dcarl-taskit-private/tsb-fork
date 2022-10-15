package tsb

import (
	"fmt"
)

type UartBaud uint16

const (
	UartBaudAuto UartBaud = iota
	UartBaud1200
	UartBaud2400
	UartBaud4800
	UartBaud9600
	UartBaud19200
	UartBaud38400
	UartBaud57600
	UartBaud115200
	UartBaud230400
	UartBaud460800
	UartBaud921600
	UartBaud1000000
	UartBaud1500000
	UartBaud3000000
)

type UartBits uint16

const (
	UartData8 UartBits = iota << 12
	UartData9
	UartData7
	UartData6
	UartData5
)

const (
	UartParityNone UartBits = iota << 10
	UartParityEven
	UartParityOdd
)

const (
	UartStopbits1 UartBits = iota << 8
	UartStopbits2
)

func (s *Server) UartInit(jack uint8, baud UartBaud, bits UartBits) (get chan []byte, put chan []byte, err error) {
	checkJack(jack)
	s.jack[jack].GetChan[TypRaw] = make(chan []byte, 10)
	get = s.jack[jack].GetChan[TypRaw]
	put = make(chan []byte, 10)
	go func(jack uint8) {
		for {
			select {
			case msg := <-put:
				{
					td := Packet{Ch: []byte{jack}, Typ: []byte{TypRaw}, Payload: msg}
					s.tdPutCh <- td
					s.redirect((td))
				}
			case <-s.done:
				{
					fmt.Printf("Uart %d closed!\n", jack)
					return
				}
			}
			td := Packet{Ch: []byte{jack}, Typ: []byte{TypRaw}, Payload: <-put}
			s.tdPutCh <- td
		}
	}(jack)
	return get, put, nil
}
