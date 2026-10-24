package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

// Thème Gio interne, uniquement pour fabriquer la barre : les couleurs sont écrasées ci-dessous.
var scrollTheme = material.NewTheme()

// ScrollList affiche une liste avec une barre de défilement aux couleurs du thème.
// state doit être conservé entre les frames, comme tout composant avec état.
func ScrollList(gtx layout.Context, th *theme.Theme, state *widget.List, length int, el layout.ListElement) layout.Dimensions {
	list := material.List(scrollTheme, state)
	list.Indicator.Color = th.Colors.Border
	list.Indicator.HoverColor = th.Colors.Foreground
	list.Indicator.CornerRadius = 0
	list.Indicator.MinorWidth = 10
	return list.Layout(gtx, length, el)
}
