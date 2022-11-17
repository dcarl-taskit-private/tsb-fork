package main

import (
	"fmt"
	"log"
	"time"

	"github.com/traulfs/tsb"
)

const MyJack byte = 5

func main() {
	//serv, err := tsb.NewSerialServer("/dev/ttyUSB0")
	serv, err := tsb.NewTcpServer("localhost:3001")
	//serv, err := tsb.NewTcpServer("10.1.108.197:3001")
	if err != nil {
		log.Fatal(err)
	}
	temp := make([]byte, 3)
	serv.I2cInit(MyJack)
	fmt.Printf("BME280 Example\n")
	serv.I2cSetAdr(MyJack, 0x76)
	serv.I2cWrite(MyJack, []byte{0xF4, 0x03})
	for i := 1; i <= 10; i++ {
		serv.I2cWrite(MyJack, []byte{0xFA})
		_, err = serv.I2cRead(MyJack, temp)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%2d: Temp: %v\n", i, temp)
		time.Sleep(time.Second)
	}
}
