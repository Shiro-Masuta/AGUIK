package components

import (
	"image"
	"math"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

// Loader est un spinner de chargement indéterminé — animation continue,
// donc premier composant du kit qui demande une invalidation de frame à
// CHAQUE Layout (pas seulement sur un événement ponctuel comme Dialog).
type Loader struct {
	Size    float32 // dp, diamètre ; 0 -> 24
	Variant ButtonVariant

	start   time.Time
	started bool
}

func (l *Loader) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Redemande une frame immédiatement — sans ça, l'animation se
	// figerait dès qu'aucun autre événement (clic, resize...) ne survient.
	gtx.Execute(op.InvalidateCmd{})

	if !l.started {
		l.start = gtx.Now
		l.started = true
	}
	elapsed := gtx.Now.Sub(l.start).Seconds()
	angle := float32(elapsed) * 4 // vitesse de rotation, en radians/seconde — ajustable

	sizeDp := l.Size
	if sizeDp == 0 {
		sizeDp = 24
	}
	d := gtx.Dp(unit.Dp(sizeDp))
	strokeW := float32(gtx.Dp(unit.Dp(3)))
	radius := float32(d)/2 - strokeW

	col, _ := variantColors(th, l.Variant)
	if l.Variant == ButtonDefault {
		col = th.Colors.Primary
	}
	center := f32.Pt(float32(d)/2, float32(d)/2)

	// Approxime un arc de ~270° par une polyligne (segments droits) — plus
	// simple et robuste à écrire correctement qu'une primitive d'arc
	// exacte, largement suffisant visuellement à cette échelle.
	const segments = 24
	const arcSpan = 1.5 * math.Pi // ~270°

	var path clip.Path
	path.Begin(gtx.Ops)
	for i := 0; i <= segments; i++ {
		t := float64(i) / float64(segments) * arcSpan
		x := center.X + radius*float32(math.Cos(t))
		y := center.Y + radius*float32(math.Sin(t))
		if i == 0 {
			path.MoveTo(f32.Pt(x, y))
		} else {
			path.LineTo(f32.Pt(x, y))
		}
	}
	stroke := clip.Stroke{Path: path.End(), Width: strokeW}

	// Rotation appliquée à TOUT ce qui est dessiné après ce Push, jusqu'au
	// Pop correspondant — c'est ce qui fait tourner l'arc dans le temps.
	rotStack := op.Affine(f32.Affine2D{}.Rotate(center, angle)).Push(gtx.Ops)
	paint.FillShape(gtx.Ops, col, stroke.Op())
	rotStack.Pop()

	return layout.Dimensions{Size: image.Pt(d, d)}
}
