// A blinkt demo that displays the CPU temperature.
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/warthog618/blinkt"
)

func getTemperature() float64 {
	targetCmd := exec.Command("vcgencmd", "measure_temp")
	var out bytes.Buffer
	targetCmd.Stdout = &out
	err := targetCmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	temp := out.String()

	// temp=35.8'C
	tempVal := temp[strings.Index(temp, "=")+1 : len(temp)-3]
	celcius, _ := strconv.ParseFloat(tempVal, 64)
	return celcius
}

func min(x float64, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func showGraph(bl *blinkt.Blinkt, v float64, r int, g int, b int) {
	v = v * 8
	one := float64(1.0)

	for i := 0; i < 8; i++ {
		if v < 0 {
			r, g, b = 0, 0, 0
		} else {
			r = int(float64(r) * min(v, one))
			g = int(float64(g) * min(v, one))
			b = int(float64(b) * min(v, one))
		}
		bl.SetPixel(i, r, g, b)
		v = v - 1
	}
	bl.Show()
}

func displayTemperature(bl *blinkt.Blinkt) {
	celcius := getTemperature()
	fmt.Printf("Temperature: %2.2f\n", celcius)

	v := celcius / 100
	showGraph(bl, v, 255, 255, 255)
}

func main() {
	bl := blinkt.New()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Press Control + C to stop")

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan struct{})

	go func() {
		<-signalChan
		bl.Close()
		close(done)
		os.Exit(1)
	}()

	displayTemperature(&bl)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				displayTemperature(&bl)
			}
		}
	}()
	<-done
}
