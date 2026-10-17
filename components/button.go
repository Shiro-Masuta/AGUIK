package components

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type ButtonVariant int

const (
	ButtonDefault ButtonVariant = iota
	ButtonPrimary
	ButtonSecondary
	ButtonSuccess
	ButtonWarning
	ButtonError
)

type ButtonSize int

const (
	SizeSM ButtonSize = iota
	SizeMD
	SizeLG
)

// Button est stateful (contrairement à Text/Icon/Card) : il porte un
// widget.Clickable qui doit survivre d'une frame à l'autre pour détecter
// les clics. D'où un struct + méthode Layout plutôt qu'une closure pure.
type Button struct {
	Clickable widget.Clickable
	Text      string
	Icon      *widget.Icon
	Variant   ButtonVariant
	Size      ButtonSize
	Disabled  bool
	Clicked   bool
}

func (b *Button) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	clicked := false
	if !b.Disabled {
		clicked = b.Clickable.Clicked(gtx)
		if clicked {
			interaction.ClearFocus(gtx)
			interaction.Invalidate(gtx)
		}
	}
	b.Clicked = clicked

	bg, fg := variantColors(th, b.Variant)
	if b.Disabled {
		bg, fg = th.Colors.Muted, th.Colors.Foreground
	}
	hovered := !b.Disabled && b.Clickable.Hovered()
	pressed := !b.Disabled && b.Clickable.Pressed()

	// Le contenu varie selon qu'un Icon est fourni ou non — le reste
	// (mesure, cadre, hitbox) ne change jamais.
	var content layout.Widget
	if b.Icon != nil {
		content = IconText(th, b.Icon, b.Text,
			[]IconOption{WithIconColor(fg)},
			[]TextOption{WithColor(fg), WithWeight(font.SemiBold)},
		)
	} else {
		content = Text(th, b.Text, WithColor(fg), WithWeight(font.SemiBold))
	}

	return layoutButtonFrame(gtx, th, content, paddingFor(th, b.Size), bg, hovered, pressed, b.Disabled, &b.Clickable)
}
func paddingFor(th *theme.Theme, s ButtonSize) float32 {
	switch s {
	case SizeMD:
		return th.Spacing.MD
	case SizeLG:
		return th.Spacing.LG
	default:
		return th.Spacing.SM
	}
}

// variantColors mappe un ButtonVariant vers ses couleurs fond/texte —
// partagé par Button et, plus tard, tout composant qui a besoin des mêmes
// variantes sémantiques (Checkbox, Alert...).
func variantColors(th *theme.Theme, v ButtonVariant) (bg, fg color.NRGBA) {
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	switch v {
	case ButtonPrimary:
		return th.Colors.Primary, white
	case ButtonSecondary:
		return th.Colors.Secondary, white
	case ButtonSuccess:
		return th.Colors.Success, white
	case ButtonWarning:
		return th.Colors.Warning, white
	case ButtonError:
		return th.Colors.Danger, white
	default:
		return th.Colors.Surface, th.Colors.Foreground
	}
}
