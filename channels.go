package event_emitter

type EventEmitter interface {
	Subscribe() (evt <-chan Event, unsubscribe func() error)
	Emit(event Event)
}

type Event struct {
	Name   string
	Params map[string]interface{}
}
