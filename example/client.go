package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/secure2work/event_emitter"
)

const (
	eventTime string = "event.time"
)

func main() {

	publisher := event_emitter.NewPubsub()

	channel1, unsubscribe1 := publisher.On("nori/plugins/*")
	channel2, unsubscribe2 := publisher.On("nori/plugins/started")

	listen(channel1, "ch1", 5)
	listen(channel2, "ch2", 5)

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				publisher.Emit(event_emitter.Event{
					Name: "nori/plugins/stopped",
					Params: map[string]interface{}{
						"time": t,
					},
				})
			}
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	if err := unsubscribe1(); err != nil {
		log.Println(err.Error())
	}
	if err := unsubscribe2(); err != nil {
		log.Println(err.Error())
	}

	fmt.Println("\r- Ctrl+C pressed in Terminal")
}

func listen(ch <-chan event_emitter.Event, chName string, div int) {
	go func() {
		for {
			e1 := <-ch
			log.Println(chName, " receive event: ", e1.Name)
		}
	}()
}
