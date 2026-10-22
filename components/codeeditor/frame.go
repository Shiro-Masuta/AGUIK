package codeeditor

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

// drawFrame dessine un cadre simple (fond + bordure arrondie) autour du
// contenu, en suivant le pattern "mesurer puis dessiner" du reste d'AGUIK :
// on enregistre les opérations du contenu, on connaît alors ses dimensions
// réelles, on dessine le fond à cette taille, puis on rejoue le contenu.
func drawFrame(gtx layout.Context, th *theme.Theme, bg, border color.NRGBA, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	rr := gtx.Dp(unit.Dp(th.Shape.Radius))
	rect := clip.RRect{
		Rect: image.Rectangle{Max: dims.Size},
		NE:   rr, NW: rr, SE: rr, SW: rr,
	}

	paint.FillShape(gtx.Ops, bg, rect.Op(gtx.Ops))

	bw := gtx.Dp(unit.Dp(th.Shape.BorderWidth))
	if bw > 0 {
		stroke := clip.Stroke{Path: rect.Path(gtx.Ops), Width: float32(bw)}
		paint.FillShape(gtx.Ops, border, stroke.Op())
	}

	call.Add(gtx.Ops)
	return dims
}
