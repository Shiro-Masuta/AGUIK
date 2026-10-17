package interaction

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/widget"
)

type Interactive interface {
	Layout(gtx layout.Context, w layout.Widget) layout.Dimensions
}

// ClickArea enregistre la zone d'écoute pointeur sur restRect (jamais sur
// tout le canvas, qui inclut la marge réservée à l'ombre/au lift) et pose
// le curseur "main". L'offset aligne la hitbox sur ce qui est réellement
// dessiné au repos (clip.Rect seul ne fait que rogner, il ne déplace pas
// l'origine locale utilisée par Clickable).
func ClickArea(gtx layout.Context, target Interactive, restRect image.Rectangle) {
	offStack := op.Offset(restRect.Min).Push(gtx.Ops)
	area := clip.Rect(image.Rectangle{Max: restRect.Size()}).Push(gtx.Ops)
	target.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: restRect.Size()}
	})
	pointer.CursorPointer.Add(gtx.Ops)
	area.Pop()
	offStack.Pop()
}

// UpdateClickable traite les events en attente AVANT toute lecture de
// Hovered()/Pressed(), pour que le rendu de cette frame reflète les events
// de cette frame (sinon on dessine l'état d'un event de retard).
func UpdateClickable(gtx layout.Context, clickable *widget.Clickable, disabled bool) {
	if !disabled {
		clickable.Update(gtx)
	}
}
