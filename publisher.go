package event_emitter

import (
	"fmt"
	"sync"
)

type publisher struct {
	counter uint64
	mu   sync.RWMutex
	subs map[uint64]chan<- Event
}

//unsubscribe func(p, nameChannel string)
//Убрать содержимое скобок или нет?
func (p *publisher) Subscribe() (evt <-chan Event, unsubscribe func()) {
	newChannel := make(chan Event)
	p.mu.Lock()
	p.counter=p.counter+1
	nameChannel := p.counter
	p.mu.Unlock()
	p.subs[nameChannel] = newChannel

	c := func() {
		_, ok := p.subs[nameChannel];
		if ok {
			close(newChannel)
			p.mu.Lock()
			delete(p.subs, nameChannel);
			p.mu.Unlock()
			fmt.Println("nameChannel is", nameChannel)
			fmt.Println("unsubscribe")
		}else{
			println("there's not channel with id ", nameChannel)
		}
	}

	return newChannel, c
}

/*func (p *publisher) UnSubscribe(p, nameChannel string) {
	_, ok := p.subs[nameChannel];
	if ok {
		delete(p.subs, nameChannel);
	}else{
		println("there's not channel with id "+nameChannel)
	}
	return
}
*/
func NewPubsub() *publisher {
	pub := &publisher{}
	pub.mu.Lock()
	pub.subs = make(map[uint64]chan<- Event)
	pub.mu.Unlock()
	return pub
}
