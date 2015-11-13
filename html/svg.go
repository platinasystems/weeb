package html

type svgNode interface {
	BodyNode
	svgNode()
}

//go:generate gentemplate -d Package=html -id svgNode -d Type=svgNode github.com/platinasystems/elib/vec.tmpl

func (n *svgNodeVec) Markup(d *Doc) string {
	s := ""
	for _, f := range *n {
		s += f.Markup(d)
	}
	return s
}

type Svg struct {
	Attrs
	X svgNodeVec
}

func (n *Svg) blockNode() {}
func (n *Svg) bodyNode()  {}
func (n *Svg) node()      {}

func (n *Svg) Markup(d *Doc) string {
	a, _ := n.Attrs.String(d)
	return wrap("svg", a, n.X.Markup(d))
}

func (n *Svg) attrs() *Attrs { return &n.Attrs }

func (n *Svg) bodyVec() (r BodyVec) {
	for i := range n.X {
		r = append(r, n.X[i])
	}
	return
}

func (d *Doc) Svg(args ...interface{}) *Svg {
	n := &Svg{}
	var ls []interface{}
	for _, a := range args {
		switch v := a.(type) {
		case string:
			d.addAttrForce(n, &n.Attrs, v)

		case svgNode:
			n.X = append(n.X, v)

		case []svgNode:
			n.X = append(n.X, v...)

		default:
			if isEventListener(a) {
				ls = append(ls, a)
			} else {
				panic(v)
			}
		}
	}

	if len(ls) > 0 {
		d.addEventListener(&n.Attrs, ls)
	}

	return n
}

type Circle Attrs

func (n *Circle) svgNode()  {}
func (n *Circle) bodyNode() {}
func (n *Circle) node()     {}

func (n *Circle) Markup(d *Doc) string { return (*Attrs)(n).markup(d, "circle") }
func (n *Circle) attrs() *Attrs        { return (*Attrs)(n) }
func (n *Circle) bodyVec() BodyVec     { return BodyVec{} }
func (d *Doc) Circle(args ...string) *Circle {
	n := &Circle{}
	for _, a := range args {
		d.addAttrForce(n, (*Attrs)(n), a)
	}
	return n
}
