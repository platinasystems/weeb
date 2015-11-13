package html

import (
	"fmt"
	"github.com/platinasystems/elib"
)

type Doc struct {
	Head               []HeadNode
	Body               BodyVec
	ClassByID          elib.StringPool
	ClassMap           map[string]ClassID
	BodyNodeById       map[string]BodyNode
	nAssignedIds       int
	EventListenersById map[string][]interface{}
}

func wrap(tag, attrs, content string) string {
	sep := ""
	if len(attrs) > 0 {
		sep = " "
	}
	return fmt.Sprintf("<%s%s%s>%s</%s>", tag, sep, attrs, content, tag)
}

func (d *Doc) Reset() { d.nAssignedIds = 0 }

func (d *Doc) Markup() string {
	s := fmt.Sprintf("<!DOCTYPE html>")
	s += fmt.Sprintf("<html>")

	s += fmt.Sprintf("<head>")
	for _, n := range d.Head {
		s += n.Markup(d)
	}
	s += fmt.Sprintf("</head>")

	s += fmt.Sprintf("<body>")
	for _, n := range d.Body {
		s += n.Markup(d)
	}
	s += fmt.Sprintf("</body>")

	s += fmt.Sprintf("</html>")
	return s
}
