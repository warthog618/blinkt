// A blinkt demo that displays the currrent time in binary.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/warthog618/blinkt"
)

var mode int
var modeTicks int

func displayTime(bl *blinkt.Blinkt, t time.Time) {
	second := t.Second()

	modeTicks++
	if modeTicks > 3 {
		modeTicks = 0
		mode++
		if mode > 2 {
			mode = 0
		}
	}
	bl.Clear()
	val := 0
	switch mode {
	case 0: // hour
		bl.SetPixel(0, 255, 0, 0)
		val = t.Hour()
	case 1: // minute
		bl.SetPixel(0, 0, 255, 0)
		val = t.Minute()
	case 2: // second
		bl.SetPixel(0, 0, 0, 255)
		val = second
	}
	if second&0x1 == 0 {
		bl.SetPixel(1, 64, 64, 64)
	}
	for i := 0; i < 6; i++ {
		if (val & (0x1 << i)) != 0 {
			bl.SetPixel(7-i, 255, 255, 255)
		}
	}
	bl.Show()
}

func main() {
	bl := blinkt.New()
	bl.SetBrightness(15)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Press Control + C to stop")

	ticker := time.NewTicker(time.Second)
	done := make(chan struct{})

	go func() {
		<-signalChan
		bl.Close()
		close(done)
		os.Exit(1)
	}()

	displayTime(&bl, time.Now())
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				displayTime(&bl, t)
			}
		}
	}()
	<-done
}
