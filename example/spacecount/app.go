package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/warthog618/blinkt"
)

func displayAstronauts(bl *blinkt.Blinkt) {
	num := getAstronautCount()

	r := 150
	g := 0
	b := 0
	bl.Clear()
	for pixel := 0; pixel < num; pixel++ {
		bl.SetPixel(pixel, r, g, b)
		bl.Show()
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {

	bl := blinkt.New()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Press Control + C to stop")

	ticker := time.NewTicker(60 * time.Second)
	done := make(chan struct{})

	go func() {
		<-signalChan
		bl.Close()
		close(done)
		os.Exit(1)
	}()

	displayAstronauts(&bl)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				displayAstronauts(&bl)
			}
		}
	}()
	<-done
}
