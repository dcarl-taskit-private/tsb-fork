package tsb

import (
	"fmt"
	"log"
	"net"

	"github.com/tarm/serial"
)

const (
	MaxJacks byte = 8
)

type JackMode uint16

const (
	JackPorts JackMode = iota
	JackUart485
	JackUart232
	JackI2c
	JackSpi8
	JackSpi16
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

type PortMode uint16

type jack struct {
	enabled    bool
	jackConfig uint16 // 4 Bits PortEnable, 4 Bits jackMode
	uartConfig uint16 // 8 Bits uartBits, 8 Bits uartBaud
	portConfig uint16 // 4 Bits Pullup, 4 Bits Direction, 4 Bits InputNotification, 4 Bits Output
	i2cConfig  uint16 // 8 Bits Address
	GetChan    [TypError - TypModbus + 1]chan []byte
	PutChan    [TypError - TypModbus + 1]chan []byte
}

type Server struct {
	address string
	typ     string
	jack    [MaxJacks]jack
	conn    net.Conn
	sport   *serial.Port
	tdPutCh chan DataTsb
	tdGetCh chan DataTsb
	done    chan struct{}
}

func NewSerialServer(address string) (Server, error) {
	var err error
	s := Server{address: address}
	s.typ = "Serial"
	s.sport, err = serial.OpenPort(&serial.Config{Name: address, Baud: 115200})
	if err != nil {
		log.Fatal(err)
	}
	s.tdPutCh = PutData(s.sport)
	s.tdGetCh, s.done = GetData(s.sport)
	s.serv()
	return s, nil
}

func NewTcpServer(address string) (server, error) {
	var err error
	s := server{address: address}
	s.typ = "TCP"
	s.conn, err = net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	s.tdPutCh = PutData(s.conn)
	s.tdGetCh, s.done = GetData(s.conn)
	s.serv()
	return s, nil
}

func (s *Server) serv() {
	fmt.Printf("TSB client connected to tsb server: %s", s.address)
	go func() {
		for {
			select {
			case td := <-s.tdGetCh:
				{
					s.redirect((td))
				}
			case <-s.done:
				{
					fmt.Printf("TSB client connection closed!\n")
					return
				}
			}
		}
	}()
}

func (s *Server) redirect(td DataTsb) {
	c := td.Typ[0]
	if c < 0 || c > (TypError-TypModbus) {
		log.Printf("Unknown Typ!\n\r")
		return
	}
	if td.Ch[0] > MaxJacks || td.Ch[0] < 1 {
		log.Printf("Invalid Jacknr!\n\r")
		return
	}
	if s.jack[td.Ch[0]].GetChan[c] == nil {
		log.Printf("Not initialized!\n\r")
		return
	}
	if len(s.jack[td.Ch[0]].GetChan[c]) >= cap(s.jack[td.Ch[0]].GetChan[c]) {

	}
	s.jack[td.Ch[0]].GetChan[c] <- td.Payload
}

func (s *Server) UartInit(jack uint8, baud UartBaud, bits UartBits) (get chan []byte, put chan []byte, err error) {
	checkJack(jack)
	get = make(chan []byte, 10)
	put = make(chan []byte, 10)
	return get, put, nil
}

func (s *Server) I2cInit(jack uint8, address uint8) (err error) {
	checkJack(jack)
	return nil
}

func (s *Server) SpiInit(jack uint8) (err error) {
	checkJack(jack)
	return nil
}

func (s *Server) PortInit(jack uint8, mode PortMode) (err error) {
	checkJack(jack)
	return nil
}

func checkJack(jack uint8) {

}
