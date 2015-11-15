//+build js

package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/gopherjs/websocket"
	"github.com/platinasystems/weeb"
	"github.com/platinasystems/weeb/canvas"
	. "github.com/platinasystems/weeb/html"
	"github.com/platinasystems/weeb/r2"

	"fmt"
	ghtml "html"
	"log"
	"strings"
	"time"
)

var jq = jquery.NewJQuery

var window = js.Global.Get("window")
var document = window.Get("document")
var location = document.Get("location")
var screenPixelsPerPoint float64

var currentPath = location.Get("pathname").String()
var doc = mySite.DocByPath[currentPath]
var page = mySite.PageByPath[currentPath]
var rpc *weeb.Rpc

func submitInput(input jquery.JQuery) {
	if currentPath != "/exec" {
		log.Println("ignore non exec ", currentPath)
		return
	}

	cmd := strings.TrimSpace(input.Val())
	if cmd == "" {
		return
	}

	go func() {
		d := mySite.DocByPath[currentPath]
		page := mySite.PageByPath[currentPath]

		var pre, result string
		args := strings.Split(cmd, " ")
		err := rpc.Call("Listener.Exec", &args, &result)
		if err != nil {
			pre = fmt.Sprintf("error: %v", err)
		} else {
			pre = result
		}
		pre = ghtml.EscapeString(pre)
		body := d.Pre(&String{pre}).Markup(d)

		jq(input).SetAttr("disabled", "yes")
		ec := jq(input).Closest("div.ExecCommand")
		res := ec.Find("div.ExecResult")

		jq(res).SetHtml(body)

		jq(res).Closest("div.row").RemoveClass("hide")

		parent := jq(ec).Parent()
		add := page.(*execPage).ExecCommand.Body(d)

		parent.Append(add.Markup(d))
		jqBind(jq(parent))

		jq(parent).Find("input:last").Focus()
	}()
}

func jqBind(j jquery.JQuery) {
	jq(j).Find("[replace]").On(jquery.CLICK, func(e jquery.Event) {
		replace(jq(e.Target))
		e.PreventDefault()
	})

	jq(j).Find(".submit_on_click").On(jquery.CLICK, func(e jquery.Event) {
		form := jq(e.Target).Closest("form")
		form.Find(":input").Each(func(index int, x interface{}) {
			i := jq(x)
			log.Printf("%s: %s", i.Attr("id"), i.Val())
		})
		e.PreventDefault()
	})

	jq(j).Find("form").On(jquery.SUBMIT, func(e jquery.Event) {
		e.PreventDefault()
	})

	jq(j).Find(".submit_on_enter > input").On(jquery.KEYUP, func(e jquery.Event) {
		if e.KeyCode == 13 {
			submitInput(jq(e.Target))
		}
		e.PreventDefault()
	})

	// Draw canvas elements for this page.
	jq(j).Find("canvas").Each(func(index int, c interface{}) {
		canvas := jq(c)
		canvasDraw(canvas)
		canvas.On(jquery.CLICK, func(e jquery.Event) {
			canvasEvent(e)
			e.PreventDefault()
		})
	})
}

func replace(t jquery.JQuery) {
	href := t.Attr("href")
	id := t.Attr("replace")

	if len(id) == 0 || len(href) == 0 {
		log.Printf("%s: replace fails id '%s' href '%s'", currentPath, id, href)
		return
	}

	if page, doc, _ := mySite.Match(href); page == nil {
		log.Printf("page not found %s", href)
	} else {
		idSelector := "#" + id
		elt := jq(idSelector)
		if elt.Length > 0 {
			currentPath = href

			bn := doc.BodyNodeById[id]
			elt.ReplaceWith(bn.Markup(doc))

			doc.AddBodyNodeEventListeners(document, bn)

			// Select again since we've just changed it.
			elt := jq(idSelector)

			jqBind(elt)
			elt.Call("foundation", "reflow")

			// Change displayed page path without reloading page.
			window.Get("history").Call("pushState", "object or string", "Title", href)
		}
	}
}

func canvasContext(c jquery.JQuery) (x *canvas.Context) {
	x = canvas.GetContext(c)

	h := float64(c.Width())
	w := float64(c.Height())
	x.Size = r2.XY(h/screenPixelsPerPoint, w/screenPixelsPerPoint)

	x.Scale(r2.XY1(screenPixelsPerPoint))

	return
}

func getDrawerListener(c jquery.JQuery) (d canvas.Drawer, l canvas.Listener) {
	id := c.Attr("id")
	page := mySite.PageByPath[currentPath]
	switch p := page.(type) {
	case canvas.Interface:
		var ok bool
		if d, ok = p.Drawer(id); !ok {
			log.Printf("%s: no canvas drawer for id '%s'", currentPath, id)
		}
		l, _ = p.Listener(id)
	default:
		log.Printf("%s: found canvas on page with no interface", currentPath)
	}
	return
}

func canvasDraw(c jquery.JQuery) {
	if c.Attr("width") == "" && c.Attr("height") == "" {
		// Determine height from aspect ratio and parent's width.
		aspect := 1.5 // 3 by 2 aspect ratio
		attr := c.Attr("aspect")
		if attr != "" {
			_, err := fmt.Sscanf(attr, "%f", &aspect)
			if err != nil {
				panic(err)
			}
		}

		parent := c.Parent()
		w := parent.Width()
		c.SetAttr("width", w)
		c.SetAttr("height", float64(w)/aspect)
	}

	d, _ := getDrawerListener(c)
	if d != nil {
		go d.Draw(canvasContext(c))
	}
}

func canvasEvent(e jquery.Event) {
	c := jq(e.Target)
	_, l := getDrawerListener(c)
	if l != nil {
		xy := c.Offset()
		px := float64(e.PageX-xy.Left) / screenPixelsPerPoint
		py := float64(e.PageY-xy.Top) / screenPixelsPerPoint
		l.Event(canvasContext(c), r2.XY(px, py))
	}
}

func findScreenPixelsPerPoint() float64 {
	id := fmt.Sprintf("_%x", time.Now().UnixNano())
	n := float64(2 * 72)
	jq("body").Append(fmt.Sprintf(`<div id="%s" style="width:%.0fpt;visible:hidden;padding:0px"></div>`, id, n))
	r := document.Call("getElementById", id).Get("offsetWidth").Float() / n
	jq("#" + id).Remove()
	return r
}

func main() {
	if true {
		wsBaseURL := fmt.Sprintf("ws://%s:%s/ws/", location.Get("hostname"), location.Get("port"))
		ws, err := websocket.Dial(wsBaseURL + "rpc/foo")
		if err != nil {
			log.Fatal("Dial ", err)
			return
		}
		l := &Listener{}
		rpc = weeb.NewRpc(ws, l)
		l.rpc = rpc

		// go HelloRpcClient(rpc)
		go rpc.Serve()
	}

	if false {
		fu()
	}

	jq().Ready(func() {
		screenPixelsPerPoint = findScreenPixelsPerPoint()
		doc.AddEventListeners(document)
		jqBind(jq(":root"))
	})
}

func err(kp *js.Object) {
	print("error")
	print(kp)
}

var tmp = js.Global.Call("eval", `({ view: function(o) { return new Uint8Array(o); } })`)

func ok(kp *js.Object) {
	print("ok")
	print(kp)
	fmt.Printf("%d\n", kp.Get("byteLength").Int())
	x := tmp.Call("view", kp)
	log.Println(x)
}

func save(kp *js.Object) {
	c := js.Global.Get("window").Get("crypto").Get("subtle")
	var b [32]byte
	x := js.Global.Get("Object").New()
	x.Set("name", "AES-GCM")
	x.Set("iv", js.NewArrayBuffer(b[:12]))
	p := c.Call("encrypt", x, kp, js.NewArrayBuffer(b[:]))
	p.Call("then", ok).Call("catch", err)
}

func fu() {
	c := js.Global.Get("window").Get("crypto").Get("subtle")
	var b [32]byte
	p := c.Call("importKey", "raw", js.NewArrayBuffer(b[:]), "AES-GCM", true, [2]string{"encrypt", "decrypt"})
	p.Call("then", save).Call("catch", err)
}
