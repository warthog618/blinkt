// +build blinkt_gpio

package blinkt

import "github.com/warthog618/gpio"

// APA102 implenments the interface to the Blinkt! hardware
type APA102 struct {
	dat *gpio.Pin
	clk *gpio.Pin
}

// Open the interface to APA102
func (a *APA102) Open() {
	gpio.Open()
	a.dat = gpio.NewPin(gpio.GPIO23)
	a.dat.Low()
	a.dat.Output()
	a.clk = gpio.NewPin(gpio.GPIO24)
	a.clk.High()
	a.clk.Output()
}

// Close the interface to the APA102
func (a *APA102) Close() {
	a.dat.Input()
	a.clk.Input()
	gpio.Close()
}

func (a *APA102) WriteBit(val int) {
	var bit gpio.Level
	if val == 0 {
		bit = gpio.Low
	} else {
		bit = gpio.High
	}
	a.dat.Write(bit)
	a.clk.Low()
	a.clk.High()
}
