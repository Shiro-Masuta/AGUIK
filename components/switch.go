package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type Switch struct {
	Bool     widget.Bool
	Variant  ButtonVariant
	Size     ButtonSize
	Disabled bool
}

func switchDimsFor(s ButtonSize) (trackWDp, trackHDp float32) {
	switch s {
	case SizeMD:
		return 44, 24
	case SizeLG:
		return 52, 28
	default:
		return 36, 20
	}
}

func (s *Switch) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if !s.Disabled {
		s.Bool.Update(gtx)
	}

	trackWDp, trackHDp := switchDimsFor(s.Size)
	trackW := gtx.Dp(unit.Dp(trackWDp))
	trackH := gtx.Dp(unit.Dp(trackHDp))
	borderW := gtx.Dp(unit.Dp(th.Shape.BorderWidth))

	onColor, _ := toggleCheckedColors(th, s.Variant)

	bg := th.Colors.Surface
	if s.Bool.Value {
		bg = onColor
	}
	if s.Disabled {
		bg = th.Colors.Muted
	}

	hovered := s.Bool.Hovered() && !s.Disabled
	pressed := s.Bool.Pressed() && !s.Disabled

	innerW := trackW - 2*borderW
	innerH := trackH - 2*borderW
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}

	frame := NeoFrame(gtx, th, innerW, innerH, bg, hovered, pressed)

	knobPad := gtx.Dp(unit.Dp(2))
	knobSize := frame.BgRect.Dy() - 2*knobPad
	if knobSize < 0 {
		knobSize = 0
	}

	knobX := frame.BgRect.Min.X + knobPad
	if s.Bool.Value {
		knobX = frame.BgRect.Max.X - knobPad - knobSize
	}
	knobY := frame.BgRect.Min.Y + knobPad

	knobColor := th.Colors.Muted
	if s.Bool.Value {
		knobColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	}
	if s.Disabled {
		knobColor = th.Colors.Foreground
	}

	knobBorderW := gtx.Dp(unit.Dp(1))
	knobRadius := gtx.Dp(unit.Dp(th.Shape.Radius))

	knobRect := image.Rectangle{
		Min: image.Pt(knobX, knobY),
		Max: image.Pt(knobX+knobSize, knobY+knobSize),
	}
	outerKnob := clip.RRect{Rect: knobRect, SE: knobRadius, SW: knobRadius, NE: knobRadius, NW: knobRadius}
	paint.FillShape(gtx.Ops, th.Colors.Border, outerKnob.Op(gtx.Ops))

	innerRadius := knobRadius - knobBorderW
	if innerRadius < 0 {
		innerRadius = 0
	}
	innerKnobRect := image.Rectangle{
		Min: knobRect.Min.Add(image.Pt(knobBorderW, knobBorderW)),
		Max: knobRect.Max.Sub(image.Pt(knobBorderW, knobBorderW)),
	}
	innerKnob := clip.RRect{Rect: innerKnobRect, SE: innerRadius, SW: innerRadius, NE: innerRadius, NW: innerRadius}
	paint.FillShape(gtx.Ops, knobColor, innerKnob.Op(gtx.Ops))

	if !s.Disabled {
		interaction.ClickArea(gtx, &s.Bool, frame.RestRect)
	}

	return layout.Dimensions{Size: frame.CanvasDim}
}
