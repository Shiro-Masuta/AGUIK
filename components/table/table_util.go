package table

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

func columnWeight(w float32) float32 {
	if w == 0 {
		return 1
	}
	return w
}

// plainRow dessine content sur un fond uni pleine largeur, avec un padding
// horizontal et vertical issu du thème. Pas de bordure ni d'ombre. Le header et
// les lignes partagent ce padding, donc leurs colonnes sont alignées.
func plainRow(gtx layout.Context, th *theme.Theme, bg color.NRGBA, content layout.Widget) layout.Dimensions {
	padX := gtx.Dp(unit.Dp(th.Spacing.MD))
	padY := gtx.Dp(unit.Dp(th.Spacing.SM))
	innerW := max(0, gtx.Constraints.Max.X-2*padX)

	inner := gtx
	inner.Constraints = layout.Constraints{
		Min: image.Pt(innerW, 0),
		Max: image.Pt(innerW, gtx.Constraints.Max.Y),
	}

	macro := op.Record(gtx.Ops)
	dims := content(inner)
	call := macro.Stop()

	size := image.Pt(gtx.Constraints.Max.X, dims.Size.Y+2*padY)
	glayout.FillRect(gtx, bg, size)

	off := op.Offset(image.Pt(padX, padY)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	off.Pop()

	return layout.Dimensions{Size: size}
}

// withRule ajoute un trait horizontal de px pixels sous dims.
func withRule(gtx layout.Context, dims layout.Dimensions, col color.NRGBA, px int) layout.Dimensions {
	defer op.Offset(image.Pt(0, dims.Size.Y)).Push(gtx.Ops).Pop()
	glayout.FillRect(gtx, col, image.Pt(dims.Size.X, px))
	return layout.Dimensions{Size: image.Pt(dims.Size.X, dims.Size.Y+px)}
}
