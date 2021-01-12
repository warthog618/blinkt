// A blinkt demo that continually sets the leds to random colours.
package main

import (
	"fmt"
	"math/rand"
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

	ticker := time.NewTicker(50 * time.Millisecond)
	done := make(chan struct{})

	go func() {
		<-signalChan
		bl.Close()
		close(done)
		os.Exit(1)
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for pixel := 0; pixel < 8; pixel++ {
					bl.SetPixel(pixel, rand.Intn(255), rand.Intn(255), rand.Intn(255))
				}
				bl.Show()
			}
		}
	}()
	<-done
}
