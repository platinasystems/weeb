//+build js

package html

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/platinasystems/weeb/r2"
)

type Event struct {
	*js.Object
	Bubbles    bool `js:"bubbles"`
	Target     *js.Object
	Timestamp  int64
	Type       string `js:"type"`
	EventPhase int    `js:"eventPhase"`
}

func (e *Event) Init() {
	e.Target = e.Object.Get("target")
	e.Timestamp = e.Object.Get("timeStamp").Int64() * 1000
}

func (e *Event) PreventDefault()  { e.Object.Call("preventDefault") }
func (e *Event) StopPropagation() { e.Object.Call("stopPropagation") }

type KeyModifier uint16

const (
	KeyControl KeyModifier = 1 << iota
	KeyShift
	KeyAlt
	KeyMeta
)

func mod(o *js.Object) (m KeyModifier) {
	if o.Get("ctrlKey").Bool() {
		m |= KeyControl
	}
	if o.Get("shiftKey").Bool() {
		m |= KeyShift
	}
	if o.Get("altKey").Bool() {
		m |= KeyAlt
	}
	if o.Get("metaKey").Bool() {
		m |= KeyMeta
	}
	return
}

type MouseEvent struct {
	Event
	ScreenX, ClientX r2.X
	KeyModifier
	Button        uint16
	RelatedTarget *js.Object
}

func xy(o *js.Object, n string) r2.X {
	return r2.XY(
		o.Get(n+"X").Float(),
		o.Get(n+"Y").Float())
}

func (e *MouseEvent) Init() {
	e.Event.Init()
	o := e.Object
	e.ScreenX = xy(o, "screen")
	e.ClientX = xy(o, "client")
	e.KeyModifier = mod(o)
	e.Button = uint16(o.Get("button").Int())
	e.RelatedTarget = o.Get("relatedTarget")
}

func addEventListenersId(document *js.Object, id string, ls []interface{}) {
	for _, l := range ls {
		if v, ok := l.(ClickInterface); ok {
			x := document.Call("getElementById", id)
			x.Call("addEventListener", "click",
				func(e *MouseEvent) {
					e.Init()
					v.Click(e)
				}, false)
		}
		if v, ok := l.(LoadInterface); ok {
			// Fixme event never fires snice main is called from load.
			w := js.Global.Get("window")
			// w := document.Call("getElementById", id)
			// w := document
			f := func(e *Event) {
				e.Init()
				v.Load(e)
			}
			w.Call("addEventListener", "load", f, false)
		}
	}
}

func (d *Doc) AddEventListeners(document *js.Object) {
	for id, listeners := range d.EventListenersById {
		addEventListenersId(document, id, listeners)
	}
}

func (d *Doc) AddBodyNodeEventListeners(document *js.Object, b BodyNode) {
	attrs := b.attrs()
	id := attrs.ID
	if len(id) > 0 {
		if ls, ok := d.EventListenersById[id]; ok {
			addEventListenersId(document, id, ls)
		}
	}
	bv := b.bodyVec()
	for i := range bv {
		d.AddBodyNodeEventListeners(document, bv[i])
	}
}
