// HTML element attributes (e.g. "id" in <a id="foo"/>)
package html

import (
	"fmt"
	"sort"
	"strings"
)

type ClassID uint32

// Normalized classes are space separated and sorted for unique hashing.
func normalizeClassName(name string) string {
	sep := " "
	if strings.Index(name, sep) < 0 {
		return name
	}
	cs := strings.Split(name, sep)
	sort.Strings(cs)
	return strings.Join(cs, sep)
}

func (d *Doc) ClassByName(name string) (id ClassID) {
	var ok bool

	name = normalizeClassName(name)

	if id, ok = d.ClassMap[name]; ok {
		return
	}
	id = ClassID(len(d.ClassByID.Strings))
	if id == 0 {
		id++
	}
	d.ClassByID.Validate(uint(id))
	d.ClassByID.Strings[id] = name

	if d.ClassMap == nil {
		d.ClassMap = make(map[string]ClassID)
	}
	d.ClassMap[name] = id
	return
}

func (d *Doc) Class(name string) Attrs {
	return Attrs{ClassID: d.ClassByName(name)}
}

// Create unique ID if element does not already have one.
func (d *Doc) assignID(attrs *Attrs) string {
	if len(attrs.ID) == 0 {
		const digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz$_"
		var ds [32]byte
		v := d.nAssignedIds
		d.nAssignedIds++
		i := len(ds)
		for {
			i--
			ds[i] = digits[v&0x3f]
			v >>= 6
			if v == 0 {
				break
			}
		}
		attrs.ID = string(ds[i:])
	}
	return attrs.ID
}

type Attrs struct {
	ID string
	ClassID
	user map[string]string
}

type I18NAttrs struct {
	Language string
}

func (a Attrs) User(n string, vals ...string) Attrs {
	if a.user == nil {
		a.user = make(map[string]string)
	}
	v := ""
	if len(vals) > 0 {
		v = vals[0]
	}
	a.user[n] = v
	return a
}

func (d *Doc) addAttrId(a *Attrs, n BodyNode) {
	if d.BodyNodeById == nil {
		d.BodyNodeById = make(map[string]BodyNode)
	}
	d.BodyNodeById[a.ID] = n
}

func (d *Doc) addAttr2(n BodyNode, a *Attrs, spec string, force bool) bool {
	// .STRING is shorthand for class=STRING but only applies once
	// For example d.Div(".foo", ".bar") => <div class=foo>.bar</div>)
	if spec[0] == '.' && a.ClassID == 0 {
		a.ClassID = d.ClassByName(spec[1:])
		return true
	}

	if spec[0] == '#' && len(a.ID) == 0 {
		a.ID = spec[1:]
		if n != nil {
			d.addAttrId(a, n)
		}
		return true
	}

	// Otherwise look for NAME=VALUE pair.
	if !force && strings.Index(spec, "=") < 0 {
		return false
	}

	e := strings.Split(spec, "=")
	if len(e) > 0 {
		switch e[0] {
		case "class":
			if len(e[1]) > 0 {
				a.ClassID = d.ClassByName(e[1])
			}
		default:
			if len(e) == 1 || len(e[1]) == 0 {
				*a = a.User(e[0])
			} else {
				*a = a.User(e[0], e[1])
			}
		}
	}
	return true
}

func (d *Doc) addAttr(n BodyNode, a *Attrs, spec string) bool {
	return d.addAttr2(n, a, spec, false)
}

func (d *Doc) addAttrForce(n BodyNode, a *Attrs, spec string) bool {
	return d.addAttr2(n, a, spec, true)
}

func (a *Attrs) Set(n BodyNode, d *Doc, spec string) {
	d.addAttrForce(n, a, spec)
}

func (d *Doc) addAttrs(n BodyNode, attrs *Attrs, args ...interface{}) (content string) {
	sep := ""
	var ls []interface{}
	for _, a := range args {
		switch v := a.(type) {
		case string:
			if !d.addAttr(n, attrs, v) {
				content += sep + v
				sep = " "
			}

		default:
			if isEventListener(v) {
				ls = append(ls, v)
			} else {
				panic(v)
			}
		}
	}
	if len(ls) > 0 {
		d.addEventListener(attrs, ls)
	}
	return
}

func (a *Attrs) String(d *Doc) (s, sep string) {
	s = ""
	sep = ""
	if len(a.ID) != 0 {
		s += fmt.Sprintf("%sid=\"%s\"", sep, a.ID)
		sep = " "
	}
	if a.ClassID != 0 {
		s += fmt.Sprintf("%sclass=\"%s\"", sep, d.ClassByID.Strings[a.ClassID])
		sep = " "
	}
	for k, v := range a.user {
		if len(v) > 0 {
			s += fmt.Sprintf("%s%s=\"%s\"", sep, k, v)
		} else {
			s += fmt.Sprintf("%s%s", sep, k)
		}
		sep = " "
	}
	return
}
