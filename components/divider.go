package components

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

// Divider est un simple trait de séparation, horizontal ou vertical.
// Non-interactif, sans état — le composant le plus simple du kit.
type Divider struct {
	Vertical  bool
	Thickness float32 // dp ; 0 -> th.Shape.BorderWidth
	Length    float32 // dp ; 0 -> occupe tout l'espace disponible sur l'axe principal
}

func (d *Divider) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	thicknessDp := d.Thickness
	if thicknessDp == 0 {
		thicknessDp = th.Shape.BorderWidth
	}
	thickness := gtx.Dp(unit.Dp(thicknessDp))

	var size image.Point
	if d.Vertical {
		length := gtx.Constraints.Max.Y
		if d.Length != 0 {
			length = gtx.Dp(unit.Dp(d.Length))
		}
		size = image.Pt(thickness, length)
	} else {
		length := gtx.Constraints.Max.X
		if d.Length != 0 {
			length = gtx.Dp(unit.Dp(d.Length))
		}
		size = image.Pt(length, thickness)
	}

	paint.FillShape(gtx.Ops, th.Colors.Border, clip.Rect{Max: size}.Op())

	return layout.Dimensions{Size: size}
}
