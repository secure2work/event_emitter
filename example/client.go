package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nori-io/common/v3/pkg/domain/event"

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

	type TestEventStruct struct {
		TimeField time.Time
	}

	event1 := event.Event{
		Name: "nori/plugins/stopped",
	}

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				event1.Params = TestEventStruct{TimeField: t}
				publisher.Emit(event1.Name, event1.Params)
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

func listen(ch <-chan event.Event, chName string, div int) {
	go func() {
		for {
			e1 := <-ch
			log.Println(chName, " receive event: ", e1.Name)
		}
	}()
}
