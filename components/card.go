package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

type cardStyle struct {
	padding    float32
	background color.NRGBA
	fillHeight bool
	fillWidth  bool
}

type CardOption func(*cardStyle)

func WithCardPadding(dp float32) CardOption {
	return func(c *cardStyle) { c.padding = dp }
}

func WithCardBackground(bg color.NRGBA) CardOption {
	return func(c *cardStyle) { c.background = bg }
}

func WithCardFillHeight() CardOption {
	return func(c *cardStyle) { c.fillHeight = true }
}

func WithCardFillWidth() CardOption {
	return func(c *cardStyle) { c.fillWidth = true }
}

// Card enveloppe un contenu quelconque dans le cadre néobrutaliste (bordure,
// ombre, fond) du Theme. Par défaut, padding et fond viennent du Theme (auto) ;
// chaque option surcharge un champ précis, exactement comme Text.
//
//	components.Card(th, content)                                    // 100% auto
//	components.Card(th, content, components.WithCardBackground(...))
func Card(th *theme.Theme, content layout.Widget, opts ...CardOption) layout.Widget {
	cs := cardStyle{
		padding:    th.Spacing.MD,
		background: th.Colors.Surface,
	}
	for _, opt := range opts {
		opt(&cs)
	}

	return func(gtx layout.Context) layout.Dimensions {
		pad := gtx.Dp(unit.Dp(cs.padding))
		overhead := NeoFrameOverhead(gtx, th)

		macro := op.Record(gtx.Ops)
		cgtx := gtx
		cgtx.Constraints.Min = image.Point{}
		if maxX := gtx.Constraints.Max.X - 2*pad - overhead.X; maxX > 0 {
			cgtx.Constraints.Max.X = maxX
		}
		if maxY := gtx.Constraints.Max.Y - 2*pad - overhead.Y; maxY > 0 {
			cgtx.Constraints.Max.Y = maxY
		}
		contentDims := content(cgtx)
		call := macro.Stop()

		innerW := contentDims.Size.X + 2*pad
		innerH := contentDims.Size.Y + 2*pad

		if cs.fillWidth && gtx.Constraints.Max.X > overhead.X {
			innerW = gtx.Constraints.Max.X - overhead.X
		}
		if cs.fillHeight && gtx.Constraints.Max.Y > overhead.Y {
			innerH = gtx.Constraints.Max.Y - overhead.Y
		}
		frame := NeoFrame(gtx, th, innerW, innerH, cs.background, false, false)

		off := frame.BgRect.Min.Add(image.Pt(pad, pad))
		stack := op.Offset(off).Push(gtx.Ops)
		call.Add(gtx.Ops)
		stack.Pop()

		return layout.Dimensions{Size: frame.CanvasDim}
	}
}
