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
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()
	bl.Clear()
	if second&0x1 == 0 {
		bl.SetPixel(0, 64, 64, 64)
	}

	for i := 0; i < 6; i++ {
		ch := 0
		cm := 0
		cs := 0
		mask := 1 << i
		if hour&mask != 0 {
			ch = 64
		}
		if minute&mask != 0 {
			cm = 64
		}
		if second&mask != 0 {
			cs = 64
		}
		bl.SetPixel(7-i, ch, cm, cs)
	}
	bl.Show()
}

func main() {
	bl := blinkt.New()
	bl.SetBrightness(20)

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
