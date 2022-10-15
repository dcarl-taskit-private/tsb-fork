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
	if err != nil {
		log.Fatal(err)
	}
	go uartExample(serv, 1)
	go uartExample(serv, 2)
	go portExample(serv, 3)
	go portExample(serv, 4)
	go i2cExample(serv, 5)
	for {

	}
}

func uartExample(s tsb.Server, jack byte) {
	GetChan, PutChan, err := s.UartInit(jack, tsb.UartBaud115200, tsb.UartData8&tsb.UartParityNone&tsb.UartStopbits1)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			PutChan <- []byte("Hello Jack" + strconv.Itoa(int(jack)))
			time.Sleep(time.Duration(jack) * time.Second)
		}
	}()
	for {
		msg := <-GetChan
		fmt.Printf("Received from Jack %d: %s\n\r", jack, msg)
	}
}

func portExample(s tsb.Server, jack byte) {
}

func i2cExample(s tsb.Server, jack byte) {
}
