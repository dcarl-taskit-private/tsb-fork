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
	GetChan    [MaxTyp + 1]chan []byte
	PutChan    [MaxTyp + 1]chan []byte
}

type Server struct {
	address string
	typ     string
	jack    [MaxJacks]jack
	conn    net.Conn
	sport   *serial.Port
	tdPutCh chan Packet
	tdGetCh chan Packet
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

func NewTcpServer(address string) (Server, error) {
	var err error
	s := Server{address: address}
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

func (s *Server) redirect(td Packet) {
	if td.Typ[0] > MaxTyp {
		log.Printf("Invalid Typ %d!\n\r", td.Typ[0])
		return
	}
	if td.Ch[0] > MaxJacks || td.Ch[0] < 1 {
		log.Printf("Invalid Jacknr %d!\n\r", td.Ch[0])
		return
	}
	if s.jack[td.Ch[0]].GetChan[td.Typ[0]] == nil {
		log.Printf("Channel: %d is not initialized!\n\r", td.Ch[0])
		return
	}
	if len(s.jack[td.Ch[0]].GetChan[td.Typ[0]]) >= cap(s.jack[td.Ch[0]].GetChan[td.Typ[0]]) {
		log.Printf("Channel Overflow! Jack: %d, Typ: %d", td.Ch[0], td.Typ[0])
	}
	s.jack[td.Ch[0]].GetChan[td.Typ[0]] <- td.Payload
}

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
