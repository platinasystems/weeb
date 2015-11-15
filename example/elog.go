//+build !js

package main

import (
	"github.com/platinasystems/weeb/canvas"
	"github.com/platinasystems/weeb/r2"
)

func (c *myElog) Draw(x *canvas.Context)          {}
func (c *myElog) Event(x *canvas.Context, p r2.X) {}
