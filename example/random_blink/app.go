// A blinkt demo that continually sets the led to random colours.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/warthog618/blinkt"
)

func shuffleAndSlice(arr []int) []int {

	t := time.Now()
	rand.Seed(int64(t.Nanosecond()))

	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i)
		arr[i], arr[j] = arr[j], arr[i]
	}

	subsetSize := rand.Intn(5) + 1 // +1 as zero based
	return arr[:subsetSize]
}

func isIn(s *[]int, e *int) bool {
	for _, a := range *s {
		if a == *e {
			return true
		}
	}
	return false
}

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

	nums := []int{0, 1, 2, 3, 4, 5, 6, 7}

	updateDisplay := func() {
		//There must be a more elegant way of doing this
		pixels := shuffleAndSlice(nums)
		for _, i := range nums {
			if isIn(&pixels, &i) {
				bl.SetPixel(i, 255, 150, 0)
			} else {
				bl.SetPixel(i, 0, 0, 0)
			}
		}
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
