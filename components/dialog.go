package components

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type DialogAction struct {
	Text    string
	Variant ButtonVariant
}

const defaultDialogWidthDp = 320

type Dialog struct {
	Open                   bool
	Title                  string
	Message                string
	Actions                []DialogAction
	Width                  float32
	DismissOnBackdropClick bool
	Clicked                int
	AutoClose              bool

	overlay      OverlayState
	actionBtns   []widget.Clickable
	cardBlockBtn widget.Clickable
}

func (d *Dialog) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	d.Clicked = -1

	d.actionBtns = glayout.Resize(d.actionBtns, len(d.Actions))

	widthDp := d.Width
	if widthDp == 0 {
		widthDp = defaultDialogWidthDp
	}
	cardW := gtx.Dp(unit.Dp(widthDp))

	return Overlay(gtx, &d.overlay, &d.Open, d.DismissOnBackdropClick, func(gtx layout.Context) layout.Dimensions {
		return d.layoutCard(gtx, th, cardW)
	})
}

func (d *Dialog) layoutCard(gtx layout.Context, th *theme.Theme, cardW int) layout.Dimensions {
	pad := gtx.Dp(unit.Dp(th.Spacing.LG))
	gap := gtx.Dp(unit.Dp(th.Spacing.SM))
	borderW := gtx.Dp(unit.Dp(th.Shape.BorderWidth))
	contentW := cardW - 2*borderW - 2*pad
	if contentW < 0 {
		contentW = 0
	}

	tgtx := gtx
	tgtx.Constraints = layout.Constraints{Min: image.Point{}, Max: image.Pt(contentW, gtx.Constraints.Max.Y)}

	var titleW, msgW layout.Widget
	var titleDims, msgDims layout.Dimensions
	if d.Title != "" {
		titleW = Text(th, d.Title, WithSize(th.Typography.FontSize*1.1), WithWeight(font.Bold))
		m := op.Record(tgtx.Ops)
		titleDims = titleW(tgtx)
		m.Stop()
	}
	if d.Message != "" {
		msgW = Text(th, d.Message, WithSize(th.Typography.FontSize*0.95), WithColor(th.Colors.Muted))
		m := op.Record(tgtx.Ops)
		msgDims = msgW(tgtx)
		m.Stop()
	}
	textH := titleDims.Size.Y + msgDims.Size.Y
	if d.Title != "" && d.Message != "" {
		textH += gap
	}

	type recordedAction struct {
		call op.CallOp
		dims layout.Dimensions
	}
	actions := make([]recordedAction, len(d.Actions))
	actionsW, actionsH := 0, 0
	for i, act := range d.Actions {
		clicked := d.actionBtns[i].Clicked(gtx)
		bg, fg := variantColors(th, act.Variant)
		hovered := d.actionBtns[i].Hovered()
		pressed := d.actionBtns[i].Pressed()
		content := Text(th, act.Text, WithColor(fg), WithWeight(font.SemiBold))

		m := op.Record(gtx.Ops)
		dims := layoutButtonFrame(gtx, th, content, paddingFor(th, SizeMD), bg, hovered, pressed, false, &d.actionBtns[i])
		call := m.Stop()

		actions[i] = recordedAction{call: call, dims: dims}
		if clicked {
			d.Clicked = i
			interaction.Invalidate(gtx)
			if d.AutoClose {
				d.Open = false
			}
		}
		if i > 0 {
			actionsW += gap
		}
		actionsW += dims.Size.X
		if dims.Size.Y > actionsH {
			actionsH = dims.Size.Y
		}
	}

	innerH := 2*pad + textH
	if len(actions) > 0 {
		innerH += gap + actionsH
	}
	innerW := cardW - 2*borderW

	frame := NeoFrame(gtx, th, innerW, innerH, th.Colors.Surface, false, false)

	blockArea := clip.Rect(image.Rectangle{Max: frame.CanvasDim}).Push(gtx.Ops)
	d.cardBlockBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: frame.CanvasDim}
	})
	blockArea.Pop()

	origin := frame.BgRect.Min.Add(image.Pt(pad, pad))
	ty := 0
	if d.Title != "" {
		off := op.Offset(origin.Add(image.Pt(0, ty))).Push(gtx.Ops)
		titleW(tgtx)
		off.Pop()
		ty += titleDims.Size.Y + gap
	}
	if d.Message != "" {
		off := op.Offset(origin.Add(image.Pt(0, ty))).Push(gtx.Ops)
		msgW(tgtx)
		off.Pop()
		ty += msgDims.Size.Y
	}

	if len(actions) > 0 {
		ty += gap
		ax := contentW - actionsW
		if ax < 0 {
			ax = 0
		}
		for _, a := range actions {
			off := op.Offset(origin.Add(image.Pt(ax, ty))).Push(gtx.Ops)
			a.call.Add(gtx.Ops)
			off.Pop()
			ax += a.dims.Size.X + gap
		}
	}

	return layout.Dimensions{Size: frame.CanvasDim}
}
