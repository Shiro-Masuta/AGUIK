package imageview

import (
	"image"
	"image/color"
	"reflect"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

func (v *ImageView) ensureOp() {
	if v.Src == nil {
		return
	}
	srcType := reflect.TypeOf(v.Src)
	if srcType.Comparable() && v.cachedSrc != nil && v.cachedSrc == v.Src {
		return
	}
	v.op = paint.NewImageOp(v.Src)
	if srcType.Comparable() {
		v.cachedSrc = v.Src
	} else {
		v.cachedSrc = nil
	}
}

func (v *ImageView) drawImage(gtx layout.Context, size image.Point) {
	viewClip := clip.Rect{Max: size}.Push(gtx.Ops)
	tr := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(v.scaleX, v.scaleY)).
		Offset(v.offset)).Push(gtx.Ops)
	v.op.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	tr.Pop()
	viewClip.Pop()
}

func (v *ImageView) drawHighlights(gtx layout.Context, th *theme.Theme, size image.Point) {
	viewClip := clip.Rect{Max: size}.Push(gtx.Ops)
	for i, h := range v.Highlights {
		r := v.toScreen(h.Rect)
		style := v.resolveStyle(th, i, h)

		if style.Fill != (color.NRGBA{}) {
			paint.FillShape(gtx.Ops, style.Fill, clip.Rect(r).Op())
		}
		if style.Stroke != (color.NRGBA{}) && style.StrokeWidth > 0 {
			paint.FillShape(gtx.Ops, style.Stroke, clip.Stroke{
				Path:  clip.Rect(r).Path(),
				Width: style.StrokeWidth,
			}.Op())
		}
	}
	viewClip.Pop()
}
