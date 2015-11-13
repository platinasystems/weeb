package html

func (d *Doc) addEventListener(attrs *Attrs, ls ...interface{}) {
	id := d.assignID(attrs)
	if d.EventListenersById == nil {
		d.EventListenersById = make(map[string][]interface{})
	}
	d.EventListenersById[id] = append(d.EventListenersById[id], ls...)
}

type EventType int

// fixme go generate
const (
	Invalid EventType = iota
	// HTML events
	Load
	Unload
	Abort
	Error
	SelectEvent
	Change
	SubmitEvent
	ResetEvent
	Focus
	Blur
	Resize
	Scroll
	// Mouse events
	Click
	MouseDown
	MouseUp
	MouseOver
	MouseMove
	MouseOut
)

type ClickInterface interface {
	Click(e *MouseEvent)
}

type LoadInterface interface {
	Load(e *Event)
}

func isEventListener(x interface{}) bool {
	if _, ok := x.(ClickInterface); ok {
		return true
	}
	if _, ok := x.(LoadInterface); ok {
		return true
	}
	return false
}
