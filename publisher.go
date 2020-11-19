package event_emitter

import (
	"sync"
)

type publisher struct {
	counter uint64
	mu      sync.RWMutex
	subs    map[uint64]chan<- Event
}

func (p *publisher) Subscribe() (evt <-chan Event, unsubscribe func() error) {
	newChannel := make(chan Event)
	p.mu.Lock()
	defer p.mu.Unlock()

	p.counter = p.counter + 1
	nameChannel := p.counter
	p.subs[nameChannel] = newChannel

	c := func() error {
		_, ok := p.subs[nameChannel]
		if ok {
			p.mu.Lock()
			defer p.mu.Unlock()
			close(newChannel)
			delete(p.subs, nameChannel)
			return nil
		} else {
			return nil
		}
	}

	return newChannel, c
}

func (p *publisher) Emit(event Event) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, value := range p.subs {
		go func(ch chan<- Event) {
			ch <- event
		}(value)
	}
}

func NewPubsub() *publisher {
	pub := &publisher{}
	pub.mu.Lock()
	defer pub.mu.Unlock()
	pub.subs = make(map[uint64]chan<- Event)
	return pub
}
