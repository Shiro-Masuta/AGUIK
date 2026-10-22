package imageview

import (
	"image/color"

	"github.com/Shiro-Masuta/AGUIK/components"
	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

// variantColor résout la couleur par défaut d'un surlignage — réutilise
// variantColors (button_core.go), avec un repli sur Primary plutôt que
// Surface pour ButtonDefault : un surlignage gris serait invisible sur la
// plupart des images.
func variantColor(th *theme.Theme, variant components.ButtonVariant) color.NRGBA {
	if variant == components.ButtonDefault {
		return th.Colors.Primary
	}
	switch variant {
	case components.ButtonPrimary:
		return th.Colors.Primary
	case components.ButtonSecondary:
		return th.Colors.Secondary
	case components.ButtonSuccess:
		return th.Colors.Success
	case components.ButtonWarning:
		return th.Colors.Warning
	case components.ButtonError:
		return th.Colors.Danger
	default:
		return th.Colors.Primary
	}
}

// resolveStyle retourne le style effectif d'un highlight : celui fourni
// explicitement, ou un style dérivé de Variant/Active sinon.
func (v *ImageView) resolveStyle(th *theme.Theme, i int, h Highlight) HighlightStyle {
	if h.Style != nil {
		return *h.Style
	}
	col := variantColor(th, v.Variant)
	if i == v.Active {
		return HighlightStyle{Fill: glayout.WithAlpha(col, 70), Stroke: col, StrokeWidth: 2}
	}
	return HighlightStyle{Fill: glayout.WithAlpha(col, 40)}
}
