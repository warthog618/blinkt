// A blinkt demo that cycles the whole display through red green and blue.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/warthog618/blinkt"
)

func main() {
	bl := blinkt.New()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Press Control + C to stop")

	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan struct{})

	go func() {
		<-signalChan
		bl.Close()
		close(done)
		os.Exit(1)
	}()

	step := 0

	updateDisplay := func() {
		step = step % 3
		switch step {
		case 0:
			bl.SetAll(128, 0, 0)
		case 1:
			bl.SetAll(0, 128, 0)
		case 2:
			bl.SetAll(0, 0, 128)
		}

		step++
		bl.Show()
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
