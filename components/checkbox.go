package components

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type Checkbox struct {
	Bool     widget.Bool
	Label    string
	Variant  ButtonVariant
	Size     ButtonSize
	Disabled bool
}

func checkboxSizeFor(s ButtonSize) float32 {
	switch s {
	case SizeMD:
		return 20
	case SizeLG:
		return 24
	default:
		return 16
	}
}

func (c *Checkbox) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if !c.Disabled {
		c.Bool.Update(gtx)
	}

	boxPx := gtx.Dp(unit.Dp(checkboxSizeFor(c.Size)))

	checkedBg, checkedFg := toggleCheckedColors(th, c.Variant)

	bg, fg := th.Colors.Surface, checkedFg
	if c.Bool.Value {
		bg = checkedBg
	}
	if c.Disabled {
		bg, fg = th.Colors.Muted, th.Colors.Foreground
	}

	hovered := c.Bool.Hovered() && !c.Disabled
	pressed := c.Bool.Pressed() && !c.Disabled

	frame := NeoFrame(gtx, th, boxPx, boxPx, bg, hovered, pressed)

	if c.Bool.Value {
		drawCheckmark(gtx, frame.BgRect, fg, gtx.Dp(unit.Dp(2)))
	}

	if !c.Disabled {
		interaction.ClickArea(gtx, &c.Bool, frame.RestRect)
	}

	boxSize := frame.CanvasDim
	if c.Label == "" {
		return layout.Dimensions{Size: boxSize}
	}

	pad := gtx.Dp(unit.Dp(th.Spacing.SM))
	lblColor := th.Colors.Foreground
	if c.Disabled {
		lblColor = th.Colors.Muted
	}

	label := Text(th, c.Label, WithColor(lblColor))

	mgtx := gtx
	mgtx.Constraints.Min = image.Point{}
	macro := op.Record(mgtx.Ops)
	lblDims := label(mgtx)
	macro.Stop()

	labelOff := image.Pt(boxSize.X+pad, (boxSize.Y-lblDims.Size.Y)/2)
	labelStack := op.Offset(labelOff).Push(gtx.Ops)
	label(gtx)
	labelStack.Pop()

	total := image.Pt(boxSize.X+pad+lblDims.Size.X, boxSize.Y)
	if lblDims.Size.Y > total.Y {
		total.Y = lblDims.Size.Y
	}
	return layout.Dimensions{Size: total}
}

func drawCheckmark(gtx layout.Context, rect image.Rectangle, col color.NRGBA, strokeW int) {
	w, h := float32(rect.Dx()), float32(rect.Dy())
	ox, oy := float32(rect.Min.X), float32(rect.Min.Y)
	pt := func(fx, fy float32) f32.Point {
		return f32.Point{X: ox + fx*w, Y: oy + fy*h}
	}

	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(pt(0.2, 0.55))
	p.LineTo(pt(0.42, 0.75))
	p.LineTo(pt(0.8, 0.25))
	stroke := clip.Stroke{Path: p.End(), Width: float32(strokeW)}

	paint.FillShape(gtx.Ops, col, stroke.Op())
}
