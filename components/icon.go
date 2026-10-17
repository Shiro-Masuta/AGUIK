package components

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

type iconStyle struct {
	size  unit.Dp
	color color.NRGBA
}

type IconOption func(*iconStyle)

func WithIconSize(dp float32) IconOption {
	return func(s *iconStyle) { s.size = unit.Dp(dp) }
}

func WithIconColor(c color.NRGBA) IconOption {
	return func(s *iconStyle) { s.color = c }
}

// Icon dessine une icône pré-chargée (*widget.Icon) à une taille et couleur
// données. Contrairement à Text, aucun Shaper requis — l'icône est un chemin
// vectoriel déjà figé, pas du texte à mettre en page.
func Icon(th *theme.Theme, src *widget.Icon, opts ...IconOption) layout.Widget {
	is := iconStyle{
		size:  unit.Dp(24), // taille par défaut raisonnable pour une UI
		color: th.Colors.Foreground,
	}
	for _, opt := range opts {
		opt(&is)
	}

	return func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = gtx.Dp(is.size)
		gtx.Constraints.Max.Y = gtx.Dp(is.size)
		gtx.Constraints.Min = gtx.Constraints.Max
		return src.Layout(gtx, is.color)
	}
}
