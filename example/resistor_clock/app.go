// A blinkt demo that displays the currrent time using resistor band colours.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/warthog618/blinkt"
)

var colours = [10][3]int{
	{0, 0, 0},       // 0 black
	{139, 69, 19},   // 1 brown
	{255, 0, 0},     // 2 red
	{255, 69, 0},    // 3 orange
	{255, 255, 0},   // 4 yellow
	{0, 255, 0},     // 5 green
	{0, 0, 255},     // 6 blue
	{128, 0, 128},   // 7 violet
	{255, 255, 100}, // 8 grey
	{255, 255, 255}, // 9 white
}

func setPixelColour(bl *blinkt.Blinkt, l, c int) {
	bl.SetPixel(l, colours[c][0], colours[c][1], colours[c][2])
}

func displayTime(bl *blinkt.Blinkt, t time.Time) {
	hour := t.Hour()
	minute := t.Minute()

	hourten := hour / 10
	hourunit := hour % 10
	minuteten := minute / 10
	minuteunit := minute % 10

	setPixelColour(bl, 0, hourten)
	setPixelColour(bl, 1, hourten)
	setPixelColour(bl, 2, hourunit)
	setPixelColour(bl, 3, hourunit)
	setPixelColour(bl, 4, minuteten)
	setPixelColour(bl, 5, minuteten)
	setPixelColour(bl, 6, minuteunit)
	setPixelColour(bl, 7, minuteunit)

	bl.Show()
}

func tick(bl *blinkt.Blinkt, t time.Time) {
	minuteunit := t.Minute() % 10
	if minuteunit == 0 {
		setPixelColour(bl, 7, 8)
	} else {
		setPixelColour(bl, 7, 0)
	}
	bl.Show()
}

func main() {
	bl := blinkt.New()
	bl.SetBrightness(15)

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

	displayTime(&bl, time.Now())
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				if t.Nanosecond() > 500000000 {
					displayTime(&bl, t)
				} else {
					tick(&bl, t)
				}
			}
		}
	}()
	<-done
}
