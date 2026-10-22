package selectwidget

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/Shiro-Masuta/AGUIK/components"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

// drawSelectRow peint une ligne de liste plate : fond de surbrillance
// optionnel (hover) + texte. Le texte est mesuré AVANT de peindre le fond
// (sinon le rectangle de surbrillance serait peint à hauteur 0, puisque
// sa taille dépendrait d'un texte pas encore layouté).
func drawSelectRow(gtx layout.Context, th *theme.Theme, label string, hovered bool, padY int) layout.Dimensions {
	content := components.Text(th, label, components.WithColor(th.Colors.Foreground))

	mgtx := gtx
	mgtx.Constraints.Min = image.Point{}
	macro := op.Record(mgtx.Ops)
	textDims := content(mgtx)
	macro.Stop() // dry-run : jeté, on garde juste la taille

	rowH := textDims.Size.Y + 2*padY

	if hovered {
		paint.FillShape(gtx.Ops, th.Colors.SurfaceAlt, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, rowH)}.Op())
	}

	off := op.Offset(image.Pt(0, padY)).Push(gtx.Ops)
	content(gtx)
	off.Pop()

	return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, rowH)}
}
