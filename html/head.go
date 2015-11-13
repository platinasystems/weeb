package html

import (
	"fmt"
)

type HeadNode interface {
	Node
	headNode()
}

type Title struct {
	X string
}

func (n *Title) headNode() {}
func (n *Title) node()     {}

func (n *Title) Markup(d *Doc) string {
	return wrap("title", "", n.X)
}

// Document base.
type Base struct {
	URI
}

func (n *Base) headNode() {}
func (n *Base) node()     {}

func (n *Base) Markup(d *Doc) (v string) {
	return fmt.Sprintf("<base href=\"%s\"/>", n.URI)
}

// A character encoding, as per [RFC2045].
type Charset string

// Media type, as per [RFC2045].
type ContentType string

// A Uniform Resource Identifier.
type URI string

// Generic metainformation.
type Meta struct {
	Charset

	I18NAttrs

	// HTTP response header name
	HttpEquiv string
	Name      string
	Content   string
	Scheme    string
}

func (n *Meta) headNode() {}
func (n *Meta) node()     {}

func (n *Meta) Markup(d *Doc) (v string) {
	a := ""
	switch {
	case len(n.Charset) != 0:
		a += fmt.Sprintf("charset=\"%s\"", n.Charset)
	case len(n.HttpEquiv) != 0:
		a += fmt.Sprintf("http-equiv=\"%s\" content=\"%s\"", n.HttpEquiv, n.Content)
	default:
		a += fmt.Sprintf("name=\"%s\" content=\"%s\"", n.Name, n.Content)
	}
	v = fmt.Sprintf("<meta %s/>", a)
	return
}

type Link struct {
	Charset
	Href URI
	Type ContentType
	Rel  string
}

func (n *Link) headNode() {}
func (n *Link) node()     {}

func (n *Link) Markup(d *Doc) (v string) {
	return fmt.Sprintf("<link rel=\"%s\" type=\"%s\" href=\"%s\"/>", n.Rel, n.Type, n.Href)
}

type Style struct {
	// Content type of style language (e.g. text/css).
	ContentType
	Title string // Advisory title
}

func (n *Style) headNode() {}
func (n *Style) node()     {}

type Script struct {
	Charset
	Type    ContentType
	Src     URI
	Async   bool
	Defer   bool
	Content string
}

func (n *Script) headNode() {}
func (n *Script) bodyNode() {}
func (n *Script) node()     {}

func (n *Script) Markup(d *Doc) (v string) {
	a := ""
	sep := ""
	if len(n.Type) != 0 {
		a += fmt.Sprintf("%stype=\"%s\"", sep, n.Type)
		sep = " "
	}
	if len(n.Src) != 0 {
		a += fmt.Sprintf("%ssrc=\"%s\"", sep, n.Src)
		sep = " "
	}
	if n.Defer {
		a += " defer"
	}
	if n.Async {
		a += " async"
	}
	return wrap("script", a, n.Content)
}

func (n *Script) attrs() *Attrs    { return &Attrs{} }
func (n *Script) bodyVec() BodyVec { return BodyVec{} }
