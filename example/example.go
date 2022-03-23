package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gpio-net/tsb"
)

func main() {
	serv, err := tsb.NewTcpServer("loaclhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	go serv.uartExample(1)
	go serv.uartExample(2)
	go serv.portExample(3)
	go serv.portExample(4)
	go serv.i2cExample(5)
}

func (s *tsb.server) uartExample(jack int) {
	GetChan, PutChan, err := s.UartInit(1, 115200, "8N1")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		PutChan <- "Hello Chan" + strconv.Itoa(jack)
		time.Sleep(time.Duration(jack) * time.Second)
	}()
	for {
		fmt.Printf("Received from Jack%d: %s\n\r", jack, <-GetChan)
	}
}
