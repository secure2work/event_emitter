package event_emitter

import (
	"fmt"
	"testing"
)

func TestPublisher_Subscribe(t *testing.T) {
	pub:=NewPubsub()
	a,b:=pub.Subscribe()
	fmt.Println("channel", a)
	//fmt.Println("func is", b)
	b()
}
