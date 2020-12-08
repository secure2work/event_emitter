package event_emitter

import (
	"path"
	"sync"

	"github.com/nori-io/common/v3/pkg/domain/event"
)

type Sub struct {
	Pattern     string
	Ch          chan<- event.Event
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

type Middleware func(event *event.Event)

func (p *publisher) Use(pattern string, middleware ...Middleware) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, value := range middleware {
		p.middlewares = append(p.middlewares, mw{middleware: value, pattern: pattern})
	}
}

func (p *publisher) On(pattern string, middleware ...Middleware) (evt <-chan event.Event, unsubscribe func() error) {
	newChannel := make(chan event.Event)
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
		}
		return nil

	}

	return newChannel, c
}

func (p *publisher) Emit(name string, params interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	evt := event.Event{
		Name:   name,
		Params: params,
	}

	for _, m := range p.middlewares {

		ok, err := path.Match(m.pattern, name)

		if err != nil || !ok {
			continue
		}

		m.middleware(&evt)
	}
	for _, value := range p.subs {
		go func(ch chan<- event.Event, patternCh string, middlewares []Middleware, evt event.Event) {

			ok, err := path.Match(patternCh, evt.Name)

			if err != nil || !ok {
				return
			}
			for _, m := range middlewares {
				m(&evt)
			}
			ch <- evt
		}(value.Ch, value.Pattern, value.middlewares, evt)
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
