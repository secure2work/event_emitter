package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bruteforce1414/event_emitter"
)

const (
	eventTime string = "event.time"
)

func main() {
	publisher := event_emitter.NewPubsub()
	channel1, unsubscribe1 := publisher.On("nori/plugins/*")
	channel2, unsubscribe2 := publisher.On("nori/plugins/*")

	listen(channel1, "ch1", 5)
	listen(channel2, "ch2", 7)

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				publisher.Emit(event_emitter.Event{
					Name: eventTime,
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
			switch e1.Name {
			case eventTime:
				t, ok := e1.Params["time"].(time.Time)
				if !ok {
					log.Println(chName, ": no time provided")
				}
				if (t.Second() % div) == 0 {
					log.Println(chName, " log number: ", t.Second())
				}
			default:
				continue
			}
		}
	}()
}
