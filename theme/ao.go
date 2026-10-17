package theme

import "image/color"

// Ao retourne une nouvelle instance du thème sombre inspiré des bleus profonds.
// Chaque appel construit une copie fraîche — aucun état partagé entre les
// consommateurs, donc aucun risque qu'un composant modifie la palette globale
// par accident.
func Ao() Theme {
	return Theme{
		Colors: Colors{
			Background: color.NRGBA{R: 13, G: 17, B: 28, A: 255},
			Surface:    color.NRGBA{R: 21, G: 28, B: 43, A: 255},
			SurfaceAlt: color.NRGBA{R: 29, G: 39, B: 58, A: 255},

			Foreground: color.NRGBA{R: 241, G: 243, B: 245, A: 255},
			Muted:      color.NRGBA{R: 155, G: 168, B: 184, A: 255},
			Border:     color.NRGBA{R: 241, G: 243, B: 245, A: 255},

			Primary:   color.NRGBA{R: 79, G: 163, B: 209, A: 255},
			Secondary: color.NRGBA{R: 120, G: 199, B: 232, A: 255},

			Danger:  color.NRGBA{R: 217, G: 74, B: 91, A: 255},
			Success: color.NRGBA{R: 94, G: 184, B: 126, A: 255},
			Warning: color.NRGBA{R: 224, G: 170, B: 72, A: 255},
		},

		Shape: Shape{
			Radius:      5,
			BorderWidth: 3,
			Shadow: Shadow{
				Color:   color.NRGBA{R: 241, G: 243, B: 245, A: 255},
				OffsetX: 4,
				OffsetY: 4,
				Blur:    0,
				Spread:  0,
			},
		},

		Spacing: Spacing{
			XS: 4,
			SM: 8,
			MD: 12,
			LG: 16,
			XL: 24,
		},

		Typography: Typography{
			FontSize: 14,
		},
	}
}
