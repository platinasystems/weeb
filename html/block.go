// Block level html nodes.
package html

import (
	"fmt"
)

// All HTML elements (Nodes) are capable of generating markup as a string.
type Node interface {
	Markup(d *Doc) string
}

type Body interface {
	Body(d *Doc) BodyVec
}

// A Node which can be inside a document body (between <body> and </body>)
type BodyNode interface {
	attrs() *Attrs
	bodyVec() BodyVec
	Node
	bodyNode()
}

// HTML has two basic content models:
//   inline     character level elements and text strings
//   block      block-like elements e.g. paragraphs and lists
type BlockNode interface {
	BodyNode
	blockNode()
}

//go:generate gentemplate -d Package=html -id Body -d Type=BodyNode github.com/platinasystems/elib/vec.tmpl
//go:generate gentemplate -d Package=html -id Block -d Type=BlockNode github.com/platinasystems/elib/vec.tmpl

func (n *BodyVec) Markup(d *Doc) string {
	s := ""
	for _, f := range *n {
		s += f.Markup(d)
	}
	return s
}

func (n *BlockVec) Markup(d *Doc) string {
	s := ""
	for _, f := range *n {
		s += f.Markup(d)
	}
	return s
}

type Flow struct {
	Attrs
	X BodyVec
}

func (d *Doc) flow(n BodyNode, px *BodyVec, attrs *Attrs, args ...interface{}) {
	sep := ""
	x := *px
	for _, a := range args {
		switch v := a.(type) {
		case string:
			if !d.addAttr(n, attrs, v) {
				x = append(x, &String{sep + v})
				sep = " "
			}

		case int:
			x = append(x, &String{fmt.Sprintf("%d", v)})
			sep = " "

		case BodyNode:
			x = append(x, v)

		case BodyVec:
			x = append(x, v...)

		case Body:
			bs := v.Body(d)
			if isEventListener(a) {
				for i := range bs {
					attrs := bs[i].attrs()
					d.addEventListener(attrs, a)
					d.addAttrId(attrs, bs[i])
				}
			}
			x = append(x, bs...)

		default:
			panic(v)
		}
	}

	*px = x
}

func (n *Flow) markup(d *Doc, tag string) string {
	a, _ := n.Attrs.String(d)
	return wrap(tag, a, n.X.Markup(d))
}

func (n *Flow) Markup(d *Doc) string {
	return n.X.Markup(d)
}

func (n *Flow) bodyvec() (r BodyVec) {
	for i := range n.X {
		r = append(r, n.X[i])
	}
	return
}

type block struct {
	Attrs
	X BlockVec
}

func (d *Doc) block(n BodyNode, px *BlockVec, attrs *Attrs, args ...interface{}) {
	sep := ""
	x := *px
	var ls []interface{}
	for _, a := range args {
		var isEvent bool
		if isEvent = isEventListener(a); isEvent {
			ls = append(ls, a)
		}
		switch v := a.(type) {
		case string:
			if !d.addAttr(n, attrs, v) {
				x = append(x, &String{sep + v})
				sep = " "
			}
		case BlockNode:
			x = append(x, v)

		default:
			if !isEvent {
				panic(v)
			}
		}
	}

	if len(ls) > 0 {
		d.addEventListener(attrs, ls)
	}

	*px = x
}

func (n *block) markup(d *Doc, tag string) string {
	a, _ := n.Attrs.String(d)
	return wrap(tag, a, n.X.Markup(d))
}

func (n *block) bodyvec() (r BodyVec) {
	for i := range n.X {
		r = append(r, n.X[i])
	}
	return
}

func (n *Attrs) markup(d *Doc, tag string) string {
	a, _ := n.String(d)
	return fmt.Sprintf("<%s %s>", tag, a)
}

func (n *Attrs) markupEndTag(d *Doc, tag string) string {
	a, _ := n.String(d)
	return fmt.Sprintf("<%s %s></%s>", tag, a, tag)
}

func (d *Doc) attrsOnly(n BodyNode, attrs *Attrs, args ...interface{}) {
	d.addAttrs(n, attrs, args...)
}

type String struct {
	X string
}

func (n *String) bodyNode()   {}
func (n *String) blockNode()  {}
func (n *String) inlineNode() {}
func (n *String) node()       {}

func (n *String) attrs() *Attrs    { return &Attrs{} }
func (n *String) bodyVec() BodyVec { return BodyVec{} }

func (n *String) Markup(d *Doc) string {
	return n.X
}

// Paragraphs
type P inline

func (n *P) blockNode() {}
func (n *P) bodyNode()  {}
func (n *P) node()      {}

func (n *P) Markup(d *Doc) string { return (*inline)(n).markup(d, "p") }

func (n *P) attrs() *Attrs    { return &n.Attrs }
func (n *P) bodyVec() BodyVec { return (*inline)(n).bodyvec() }

func (d *Doc) P(args ...interface{}) (n *P) {
	n = &P{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}

// Headings
type H1 inline
type H2 inline
type H3 inline
type H4 inline
type H5 inline
type H6 inline

func (n *H1) bodyNode()  {}
func (n *H1) blockNode() {}
func (n *H1) node()      {}
func (n *H2) bodyNode()  {}
func (n *H2) blockNode() {}
func (n *H2) node()      {}
func (n *H3) bodyNode()  {}
func (n *H3) blockNode() {}
func (n *H3) node()      {}
func (n *H4) bodyNode()  {}
func (n *H4) blockNode() {}
func (n *H4) node()      {}
func (n *H5) bodyNode()  {}
func (n *H5) blockNode() {}
func (n *H5) node()      {}
func (n *H6) bodyNode()  {}
func (n *H6) blockNode() {}
func (n *H6) node()      {}

func (n *H1) Markup(d *Doc) string { return (*inline)(n).markup(d, "h1") }
func (n *H2) Markup(d *Doc) string { return (*inline)(n).markup(d, "h2") }
func (n *H3) Markup(d *Doc) string { return (*inline)(n).markup(d, "h3") }
func (n *H4) Markup(d *Doc) string { return (*inline)(n).markup(d, "h4") }
func (n *H5) Markup(d *Doc) string { return (*inline)(n).markup(d, "h5") }
func (n *H6) Markup(d *Doc) string { return (*inline)(n).markup(d, "h6") }

func (n *H1) attrs() *Attrs { return &n.Attrs }
func (n *H2) attrs() *Attrs { return &n.Attrs }
func (n *H3) attrs() *Attrs { return &n.Attrs }
func (n *H4) attrs() *Attrs { return &n.Attrs }
func (n *H5) attrs() *Attrs { return &n.Attrs }
func (n *H6) attrs() *Attrs { return &n.Attrs }

func (n *H1) bodyVec() BodyVec { return (*inline)(n).bodyvec() }
func (n *H2) bodyVec() BodyVec { return (*inline)(n).bodyvec() }
func (n *H3) bodyVec() BodyVec { return (*inline)(n).bodyvec() }
func (n *H4) bodyVec() BodyVec { return (*inline)(n).bodyvec() }
func (n *H5) bodyVec() BodyVec { return (*inline)(n).bodyvec() }
func (n *H6) bodyVec() BodyVec { return (*inline)(n).bodyvec() }

func (d *Doc) H1(args ...interface{}) (n *H1) {
	n = &H1{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) H2(args ...interface{}) (n *H2) {
	n = &H2{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) H3(args ...interface{}) (n *H3) {
	n = &H3{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) H4(args ...interface{}) (n *H4) {
	n = &H4{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) H5(args ...interface{}) (n *H5) {
	n = &H5{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) H6(args ...interface{}) (n *H6) {
	n = &H6{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}

// Lists
type UL Flow
type OL Flow
type LI Flow

func (n *UL) blockNode() {}
func (n *UL) bodyNode()  {}
func (n *UL) node()      {}
func (n *OL) blockNode() {}
func (n *OL) bodyNode()  {}
func (n *OL) node()      {}
func (n *LI) blockNode() {}
func (n *LI) bodyNode()  {}
func (n *LI) node()      {}

func (n *UL) Markup(d *Doc) string { return (*Flow)(n).markup(d, "ul") }
func (n *OL) Markup(d *Doc) string { return (*Flow)(n).markup(d, "ol") }
func (n *LI) Markup(d *Doc) string { return (*Flow)(n).markup(d, "li") }

func (n *UL) attrs() *Attrs { return &n.Attrs }
func (n *OL) attrs() *Attrs { return &n.Attrs }
func (n *LI) attrs() *Attrs { return &n.Attrs }

func (n *UL) bodyVec() BodyVec { return (*Flow)(n).bodyvec() }
func (n *OL) bodyVec() BodyVec { return (*Flow)(n).bodyvec() }
func (n *LI) bodyVec() BodyVec { return (*Flow)(n).bodyvec() }

func (d *Doc) UL(args ...interface{}) (n *UL) {
	n = &UL{}
	d.flow(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) OL(args ...interface{}) (n *OL) {
	n = &OL{}
	d.flow(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) LI(args ...interface{}) (n *LI) {
	n = &LI{}
	d.flow(n, &n.X, &n.Attrs, args...)
	return
}

// Preformatted text
type Pre inline

func (n *Pre) blockNode() {}
func (n *Pre) bodyNode()  {}
func (n *Pre) node()      {}

func (n *Pre) Markup(d *Doc) string { return (*inline)(n).markup(d, "pre") }
func (n *Pre) attrs() *Attrs        { return &n.Attrs }
func (n *Pre) bodyVec() BodyVec     { return (*inline)(n).bodyvec() }

func (d *Doc) Pre(args ...interface{}) (n *Pre) {
	n = &Pre{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}

// Definition lists: DT for term, DD for its definition.
type DL struct {
	Attrs
	X []DLNode
}

type DLNode interface {
	BodyNode
	dlNode()
}

type DT struct {
	Attrs
	inline
}

type DD struct {
	Flow
	inline
}

func (n *DD) bodyNode()  {}
func (n *DD) dlNode()    {}
func (n *DD) node()      {}
func (n *DL) blockNode() {}
func (n *DL) bodyNode()  {}
func (n *DL) node()      {}
func (n *DT) bodyNode()  {}
func (n *DT) dlNode()    {}
func (n *DT) node()      {}

type Div Flow

func (n *Div) blockNode() {}
func (n *Div) bodyNode()  {}
func (n *Div) node()      {}

func (n *Div) Markup(d *Doc) string { return (*Flow)(n).markup(d, "div") }
func (n *Div) attrs() *Attrs        { return &n.Attrs }
func (n *Div) bodyVec() BodyVec     { return (*Flow)(n).bodyvec() }
func (d *Doc) Div(args ...interface{}) (n *Div) {
	n = &Div{}
	d.flow(n, &n.X, &n.Attrs, args...)
	return
}

type Nav Flow

func (n *Nav) blockNode() {}
func (n *Nav) bodyNode()  {}
func (n *Nav) node()      {}

func (n *Nav) Markup(d *Doc) string { return (*Flow)(n).markup(d, "nav") }
func (n *Nav) attrs() *Attrs        { return &n.Attrs }
func (n *Nav) bodyVec() BodyVec     { return (*Flow)(n).bodyvec() }
func (d *Doc) Nav(args ...interface{}) (n *Nav) {
	n = &Nav{}
	d.flow(n, &n.X, &n.Attrs, args...)
	return
}

type Section Flow

func (n *Section) blockNode() {}
func (n *Section) bodyNode()  {}
func (n *Section) node()      {}

func (n *Section) Markup(d *Doc) string { return (*Flow)(n).markup(d, "section") }
func (n *Section) attrs() *Attrs        { return &n.Attrs }
func (n *Section) bodyVec() BodyVec     { return (*Flow)(n).bodyvec() }
func (d *Doc) Section(args ...interface{}) (n *Section) {
	n = &Section{}
	d.flow(n, &n.X, &n.Attrs, args...)
	return
}

type Blockquote block

func (n *Blockquote) blockNode() {}
func (n *Blockquote) bodyNode()  {}
func (n *Blockquote) node()      {}

func (n *Blockquote) Markup(d *Doc) string { return (*block)(n).markup(d, "blockquote") }
func (n *Blockquote) attrs() *Attrs        { return &n.Attrs }
func (n *Blockquote) bodyVec() BodyVec     { return (*block)(n).bodyvec() }
func (d *Doc) Blockquote(args ...interface{}) (n *Blockquote) {
	n = &Blockquote{}
	d.block(n, &n.X, &n.Attrs, args...)
	return
}

type HR Attrs

func (n *HR) blockNode() {}
func (n *HR) bodyNode()  {}
func (n *HR) node()      {}

func (n *HR) Markup(d *Doc) string { return (*Attrs)(n).markup(d, "hr") }
func (n *HR) attrs() *Attrs        { return (*Attrs)(n) }
func (n *HR) bodyVec() BodyVec     { return BodyVec{} }
func (d *Doc) HR(args ...interface{}) (n *HR) {
	n = &HR{}
	d.attrsOnly(n, (*Attrs)(n), args...)
	return
}

// Information on author.
type Address inline

func (n *Address) blockNode() {}
func (n *Address) bodyNode()  {}
func (n *Address) node()      {}

func (n *Address) Markup(d *Doc) string { return (*inline)(n).markup(d, "address") }
func (n *Address) attrs() *Attrs        { return &n.Attrs }
func (n *Address) bodyVec() BodyVec     { return (*inline)(n).bodyvec() }
func (d *Doc) Address(args ...interface{}) (n *Address) {
	n = &Address{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}

// HTML5 canvas
type Canvas Attrs

func (n *Canvas) blockNode() {}
func (n *Canvas) bodyNode()  {}
func (n *Canvas) node()      {}

func (n *Canvas) Markup(d *Doc) string { return (*Attrs)(n).markupEndTag(d, "canvas") }
func (n *Canvas) attrs() *Attrs        { return (*Attrs)(n) }
func (n *Canvas) bodyVec() BodyVec     { return BodyVec{} }
func (d *Doc) Canvas(args ...interface{}) (n *Canvas) {
	n = &Canvas{}
	d.attrsOnly(n, (*Attrs)(n), args...)
	return
}
