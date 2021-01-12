// +build blinkt_wiringpi

package blinkt

import "github.com/alexellis/rpi"

const DAT int = 23
const CLK int = 24

// APA102 implenments the interface to the Blinkt! hardware
type APA102 struct {
}

// Open the interface to APA102
func (a *APA102) Open() {
	rpi.WiringPiSetup()
	rpi.PinMode(rpi.GpioToPin(DAT), rpi.OUTPUT)
	rpi.PinMode(rpi.GpioToPin(CLK), rpi.OUTPUT)
}

// Close the interface to the APA102
func (a *APA102) Close() {
	rpi.PinMode(rpi.GpioToPin(DAT), rpi.INPUT)
	rpi.PinMode(rpi.GpioToPin(CLK), rpi.INPUT)
}

func (a *APA102) WriteBit(bit int) {
	rpi.DigitalWrite(rpi.GpioToPin(DAT), bit)
	rpi.DigitalWrite(rpi.GpioToPin(CLK), 0)
	rpi.DigitalWrite(rpi.GpioToPin(CLK), 1)
}
