// A blinkt demo that displays the currrent time in binary.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/warthog618/blinkt"
)

type block struct {
	size     int
	position int
	r        int
	g        int
	b        int
}

var grid = [4][16]int{
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0},
}

func displayBlock(bl *blinkt.Blinkt, b block) {
	row := grid[b.size]
	bl.Clear()
	switch {
	case b.position < 8:
		for i := 0; i < 8; i++ {
			if row[b.position+i] != 0 {
				bl.SetPixel(i, b.r, b.g, b.b)
			}
		}
	case b.position == 8:
		// hide
	case b.position == 9:
		for i := 0; i < 8; i++ {
			if row[7+i] != 0 {
				bl.SetPixel(i, b.r, b.g, b.b)
			}
		}
	case b.position > 9:
		// hide
	}
	bl.Show()
}

func newBlock() block {
	return block{
		size:     int(rand.Float64() * 4),
		position: 0,
		r:        int(rand.Float64() * 255),
		g:        int(rand.Float64() * 255),
		b:        int(rand.Float64() * 50),
	}
}

func main() {
	bl := blinkt.New()
	bl.SetBrightness(10)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Press Control + C to stop")

	ticker := time.NewTicker(250 * time.Millisecond)
	done := make(chan struct{})

	go func() {
		<-signalChan
		bl.Close()
		close(done)
		os.Exit(1)
	}()

	rand.Seed(time.Now().UnixNano())

	b := newBlock()
	displayBlock(&bl, b)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				b.position++
				if b.position == 13 {
					b = newBlock()
				}
				displayBlock(&bl, b)
			}
		}
	}()
	<-done
}
