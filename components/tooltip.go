package components

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

type TooltipPlacement int

const (
	TooltipTop TooltipPlacement = iota
	TooltipBottom
	TooltipRight
)

// Tooltip enveloppe un widget NON-INTERACTIF (Icon, Text...) et affiche une
// bulle au survol. Pour un widget déjà cliquable (Button, IconButton), ne
// pas utiliser ce wrapper — appeler directement LayoutTooltipBubble avec
// le hover que le composant suit déjà lui-même (voir la note dans ce
// fichier), pour éviter de superposer deux détections de survol sur la
// même zone.
type Tooltip struct {
	Text      string
	Placement TooltipPlacement

	hover gesture.Hover
}

func (t *Tooltip) Layout(gtx layout.Context, th *theme.Theme, content layout.Widget) layout.Dimensions {
	dims := content(gtx)

	t.hover.Add(gtx.Ops)
	hovered := t.hover.Update(gtx.Source)

	if hovered && t.Text != "" {
		LayoutTooltipBubble(gtx, th, t.Text, t.Placement, dims.Size)
	}

	return dims
}

// LayoutTooltipBubble dessine la bulle elle-même, différée via op.Defer
// (même technique que le dropdown de Select) pour passer au-dessus de
// tout le reste. anchorSize = taille du widget déclencheur, pour savoir
// où positionner la bulle par rapport à lui.
func LayoutTooltipBubble(gtx layout.Context, th *theme.Theme, text string, placement TooltipPlacement, anchorSize image.Point) {
	pad := gtx.Dp(unit.Dp(th.Spacing.SM))

	content := Text(th, text, WithColor(th.Colors.Foreground), WithSize(th.Typography.FontSize*0.85))

	mgtx := gtx
	mgtx.Constraints.Min = image.Point{}
	macro := op.Record(mgtx.Ops)
	textDims := content(mgtx)
	macro.Stop()

	innerW := textDims.Size.X + 2*pad
	innerH := textDims.Size.Y + 2*pad

	bubbleMacro := op.Record(gtx.Ops)
	frame := NeoFrame(gtx, th, innerW, innerH, th.Colors.SurfaceAlt, false, false)
	off := frame.BgRect.Min.Add(image.Pt(pad, pad))
	textStack := op.Offset(off).Push(gtx.Ops)
	content(gtx)
	textStack.Pop()
	bubbleCall := bubbleMacro.Stop()

	var posX, posY int
	switch placement {
	case TooltipBottom:
		posY = anchorSize.Y
	case TooltipRight:
		posX = anchorSize.X
		posY = (anchorSize.Y - innerH) / 2
	default:
		posY = -innerH
	}

	posStack := op.Offset(image.Pt(posX, posY)).Push(gtx.Ops)
	op.Defer(gtx.Ops, bubbleCall)
	posStack.Pop()
}
