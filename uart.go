package tsb

type UartBaud uint16

const (
	UartBaudAuto uint16 = iota
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
	UartData8 uint16 = iota << 12
	UartData9
	UartData7
	UartData6
	UartData5
)

const (
	UartParityNone uint16 = iota << 10
	UartParityEven
	UartParityOdd
)

const (
	UartStopbits1 uint16 = iota << 8
	UartStopbits2
)

func (s *Server) UartInit(jack byte, baud uint16, bits uint16) (err error) {
	CheckJack(jack)
	s.I2cWrite(jack, 130, []byte{byte(baud), byte(baud >> 8), byte(bits), byte(bits >> 8)})
	return nil
}

func (s *Server) UartWrite(jack byte, b []byte) (n int, err error) {
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypRaw}, Payload: b}
	s.tdPutCh <- td
	return len(b), nil
}

func (s *Server) UartRead(jack byte, b []byte) (n int, err error) {
	b[0] = <-s.jack[jack].ReadChan[TypRaw]
	n = len(s.jack[jack].ReadChan[TypRaw]) + 1
	if n > len(b) {
		n = len(b)
	}
	for i := 1; i < n; i++ {
		b[i] = <-s.jack[jack].ReadChan[TypRaw]
	}
	return n, nil
}
