package interaction

import (
	"image"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/widget"
)

func ClearFocus(gtx layout.Context) {
	gtx.Execute(key.FocusCmd{Tag: nil})
}

type FocusClearer struct {
	clickable widget.Clickable
}

func (f *FocusClearer) Layout(gtx layout.Context) layout.Dimensions {
	if f.clickable.Clicked(gtx) {
		ClearFocus(gtx)
	}

	size := gtx.Constraints.Max
	area := clip.Rect{Max: size}.Push(gtx.Ops)
	f.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}
	})
	area.Pop()
	return layout.Dimensions{Size: size}
}
