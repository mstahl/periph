package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
	"periph.io/x/periph/host"
)

const bme680Address uint16 = 0x77

func main() {
	fmt.Println("Testing a BME680 chip")

	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal("Error initializing", err)
	}

	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal("Error opening I2C", err)
	}
	defer bus.Close()

	// Open a handle to a bme280/bmp280 connected on the I²C bus using default
	// settings:
	dev, err := bmxx80.NewI2C(bus, bme680Address, &bmxx80.Opts{
		Temperature: bmxx80.O8x,
		Humidity:    bmxx80.O1x,
		Pressure:    bmxx80.O1x,
	})
	if err != nil {
		log.Fatal("Error initializing bmxx80 driver", err)
	}
	defer dev.Halt()

	var env physic.Env

	ticker := time.NewTicker(time.Second)

	for {
		if err := dev.Sense(&env); err != nil {
			panic(err)
		}
		fmt.Println(env)

		<-ticker.C
	}

	fmt.Println("Woot")
}
