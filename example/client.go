package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/bruteforce1414/event_emitter"
)

func main() {
	publisher1 := event_emitter.NewPubsub()
	channel1, funcUnsubscribe1 := publisher1.Subscribe()
	channel2, funcUnsubscribe2 := publisher1.Subscribe()

	wg := sync.WaitGroup{}
	wg.Add(1)

	var e1, e2 event_emitter.Event

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		time.Sleep(5000)
		os.Exit(0)
	}()

	go func() {
		for {
			wg.Add(1)
			e1 = <-channel1
			i, _ := strconv.Atoi(e1.Name)
			if (i % 2) == 0 {
				fmt.Println("1 channel log number", e1.Name)
			}
			fmt.Println(e1.Name)
			wg.Done()
		}
	}()

	go func() {
		for {
			wg.Add(1)
			e2 = <-channel2
			i, _ := strconv.Atoi(e2.Name)
			if (i % 3) == 0 {
				fmt.Println("2 channel log number", e2.Name)
			}
			wg.Done()

		}
	}()

	go func() {
		for {
			event1 := event_emitter.Event{
				Name:   fmt.Sprint(time.Now().Second()),
				Params: make(map[string]interface{}),
			}
			publisher1.Emit(event1)
			time.Sleep(2 * time.Second)
		}
	}()

	wg.Wait()

	funcUnsubscribe1()
	funcUnsubscribe2()
}
