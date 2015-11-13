//+build js

package canvas

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"weeb/r2"
)

type CompositeOperation int

const (
	SrcOver CompositeOperation = iota // A over B (default)
	SrcAtop                           // A atop B
	SrcIn                             // A in B
	SrcOut                            // A out B
	DstOver                           // B over A
	DstAtop                           // B atop A
	DstIn                             // B in A
	DstOut                            // B out A
	Lighter                           // A plus B
	Copy                              // A (B is ignored)
	Xor                               // A xor B
)

var compositionStrings = []string{
	SrcOver: "source-over",
	SrcAtop: "source-atop",
	SrcIn:   "source-in",
	SrcOut:  "source-out",
	DstOver: "destination-over",
	DstAtop: "destination-atop",
	DstIn:   "destination-in",
	DstOut:  "destination-out",
	Lighter: "lighter",
	Copy:    "copy",
	Xor:     "xor",
}

type Context struct {
	*js.Object
	Font                     string  `js:"font"` // "10px sans-serif"
	GlobalAlpha              float64 `js:"globalAlpha"`
	GlobalCompositeOperation string  `js:"globalCompositeOperation"` // "source-over"
	ImageSmoothingEnabled    bool    `js:"imageSmoothingEnabled"`
	LineCap                  string  `js:"lineCap"` // "butt" round square
	LineDashOffset           float64 `js:"lineDashOffset"`
	LineJoin                 string  `js:"lineJoin"` // "miter" round miter
	LineWidth                float64 `js:"lineWidth"`
	MiterLimit               float64 `js:"miterLimit"`
	ShadowBlur               float64 `js:"shadowBlur"`
	ShadowColor              string  `js:"shadowColor"`
	ShadowOffsetX            float64 `js:"shadowOffsetX"`
	ShadowOffsetY            float64 `js:"shadowOffsetY"`
	FillStyle                string  `js:"fillStyle"` // css color, CanvasGradient, CanvasPattern
	StrokeStyle              string  `js:"strokeStyle"`
	TextAlign                string  `js:"textAlign"`    // "start" (default), "end", "left", "right", "center"
	TextBaseline             string  `js:"textBaseline"` // "top", "hanging", "middle", "alphabetic" (default), "ideographic", "bottom"
	Size                     r2.X
}

// (x,y) coordinate as complex number for easy arithmetic.
// X increasing to the right; Y increasing down on screen.
type Xy complex128

func GetContext(c interface{}) *Context {
	return &Context{Object: jquery.NewJQuery(c).Underlying().Index(0).Call("getContext", "2d")}
}

// Push/pop graphics state on stack.
func (c *Context) Save()    { c.Call("save") }
func (c *Context) Restore() { c.Call("restore") }

func (c *Context) BeginPath() { c.Call("beginPath") }
func (c *Context) ClosePath() { c.Call("closePath") }
func (c *Context) Fill()      { c.Call("fill") }
func (c *Context) Stroke()    { c.Call("stroke") }
func (c *Context) Clip()      { c.Call("clip") }

func (c *Context) IsPointInPath(x r2.X) { c.Call("isPointInPath", x.X(), x.Y()) }

func (c *Context) FillRect(x, s r2.X)   { c.Call("fillRect", x.X(), x.Y(), s.X(), s.Y()) }
func (c *Context) StrokeRect(x, s r2.X) { c.Call("strokeRect", x.X(), x.Y(), s.X(), s.Y()) }
func (c *Context) ClearRect(x, s r2.X)  { c.Call("clearRect", x.X(), x.Y(), s.X(), s.Y()) }

func (c *Context) FillText(text string, x r2.X, maxWidth ...float64) {
	if len(maxWidth) > 0 {
		c.Call("fillText", text, x.X(), x.Y(), maxWidth)
	} else {
		c.Call("fillText", text, x.X(), x.Y())
	}
}
func (c *Context) StrokeText(text string, x r2.X, maxWidth ...float64) {
	if len(maxWidth) > 0 {
		c.Call("strokeText", text, x.X(), x.Y(), maxWidth)
	} else {
		c.Call("strokeText", text, x.X(), x.Y())
	}
}

func (c *Context) MeasureText(text string) r2.X {
	dx := c.Call("measureText", text).Get("width").Float()
	return r2.X(complex(dx, 0))
}

func (c *Context) MoveTo(x r2.X)  { c.Call("moveTo", x.X(), x.Y()) }
func (c *Context) LineTo(x r2.X)  { c.Call("lineTo", x.X(), x.Y()) }
func (c *Context) Rect(x, s r2.X) { c.Call("rect", x.X(), x.Y(), s.X(), s.Y()) }

func (c *Context) QuadraticCurveTo(c1, x1 r2.X) {
	c.Call("quadraticCurveTo", c1.X(), c1.Y(), x1.X(), x1.Y())
}

func (c *Context) BezierCurveTo(c1, c2, x1 r2.X) {
	c.Call("bezierCurveTo", c1.X(), c1.Y(), c2.X(), c2.Y(), x1.X(), x1.Y())
}

func (c *Context) Arc(x r2.X, r float64, θ0, θ1 r2.Angle, ccw ...bool) {
	c.Call("arc", x.X(), x.Y(), r, θ0.Radians(), θ1.Radians(), ccw)
}

func (c *Context) Ellipse(x, r r2.X, θ0, θ1, rotation r2.Angle, ccw ...bool) {
	c.Call("ellipse", x.X(), x.Y(), r.X(), r.Y(), rotation.Radians(), θ0.Radians(), θ1.Radians(), ccw)
}

func (c *Context) ArcTo(x1, x2 r2.X, r float64) { c.Call("arcTo", x1.X(), x1.Y(), x2.X(), x2.Y(), r) }

// Transforms
func (c *Context) Scale(x r2.X)      { c.Call("scale", x.X(), x.Y()) }
func (c *Context) Translate(x r2.X)  { c.Call("translate", x.X(), x.Y()) }
func (c *Context) Rotate(θ r2.Angle) { c.Call("rotate", θ.Radians()) }

// Applies to current transform.
func (c *Context) Transform(m00, m01, m10, m11 float64, dx r2.X) {
	c.Call("transform", m00, m01, m10, m11, dx.X(), dx.Y())
}

// Resets transform to DX 0 and identity matrix then applies given transform.
func (c *Context) SetTransform(m00, m01, m10, m11 float64, dx r2.X) {
	c.Call("setTransform", m00, m01, m10, m11, dx.X(), dx.Y())
}

type RGBA struct {
	R, G, B, A float32
}

// Saturate from 0 to 255
func sat(x float32) int {
	i := int(x * 256)
	switch {
	case i < 0:
		i = 0
	case i > 255:
		i = 255
	}
	return i
}

// For {Fill,Stroke}Style
func (c *Context) Style(x interface{}) string {
	var s string
	switch v := x.(type) {
	case RGBA:
		s = fmt.Sprintf("rgba(%d,%d,%d,%f)", sat(v.R), sat(v.G), sat(v.B), v.A)
	case string:
		s = v
	default:
		panic(v)
	}
	return s
}

func (c *Context) SetFillStyle(x interface{}) {
	c.FillStyle = c.Style(x)
}

func (c *Context) SetStrokeStyle(x interface{}) {
	c.StrokeStyle = c.Style(x)
}

// Drawer is an interface for types which know how to draw with a Context.
type Drawer interface {
	Draw(c *Context)
}

type Listener interface {
	Event(c *Context, x r2.X)
}

type Interface interface {
	Drawer(id string) (f Drawer, ok bool)
	Listener(id string) (f Listener, ok bool)
}

type Page struct {
	Interface
	DrawerById   map[string]Drawer
	ListenerById map[string]Listener
}

func (p *Page) Drawer(id string) (d Drawer, ok bool) {
	d, ok = p.DrawerById[id]
	return
}

func (p *Page) SetDrawer(id string, d Drawer) {
	if p.DrawerById == nil {
		p.DrawerById = make(map[string]Drawer)
	}
	p.DrawerById[id] = d
}

func (p *Page) Listener(id string) (d Listener, ok bool) {
	d, ok = p.ListenerById[id]
	return
}

func (p *Page) SetListener(id string, d Listener) {
	if p.ListenerById == nil {
		p.ListenerById = make(map[string]Listener)
	}
	p.ListenerById[id] = d
}

/*
Not yet:

createImageData: createImageData()
drawImage: drawImage()
getImageData: getImageData()
putImageData: putImageData()

createLinearGradient: createLinearGradient()
createPattern: createPattern()
createRadialGradient: createRadialGradient()

drawFocusIfNeeded: drawFocusIfNeeded()
getContextAttributes: getContextAttributes()
*/
