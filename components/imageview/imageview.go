package imageview

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op/paint"

	"github.com/Shiro-Masuta/AGUIK/components"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

// ImageFit contrôle comment l'image source est ajustée aux contraintes,
// même sémantique que CSS object-fit.
type ImageFit int

const (
	FitContain ImageFit = iota // préserve le ratio, image entière visible (défaut)
	FitCover                   // préserve le ratio, remplit tout l'espace, peut rogner
	FitFill                    // étire l'image pour remplir exactement, ratio non préservé
)

// HighlightStyle décrit l'apparence d'un surlignage précis. Un Highlight
// dont Style est nil retombe sur la couleur dérivée de ImageView.Variant
// (comportement générique par défaut, sans obliger l'appelant à préciser
// un style pour chaque rectangle).
type HighlightStyle struct {
	Fill        color.NRGBA
	Stroke      color.NRGBA
	StrokeWidth float32
}

// Highlight est un rectangle surligné, en coordonnées pixel de l'image
// source (pas de l'écran). Style permet de personnaliser chaque
// surlignage individuellement — utile pour distinguer plusieurs
// catégories dans l'image (pas seulement un état actif/inactif).
type Highlight struct {
	Rect  image.Rectangle
	Style *HighlightStyle
}

// ImageView affiche une image avec un ensemble de rectangles surlignés.
// Générique : ne connaît rien d'OCR ou de tout autre domaine, juste des
// rectangles indexés et des styles optionnels.
type ImageView struct {
	// Src est l'image à afficher. La remplacer déclenche la reconstruction
	// de la texture au prochain Layout.
	Src image.Image

	// Fit contrôle l'ajustement de l'image aux contraintes disponibles.
	Fit ImageFit

	// Highlights sont les rectangles à dessiner, en coordonnées image source.
	Highlights []Highlight

	// Active est l'index du surlignage mis en accent (bordure + remplissage
	// plus opaque), utilisé seulement pour les Highlights sans Style
	// explicite. -1 = aucun.
	Active int

	// Variant contrôle la couleur par défaut des surlignages sans Style
	// explicite.
	Variant components.ButtonVariant

	// Disabled désactive la zone de clic.
	Disabled bool

	clicked int

	click gesture.Click

	scaleX float32
	scaleY float32
	offset f32.Point

	cachedSrc image.Image
	op        paint.ImageOp
}

func (v *ImageView) Clicked() int { return v.clicked }

func (v *ImageView) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	v.clicked = -1
	size := gtx.Constraints.Max

	if v.Src == nil {
		return layout.Dimensions{Size: size}
	}
	v.ensureOp()

	b := v.Src.Bounds().Size()
	if b.X == 0 || b.Y == 0 {
		return layout.Dimensions{Size: size}
	}

	v.computeFit(size, b)

	v.updateClick(gtx)
	v.drawImage(gtx, size)
	v.drawHighlights(gtx, th, size)
	v.registerClickArea(gtx, size)

	return layout.Dimensions{Size: size}
}
