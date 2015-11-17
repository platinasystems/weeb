//+build js

package main

import (
	"github.com/platinasystems/elib/elog"
	"github.com/platinasystems/weeb/canvas"
	"github.com/platinasystems/weeb/r2"

	"fmt"
	"math"
)

func (m *myElog) Draw(ctx *canvas.Context) {
	vw := &m.view
	err := rpc.Call("Listener.GetEventView", &struct{}{}, vw)
	if err != nil {
		panic(err)
	}

	ctx.Save()
	ctx.BeginPath()
	ctx.SetFillStyle(canvas.RGBA{A: 1, R: .95, G: .95, B: .95})
	ctx.Rect(0, ctx.Size)
	ctx.Fill()
	ctx.Restore()

	vw.GetTimeBounds(&m.tb)

	// Switch to Cartesian coordinates where (0,0) is at the lower left and S is the upper right.
	m.margin = r2.XY(16., 24.)
	m.max = ctx.Size - 2*m.margin
	ctx.Translate(r2.XY(m.margin.X(), ctx.Size.Y()-m.margin.Y()))

	/* Draw and label time axis. */
	{
		ctx.Save()

		ctx.SetStrokeStyle(canvas.RGBA{A: 1})
		ctx.LineWidth = 1

		ctx.BeginPath()
		ctx.MoveTo(0)
		ctx.LineTo(r2.XY(m.max.X(), 0))
		ctx.Stroke()

		ctx.Font = "10px Ariel"
		ctx.TextBaseline = "top"
		ctx.TextAlign = "center"
		nUnit := int(.5 + m.tb.Dt/(m.tb.Round*m.tb.Unit))
		x, dx := r2.XY1(0), r2.XY(m.max.X()/float64(nUnit), 0)
		t, dt := m.tb.Min, m.tb.Round*m.tb.Unit
		for i := 0; i <= nUnit; i++ {
			str := fmt.Sprintf("%.0f", t/m.tb.Unit)
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
		label := fmt.Sprintf("Time in %s from %s", m.tb.UnitName, m.tb.Start.Format("2006-01-02 15:04:05"))
		ctx.FillText(label, r2.XY(m.max.X()/2, 14))
		ctx.Restore()
	}

	{
		ctx.Save()
		dy := 7.
		ctx.Font = fmt.Sprintf("%.0fpt monospace", dy)
		ctx.TextAlign = "center"
		ctx.TextBaseline = "middle"
		black := canvas.RGBA{A: 1}
		stroke := canvas.RGBA{A: 1, R: 151 / 256., G: 187 / 256., B: 205 / 256.}
		fill := stroke
		fill.A = .2
		labelStyle := canvas.RGBA{A: 1}
		ctx.SetStrokeStyle(black)
		ctx.SetFillStyle(fill)
		lastXmax := 0.
		idy := 0
		for i := range vw.Events {
			e := &vw.Events[i]
			t := vw.Time(e).Sub(m.tb.Start).Seconds()

			// label := e.Type().Name
			label := e.String()
			if len(label) > 20 {
				label = label[:20]
			}

			dx := 1.1*ctx.MeasureText(label) + r2.XY(0, 1.5*dy)
			x := r2.XY(m.max.X()*(t-m.tb.Min)/m.tb.Dt, -m.max.Y()/2)

			if i > 0 && x.X()-.5*dx.X() < lastXmax {
				idy++
			} else {
				idy = 0
			}

			thisDy := 0.
			if idy != 0 {
				// Even stacked events go down; odd events up.
				if idy%2 == 0 {
					thisDy = float64(idy/2) * dy
				} else {
					thisDy = -float64((1+idy)/2) * dy
				}
			}
			// fmt.Println(idy, thisDy)

			x += r2.XY(0, thisDy*2)
			lastXmax = x.X() + .5*dx.X()

			ctx.CenteredRoundedRect(x, dx, 2)
			ctx.SetFillStyle(fill)
			ctx.Fill()
			ctx.Stroke()
			ctx.SetFillStyle(labelStyle)
			ctx.FillText(label, x)
		}
		ctx.Restore()
	}
}

func (m *myElog) findNearby(tol float64, p r2.X) (es []*elog.Event) {
	px := p.X() - m.margin.X()
	dx := px / m.max.X()
	pt := m.tb.Min + m.tb.Dt*dx
	vw := &m.view
	for ei := range vw.Events {
		e := &vw.Events[ei]
		dt := math.Abs(m.view.Time(e).Sub(m.tb.Start).Seconds() - pt)
		if dt < .05*m.tb.Dt {
			es = append(es, e)
		}
	}
	return
}

func (m *myElog) Event(x *canvas.Context, p r2.X) {
	es := m.findNearby(.05, p)
	fmt.Printf("%d events\n", len(es))
	for _, e := range es {
		fmt.Printf("%s\n", m.view.EventString(e))
	}
}
