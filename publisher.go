package event_emitter

import (
	"sync"

	e "github.com/nori-io/common/v3/pkg/domain/enum/event"
)

type Sub struct {
	Pattern string
	Ch      chan<- Event
}

type publisher struct {
	counter uint64
	mu      sync.RWMutex
	subs    map[uint64]Sub
}

type Middleware func(event *Event)

func (p *publisher) Use(pattern string, middleware ...Middleware) {

}

func (p *publisher) On(pattern string, middleware ...Middleware) (evt <-chan Event, unsubscribe func() error) {
	newChannel := make(chan Event)
	p.mu.Lock()
	defer p.mu.Unlock()

	p.counter = p.counter + 1
	nameChannel := p.counter
	p.subs[nameChannel] = Sub{
		Pattern: pattern,
		Ch:      newChannel,
	}

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
			if (event.Name == value.Pattern) || (value.Pattern == string(e.NoriPlugins)) {
				ch <- event
			}

		}(value.Ch)
	}
}

func NewPubsub() *publisher {
	pub := &publisher{}
	pub.mu.Lock()
	defer pub.mu.Unlock()
	pub.subs = make(map[uint64]Sub)
	return pub
}
