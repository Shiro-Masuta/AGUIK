package codeeditor

import (
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

// layoutScroll place la gouttière et l'éditeur côte à côte, sans hauteur
// bornée (donc sans scroll interne propre à l'éditeur), puis fait défiler
// ce bloc entier via une unique liste verticale native. Gouttière et texte
// bougent donc toujours ensemble, car ils appartiennent à la même passe de
// layout.
func (c *CodeEditor) layoutScroll(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	c.list.Axis = layout.Vertical
	return c.list.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
		return c.layoutContent(gtx, th)
	})
}

func (c *CodeEditor) layoutContent(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	c.Editor.ReadOnly = c.Disabled
	c.Editor.SingleLine = false
	c.Editor.Submit = false

	textSize := c.TextSize
	if textSize <= 0 {
		textSize = th.Typography.FontSize
	}
	fnt := c.Font
	if fnt == (font.Font{}) {
		fnt = th.Typography.Font
	}

	// Thème "material" minimal, requis par material.Editor, qui délègue à
	// AGUIK la couleur, la bordure et le fond réels.
	mth := material.NewTheme()
	mth.Shaper = th.Shaper

	children := []glayout.Child{}

	if !c.HideLineNumbers {
		lineCount := strings.Count(c.Editor.Text(), "\n") + 1
		children = append(children,
			glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return c.layoutGutter(gtx, th, fnt, textSize, lineCount)
			}),
			glayout.Rigid(glayout.Spacer(th.Spacing.SM)),
		)
	}

	children = append(children, glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
		ed := material.Editor(mth, &c.Editor, c.Placeholder)
		ed.Font = fnt
		ed.TextSize = unit.Sp(textSize)
		ed.Color = th.Colors.Foreground
		ed.HintColor = th.Colors.Muted
		// Couleur de sélection bien visible : c'est ce qui "montre" le
		// code sélectionné.
		ed.SelectionColor = glayout.WithAlpha(th.Colors.Primary, 110)
		return ed.Layout(gtx)
	}))

	return glayout.Row(0, layout.Start, children...)(gtx)
}
