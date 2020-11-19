package event_emitter

import (
	"fmt"
	"testing"
)

func TestPublisher_Subscribe(t *testing.T) {
	pub:=NewPubsub()
	ch,funcUnsubscribe:=pub.Subscribe()
	fmt.Println("channel", ch)
	//fmt.Println("func is", b)
	funcUnsubscribe()
}
