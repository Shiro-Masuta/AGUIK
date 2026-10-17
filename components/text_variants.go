package components

import (
	"gioui.org/font"
	"gioui.org/layout"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

// Title — texte de titre principal, grande taille, poids fort.
func Title(th *theme.Theme, str string, opts ...TextOption) layout.Widget {
	base := []TextOption{
		WithSize(th.Typography.FontSize * 1.75), // ex: 14*1.75 = 24.5
		WithWeight(font.Bold),
	}
	return Text(th, str, append(base, opts...)...)
}

// Subtitle — intermédiaire, poids medium.
func Subtitle(th *theme.Theme, str string, opts ...TextOption) layout.Widget {
	base := []TextOption{
		WithSize(th.Typography.FontSize * 1.25),
		WithWeight(font.Medium),
	}
	return Text(th, str, append(base, opts...)...)
}

// Body — texte courant, taille et poids du theme, aucune surcharge.
func Body(th *theme.Theme, str string, opts ...TextOption) layout.Widget {
	return Text(th, str, opts...)
}

// Caption — petit texte secondaire, couleur atténuée par défaut.
func Caption(th *theme.Theme, str string, opts ...TextOption) layout.Widget {
	base := []TextOption{
		WithSize(th.Typography.FontSize * 0.85),
		WithColor(th.Colors.Muted),
	}
	return Text(th, str, append(base, opts...)...)
}
