package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

const defaultInputWidthDp = 200

type Input struct {
	Editor      widget.Editor
	Placeholder string
	Width       float32
	Size        ButtonSize
	Disabled    bool
	Password    bool // masque la saisie (•••) si true
}

func (in *Input) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	in.Editor.SingleLine = true
	in.Editor.ReadOnly = in.Disabled
	if in.Password {
		in.Editor.Mask = '•' // gioui.org/widget.Editor expose Mask rune pour ça
	} else {
		in.Editor.Mask = 0
	}

	if !in.Disabled {
		for {
			if _, ok := in.Editor.Update(gtx); !ok {
				break
			}
		}
	}

	bg, fg := th.Colors.Surface, th.Colors.Foreground
	if in.Disabled {
		bg, fg = th.Colors.Muted, th.Colors.Foreground
	}
	focused := gtx.Focused(&in.Editor) && !in.Disabled

	pad := gtx.Dp(unit.Dp(paddingFor(th, in.Size)))
	width := in.Width
	if width == 0 {
		width = defaultInputWidthDp
	}
	fieldW := gtx.Dp(unit.Dp(width))
	lineH := gtx.Dp(unit.Dp(th.Typography.FontSize * 1.3))

	ed := material.Editor(&material.Theme{Shaper: th.Shaper}, &in.Editor, in.Placeholder)
	ed.Color = fg
	ed.HintColor = th.Colors.Muted
	ed.SelectionColor = th.Colors.Primary
	ed.TextSize = unit.Sp(th.Typography.FontSize)

	return layoutInputFrame(gtx, th, ed, inputDims{width: fieldW, height: lineH}, pad, bg, focused)
}
