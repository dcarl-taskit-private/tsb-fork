package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/traulfs/tsb"
)

const MyJack byte = 5 // select Jack 1-8

func main() {
	var buf []byte = make([]byte, 256)
	//serv, err := tsb.NewSerialServer("/dev/ttyUSB0")
	s, err := tsb.NewTcpServer("localhost:3001")
	if err != nil {
		log.Fatal(err)
	}
	err = s.UartInit(MyJack, tsb.UartBaud115200, tsb.UartData8, tsb.UartParityNone, tsb.UartStopbits1)
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
	}(MyJack)
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
	}(MyJack)
}
