// A blinkt demo that cycles a rainbow of colours.
package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/warthog618/blinkt"
)

// the display covers half the colour wheel at any time...
var spacing = 360.0 / 16

func hueToIntensity(hue float64) int {
	h := math.Mod(hue, 360)
	switch {
	case h <= 120:
		return 0
	case h < 180:
		return int(math.Round(((h - 120) * 255) / 60))
	case h <= 300:
		return 255
	default:
		return int(math.Round(((360 - h) * 255) / 60))
	}
}

func hueToRGB(h float64) (int, int, int) {
	r := hueToIntensity(h)
	g := hueToIntensity(h + 120)
	b := hueToIntensity(h + 240)
	return r, g, b
}

func display(bl *blinkt.Blinkt, hue float64) {
	bl.Clear()
	for i := 0; i < bl.Pixels; i++ {
		r, g, b := hueToRGB(hue)
		bl.SetPixel(i, r, g, b)
		hue += spacing
	}
	bl.Show()
}

func main() {
	bl := blinkt.New()
	bl.SetBrightness(10)

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

	go func() {
		h := 0.0
		display(&bl, h)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				h += spacing
				if h > 360 {
					h -= 360
				}
				display(&bl, h)
			}
		}
	}()
	<-done
}
