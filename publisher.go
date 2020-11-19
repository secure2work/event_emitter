package event_emitter

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type publisher struct {
	mu   sync.RWMutex
	subs map[string]chan<- Event
}

//unsubscribe func(p, nameChannel string)
//Убрать содержимое скобок или нет?
func (p *publisher) Subscribe() (evt <-chan Event, unsubscribe func()) {
	newChannel := make(chan Event)
	nameChannel := uuid.New().String()
	p.subs[nameChannel] = newChannel


	c := func() {
		_, ok := p.subs[nameChannel];
		if ok {
			delete(p.subs, nameChannel);
			fmt.Println("nameChannel is", nameChannel)
			fmt.Println("unsubscribe")
		}else{
			println("there's not channel with id "+nameChannel)
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
	pub.subs = make(map[string]chan<- Event)
	return pub
}
