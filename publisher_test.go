package event_emitter

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublisher_Subscribe(t *testing.T) {
	pub := NewPubsub()
	ch1, funcUnsubscribe1 := pub.Subscribe()
	ch2, funcUnsubscribe2 := pub.Subscribe()

	wg := sync.WaitGroup{}
	wg.Add(2)

	var e1, e2 Event
	go func() {
		select {
		case e1 = <-ch1:
			wg.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			wg.Done()
		}
	}()

	eventTest := Event{
		Name:   "all plugins inited",
		Params: nil,
	}

	pub.Emit(eventTest)

	wg.Wait()

	assert.Equal(t, e1, e2)

	funcUnsubscribe1()
	funcUnsubscribe2()
}
