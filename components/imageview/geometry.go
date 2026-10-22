package imageview

import (
	"image"

	"gioui.org/f32"
)

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

func (v *ImageView) toScreen(r image.Rectangle) image.Rectangle {
	return image.Rect(
		int(float32(r.Min.X)*v.scaleX+v.offset.X),
		int(float32(r.Min.Y)*v.scaleY+v.offset.Y),
		int(float32(r.Max.X)*v.scaleX+v.offset.X),
		int(float32(r.Max.Y)*v.scaleY+v.offset.Y),
	)
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
