package interaction

import (
	"gioui.org/layout"
	"gioui.org/op"
)

// Invalidate demande explicitement une nouvelle frame. Gio ne redessine
// jamais tout seul sur une simple mutation de state Go — tout composant
// qui fait transitionner un état "logique" (Clicked, Toggled, Open...)
// suite à un event pointeur DOIT l'appeler ici, une seule fois, plutôt
// que de compter sur l'appelant pour y penser à chaque site d'usage.
func Invalidate(gtx layout.Context) {
	gtx.Execute(op.InvalidateCmd{})
}
