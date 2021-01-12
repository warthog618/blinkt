// +build blinkt_sysfs

package blinkt

import "github.com/alexellis/blinkt_go/sysfs/gpio"

// APA102 implenments the interface to the Blinkt! hardware
type APA102 struct {
}

const DAT int = 23
const CLK int = 24

// Open the interface to APA102
func (a *APA102) Open() {
	gpio.Setup()
	gpio.PinMode(gpio.GpioToPin(DAT), gpio.OUTPUT)
	gpio.PinMode(gpio.GpioToPin(CLK), gpio.OUTPUT)
}

// Close the interface to the APA102
func (a *APA102) Close() {
	gpio.Cleanup()
}

func (a *APA102) WriteBit(bit int) {
	gpio.DigitalWrite(gpio.GpioToPin(DAT), bit)
	gpio.DigitalWrite(gpio.GpioToPin(CLK), 0)
	gpio.DigitalWrite(gpio.GpioToPin(CLK), 1)
}
