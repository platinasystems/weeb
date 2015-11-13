// Inline level html nodes.
package html

import (
	"fmt"
)

type InlineNode interface {
	BodyNode
	inlineNode()
}

//go:generate gentemplate -d Package=html -id Inline -d Type=InlineNode github.com/platinasystems/elib/vec.tmpl

func (n *InlineVec) Markup(d *Doc) string {
	s := ""
	for _, f := range *n {
		s += f.Markup(d)
	}
	return s
}

type inline struct {
	Attrs
	X InlineVec
}

func (d *Doc) inline(n BodyNode, px *InlineVec, attrs *Attrs, args ...interface{}) {
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
		case InlineNode:
			x = append(x, v)

		case InlineVec:
			x = append(x, v...)

		case []InlineNode:
			x = append(x, v...)

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

func (n *inline) markup(d *Doc, tag string) string {
	a, _ := n.Attrs.String(d)
	return wrap(tag, a, n.X.Markup(d))
}

func (n *inline) bodyvec() (r BodyVec) {
	for i := range n.X {
		r = append(r, n.X[i])
	}
	return
}

// Anchor
type A struct {
	inline
	Href URI
}

func (n *A) bodyNode()   {}
func (n *A) inlineNode() {}
func (n *A) node()       {}

func (n *A) Markup(d *Doc) string {
	a, sep := n.Attrs.String(d)
	if len(n.Href) > 0 {
		a += fmt.Sprintf("%shref=\"%s\"", sep, n.Href)
		sep = " "
	}
	return wrap("a", a, n.X.Markup(d))
}

func (n *A) attrs() *Attrs    { return &n.Attrs }
func (n *A) bodyVec() BodyVec { return n.inline.bodyvec() }

func (d *Doc) A(args ...interface{}) (n *A) {
	n = &A{}
	d.inline(n, &n.X, &n.Attrs, args...)
	var ok bool
	var u string
	if u, ok = n.Attrs.user["href"]; ok {
		n.Href = URI(u)
		delete(n.Attrs.user, "href")
	}
	return
}

type Span inline

func (n *Span) bodyNode()   {}
func (n *Span) inlineNode() {}
func (n *Span) node()       {}

func (n *Span) Markup(d *Doc) string { return (*inline)(n).markup(d, "span") }
func (n *Span) attrs() *Attrs        { return &n.Attrs }
func (n *Span) bodyVec() BodyVec     { return (*inline)(n).bodyvec() }
func (d *Doc) Span(args ...interface{}) (n *Span) {
	n = &Span{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}

// <!ENTITY % fontstyle
//  "TT | I | B | BIG | SMALL">

type TT inline
type I inline
type B inline
type Big inline
type Small inline

func (n *TT) bodyNode()      {}
func (n *TT) inlineNode()    {}
func (n *TT) node()          {}
func (n *I) bodyNode()       {}
func (n *I) inlineNode()     {}
func (n *I) node()           {}
func (n *B) bodyNode()       {}
func (n *B) inlineNode()     {}
func (n *B) node()           {}
func (n *Big) bodyNode()     {}
func (n *Big) inlineNode()   {}
func (n *Big) node()         {}
func (n *Small) bodyNode()   {}
func (n *Small) inlineNode() {}
func (n *Small) node()       {}

func (n *TT) Markup(d *Doc) string    { return (*inline)(n).markup(d, "tt") }
func (n *I) Markup(d *Doc) string     { return (*inline)(n).markup(d, "i") }
func (n *B) Markup(d *Doc) string     { return (*inline)(n).markup(d, "b") }
func (n *Big) Markup(d *Doc) string   { return (*inline)(n).markup(d, "big") }
func (n *Small) Markup(d *Doc) string { return (*inline)(n).markup(d, "small") }

func (n *TT) attrs() *Attrs    { return &n.Attrs }
func (n *I) attrs() *Attrs     { return &n.Attrs }
func (n *B) attrs() *Attrs     { return &n.Attrs }
func (n *Big) attrs() *Attrs   { return &n.Attrs }
func (n *Small) attrs() *Attrs { return &n.Attrs }

func (n *TT) bodyVec() BodyVec    { return (*inline)(n).bodyvec() }
func (n *I) bodyVec() BodyVec     { return (*inline)(n).bodyvec() }
func (n *B) bodyVec() BodyVec     { return (*inline)(n).bodyvec() }
func (n *Big) bodyVec() BodyVec   { return (*inline)(n).bodyvec() }
func (n *Small) bodyVec() BodyVec { return (*inline)(n).bodyvec() }

func (d *Doc) TT(args ...interface{}) (n *TT) {
	n = &TT{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) I(args ...interface{}) (n *I) {
	n = &I{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) B(args ...interface{}) (n *B) {
	n = &B{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) Big(args ...interface{}) (n *Big) {
	n = &Big{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) Small(args ...interface{}) (n *Small) {
	n = &Small{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}

// <!ENTITY % phrase "EM | STRONG | DFN | CODE |
//                    SAMP | KBD | VAR | CITE | ABBR | ACRONYM" >

// Forced line break.
type BR Attrs

func (n *BR) bodyNode()   {}
func (n *BR) inlineNode() {}
func (n *BR) node()       {}

func (n *BR) Markup(d *Doc) string { return (*Attrs)(n).markup(d, "br") }
func (n *BR) attrs() *Attrs        { return (*Attrs)(n) }
func (n *BR) bodyVec() BodyVec     { return BodyVec{} }
func (d *Doc) BR(args ...interface{}) (n *BR) {
	n = &BR{}
	d.attrsOnly(n, (*Attrs)(n), args...)
	return
}
