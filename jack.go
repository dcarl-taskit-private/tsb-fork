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

func (s *Server) SpiInit(jack uint8) (err error) {
	checkJack(jack)
	return nil
}

func checkJack(jack uint8) {
	if jack > MaxJacks {
		log.Fatalf("Illegal Jack nr: %d", jack)
	}
}
