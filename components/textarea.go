package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

type TextArea struct {
	Editor      widget.Editor
	Placeholder string
	Width       float32
	Lines       int // nombre de lignes visibles ; 0 -> 4 par défaut
	Size        ButtonSize
	Disabled    bool
}

func (ta *TextArea) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	ta.Editor.SingleLine = false
	ta.Editor.ReadOnly = ta.Disabled

	if !ta.Disabled {
		for {
			if _, ok := ta.Editor.Update(gtx); !ok {
				break
			}
		}
	}

	bg, fg := th.Colors.Surface, th.Colors.Foreground
	if ta.Disabled {
		bg, fg = th.Colors.Muted, th.Colors.Foreground
	}
	focused := gtx.Focused(&ta.Editor) && !ta.Disabled

	pad := gtx.Dp(unit.Dp(paddingFor(th, ta.Size)))
	width := ta.Width
	if width == 0 {
		width = defaultInputWidthDp
	}
	fieldW := gtx.Dp(unit.Dp(width))

	lines := ta.Lines
	if lines == 0 {
		lines = 4
	}
	lineH := gtx.Dp(unit.Dp(th.Typography.FontSize * 1.3))

	ed := material.Editor(&material.Theme{Shaper: th.Shaper}, &ta.Editor, ta.Placeholder)
	ed.Color = fg
	ed.HintColor = th.Colors.Muted
	ed.SelectionColor = th.Colors.Primary
	ed.TextSize = unit.Sp(th.Typography.FontSize)

	return layoutInputFrame(gtx, th, ed, inputDims{width: fieldW, height: lineH * lines}, pad, bg, focused)
}
