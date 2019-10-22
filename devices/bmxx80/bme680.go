package bmxx80

import (
	"fmt"
	"time"

	"periph.io/x/periph/conn/physic"
)

// GasTemp represents a temperature setting for the BME680 gas sensor's heater.
// These are set using a multiplier and a value, and the temperature setting is
// the product of these.
type GasTemp uint8

// Use these multipliers for setting temperatures for the gas sensor.
const (
	GTx1  uint8 = 0x0
	GTx4  uint8 = 0x1
	GTx16 uint8 = 0x2
	GTx64 uint8 = 0x3
)

func (g GasTempSetPoint) Int() {
}

func (g GasTempSetPoint) String() {
}

// configureHeater680 sets the heating duration and temperature for the gas
// sensor.
//
// It must be called with d.mu lock held.
func (d *Dev) configureHeater680(gasWait time.Duration, gasTemp physic.Temperature) error {
	return nil
}

// sense680 reads the device's registers for bme680
//
// It must be called with d.mu lock held.
func (d *Dev) sense680(e *physic.Env) error {
	// All registers must be read in a single pass, as noted at page 21, section
	// 4.1.
	// Pressure: 0x1F-0x20
	// Temperature: 0x22-0x24
	// Humidity: 0x25-0x26
	buf := [8]byte{}
	b := buf[:]

	if err := d.readReg(0x1F, b); err != nil {
		return err
	}

	// These values are 20 bits as per doc.
	pRaw := int32(buf[0])<<12 | int32(buf[1])<<4 | int32(buf[2])>>4
	tRaw := int32(buf[3])<<12 | int32(buf[4])<<4 | int32(buf[5])>>4

	fmt.Println("raw temperature:", tRaw)
	// TODO: Also need to read `gas_r` register here

	t, tFine := d.cal280.compensateTempInt(tRaw)
	fmt.Println("t = ", t, "tFine = ", tFine)

	// Convert CentiCelsius to Kelvin.
	e.Temperature = physic.Temperature(t)*10*physic.MilliCelsius + physic.ZeroCelsius

	if d.opts.Pressure != Off {
		p := d.cal280.compensatePressureInt64(pRaw, tFine)
		// It has 8 bits of fractional Pascal.
		e.Pressure = physic.Pressure(p) * 15625 * physic.MicroPascal / 4
	}

	if d.opts.Humidity != Off {
		// This value is 16 bits as per doc.
		hRaw := int32(buf[6])<<8 | int32(buf[7])
		h := physic.RelativeHumidity(d.cal280.compensateHumidityInt(hRaw, tFine))
		// Convert base 1024 to base 1000.
		e.Humidity = h * 10000 / 1024 * physic.MicroRH
	}
	return nil
}
