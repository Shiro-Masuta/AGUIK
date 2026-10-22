package imageview

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

func (v *ImageView) updateClick(gtx layout.Context) {
	if v.Disabled {
		return
	}
	for {
		e, ok := v.click.Update(gtx.Source)
		if !ok {
			return
		}
		if e.Kind == gesture.KindClick {
			if i := v.HitTest(e.Position); i >= 0 {
				v.clicked = i
			}
		}
	}
}

func (v *ImageView) registerClickArea(gtx layout.Context, size image.Point) {
	if !v.Disabled {
		area := clip.Rect{Max: size}.Push(gtx.Ops)
		v.click.Add(gtx.Ops)
		area.Pop()
	}
}
