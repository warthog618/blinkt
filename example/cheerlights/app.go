// A blinkt demo that displays the current CheerLights colour.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/warthog618/blinkt"
)

func getEnv(envVar string, assumed int) int {
	if value, exists := os.LookupEnv(envVar); exists {
		if period, err := strconv.Atoi(value); err == nil {
			return period
		}
	}
	return assumed
}

const defaultRefreshSeconds = 60
const envRefreshSeconds = "refresh_seconds"

func main() {
	bl := blinkt.New()
	bl.SetBrightness(15)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Press Control + C to stop")

	checkPeriod := time.Duration(getEnv(envRefreshSeconds, defaultRefreshSeconds))
	ticker := time.NewTicker(checkPeriod * time.Second)
	done := make(chan struct{})

	go func() {
		<-signalChan
		bl.Close()
		close(done)
		os.Exit(1)
	}()

	updateDisplay := func() {
		r, g, b := getCheerlightColours()

		bl.Clear()
		bl.SetAll(r, g, b)
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
