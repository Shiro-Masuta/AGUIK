package components

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

const defaultProgressWidthDp = 200

// ProgressBar affiche une progression déterminée (0.0 à 1.0), non
type ProgressBar struct {
	Value   float32 // 0..1, clampé automatiquement
	Variant ButtonVariant
	Width   float32 // dp ; 0 -> defaultProgressWidthDp
	Height  float32 // dp ; 0 -> 12
}

func (p *ProgressBar) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	v := p.Value
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}

	widthDp := p.Width
	if widthDp == 0 {
		widthDp = defaultProgressWidthDp
	}
	heightDp := p.Height
	if heightDp == 0 {
		heightDp = 12
	}

	trackW := gtx.Dp(unit.Dp(widthDp))
	trackH := gtx.Dp(unit.Dp(heightDp))
	borderW := gtx.Dp(unit.Dp(th.Shape.BorderWidth))

	innerW := trackW - 2*borderW
	innerH := trackH - 2*borderW
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}

	fillColor, _ := variantColors(th, p.Variant)
	if p.Variant == ButtonDefault {
		fillColor = th.Colors.Primary
	}

	frame := NeoFrame(gtx, th, innerW, innerH, th.Colors.Surface, false, false)

	fillW := int(float32(frame.BgRect.Dx()) * v)
	if fillW > 0 {
		fillRect := image.Rectangle{
			Min: frame.BgRect.Min,
			Max: image.Pt(frame.BgRect.Min.X+fillW, frame.BgRect.Max.Y),
		}
		radius := gtx.Dp(unit.Dp(th.Shape.Radius))
		fillClip := clip.RRect{Rect: fillRect, SE: radius, SW: radius, NE: radius, NW: radius}
		paint.FillShape(gtx.Ops, fillColor, fillClip.Op(gtx.Ops))
	}

	return layout.Dimensions{Size: frame.CanvasDim}
}
