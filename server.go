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

type jack struct {
	ReadChan [MaxTyp + 1]chan byte
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
	for i := 0; i < int(MaxJacks); i++ {
		s.jack[i].ReadChan[TypI2c] = make(chan byte, 1024)
	}
	fmt.Printf("TSB client connected to tsb server: %s\n", s.address)
	go func() {
		for {
			select {
			case td := <-s.tdGetCh:
				{
					if td.Typ[0] > MaxTyp {
						log.Printf("Invalid Typ %d!\n\r", td.Typ[0])
						return
					}
					if td.Ch[0] > byte(MaxJacks) || td.Ch[0] < 1 {
						log.Printf("Invalid Jacknr %d!\n\r", td.Ch[0])
						return
					}
					if s.jack[td.Ch[0]].ReadChan[td.Typ[0]] == nil {
						log.Printf("Channel: %d is not initialized!\n\r", td.Ch[0])
						return
					}
					if len(s.jack[td.Ch[0]].ReadChan[td.Typ[0]]) >= cap(s.jack[td.Ch[0]].ReadChan[td.Typ[0]]) {
						log.Printf("Channel Overflow! Jack: %d, Typ: %d", td.Ch[0], td.Typ[0])
					}
					for i := range td.Payload {
						s.jack[td.Ch[0]].ReadChan[td.Typ[0]] <- td.Payload[i]
					}
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

func (s *Server) SpiInit(jack byte) (err error) {
	CheckJack(jack)
	return nil
}

func CheckJack(jack byte) {
	if jack > MaxJacks {
		log.Fatalf("Illegal Jack nr: %d", jack)
	}
}