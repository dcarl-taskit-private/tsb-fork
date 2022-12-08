package tsb

const (
	UartBaudAuto byte = iota
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

const (
	UartStopbits1 byte = iota
	UartStopbits15
	UartStopbits2
)

const (
	UartParityNone byte = iota << 2
	UartParityEven
	UartParityOdd
)

const (
	UartData8 byte = iota << 4
	UartData9
	UartData7
	UartData6
	UartData5
)

func (s *Server) UartInit(jack byte, baud byte, parity byte, datalen byte, databits byte) (err error) {
	CheckJack(jack)
	s.I2cSetAdr(jack, JackModeReg)
	s.I2cWrite(jack, []byte{JackUart})
	s.I2cSetAdr(jack, JackUartReg)
	s.I2cWrite(jack, []byte{baud, parity | datalen | databits})
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
