package overlay

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

// OverlayState porte le Clickable de fond — doit être un champ persistant
// du composant appelant (Dialog, futur Toast...), jamais recréé à chaque
// frame, comme tout state Gio.
type OverlayState struct {
	backdrop widget.Clickable
}

// Overlay centralise tout ce qui est commun à un popup modal centré : fond
// assombri qui bloque les clics vers ce qu'il y a derrière, mesure du
// contenu, centrage, fermeture au clic sur le fond, invalidation de frame.
// Le contenu (dialog, alerte, toast...) n'a besoin de connaître QUE sa
// propre apparence — jamais la mécanique de superposition.
func Overlay(gtx layout.Context, state *OverlayState, open *bool, dismissOnBackdropClick bool, content layout.Widget) layout.Dimensions {
	if !*open {
		return layout.Dimensions{}
	}

	full := gtx.Constraints.Max

	backdropClicked := state.backdrop.Clicked(gtx)
	paint.FillShape(gtx.Ops, color.NRGBA{A: 160}, clip.Rect(image.Rectangle{Max: full}).Op())
	area := clip.Rect(image.Rectangle{Max: full}).Push(gtx.Ops)
	state.backdrop.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: full}
	})
	area.Pop()

	if backdropClicked && dismissOnBackdropClick {
		*open = false
		gtx.Execute(op.InvalidateCmd{}) // Gio ne redessine jamais sur un simple changement d'état interne
		return layout.Dimensions{Size: full}
	}

	// Contenu mesuré UNE SEULE FOIS (clics internes compris), puis rejoué
	// centré — même pattern que Card, appliqué ici à l'échelle de toute
	// la fenêtre.
	mgtx := gtx
	macro := op.Record(mgtx.Ops)
	dims := content(mgtx)
	call := macro.Stop()

	x := (full.X - dims.Size.X) / 2
	y := (full.Y - dims.Size.Y) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	stack := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	stack.Pop()

	return layout.Dimensions{Size: full}
}
