package html

import (
	"fmt"
)

type FormNode interface {
	BodyNode
	formNode()
}

type Form struct {
	Flow
	Action URI
	FormMethod
}

func (n *Form) blockNode() {}
func (n *Form) bodyNode()  {}
func (n *Form) node()      {}

type FormMethod int

const (
	FormMethodNone FormMethod = iota
	FormGet
	FormPost
)

var formMethodStrings = []string{
	FormGet:  "GET",
	FormPost: "POST",
}

var formMethodMap = map[string]FormMethod{
	"GET":  FormGet,
	"POST": FormPost,
}

func (n *Form) Markup(d *Doc) string {
	a, sep := n.Attrs.String(d)

	if n.FormMethod != FormMethodNone {
		a += fmt.Sprintf("%smethod=\"%s\"", sep, formMethodStrings[n.FormMethod])
		sep = " "
	}
	if len(n.Action) > 0 {
		a += fmt.Sprintf("%saction=\"%s\"", sep, n.Action)
		sep = " "
	}
	return wrap("form", a, n.X.Markup(d))
}

func (n *Form) attrs() *Attrs    { return &n.Attrs }
func (n *Form) bodyVec() BodyVec { return n.Flow.bodyvec() }

func (d *Doc) Form(args ...interface{}) (n *Form) {
	n = &Form{}
	d.flow(n, &n.X, &n.Attrs, args...)
	var ok bool
	var u string

	if u, ok = n.Attrs.user["method"]; ok {
		n.FormMethod = formMethodMap[u]
		delete(n.Attrs.user, "method")
	}

	if u, ok = n.Attrs.user["action"]; ok {
		n.Action = URI(u)
		delete(n.Attrs.user, "action")
	}

	return
}

// Label
type Label struct {
	inline
	For string
}

func (n *Label) bodyNode()   {}
func (n *Label) inlineNode() {}
func (n *Label) node()       {}

func (n *Label) Markup(d *Doc) string {
	a, sep := n.Attrs.String(d)
	if len(n.For) > 0 {
		a += fmt.Sprintf("%sfor=\"%s\"", sep, n.For)
		sep = " "
	}
	return wrap("label", a, n.X.Markup(d))
}

func (n *Label) attrs() *Attrs    { return &n.Attrs }
func (n *Label) bodyVec() BodyVec { return n.inline.bodyvec() }

func (d *Doc) Label(args ...interface{}) (n *Label) {
	n = &Label{}
	d.inline(n, &n.X, &n.Attrs, args...)
	var ok bool
	var v string
	if v, ok = n.Attrs.user["for"]; ok {
		n.For = v
		delete(n.Attrs.user, "for")
	}
	return
}

type Input struct {
	Attrs
	InputType
}

type InputType int

const (
	Text InputType = iota + 1
	Password
	Checkbox
	Radio
	Submit
	Reset
	File
	Hidden
	Image
	Button
	Number
)

var inputTypeStrings = []string{
	Text:     "text",
	Password: "password",
	Checkbox: "checkbox",
	Radio:    "radio",
	Submit:   "submit",
	Reset:    "reset",
	File:     "file",
	Hidden:   "hidden",
	Image:    "image",
	Button:   "button",
	Number:   "number",
}

var inputTypeMap = map[string]InputType{
	"text":     Text,
	"password": Password,
	"checkbox": Checkbox,
	"radio":    Radio,
	"submit":   Submit,
	"reset":    Reset,
	"file":     File,
	"hidden":   Hidden,
	"image":    Image,
	"button":   Button,
	"number":   Number,
}

func (n *Input) formNode()   {}
func (n *Input) bodyNode()   {}
func (n *Input) inlineNode() {}
func (n *Input) node()       {}

func (n *Input) Markup(d *Doc) string {
	a, sep := n.Attrs.String(d)
	a += fmt.Sprintf("%stype=\"%s\"", sep, inputTypeStrings[n.InputType])
	return fmt.Sprintf("<input %s/>", a)
}

func (n *Input) attrs() *Attrs    { return &n.Attrs }
func (n *Input) bodyVec() BodyVec { return BodyVec{} }

func (d *Doc) Input(args ...interface{}) (n *Input) {
	n = &Input{}
	d.addAttrs(n, &n.Attrs, args...)
	var ok bool
	var u string
	if u, ok = n.Attrs.user["type"]; ok {
		n.InputType = inputTypeMap[u]
		delete(n.Attrs.user, "type")
	} else {
		n.InputType = Text
	}
	return
}

type Select struct {
	Attrs
	Options []OptionNode
}

func (n *Select) formNode()   {}
func (n *Select) bodyNode()   {}
func (n *Select) inlineNode() {}
func (n *Select) node()       {}

func (n *Select) Markup(d *Doc) string {
	a, _ := n.Attrs.String(d)
	options := ""
	for _, o := range n.Options {
		options += o.Markup(d)
	}
	return wrap("select", a, options)
}

func (n *Select) attrs() *Attrs    { return &n.Attrs }
func (n *Select) bodyVec() BodyVec { return BodyVec{} }

func (d *Doc) Select(args ...interface{}) (n *Select) {
	n = &Select{}
	for _, a := range args {
		switch v := a.(type) {
		case string:
			if !d.addAttr(n, &n.Attrs, v) {
				panic(v)
			}

		case OptionNode:
			n.Options = append(n.Options, v)

		default:
			panic(v)
		}
	}
	return
}

type OptionNode interface {
	FormNode
	optionNode()
}

type Option struct {
	Attrs
	Selected bool
	Disabled bool
	Value    string
}

func (n *Option) optionNode() {}
func (n *Option) formNode()   {}
func (n *Option) bodyNode()   {}
func (n *Option) node()       {}

func (n *Option) Markup(d *Doc) string {
	a, _ := n.Attrs.String(d)
	return wrap("option", a, n.Value)
}

func (n *Option) attrs() *Attrs    { return &n.Attrs }
func (n *Option) bodyVec() BodyVec { return BodyVec{} }

func (d *Doc) Option(args ...interface{}) OptionNode {
	n := &Option{}
	n.Value = d.addAttrs(n, &n.Attrs, args...)
	return n
}

type Textarea struct {
	Attrs
	Content string
}

func (n *Textarea) formNode()   {}
func (n *Textarea) bodyNode()   {}
func (n *Textarea) inlineNode() {}
func (n *Textarea) node()       {}

func (n *Textarea) Markup(d *Doc) string {
	a, _ := n.Attrs.String(d)
	return wrap("textarea", a, n.Content)
}

func (n *Textarea) attrs() *Attrs    { return &n.Attrs }
func (n *Textarea) bodyVec() BodyVec { return BodyVec{} }

func (d *Doc) Textarea(args ...interface{}) (n *Textarea) {
	n = &Textarea{}
	for _, a := range args {
		switch v := a.(type) {
		case string:
			if !d.addAttr(n, &n.Attrs, v) {
				panic(v)
			}

		default:
			panic(v)
		}
	}
	return
}

type Fieldset struct {
	Attrs
}

func (n *Fieldset) blockNode() {}
func (n *Fieldset) bodyNode()  {}
func (n *Fieldset) node()      {}

func (n *Fieldset) attrs() *Attrs    { return &n.Attrs }
func (n *Fieldset) bodyVec() BodyVec { return BodyVec{} }
