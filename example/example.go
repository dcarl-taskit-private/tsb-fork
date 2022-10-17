package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/traulfs/tsb"
)

func main() {
	//serv, err := tsb.NewSerialServer("/dev/ttyUSB0")
	serv, err := tsb.NewTcpServer("localhost:3001")
	//serv, err := tsb.NewTcpServer("10.1.108.197:3001")
	if err != nil {
		log.Fatal(err)
	}
	go uartExample(serv, 1)
	go uartExample(serv, 2)
	go portExample(serv, 3)
	go portExample(serv, 4)
	go i2cExample(serv, 5)
	go i2cExample(serv, 6)
	for {

	}
}

func uartExample(s tsb.Server, jack byte) {
	var buf []byte = make([]byte, 256)
	err := s.UartInit(jack, tsb.UartBaud115200, tsb.UartData8&tsb.UartParityNone&tsb.UartStopbits1)
	if err != nil {
		log.Fatal(err)
	}
	go func(jack byte) {
		for {
			_, err := s.UartWrite(jack, []byte("Hello Jack"+strconv.Itoa(int(jack))+"\n"))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(5 * time.Duration(time.Second))
		}
	}(jack)
	go func(jack byte) {
		for {
			//fmt.Printf("%d", jack)
			n, err := s.UartRead(jack, buf)
			if err != nil {
				log.Fatal(err)
			}
			if n > 0 {
				fmt.Printf("Received from Jack %d: %s\n", jack, buf)
			}
		}
	}(jack)
}

func portExample(s tsb.Server, jack byte) {
	err := s.PortInit(jack)
	if err != nil {
		log.Fatal(err)
	}
	go func(jack byte) {
		for {
			s.PortPutc(jack, 0x31)
			time.Sleep(time.Duration(time.Second))
		}
	}(jack)
}

func i2cExample(s tsb.Server, jack byte) {
}
