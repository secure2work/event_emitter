package event_emitter

import (
	"fmt"
	"sync"
)

type publisher struct {
	counter uint64
	mu      sync.RWMutex
	subs    map[uint64]chan<- Event
}

func (p *publisher) Subscribe() (evt <-chan Event, unsubscribe func()) {
	newChannel := make(chan Event)
	p.mu.Lock()
	p.counter = p.counter + 1
	nameChannel := p.counter
	p.mu.Unlock()
	p.subs[nameChannel] = newChannel

	c := func() {
		_, ok := p.subs[nameChannel]
		if ok {
			close(newChannel)
			p.mu.Lock()
			delete(p.subs, nameChannel)
			p.mu.Unlock()
			fmt.Println("nameChannel is", nameChannel)
			fmt.Println("unsubscribe")
		} else {
			println("there's not channel with id ", nameChannel)
		}
	}

	return newChannel, c
}

func (p *publisher) Emit(event Event) {
	p.mu.Lock()

	for _, value := range p.subs {

		go func(evt chan<- Event) {
			value <- event
		}(value)
		fmt.Println("channel is", value)
		fmt.Println("event is", event)
	}
	p.mu.Unlock()
}

func NewPubsub() *publisher {
	pub := &publisher{}
	pub.mu.Lock()
	pub.subs = make(map[uint64]chan<- Event)

	pub.mu.Unlock()
	return pub
}
