package event_emitter

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublisher_On_1subscriber_1receiving(t *testing.T) {
	pub := NewPubsub()
	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*")

	wg := sync.WaitGroup{}
	wg.Add(1)

	var e1, e2 Event
	go func() {
		select {
		case e1 = <-ch1:
			wg.Done()
		}
	}()

	eventTest := Event{
		Name:   "nori/plugins/inited",
		Params: nil,
	}

	pub.Emit(eventTest)

	wg.Wait()

	t.Log("event in 1st channel: ", e1.Name)
	t.Log("event in 2nd channel: ", e2.Name)

	assert.Equal(t, e1.Name, "nori/plugins/inited")

	funcUnsubscribe1()
}

func TestPublisher_On_2subscribers_2receiving(t *testing.T) {
	pub := NewPubsub()
	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*")
	ch2, funcUnsubscribe2 := pub.On("nori/plugins/started")

	wg := sync.WaitGroup{}
	wg.Add(2)

	var e1, e2 Event
	go func() {
		select {
		case e1 = <-ch1:
			wg.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			wg.Done()
		}
	}()

	eventTest := Event{
		Name:   "nori/plugins/started",
		Params: nil,
	}

	pub.Emit(eventTest)

	wg.Wait()

	t.Log("event in 1st channel: ", e1.Name)
	t.Log("event in 2nd channel: ", e2.Name)

	assert.Equal(t, e1.Name, e2.Name)

	funcUnsubscribe1()
	funcUnsubscribe2()
}

func TestPublisher_On_2subscribers_1receiving(t *testing.T) {
	pub := NewPubsub()
	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*")
	ch2, funcUnsubscribe2 := pub.On("nori/plugins/started")

	wg := sync.WaitGroup{}
	wg.Add(2)

	var e1, e2 Event
	go func() {
		select {
		case e1 = <-ch1:
			wg.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			wg.Done()

		}

	}()
	wg.Done()

	eventTest := Event{
		Name:   "nori/plugins/stopped",
		Params: nil,
	}

	pub.Emit(eventTest)

	wg.Wait()

	t.Log("event in 1st channel: ", e1.Name)
	t.Log("event in 2nd channel: ", e2.Name)
	assert.Equal(t, e1.Name, "nori/plugins/stopped")
	assert.Equal(t, e2.Name, "")
	funcUnsubscribe1()
	funcUnsubscribe2()
}

func TestPublisher_On_2subscribers_2receiving_2globalmiddleware(t *testing.T) {
	pub := NewPubsub()

	eventTest := Event{
		Name:   "nori/plugins/started",
		Params: make(map[string]interface{}),
	}

	m1 := func(event *Event) {
		event.Params["global_m1_key"] = "global_m1_value"
	}

	m2 := func(event *Event) {
		event.Params["global_m2_key"] = "global_m2_value"
	}

	pub.Use("nori/plugins/*", m1, m2)

	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*")
	ch2, funcUnsubscribe2 := pub.On("nori/plugins/started")

	wg := sync.WaitGroup{}
	wg.Add(2)

	var e1, e2 Event
	go func() {
		select {
		case e1 = <-ch1:
			wg.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			wg.Done()
		}
	}()

	pub.Emit(eventTest)

	wg.Wait()

	t.Log("params in 1st channel: ", e1.Params)
	t.Log("params in 2nd channel: ", e2.Params)
	assert.Equal(t, e1.Params, e2.Params)

	funcUnsubscribe1()
	funcUnsubscribe2()
}

func TestPublisher_On_2subscribers_2receiving_2localmiddleware(t *testing.T) {
	pub := NewPubsub()

	eventTest := Event{
		Name:   "nori/plugins/started",
		Params: make(map[string]interface{}),
	}

	m1 := func(event *Event) {
		event.Params["local_m1key"] = "local_m1value"
	}

	m2 := func(event *Event) {
		event.Params["local_m2key"] = "local_m2value"
	}

	m3 := func(event *Event) {
		event.Params["local_m3key"] = "local_m3value"
	}

	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*", m1, m2)
	ch2, funcUnsubscribe2 := pub.On("nori/plugins/started", m1, m2, m3)

	wg := sync.WaitGroup{}
	wg.Add(2)

	var e1, e2 Event
	go func() {
		select {
		case e1 = <-ch1:
			wg.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			wg.Done()
		}
	}()

	pub.Emit(eventTest)

	wg.Wait()

	t.Log("params in 1st channel: ", e1.Params)
	t.Log("params in 2nd channel: ", e2.Params)
	assert.NotEqual(t, e1.Params, e2.Params)

	eventForComparing1 := Event{
		Name:   "nori/plugins/started",
		Params: make(map[string]interface{}),
	}
	eventForComparing1.Params["local_m1key"] = "local_m1value"
	eventForComparing1.Params["local_m2key"] = "local_m2value"
	assert.Equal(t, e1.Params, eventForComparing1.Params)

	eventForComparing2 := Event{
		Name:   "nori/plugins/started",
		Params: make(map[string]interface{}),
	}
	eventForComparing2.Params["local_m1key"] = "local_m1value"
	eventForComparing2.Params["local_m2key"] = "local_m2value"
	eventForComparing2.Params["local_m3key"] = "local_m3value"

	assert.Equal(t, e2.Params, eventForComparing2.Params)

	funcUnsubscribe1()
	funcUnsubscribe2()
}

func TestPublisher_On_2subscribers_2receiving_2_globalmiddleware_2localmiddleware(t *testing.T) {
	pub := NewPubsub()

	eventTest := Event{
		Name:   "nori/plugins/started",
		Params: make(map[string]interface{}),
	}

	mG1 := func(event *Event) {
		event.Params["global_m1key"] = "global_m1value"
	}

	mG2 := func(event *Event) {
		event.Params["global_m2key"] = "global_m2value"
	}

	mL1 := func(event *Event) {
		event.Params["local_m1key"] = "local_m1value"
	}

	mL2 := func(event *Event) {
		event.Params["local_m2key"] = "local_m2value"
	}

	pub.Use("nori/plugins/*", mG1, mG2)

	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*", mL1)
	ch2, funcUnsubscribe2 := pub.On("nori/plugins/started", mL1, mL2)

	wg := sync.WaitGroup{}
	wg.Add(2)

	var e1, e2 Event
	go func() {
		select {
		case e1 = <-ch1:
			wg.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			wg.Done()
		}
	}()

	pub.Emit(eventTest)

	wg.Wait()

	t.Log("params in 1st channel: ", e1.Params)
	t.Log("params in 2nd channel: ", e2.Params)
	assert.NotEqual(t, e1.Params, e2.Params)

	eventForComparing1 := Event{
		Name:   "nori/plugins/started",
		Params: make(map[string]interface{}),
	}
	eventForComparing1.Params["global_m1key"] = "global_m1value"
	eventForComparing1.Params["global_m2key"] = "global_m2value"
	eventForComparing1.Params["local_m1key"] = "local_m1value"
	assert.Equal(t, e1.Params, eventForComparing1.Params)

	eventForComparing2 := Event{
		Name:   "nori/plugins/started",
		Params: make(map[string]interface{}),
	}
	eventForComparing2.Params["global_m1key"] = "global_m1value"
	eventForComparing2.Params["global_m2key"] = "global_m2value"
	eventForComparing2.Params["local_m1key"] = "local_m1value"
	eventForComparing2.Params["local_m2key"] = "local_m2value"

	assert.Equal(t, e2.Params, eventForComparing2.Params)

	funcUnsubscribe1()
	funcUnsubscribe2()
}
