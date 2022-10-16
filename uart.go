package tsb

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

func (s *Server) UartInit(jack byte, baud UartBaud, bits UartBits) (err error) {
	CheckJack(jack)
	s.jack[jack].ReadChan[TypRaw] = make(chan byte, 1024)
	/*
		get = s.jack[jack].ReadChan[TypRaw]
		put = make(chan []byte, 10)
		go func(jack uint8) {
			for {
				select {
				case msg := <-put:
					{
						td := Packet{Ch: []byte{jack}, Typ: []byte{TypRaw}, Payload: msg}
						s.tdPutCh <- td
						//s.redirect((td))
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
	*/
	return nil
}

func (s *Server) UartWrite(jack byte, b []byte) (n int, err error) {
	td := Packet{Ch: []byte{byte(jack)}, Typ: []byte{TypRaw}, Payload: b}
	s.tdPutCh <- td
	return len(b), nil
}

func (s *Server) UartRead(jack byte, b []byte) (n int, err error) {
	n = len(s.jack[jack].ReadChan[TypRaw])
	if n > len(b) {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		b[i] = <-s.jack[jack].ReadChan[TypRaw]
	}
	return n, nil
}
