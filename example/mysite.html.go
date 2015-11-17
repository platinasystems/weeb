package main

import (
	"fmt"
	"math/rand"

	"github.com/platinasystems/elib/elog"
	"github.com/platinasystems/weeb"
	"github.com/platinasystems/weeb/canvas"
	. "github.com/platinasystems/weeb/html"
	"github.com/platinasystems/weeb/r2"
)

type ExecCommand struct {
	Prompt string
	Cmd    string
	Result string
}

func (c *ExecCommand) Body(d *Doc) BodyVec {
	i := d.Input("type=text")
	if len(c.Prompt) > 0 {
		i.Attrs.Set(i, d, "placeholder="+c.Prompt)
	}
	return BodyVec{
		d.Div(".ExecCommand",
			d.Div(".row",
				d.Div(".large-12 columns submit_on_enter", i)),
			d.Div(".row hide",
				d.Div(".large-12 columns",
					d.Div(".ExecResult panel")),
			),
		),
	}
}

var mySite *weeb.Site = initSite()

func initSite() *weeb.Site {
	s := &weeb.Site{
		PageByPath: map[string]weeb.Page{
			"/":       &rootPage{},
			"/page/":  &pathPage{},
			"/exec":   &execPage{},
			"/canvas": &myCanvasPage{},
			"/elog":   &myElogPage{},
		},
	}

	s.DocByPath = make(map[string]*Doc)
	for path, p := range s.PageByPath {
		d := &Doc{
			Head: head,
		}

		d.Body = p.PageBody(path, d)

		s.DocByPath[path] = d
	}

	return s
}

var head = []HeadNode{
	&Meta{Charset: "utf-8"},
	&Meta{HttpEquiv: "X-UA-Compatible", Content: "IE=edge"},
	&Meta{Name: "viewport", Content: "width=device-width, initial-scale=1"},
	&Link{Rel: "stylesheet", Type: "text/css", Href: "/css/eg.min.css"},
	&Script{Type: "text/javascript", Src: "/js/foundation_deps.min.js"},
	&Script{Type: "text/javascript", Src: "/js/foundation.min.js"},
	&Script{Type: "text/javascript", Src: "/js/js.min.js"},
	&Title{"Weeb Title"},
}

type standardBody struct {
	BodyVec
	noSideBar bool
}

func (s *standardBody) Body(d *Doc) BodyVec {
	pb := d.Div("#page_body", s.BodyVec)
	var cols BodyVec
	if !s.noSideBar {
		cols = BodyVec{
			d.Div(".large-3 columns",
				d.Div(".hide-for-small",
					d.Div(".sidebar",
						d.UL(".side-nav",
							d.LI(".heading", "Heading One"),
							d.LI(d.A("href=/page/1", "replace=page_body", "Link 1")),
							d.LI(d.A("href=/page/2", "replace=page_body", "Link 2")),
							d.LI(d.A("href=/notfound", "replace=page_body", "Link 3")),
							d.LI(d.A("href=#", "Link 4")),
							d.LI(".divider"),
							d.LI(".heading", "Heading Two"),
							d.LI(d.A("href=/exec", "replace=page_body", "Exec")),
							d.LI(d.A("href=/canvas", "replace=page_body", "Canvas")),
							d.LI(d.A("href=/elog", "replace=page_body", "Event Log")),
							d.LI(d.A("href=#", "Link 3")),
							d.LI(d.A("href=#", "Link 4")),
						)))),
			d.Div(".large-9 columns", pb),
		}
	} else {
		cols = BodyVec{
			d.Div(".small-12 columns", pb),
		}
	}
	return BodyVec{
		d.Div(".contain-to-grid fixed",
			d.Nav(".top-bar", "data-topbar=", "role=navigation",
				d.UL(".title-area",
					d.LI(".name", d.H1(d.A("href=/", "replace=page_body", "My Site"))),
					d.LI(".toggle-topbar menu-icon", d.A("href=#", d.Span("My Site"))),
				),
				d.Section(".top-bar-section",
					d.UL(".right",
						d.LI(".active", d.A("href=#", "Right Button Active")),
						d.LI(".has-dropdown",
							d.A("href=#", "Right Button Dropdown"),
							d.UL(".dropdown",
								d.LI(d.A("href=#", "First link in dropdown")),
								d.LI(".active", d.A("href=#", "Active link in dropdown"))))),
					d.UL(".left",
						d.LI(d.A("href=#", "Left Nav Button"))),
				),
			),
		),
		d.Div(".row", cols),

		// Initialize Zurb foundation javascript.
		&Script{Content: "$(document).foundation();"},
	}
}

type inlineLabel struct {
	Id   string
	Type InputType
}

func (t *inlineLabel) Body(d *Doc) BodyVec {
	p := t.Type
	if p == InputType(0) {
		p = Text
	}
	i := d.Input("#" + t.Id)
	i.InputType = p
	return BodyVec{
		d.Div(".row",
			d.Div(".large-3 columns", d.Label(".right inline", "for="+t.Id, t.Id)),
			d.Div(".large-9 columns", i)),
	}
}

type T struct {
	A int
	B int
	C string
	D float64
}

func (t *T) Body(d *Doc) BodyVec {
	if false {
		return BodyVec{
			d.Div(".row",
				d.Div(".large-4 columns", d.Label("A", d.Input("#A", "type=number"))),
				d.Div(".large-4 columns", d.Label("B", d.Input("#B", "type=number"))),
				d.Div(".large-4 columns", d.Label("C", d.Input("#C", "type=text"))),
			),
		}
	} else {
		return BodyVec{
			d.Div(".row",
				d.Div(".large-4 columns", &inlineLabel{Id: "A", Type: Number}),
				d.Div(".large-4 columns",
					d.Div(".row",
						d.Div(".large-3 columns", d.Label(".right inline", "for=B", "B")),
						d.Div(".large-9 columns", d.Input("#B", "type=number")))),
				d.Div(".large-4 columns",
					d.Div(".row",
						d.Div(".large-3 columns", d.Label(".right inline", "for=C", "C")),
						d.Div(".large-9 columns", d.Input("#C", "type=text")))),
			),
		}
	}
}

type foo struct{}

func (f *foo) Body(d *Doc) BodyVec {
	return BodyVec{
		d.Svg("width=100", "height=100",
			d.Circle("cx=50", "cy=50", "r=40", "stroke=black", "stroke-width=3", "fill=red")),
	}
}

func (f *foo) Click(e *MouseEvent) {
	fmt.Printf("here %v\n", e)
}

type rootPage struct{}

func (r *rootPage) PageBody(path string, d *Doc) BodyVec {
	return (&standardBody{BodyVec: r.Body(path, d)}).Body(d)
}

func (r *rootPage) Body(path string, d *Doc) BodyVec {
	return BodyVec{
		d.H1("Root Title"),
		d.H3(".subheader", `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam
tristique lacus sed ex viverra, quis bibendum mauris vehicula.`),
		d.HR(),
		d.Div(".row",
			d.Div(".large-12 columns", &foo{}),
			d.Div(".large-12 columns",
				d.Blockquote("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas rutrum blandit purus, sed sollicitudin augue elementum at. Donec congue enim in mauris finibus congue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris et tempor ligula. Phasellus porttitor consequat nisl sed viverra. Nunc elit ligula, tempus id sagittis at, iaculis et magna. Ut blandit erat nisi, vitae blandit eros vestibulum sed. Ut malesuada pellentesque turpis quis ullamcorper. Nam tincidunt, mauris et laoreet auctor, enim turpis rutrum urna, vel mollis diam metus sit amet eros. Quisque accumsan sapien turpis, non tincidunt nulla imperdiet ut."),
				d.HR(),
				d.Address(
					"Written by ", d.A("href=mailto:webmaster@example.com", "Jon Doe"), ".", d.BR(),
					"Visit us at:", d.BR(),
					"Example.com", d.BR(),
					"Box 564, Disneyland", d.BR(),
					"USA"),
				d.HR(),
				d.Table(
					d.Colgroup(
						d.Col("a=1", "b=2"),
						d.Col("a=1", "b=2"),
					),
					d.Thead(
						d.TR(d.TH("Foo"), d.TH("Bar")),
					),
					d.Tbody(
						d.TR(d.TH("1"), d.TH("2")),
						d.TR(d.TH(3), d.TH(4)),
					),
				),
				d.HR(),
				d.Div(".panel",
					d.Form(
						d.Div(".row", &T{A: 1, C: "value"}),
						d.Div(".row",
							d.Div(".large-12 columns",
								d.A(".large radius button submit_on_click", "SUBMIT"))),
					)))),
		d.HR(),
		d.Blockquote("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas rutrum blandit purus, sed sollicitudin augue elementum at. Donec congue enim in mauris finibus congue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris et tempor ligula. Phasellus porttitor consequat nisl sed viverra. Nunc elit ligula, tempus id sagittis at, iaculis et magna. Ut blandit erat nisi, vitae blandit eros vestibulum sed. Ut malesuada pellentesque turpis quis ullamcorper. Nam tincidunt, mauris et laoreet auctor, enim turpis rutrum urna, vel mollis diam metus sit amet eros. Quisque accumsan sapien turpis, non tincidunt nulla imperdiet ut."),
		d.HR(),
		d.Blockquote("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas rutrum blandit purus, sed sollicitudin augue elementum at. Donec congue enim in mauris finibus congue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris et tempor ligula. Phasellus porttitor consequat nisl sed viverra. Nunc elit ligula, tempus id sagittis at, iaculis et magna. Ut blandit erat nisi, vitae blandit eros vestibulum sed. Ut malesuada pellentesque turpis quis ullamcorper. Nam tincidunt, mauris et laoreet auctor, enim turpis rutrum urna, vel mollis diam metus sit amet eros. Quisque accumsan sapien turpis, non tincidunt nulla imperdiet ut."),
		d.HR(),
		d.Blockquote("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas rutrum blandit purus, sed sollicitudin augue elementum at. Donec congue enim in mauris finibus congue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris et tempor ligula. Phasellus porttitor consequat nisl sed viverra. Nunc elit ligula, tempus id sagittis at, iaculis et magna. Ut blandit erat nisi, vitae blandit eros vestibulum sed. Ut malesuada pellentesque turpis quis ullamcorper. Nam tincidunt, mauris et laoreet auctor, enim turpis rutrum urna, vel mollis diam metus sit amet eros. Quisque accumsan sapien turpis, non tincidunt nulla imperdiet ut."),
		d.HR(),
		d.Blockquote("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas rutrum blandit purus, sed sollicitudin augue elementum at. Donec congue enim in mauris finibus congue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris et tempor ligula. Phasellus porttitor consequat nisl sed viverra. Nunc elit ligula, tempus id sagittis at, iaculis et magna. Ut blandit erat nisi, vitae blandit eros vestibulum sed. Ut malesuada pellentesque turpis quis ullamcorper. Nam tincidunt, mauris et laoreet auctor, enim turpis rutrum urna, vel mollis diam metus sit amet eros. Quisque accumsan sapien turpis, non tincidunt nulla imperdiet ut."),
		d.HR(),
		d.Blockquote("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas rutrum blandit purus, sed sollicitudin augue elementum at. Donec congue enim in mauris finibus congue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris et tempor ligula. Phasellus porttitor consequat nisl sed viverra. Nunc elit ligula, tempus id sagittis at, iaculis et magna. Ut blandit erat nisi, vitae blandit eros vestibulum sed. Ut malesuada pellentesque turpis quis ullamcorper. Nam tincidunt, mauris et laoreet auctor, enim turpis rutrum urna, vel mollis diam metus sit amet eros. Quisque accumsan sapien turpis, non tincidunt nulla imperdiet ut."),
		d.HR(),
	}
}

type pathPage struct{}

func (p *pathPage) PageBody(path string, d *Doc) BodyVec {
	s := &standardBody{BodyVec: p.Body(path, d)}
	return s.Body(d)
}

func (p *pathPage) Body(path string, d *Doc) BodyVec {
	return BodyVec{
		d.H1(fmt.Sprintf("Path %s", path)),
		d.H3(".subheader", `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam
tristique lacus sed ex viverra, quis bibendum mauris vehicula.`),
	}
}

type execPage struct {
	ExecCommand
}

func (p *execPage) PageBody(path string, d *Doc) BodyVec {
	s := &standardBody{
		BodyVec: p.Body(path, d),
		// noSideBar: true,
	}
	return s.Body(d)
}

func (p *execPage) Body(path string, d *Doc) BodyVec {
	p.ExecCommand = ExecCommand{Prompt: "Enter Command"}
	return BodyVec{
		d.H2("Exec Command"),
		d.Div(".panel",
			d.Form(
				&p.ExecCommand,
			),
		),
	}
}

type myCanvas struct {
	id, greeting string
}

func (c *myCanvas) Body(d *Doc) BodyVec {
	return BodyVec{
		d.Div(".panel", d.Canvas("#"+c.id, "aspect=2")),
	}
}

func (c *myCanvas) Load(e *Event) {
	print(e)
}

func (c *myCanvas) Draw(x *canvas.Context) {
	if true {
		x.SetFillStyle(canvas.RGBA{R: 1, A: 1})
		r := r2.Rect{X: 0, Size: 1 + 1i}
		scale := float64(72)

		x.Save()
		x.Scale(r2.XY(scale, scale))
		x.FillRect(r.X, r.Size)
		x.BeginPath()
		x.Arc(r.Size, .5, 0, r2.AngleMax, true)
		x.SetFillStyle(canvas.RGBA{G: 1, A: 1})
		x.Fill()
		x.LineWidth = .1
		x.SetStrokeStyle(canvas.RGBA{A: .25})
		x.Stroke()
		x.Restore()

		x.Save()
		x.SetFillStyle(canvas.RGBA{A: 1})
		x.Font = "24pt Ariel"
		x.Translate(x.Size / 2)
		x.Rotate(r2.AngleMax / 8)
		w := x.MeasureText(c.greeting)
		x.FillText(c.greeting, -w/2)
		x.Restore()
	} else {
		x.Save()
		x.SetStrokeStyle(canvas.RGBA{A: .8})
		x.LineWidth = 2
		for i := 0; i < 1000; i++ {
			z := r2.XY(
				rand.Float64()*x.Size.X(),
				rand.Float64()*x.Size.Y(),
			)
			x.SetFillStyle(canvas.RGBA{G: rand.Float32(), A: 1})
			x.BeginPath()
			x.Arc(z, 4, 0, r2.AngleMax)
			x.Fill()
			x.Stroke()
		}
		x.Restore()
	}
}

func (c *myCanvas) Event(x *canvas.Context, p r2.X) {
	fmt.Println(c.id, p)
}

type myCanvasPage struct {
	canvas.Page
}

func (p *myCanvasPage) PageBody(path string, d *Doc) BodyVec {
	s := &standardBody{BodyVec: p.Body(path, d)}
	return s.Body(d)
}

func (p *myCanvasPage) Body(path string, d *Doc) BodyVec {
	cs := []myCanvas{
		{id: "canvas1", greeting: "Hello 1"},
		{id: "canvas2", greeting: "Hello 2"},
	}
	bv := BodyVec{d.H2("Canvas")}
	for i := range cs {
		p.Page.SetDrawer(cs[i].id, &cs[i])
		p.Page.SetListener(cs[i].id, &cs[i])
		bv = append(bv, d.Div(&cs[i]))
	}
	return bv
}

type myElog struct {
	view        elog.View
	tb          elog.TimeBounds
	max, margin r2.X
	id          string
}

func (c *myElog) Body(d *Doc) BodyVec {
	return BodyVec{
		d.Canvas("#"+c.id, "aspect=1.5"),
	}
}

type myElogPage struct {
	canvas.Page
}

func (p *myElogPage) PageBody(path string, d *Doc) BodyVec {
	s := &standardBody{BodyVec: p.Body(path, d)}
	return s.Body(d)
}

func (p *myElogPage) Body(path string, d *Doc) BodyVec {
	cs := []myElog{
		{id: "elog_canvas1"},
	}
	bv := BodyVec{d.H2("Event Log")}
	for i := range cs {
		p.Page.SetDrawer(cs[i].id, &cs[i])
		p.Page.SetListener(cs[i].id, &cs[i])
		bv = append(bv, d.Div(&cs[i]))
	}
	return bv
}
