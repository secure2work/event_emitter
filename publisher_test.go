package event_emitter

import (
	"fmt"
	"testing"
)

func TestPublisher_Subscribe(t *testing.T) {
	pub:=NewPubsub()
	ch1,funcUnsubscribe1:=pub.Subscribe()
	ch2,funcUnsubscribe2:=pub.Subscribe()

	fmt.Println("channel", ch1)

	fmt.Println("channel2", ch2)

	//fmt.Println("func is", b)
	funcUnsubscribe1()
	funcUnsubscribe2()
}
