//+build js

package main

import (
	"github.com/platinasystems/elib/elog"
	"github.com/platinasystems/weeb/canvas"
	"github.com/platinasystems/weeb/r2"

	"fmt"
)

func (c *myElog) Draw(ctx *canvas.Context) {
	log := &elog.Log{}
	err := rpc.Call("Listener.GetEventLog", &struct{}{}, log)
	if err != nil {
		panic(err)
	}

	if false {
		ctx.Save()
		lw := r2.XY1(2)
		ctx.LineWidth = lw.X()
		ctx.SetStrokeStyle(canvas.RGBA{A: 1})
		ctx.SetFillStyle(canvas.RGBA{G: 1, A: 1})
		ctx.RoundedRect(lw/2, ctx.Size-lw, 8)
		ctx.Fill()
		ctx.Stroke()
		ctx.Restore()
	}

	tb := log.TimeBounds()

	// Switch to Cartesian coordinates where (0,0) is at the lower left and S is the upper right.
	margin := r2.XY(12., 24.)
	max := ctx.Size - 2*margin
	ctx.Transform(1, 0, 0, 1, r2.XY(margin.X(), ctx.Size.Y()-margin.Y()))

	/* Draw and label time axis. */
	{
		ctx.Save()

		ctx.SetStrokeStyle(canvas.RGBA{A: 1})
		ctx.LineWidth = 1

		ctx.BeginPath()
		ctx.MoveTo(0)
		ctx.LineTo(r2.XY(max.X(), 0))
		ctx.Stroke()

		ctx.Font = "10px Ariel"
		ctx.TextBaseline = "top"
		ctx.TextAlign = "center"
		nUnit := int(.5 + (tb.Max-tb.Min)/(tb.Round*tb.Unit))
		x, dx := r2.XY1(0), r2.XY(max.X()/float64(nUnit), 0)
		t, dt := tb.Min, tb.Round*tb.Unit
		for i := 0; i <= nUnit; i++ {
			str := fmt.Sprintf("%.0f", t/tb.Unit)
			dy := r2.XY(0, 2)
			ctx.FillText(str, x+dy)
			ctx.BeginPath()
			ctx.MoveTo(x + dy)
			ctx.LineTo(x - dy)
			ctx.Stroke()
			x += dx
			t += dt
		}
		ctx.Font = "11px Ariel"
		label := fmt.Sprintf("Time in %s from %s", tb.UnitName, tb.Start.Format("2006-01-02 15:04:05"))
		ctx.FillText(label, r2.XY(max.X()/2, 14))
		ctx.Restore()
	}

	{
		ctx.Save()
		ctx.SetStrokeStyle(canvas.RGBA{A: 1})
		ctx.SetFillStyle(canvas.RGBA{A: .8, G: 1})
		log.ForeachEvent(func(e *elog.Event) {
			t := e.Time(log).Sub(tb.Start).Seconds()
			ctx.BeginPath()
			x := r2.XY(max.X()*(t-tb.Min)/(tb.Max-tb.Min), -max.Y()/2)
			ctx.Arc(x, 4, 0, r2.AngleMax)
			ctx.Fill()
			ctx.Stroke()
		})
		ctx.Restore()
	}
}

func (c *myElog) Event(x *canvas.Context, p r2.X) {
	fmt.Println(c.id, p)
}
