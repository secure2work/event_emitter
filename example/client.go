package main

import (
	"github.com/bruteforce1414/event_emitter"
)

func main() {
	publisher1 := event_emitter.NewPubsub()
	channel1, funcUnsubscribe1 := publisher1.Subscribe()
	channel2, funcUnsubscribe2 := publisher1.Subscribe()
}
