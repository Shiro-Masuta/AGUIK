package components

import (
	"image"

	"gioui.org/io/input"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type SidebarItem struct {
	Icon    *widget.Icon
	Tooltip string
}

type Sidebar struct {
	Items   []SidebarItem
	Active  int
	Variant ButtonVariant
	Size    ButtonSize

	btns []IconButton
}

func (s *Sidebar) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if s.Active < 0 || s.Active >= len(s.Items) {
		s.Active = 0
	}
	if len(s.btns) != len(s.Items) {
		btns := make([]IconButton, len(s.Items))
		copy(btns, s.btns)
		s.btns = btns
	}

	pad := gtx.Dp(unit.Dp(th.Spacing.SM))

	children := make([]glayout.Child, len(s.Items))
	for i, item := range s.Items {
		i := i
		btn := &s.btns[i]
		btn.Icon = item.Icon
		btn.Size = s.Size
		btn.TooltipText = item.Tooltip
		btn.TooltipPlacement = TooltipRight
		if i == s.Active {
			btn.Variant = s.Variant
		} else {
			btn.Variant = ButtonDefault
		}

		children[i] = glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
			dims := btn.Layout(gtx, th)
			if btn.Clicked && s.Active != i {
				s.Active = i
				interaction.Invalidate(gtx)
			}
			return dims
		})
	}
	content := glayout.Column(th.Spacing.SM, layout.Start, children...)

	// Mesure la largeur : passe "à blanc", sans Source d'événements pour
	// qu'un clic en attente ne soit pas consommé ici plutôt que pendant
	// le vrai passage de dessin plus bas.
	mgtx := gtx
	mgtx.Constraints.Min = image.Point{}
	mgtx.Source = input.Source{}
	macro := op.Record(mgtx.Ops)
	contentDims := content(mgtx)
	macro.Stop()

	innerW := contentDims.Size.X + 2*pad
	innerH := gtx.Constraints.Max.Y - 2*pad
	if innerH < contentDims.Size.Y {
		innerH = contentDims.Size.Y
	}

	frame := NeoFrame(gtx, th, innerW, innerH, th.Colors.Surface, false, false)

	off := frame.BgRect.Min.Add(image.Pt(pad, pad))
	stack := op.Offset(off).Push(gtx.Ops)
	content(gtx) // seule passe qui reçoit vraiment les événements
	stack.Pop()

	return layout.Dimensions{Size: frame.CanvasDim}
}
