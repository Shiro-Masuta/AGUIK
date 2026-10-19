package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

const (
	hoverLiftDp = 2
	hoverGrowDp = 2
)

type Frame struct {
	BgRect    image.Rectangle
	RestRect  image.Rectangle
	CanvasDim image.Point
}

// neoFrame dessine le cadre néobrutaliste commun (ombre, bordure, fond,
// translation hover/pressed). Utilisé en interne par Card, Button,
// IconButton — jamais appelé directement par l'utilisateur de la lib.

func NeoFrame(gtx layout.Context, th *theme.Theme, innerW, innerH int, bg color.NRGBA, hovered, pressed bool) Frame {
	dp := func(v float32) int { return gtx.Dp(unit.Dp(v)) }

	baseX, baseY := dp(th.Shape.Shadow.OffsetX), dp(th.Shape.Shadow.OffsetY)
	lift := dp(hoverLiftDp)
	grow := dp(hoverGrowDp)
	borderW := dp(th.Shape.BorderWidth)
	radius := dp(th.Shape.Radius)

	var contentOffX, contentOffY, shadowOffX, shadowOffY int
	shadowVisible := true
	switch {
	case pressed:
		contentOffX, contentOffY = baseX, baseY
		shadowVisible = false
	case hovered:
		contentOffX, contentOffY = -lift, -lift
		shadowOffX, shadowOffY = baseX+grow, baseY+grow
	default:
		shadowOffX, shadowOffY = baseX, baseY
	}

	padLeft, padTop := lift, lift
	padRight, padBottom := baseX+grow, baseY+grow

	outerW := innerW + 2*borderW
	outerH := innerH + 2*borderW
	canvasW := padLeft + outerW + padRight
	canvasH := padTop + outerH + padBottom

	if shadowVisible {
		shadowRect := image.Rectangle{
			Min: image.Pt(padLeft+shadowOffX, padTop+shadowOffY),
			Max: image.Pt(padLeft+shadowOffX+outerW, padTop+shadowOffY+outerH),
		}
		shadowClip := clip.RRect{Rect: shadowRect, SE: radius, SW: radius, NE: radius, NW: radius}
		paint.FillShape(gtx.Ops, th.Shape.Shadow.Color, shadowClip.Op(gtx.Ops))
	}

	borderRect := image.Rectangle{
		Min: image.Pt(padLeft+contentOffX, padTop+contentOffY),
		Max: image.Pt(padLeft+contentOffX+outerW, padTop+contentOffY+outerH),
	}
	borderClip := clip.RRect{Rect: borderRect, SE: radius, SW: radius, NE: radius, NW: radius}
	paint.FillShape(gtx.Ops, th.Colors.Border, borderClip.Op(gtx.Ops))

	bgRect := image.Rectangle{
		Min: borderRect.Min.Add(image.Pt(borderW, borderW)),
		Max: borderRect.Max.Sub(image.Pt(borderW, borderW)),
	}
	innerRadius := radius - borderW
	if innerRadius < 0 {
		innerRadius = 0
	}
	bgClip := clip.RRect{Rect: bgRect, SE: innerRadius, SW: innerRadius, NE: innerRadius, NW: innerRadius}
	paint.FillShape(gtx.Ops, bg, bgClip.Op(gtx.Ops))

	restRect := image.Rectangle{
		Min: image.Pt(padLeft, padTop),
		Max: image.Pt(padLeft+outerW, padTop+outerH),
	}

	return Frame{
		BgRect:    bgRect,
		RestRect:  restRect,
		CanvasDim: image.Pt(canvasW, canvasH),
	}
}

// neoFrameOverheadX est la largeur que NeoFrame ajoute autour du contenu
// (marges de lift/ombre + bordures). Garder synchronisé avec NeoFrame.
func neoFrameOverheadX(gtx layout.Context, th *theme.Theme) int {
	dp := func(v float32) int { return gtx.Dp(unit.Dp(v)) }
	baseX := dp(th.Shape.Shadow.OffsetX)
	lift := dp(hoverLiftDp)
	grow := dp(hoverGrowDp)
	borderW := dp(th.Shape.BorderWidth)
	return lift + baseX + grow + 2*borderW
}

func neoFrameOverheadY(gtx layout.Context, th *theme.Theme) int {
	dp := func(v float32) int { return gtx.Dp(unit.Dp(v)) }
	baseY := dp(th.Shape.Shadow.OffsetY)
	lift := dp(hoverLiftDp)
	grow := dp(hoverGrowDp)
	borderW := dp(th.Shape.BorderWidth)
	return lift + baseY + grow + 2*borderW
}
