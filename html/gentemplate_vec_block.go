// autogenerated: do not edit!
// generated from gentemplate [gentemplate -d Package=html -id Block -d Type=BlockNode github.com/platinasystems/elib/vec.tmpl]

package html

import (
	. "github.com/platinasystems/elib"
)

type BlockVec []BlockNode

func (p *BlockVec) Resize(n uint) {
	c := Index(cap(*p))
	l := Index(len(*p)) + Index(n)
	if l > c {
		c = NextResizeCap(l)
		q := make([]BlockNode, l, c)
		copy(q, *p)
		*p = q
	}
	*p = (*p)[:l]
}

func (p *BlockVec) Validate(i uint) {
	c := Index(cap(*p))
	l := Index(i) + 1
	if l > c {
		c = NextResizeCap(l)
		q := make([]BlockNode, l, c)
		copy(q, *p)
		*p = q
	}
	if l > Index(len(*p)) {
		*p = (*p)[:l]
	}
}
