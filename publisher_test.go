package event_emitter

import (
	"sync"
	"testing"

	"github.com/nori-io/common/v3/pkg/domain/event"
	"github.com/stretchr/testify/assert"
)

func TestPublisher_On_1subscriber_1receiving(t *testing.T) {
	pub := NewPubsub()
	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*")

	wg1 := sync.WaitGroup{}
	wg1.Add(1)

	var e1, e2 event.Event
	go func() {
		select {
		case e1 = <-ch1:
			wg1.Done()
		}
	}()

	eventTest := event.Event{
		Name:   "nori/plugins/inited",
		Params: nil,
	}

	pub.Emit(eventTest.Name, eventTest.Params)

	wg1.Wait()

	t.Log("event in 1st channel: ", e1.Name)
	t.Log("event in 2nd channel: ", e2.Name)

	assert.Equal(t, e1.Name, "nori/plugins/inited")

	funcUnsubscribe1()
}

func TestPublisher_On_2subscribers_2receiving(t *testing.T) {
	pub := NewPubsub()
	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*")
	ch2, funcUnsubscribe2 := pub.On("nori/plugins/started")

	wg2 := sync.WaitGroup{}
	wg2.Add(2)

	var e1, e2 event.Event
	go func() {
		select {
		case e1 = <-ch1:
			wg2.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			wg2.Done()
		}
	}()

	eventTest := event.Event{
		Name:   "nori/plugins/started",
		Params: nil,
	}

	pub.Emit(eventTest.Name, eventTest.Params)

	wg2.Wait()

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

	wg3 := sync.WaitGroup{}
	wg3.Add(2)

	var e1, e2 event.Event
	go func() {
		select {
		case e1 = <-ch1:
			wg3.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			//wg3.Done()

		}

	}()
	t.Log("e2 is", e2)
	t.Log("wg3 is", wg3)
	wg3.Done()

	eventTest := event.Event{
		Name:   "nori/plugins/stopped",
		Params: nil,
	}

	pub.Emit(eventTest.Name, eventTest.Params)

	wg3.Wait()

	t.Log("event in 1st channel: ", e1.Name)
	t.Log("event in 2nd channel: ", e2.Name)
	assert.Equal(t, e1.Name, "nori/plugins/stopped")
	assert.Equal(t, e2.Name, "")
	funcUnsubscribe1()
	funcUnsubscribe2()
}

func TestPublisher_On_2subscribers_2receiving_2globalmiddleware(t *testing.T) {

	pub := NewPubsub()

	eventTest := event.Event{
		Name:   "nori/plugins/started",
		Params: make(map[string]interface{}),
	}

	m1 := func(event *event.Event) {
		//event.Params["global_m1_key"] = "global_m1_value"
	}

	m2 := func(event *event.Event) {
		//event.Params["global_m2_key"] = "global_m2_value"
	}

	pub.Use("nori/plugins/*", m1, m2)

	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*")
	ch2, funcUnsubscribe2 := pub.On("nori/plugins/started")

	wg4 := sync.WaitGroup{}
	wg4.Add(2)

	var e1, e2 event.Event
	go func() {
		select {
		case e1 = <-ch1:
			wg4.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			wg4.Done()
		}
	}()

	pub.Emit(eventTest.Name, eventTest.Params)

	wg4.Wait()

	t.Log("params in 1st channel: ", e1.Params)
	t.Log("params in 2nd channel: ", e2.Params)
	assert.Equal(t, e1.Params, e2.Params)

	funcUnsubscribe1()
	funcUnsubscribe2()
}

func TestPublisher_On_2subscribers_2receiving_2localmiddleware(t *testing.T) {
	pub := NewPubsub()

	type TestEventStruct struct {
		Key string
	}

	eventTest := event.Event{
		Name: "nori/plugins/started",
	}

	m1 := func(event *event.Event) {
		event.Params = TestEventStruct{
			Key: "local_m1value",
		}
	}

	m2 := func(event *event.Event) {
		event.Params = TestEventStruct{
			Key: "local_m2value",
		}
	}

	m3 := func(event *event.Event) {
		event.Params = TestEventStruct{
			Key: "local_m3value",
		}
	}

	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*", m1, m2)
	ch2, funcUnsubscribe2 := pub.On("nori/plugins/started", m1, m2, m3)

	wg5 := sync.WaitGroup{}
	wg5.Add(2)

	var e1, e2 event.Event
	go func() {
		select {
		case e1 = <-ch1:
			wg5.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			wg5.Done()
		}
	}()

	pub.Emit(eventTest.Name, eventTest.Params)

	wg5.Wait()

	t.Log("params in 1st channel: ", e1.Params)
	t.Log("params in 2nd channel: ", e2.Params)
	assert.NotEqual(t, e1.Params, e2.Params)

	eventForComparing1 := event.Event{
		Name: "nori/plugins/started",
		Params: TestEventStruct{
			Key: "local_m2value",
		},
	}

	assert.Equal(t, e1.Params, eventForComparing1.Params)

	eventForComparing2 := event.Event{
		Name: "nori/plugins/started",
		Params: TestEventStruct{
			Key: "local_m3value",
		},
	}

	assert.Equal(t, e2.Params, eventForComparing2.Params)

	funcUnsubscribe1()
	funcUnsubscribe2()
}

func TestPublisher_On_2subscribers_2receiving_2_globalmiddleware_2localmiddleware(t *testing.T) {
	pub := NewPubsub()

	eventTest := event.Event{
		Name: "nori/plugins/started",
	}

	type TestEventStruct struct {
		Key_Primary   string
		Key_Secondary string
	}
	mG1 := func(event *event.Event) {
		event.Params = TestEventStruct{
			Key_Primary: "key_primary_1",
		}

	}

	mG2 := func(event *event.Event) {
		event.Params = TestEventStruct{
			Key_Primary: "key_primary_2",
		}
	}

	mL1 := func(event *event.Event) {
		event.Params = TestEventStruct{
			Key_Primary:   "key_primary_1",
			Key_Secondary: "key_secondary_2",
		}
	}

	mL2 := func(event *event.Event) {
		event.Params = TestEventStruct{
			Key_Primary:   "key_primary_2",
			Key_Secondary: "key_secondary_2",
		}
	}

	pub.Use("nori/plugins/*", mG1, mG2)

	ch1, funcUnsubscribe1 := pub.On("nori/plugins/*", mL1)
	ch2, funcUnsubscribe2 := pub.On("nori/plugins/started", mL1, mL2)

	wg6 := sync.WaitGroup{}
	wg6.Add(2)

	var e1, e2 event.Event
	go func() {
		select {
		case e1 = <-ch1:
			wg6.Done()
		}
	}()

	go func() {
		select {
		case e2 = <-ch2:
			wg6.Done()
		}
	}()

	pub.Emit(eventTest.Name, eventTest.Params)

	wg6.Wait()

	t.Log("params in 1st channel: ", e1.Params)
	t.Log("params in 2nd channel: ", e2.Params)
	assert.NotEqual(t, e1.Params, e2.Params)

	eventForComparing1 := event.Event{
		Name:   "nori/plugins/started",
		Params: make(map[string]interface{}),
	}
	eventForComparing1.Params = TestEventStruct{
		Key_Primary:   "key_primary_1",
		Key_Secondary: "key_secondary_2",
	}

	assert.Equal(t, e1.Params, eventForComparing1.Params)

	eventForComparing2 := event.Event{
		Name:   "nori/plugins/started",
		Params: make(map[string]interface{}),
	}
	eventForComparing2.Params = TestEventStruct{
		Key_Primary:   "key_primary_2",
		Key_Secondary: "key_secondary_2",
	}

	assert.Equal(t, e2.Params, eventForComparing2.Params)

	funcUnsubscribe1()
	funcUnsubscribe2()
}
