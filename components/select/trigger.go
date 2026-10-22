package selectwidget

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/components"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

func (s *Select) layoutTrigger(gtx layout.Context, th *theme.Theme, label string) components.Frame {

	bg, fg := th.Colors.Surface, th.Colors.Foreground
	if s.Disabled {
		bg, fg = th.Colors.Muted, th.Colors.Foreground
	}
	hovered := !s.Disabled && s.trigger.Hovered()
	// Le style "actif" reflète simplement l'ouverture du dropdown — pas
	// de focus clavier à suivre ici : contrairement à widget.Editor
	// (Input), widget.Clickable ne gère que les clics pointeur, jamais
	// le focus. gtx.Focused ne se déclenche jamais pour lui, ce qui
	// causait l'ouverture/fermeture instantanée qu'on a corrigée.
	active := !s.Disabled && s.open

	pad := gtx.Dp(unit.Dp(selectPaddingFor(th, s.Size)))
	widthDp := s.Width
	if widthDp == 0 {
		widthDp = defaultSelectWidthDp
	}
	fieldW := gtx.Dp(unit.Dp(widthDp))
	maxFieldW := gtx.Constraints.Max.X - components.NeoFrameOverhead(gtx, th).X
	if maxFieldW < 0 {
		maxFieldW = 0
	}
	if fieldW > maxFieldW {
		fieldW = maxFieldW
	}

	content := components.Text(th, label, components.WithColor(fg))

	mgtx := gtx
	mgtx.Constraints.Min = image.Point{}
	mgtx.Constraints.Max.X = fieldW - 2*pad
	if mgtx.Constraints.Max.X < 0 {
		mgtx.Constraints.Max.X = 0
	}
	macro := op.Record(mgtx.Ops)
	contentDims := content(mgtx)
	macro.Stop()

	innerW := fieldW
	innerH := contentDims.Size.Y + 2*pad
	frame := components.NeoFrame(gtx, th, innerW, innerH, bg, hovered, active)

	off := frame.BgRect.Min.Add(image.Pt(pad, pad))
	stack := op.Offset(off).Push(gtx.Ops)
	content(gtx)
	stack.Pop()

	return frame
}

func selectPaddingFor(th *theme.Theme, size components.ButtonSize) float32 {
	switch size {
	case components.SizeMD:
		return th.Spacing.MD
	case components.SizeLG:
		return th.Spacing.LG
	default:
		return th.Spacing.SM
	}
}
