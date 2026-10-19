package components

import (
	"image"
	"image/color"
	"reflect"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/Shiro-Masuta/AGUIK/glayout"
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
	Variant ButtonVariant

	// Disabled désactive la zone de clic.
	Disabled bool

	// Clicked contient l'index du surlignage cliqué cette frame, -1 sinon
	// — à lire juste après Layout, même convention que Dialog.Clicked.
	Clicked int

	click gesture.Click

	scaleX float32
	scaleY float32
	offset f32.Point

	cachedSrc image.Image
	op        paint.ImageOp
}

// HitTest retourne l'index du surlignage sous le point p (coordonnées
// écran), ou -1. Parcourt en ordre inverse : le dernier dessiné (donc
// visuellement au-dessus en cas de chevauchement) est testé en premier.
func (v *ImageView) HitTest(p image.Point) int {
	if v.scaleX == 0 || v.scaleY == 0 {
		return -1
	}
	ip := image.Pt(
		int((float32(p.X)-v.offset.X)/v.scaleX),
		int((float32(p.Y)-v.offset.Y)/v.scaleY),
	)
	for i := len(v.Highlights) - 1; i >= 0; i-- {
		if ip.In(v.Highlights[i].Rect) {
			return i
		}
	}
	return -1
}

func (v *ImageView) toScreen(r image.Rectangle) image.Rectangle {
	return image.Rect(
		int(float32(r.Min.X)*v.scaleX+v.offset.X),
		int(float32(r.Min.Y)*v.scaleY+v.offset.Y),
		int(float32(r.Max.X)*v.scaleX+v.offset.X),
		int(float32(r.Max.Y)*v.scaleY+v.offset.Y),
	)
}

func (v *ImageView) ensureOp() {
	if v.Src == nil {
		return
	}
	srcType := reflect.TypeOf(v.Src)
	if srcType.Comparable() && v.cachedSrc != nil && v.cachedSrc == v.Src {
		return
	}
	v.op = paint.NewImageOp(v.Src)
	if srcType.Comparable() {
		v.cachedSrc = v.Src
	} else {
		v.cachedSrc = nil
	}
}

// variantColor résout la couleur par défaut d'un surlignage — réutilise
// variantColors (button_core.go), avec un repli sur Primary plutôt que
// Surface pour ButtonDefault : un surlignage gris serait invisible sur la
// plupart des images.
func variantColor(th *theme.Theme, variant ButtonVariant) color.NRGBA {
	if variant == ButtonDefault {
		return th.Colors.Primary
	}
	bg, _ := variantColors(th, variant)
	return bg
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

// computeFit calcule scale/offset selon Fit — factorisé pour ne pas
// dupliquer la logique entre les trois modes.
func (v *ImageView) computeFit(size, b image.Point) {
	sx := float32(size.X) / float32(b.X)
	sy := float32(size.Y) / float32(b.Y)

	switch v.Fit {
	case FitCover:
		scale := sx
		if sy > sx {
			scale = sy
		}
		v.scaleX, v.scaleY = scale, scale
	case FitFill:
		v.scaleX, v.scaleY = sx, sy
	default: // FitContain
		scale := sx
		if sy < sx {
			scale = sy
		}
		v.scaleX, v.scaleY = scale, scale
	}

	v.offset = f32.Pt(
		(float32(size.X)-float32(b.X)*v.scaleX)/2,
		(float32(size.Y)-float32(b.Y)*v.scaleY)/2,
	)
}

func (v *ImageView) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	v.Clicked = -1
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

	if !v.Disabled {
		for {
			e, ok := v.click.Update(gtx.Source)
			if !ok {
				break
			}
			if e.Kind == gesture.KindClick {
				if i := v.HitTest(e.Position); i >= 0 {
					v.Clicked = i
				}
			}
		}
	}

	viewClip := clip.Rect{Max: size}.Push(gtx.Ops)
	tr := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(v.scaleX, v.scaleY)).
		Offset(v.offset)).Push(gtx.Ops)
	v.op.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	tr.Pop()

	for i, h := range v.Highlights {
		r := v.toScreen(h.Rect)
		style := v.resolveStyle(th, i, h)

		if style.Fill != (color.NRGBA{}) {
			paint.FillShape(gtx.Ops, style.Fill, clip.Rect(r).Op())
		}
		if style.Stroke != (color.NRGBA{}) && style.StrokeWidth > 0 {
			paint.FillShape(gtx.Ops, style.Stroke, clip.Stroke{
				Path:  clip.Rect(r).Path(),
				Width: style.StrokeWidth,
			}.Op())
		}
	}
	viewClip.Pop()

	if !v.Disabled {
		area := clip.Rect{Max: size}.Push(gtx.Ops)
		v.click.Add(gtx.Ops)
		area.Pop()
	}

	return layout.Dimensions{Size: size}
}
