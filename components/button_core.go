package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

// layoutButtonFrame mesure content, dessine le NeoFrame autour, rejoue le
// contenu dedans, et enregistre la hitbox si le bouton est actif. Partagé
// par Button et IconButton — un seul endroit qui sait "fabriquer un
// bouton", peu importe ce qu'il affiche dedans (texte, icône, les deux).
func layoutButtonFrame(gtx layout.Context, th *theme.Theme, content layout.Widget, padding float32, bg color.NRGBA, hovered, pressed, disabled bool, clickable interaction.Interactive) layout.Dimensions {
	pad := gtx.Dp(unit.Dp(padding))

	mgtx := gtx
	mgtx.Constraints.Min = image.Point{}
	macro := op.Record(mgtx.Ops)
	contentDims := content(mgtx)
	macro.Stop() // dry-run : on jette ces opérations, on garde juste la taille

	innerW := contentDims.Size.X + 2*pad
	innerH := contentDims.Size.Y + 2*pad

	frame := NeoFrame(gtx, th, innerW, innerH, bg, hovered, pressed)

	off := frame.BgRect.Min.Add(image.Pt(pad, pad))
	stack := op.Offset(off).Push(gtx.Ops)
	content(gtx)
	stack.Pop()

	if !disabled {
		interaction.ClickArea(gtx, clickable, frame.RestRect)
	}

	return layout.Dimensions{Size: frame.CanvasDim}
}
