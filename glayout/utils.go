package glayout

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Resize retourne s avec exactement n éléments, en conservant les existants.
func Resize[E any](s []E, n int) []E {
	if n <= len(s) {
		return s[:n]
	}
	return append(s, make([]E, n-len(s))...)
}

// BlendColor mélange base vers tint (ratio 0..1) et retourne une couleur opaque.
func BlendColor(base, tint color.NRGBA, ratio float32) color.NRGBA {
	lerp := func(a, b uint8) uint8 {
		return uint8(float32(a) + (float32(b)-float32(a))*ratio)
	}
	return color.NRGBA{
		R: lerp(base.R, tint.R),
		G: lerp(base.G, tint.G),
		B: lerp(base.B, tint.B),
		A: 255,
	}
}

// WithAlpha retourne c avec le canal alpha a.
func WithAlpha(c color.NRGBA, a uint8) color.NRGBA {
	c.A = a
	return c
}

// FillRect remplit un rectangle de size à l'origine courante.
func FillRect(gtx layout.Context, col color.NRGBA, size image.Point) {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
