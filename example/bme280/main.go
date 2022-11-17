package main

import (
	"fmt"
	"log"
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
	temp := make([]byte, 3)
	serv.I2cInit(1)
	fmt.Printf("BME280 Example\n")
	for i := 0; i < 10; i++ {
		serv.I2cWrite(1, 0xF4, []byte{0x03})
		_, err = serv.I2cRead(1, 0xFA, temp)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d: Temp: %v\n", i, temp)
		time.Sleep(time.Second)
	}
}
