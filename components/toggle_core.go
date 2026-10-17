package components

import (
	"image/color"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

// toggleCheckedColors résout la couleur "coché/activé" d'un composant à
// bascule (Checkbox, Switch, RadioGroup...). ButtonDefault (gris neutre)
// serait visuellement trop proche de l'état "non coché" — on retombe sur
// Primary/blanc dans ce cas précis plutôt que d'obliger l'appelant à
// toujours préciser un Variant.
func toggleCheckedColors(th *theme.Theme, v ButtonVariant) (bg, fg color.NRGBA) {
	bg, fg = variantColors(th, v)
	if v == ButtonDefault {
		bg, fg = th.Colors.Primary, color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	}
	return bg, fg
}
