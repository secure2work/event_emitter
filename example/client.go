package main

import (
	"fmt"
	"strconv"
	"sync"
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
	fmt.Println("e1 is", e1)
	fmt.Println("e2 is", e2)

	funcUnsubscribe1()
	funcUnsubscribe2()
}
