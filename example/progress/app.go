// A blinkt demo that cycles a red led backwards and forwards across the
// display, Cylon style.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/warthog618/blinkt"
)

var pixel int
var countDown bool

func nextPixel() {
	switch countDown {
	case true:
		pixel--
		if pixel == 0 {
			countDown = false
		}
	default:
		pixel++
		if pixel >= 7 {
			countDown = true
		}
	}
}

func main() {
	bl := blinkt.New()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Press Control + C to stop")

	ticker := time.NewTicker(100 * time.Millisecond)
	done := make(chan struct{})

	go func() {
		<-signalChan
		bl.Close()
		close(done)
		os.Exit(1)
	}()

	updateDisplay := func() {
		bl.SetPixel(pixel, 150, 0, 0)
		bl.Show()
		bl.ClearPixel(pixel)
		nextPixel()
	}
	updateDisplay()
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				updateDisplay()
			}
		}
	}()
	<-done
}
