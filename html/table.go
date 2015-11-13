package html

type tableNode interface {
	BodyNode
	tableNode()
}

//go:generate gentemplate -d Package=html -id tableNode -d Type=tableNode github.com/platinasystems/elib/vec.tmpl

func (n *tableNodeVec) Markup(d *Doc) string {
	s := ""
	for _, f := range *n {
		s += f.Markup(d)
	}
	return s
}

type Table struct {
	Attrs
	X tableNodeVec
}

func (n *Table) blockNode() {}
func (n *Table) bodyNode()  {}
func (n *Table) node()      {}

func (n *Table) Markup(d *Doc) string {
	a, _ := n.Attrs.String(d)
	return wrap("table", a, n.X.Markup(d))
}

func (n *Table) attrs() *Attrs { return &n.Attrs }

func (n *Table) bodyVec() (r BodyVec) {
	for i := range n.X {
		r = append(r, n.X[i])
	}
	return
}

func (d *Doc) Table(args ...interface{}) *Table {
	n := &Table{}
	for _, a := range args {
		switch v := a.(type) {
		case string:
			d.addAttrForce(n, &n.Attrs, v)

		case tableNode:
			n.X = append(n.X, v)

		case []tableNode:
			n.X = append(n.X, v...)

		default:
			panic(v)
		}
	}
	return n
}

// Table caption.
type Caption inline

func (n *Caption) tableNode() {}
func (n *Caption) bodyNode()  {}
func (n *Caption) node()      {}

func (n *Caption) Markup(d *Doc) string { return (*inline)(n).markup(d, "caption") }
func (n *Caption) attrs() *Attrs        { return &n.Attrs }
func (n *Caption) bodyVec() BodyVec     { return (*inline)(n).bodyvec() }
func (d *Doc) Caption(args ...interface{}) (n *Caption) {
	n = &Caption{}
	d.inline(n, &n.X, &n.Attrs, args...)
	return
}

type Colgroup struct {
	Attrs
	X ColVec
}

func (n *Colgroup) tableNode() {}
func (n *Colgroup) bodyNode()  {}
func (n *Colgroup) node()      {}

func (n *Colgroup) Markup(d *Doc) string {
	a, _ := n.Attrs.String(d)
	s := ""
	for _, f := range n.X {
		s += f.Markup(d)
	}
	return wrap("colgroup", a, s)
}

func (n *Colgroup) attrs() *Attrs { return &n.Attrs }
func (n *Colgroup) bodyVec() (r BodyVec) {
	for i := range n.X {
		r = append(r, &n.X[i])
	}
	return
}

func (d *Doc) Colgroup(args ...interface{}) *Colgroup {
	n := &Colgroup{}
	for _, a := range args {
		switch v := a.(type) {
		case Col:
			n.X = append(n.X, v)

		case *Col:
			n.X = append(n.X, *v)

		case []Col:
			n.X = append(n.X, v...)

		default:
			panic(v)
		}
	}
	return n
}

type Col Attrs

func (c *Col) node()     {}
func (c *Col) bodyNode() {}

func (n *Col) Markup(d *Doc) string { return (*Attrs)(n).markup(d, "col") }
func (n *Col) attrs() *Attrs        { return (*Attrs)(n) }
func (n *Col) bodyVec() BodyVec     { return BodyVec{} }
func (d *Doc) Col(args ...string) *Col {
	n := &Col{}
	for _, a := range args {
		d.addAttrForce(n, (*Attrs)(n), a)
	}
	return n
}

//go:generate gentemplate -d Package=html -id Col -d Type=Col github.com/platinasystems/elib/vec.tmpl

type tableRows struct {
	Attrs
	X tableRowVec
}

func (n *tableRows) markup(d *Doc, tag string) string {
	a, _ := n.Attrs.String(d)
	s := ""
	for _, f := range n.X {
		s += f.Markup(d)
	}
	return wrap(tag, a, s)
}

func (n *tableRows) bodyvec() (r BodyVec) {
	for i := range n.X {
		r = append(r, &n.X[i])
	}
	return
}

func (d *Doc) tableRows(n BodyNode, px *tableRowVec, attrs *Attrs, args ...interface{}) {
	x := *px
	for _, a := range args {
		switch v := a.(type) {
		case string:
			d.addAttrForce(n, attrs, v)

		case TR:
			x = append(x, v)

		case *TR:
			x = append(x, *v)

		case []TR:
			x = append(x, v...)

		default:
			panic(v)
		}
	}
	*px = x
}

type tableRowNode interface {
	BodyNode
	tableRowNode()
}

type TR struct {
	Attrs
	X tableRowNodeVec
}

//go:generate gentemplate -d Package=html -id tableRow -d Type=TR github.com/platinasystems/elib/vec.tmpl
//go:generate gentemplate -d Package=html -id tableRowNode -d Type=tableRowNode github.com/platinasystems/elib/vec.tmpl

// Table header, body and footer.
type Thead tableRows
type Tbody tableRows
type Tfoot tableRows

func (n *Thead) tableNode() {}
func (n *Thead) bodyNode()  {}
func (n *Thead) node()      {}
func (n *Tbody) tableNode() {}
func (n *Tbody) bodyNode()  {}
func (n *Tbody) node()      {}
func (n *Tfoot) tableNode() {}
func (n *Tfoot) bodyNode()  {}
func (n *Tfoot) node()      {}

func (n *Thead) Markup(d *Doc) string { return (*tableRows)(n).markup(d, "thead") }
func (n *Tbody) Markup(d *Doc) string { return (*tableRows)(n).markup(d, "tbody") }
func (n *Tfoot) Markup(d *Doc) string { return (*tableRows)(n).markup(d, "tfoot") }

func (n *Thead) attrs() *Attrs { return &n.Attrs }
func (n *Tbody) attrs() *Attrs { return &n.Attrs }
func (n *Tfoot) attrs() *Attrs { return &n.Attrs }

func (n *Thead) bodyVec() BodyVec { return (*tableRows)(n).bodyvec() }
func (n *Tbody) bodyVec() BodyVec { return (*tableRows)(n).bodyvec() }
func (n *Tfoot) bodyVec() BodyVec { return (*tableRows)(n).bodyvec() }

func (d *Doc) Thead(args ...interface{}) (n *Thead) {
	n = &Thead{}
	d.tableRows(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) Tbody(args ...interface{}) (n *Tbody) {
	n = &Tbody{}
	d.tableRows(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) Tfoot(args ...interface{}) (n *Tfoot) {
	n = &Tfoot{}
	d.tableRows(n, &n.X, &n.Attrs, args...)
	return
}

func (n *TR) tableRowNode() {}
func (n *TR) bodyNode()     {}
func (n *TR) node()         {}

func (n *TR) Markup(d *Doc) string {
	a, _ := n.Attrs.String(d)
	s := ""
	for _, f := range n.X {
		s += f.Markup(d)
	}
	return wrap("tr", a, s)
}

func (n *TR) attrs() *Attrs { return &n.Attrs }
func (n *TR) bodyVec() (r BodyVec) {
	for i := range n.X {
		r = append(r, n.X[i])
	}
	return
}

func (d *Doc) TR(args ...interface{}) *TR {
	n := &TR{}
	for _, a := range args {
		switch v := a.(type) {
		case string:
			d.addAttrForce(n, &n.Attrs, v)

		case tableRowNode:
			n.X = append(n.X, v)

		case *tableRowNode:
			n.X = append(n.X, *v)

		case []tableRowNode:
			n.X = append(n.X, v...)

		default:
			panic(v)
		}
	}
	return n
}

type TH Flow
type TD Flow

func (n *TH) tableRowNode() {}
func (n *TH) bodyNode()     {}
func (n *TH) node()         {}
func (n *TD) tableRowNode() {}
func (n *TD) bodyNode()     {}
func (n *TD) node()         {}

func (n *TH) Markup(d *Doc) string { return (*Flow)(n).markup(d, "th") }
func (n *TD) Markup(d *Doc) string { return (*Flow)(n).markup(d, "td") }

func (n *TH) attrs() *Attrs { return &n.Attrs }
func (n *TD) attrs() *Attrs { return &n.Attrs }

func (n *TH) bodyVec() BodyVec { return (*Flow)(n).bodyvec() }
func (n *TD) bodyVec() BodyVec { return (*Flow)(n).bodyvec() }

func (d *Doc) TH(args ...interface{}) (n *TH) {
	n = &TH{}
	d.flow(n, &n.X, &n.Attrs, args...)
	return
}
func (d *Doc) TD(args ...interface{}) (n *TD) {
	n = &TD{}
	d.flow(n, &n.X, &n.Attrs, args...)
	return
}
