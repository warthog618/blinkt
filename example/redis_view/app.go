// A blinkt demo that sets the display using values provided by a Redis pubsub
// channel.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/warthog618/blinkt"
	"gopkg.in/redis.v5"
)

type ledColor struct {
	Red   int `json:"r"`
	Green int `json:"g"`
	Blue  int `json:"b"`
}

type ledMsg struct {
	Leds []ledColor `json:"leds"`
}

func newClient() *redis.Client {
	addr := os.Getenv("ADDR")
	client := redis.NewClient(&redis.Options{
		Addr:     addr + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Output: PONG <nil>
	return client
}

func main() {
	bl := blinkt.New()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Press Control + C to stop")

	go func() {
		<-signalChan
		bl.Close()
		os.Exit(1)
	}()

	client := newClient()
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	pubsub, err := client.Subscribe("lights")
	if err != nil {
		panic(err)
	}
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}
		if msg.Channel == "lights" {
			ledMsg := ledMsg{}
			jsonErr := json.Unmarshal([]byte(msg.Payload), &ledMsg)
			if jsonErr != nil {
				fmt.Println(jsonErr)
			}

			bl.Clear()

			var r, g, b int
			fmt.Println("Setting LEDs")

			for pixel := 0; pixel < 8; pixel++ {
				if pixel > len(ledMsg.Leds) {
					r = 0
					g = 0
					b = 0
				} else {
					r = ledMsg.Leds[pixel].Red
					g = ledMsg.Leds[pixel].Green
					b = ledMsg.Leds[pixel].Blue
				}
				bl.SetPixel(pixel, r, g, b)
			}
			bl.Show()
		}
	}
}
