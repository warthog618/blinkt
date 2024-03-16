// +build !blinkt_sysfs
// +build !blinkt_wiringpi
// +build !blinkt_gpio

package blinkt

import (
	"log"

	"github.com/warthog618/go-gpiocdev"
	"github.com/warthog618/go-gpiocdev/device/rpi"
)

// APA102 implements the interface to the Blinkt! hardware
type APA102 struct {
	dat *gpiocdev.Line
	clk *gpiocdev.Line
}

// Open the interface to APA102
func (a *APA102) Open() {
	c, err := gpiocdev.NewChip("gpiochip0", gpiocdev.WithConsumer("blinkt!"))
	if err != nil {
		log.Fatalf("Error opening gpiochip0: %s", err)
	}
	a.dat, err = c.RequestLine(rpi.GPIO23, gpiocdev.AsOutput(0))
	if err != nil {
		log.Fatalf("Error requesting data line: %s", err)
	}
	a.clk, err = c.RequestLine(rpi.GPIO24, gpiocdev.AsOutput(1))
	if err != nil {
		log.Fatalf("Error requesting clock line: %s", err)
	}
	c.Close()
}

// Close the interface to the APA102
func (a *APA102) Close() {
	a.dat.Reconfigure(gpiocdev.AsInput)
	a.clk.Reconfigure(gpiocdev.AsInput)
	a.dat.Close()
	a.clk.Close()
}

// WriteBit writes a single data bit to the APA102
func (a *APA102) WriteBit(bit int) {
	a.dat.SetValue(bit)
	a.clk.SetValue(0)
	a.clk.SetValue(1)
}
