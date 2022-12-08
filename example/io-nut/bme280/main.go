package main

import (
	"fmt"
	"os"

	"github.com/traulfs/tsb"
)

var verbose bool

const MyJack byte = 5 // select Jack 1-8

func main() {
	server, err := tsb.NewTcpServer("localhost:3010")
	//server, err := tsb.NewSerialServer("/dev/tty.usbmodem11401")
	if err != nil {
		Fatal(err)
	}
	// Create new connection to i2c-bus on 1 line with address 0x76.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := tsb.NewI2c(0x76, MyJack, server)
	if err != nil {
		Fatal(err)
	}
	defer i2c.Close()

	//Notify("***************************************************************************************************")

	// sensor, err := NewBMP(BMP180, i2c) // signature=0x55
	// sensor, err := NewBMP(BMP280, i2c) // signature=0x58
	sensor, err := NewBMP(BME280, i2c) // signature=0x60
	// sensor, err := NewBMP(BMP388, i2c) // signature=0x50
	if err != nil {
		Fatal(err)
	}

	id, err := sensor.ReadSensorID()
	if err != nil {
		Fatal(err)
	}
	Infof("This Bosch Sensortec sensor has signature: 0x%x", id)

	err = sensor.IsValidCoefficients()
	if err != nil {
		Fatal(err)
	}

	// Read temperature in celsius degree
	t, err := sensor.ReadTemperatureC(ACCURACY_STANDARD)
	if err != nil {
		Fatal(err)
	}
	Infof("Temperature = %v°C", t)

	// Read atmospheric pressure in pascal
	p, err := sensor.ReadPressurePa(ACCURACY_LOW)
	if err != nil {
		Fatal(err)
	}
	Infof("Pressure = %v Pa", p)

	// Read atmospheric pressure in mmHg
	p, err = sensor.ReadPressureMmHg(ACCURACY_LOW)
	if err != nil {
		Fatal(err)
	}
	Infof("Pressure = %v mmHg", p)

	// Read atmospheric pressure in mmHg
	supported, h1, err := sensor.ReadHumidityRH(ACCURACY_LOW)
	if supported {
		if err != nil {
			Fatal(err)
		}
		Infof("Humidity = %v %%", h1)
	}

	// Read atmospheric altitude in meters above sea level, if we assume
	// that pressure at see level is equal to 101325 Pa.
	a, err := sensor.ReadAltitude(ACCURACY_LOW)
	if err != nil {
		Fatal(err)
	}
	Infof("Altitude = %v m", a)
}

func Debugf(format string, values ...interface{}) {
	if verbose {
		fmt.Printf("Debug: "+format+"\r\n", values...)
	}
}

func Fatal(error error) {
	fmt.Printf("Fatal: %s\r\n", error)
	os.Exit(1)
}

func Notify(message string) {
	fmt.Printf("Notify: %s\r\n", message)
}

func Infof(format string, values ...interface{}) {
	fmt.Printf("Info: "+format+"\r\n", values...)
}
