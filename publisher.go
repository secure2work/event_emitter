package event_emitter

import (
	"log"
	"path"
	"sync"
)

type Sub struct {
	Pattern     string
	Ch          chan<- Event
	middlewares []Middleware
}

type mw struct {
	middleware Middleware
	pattern    string
}

type publisher struct {
	counter     uint64
	mu          sync.RWMutex
	subs        map[uint64]Sub
	middlewares []mw
}

type Middleware func(event *Event)

func (p *publisher) Use(pattern string, middleware ...Middleware) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, value := range middleware {
		p.middlewares = append(p.middlewares, mw{middleware: value, pattern: pattern})
	}
}

func (p *publisher) On(pattern string, middleware ...Middleware) (evt <-chan Event, unsubscribe func() error) {
	newChannel := make(chan Event)
	p.mu.Lock()
	defer p.mu.Unlock()

	p.counter = p.counter + 1
	nameChannel := p.counter

	middlewaresArr := make([]Middleware, 0)
	for _, value := range middleware {
		middlewaresArr = append(middlewaresArr, value)
	}
	p.subs[nameChannel] = Sub{
		Pattern:     pattern,
		Ch:          newChannel,
		middlewares: middlewaresArr,
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

	for _, m := range p.middlewares {

		ok, err := path.Match(m.pattern, event.Name)

		if err != nil || !ok {
			continue
		}

		m.middleware(&event)
	}
	log.Println("p.middlewares", p.middlewares)
	for _, value := range p.subs {

		go func(ch chan<- Event, patternCh string, middlewaresSub []Middleware, event2 Event) {

			ok, err := path.Match(patternCh, event.Name)

			if err != nil || !ok {
				return
			}
			for _, m := range middlewaresSub {
				m(&event)
			}
			ch <- event
		}(value.Ch, value.Pattern, value.middlewares, event)
	}
}

func NewPubsub() *publisher {
	pub := &publisher{}
	pub.mu.Lock()
	defer pub.mu.Unlock()
	pub.subs = make(map[uint64]Sub)
	pub.middlewares = make([]mw, 0)
	return pub
}
