// autogenerated: do not edit!
// generated from gentemplate [gentemplate -d Package=html -id Col -d Type=Col github.com/platinasystems/elib/vec.tmpl]

package html

import (
	. "github.com/platinasystems/elib"
)

type ColVec []Col

func (p *ColVec) Resize(n uint) {
	c := Index(cap(*p))
	l := Index(len(*p)) + Index(n)
	if l > c {
		c = NextResizeCap(l)
		q := make([]Col, l, c)
		copy(q, *p)
		*p = q
	}
	*p = (*p)[:l]
}

func (p *ColVec) Validate(i uint) {
	c := Index(cap(*p))
	l := Index(i) + 1
	if l > c {
		c = NextResizeCap(l)
		q := make([]Col, l, c)
		copy(q, *p)
		*p = q
	}
	if l > Index(len(*p)) {
		*p = (*p)[:l]
	}
}
