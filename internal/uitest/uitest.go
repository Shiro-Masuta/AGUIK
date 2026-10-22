// Package uitest simule des frames et des clics Gio pour tester les composants
// sans fenêtre.
package uitest

import (
	"image"
	"time"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

// Harness rejoue des frames successives en conservant le routeur d'événements,
// comme une vraie fenêtre : un clic est livré à la frame suivante.
type Harness struct {
	Theme   theme.Theme
	Max     image.Point
	PxPerDp float32

	router input.Router
}

// New retourne un Harness avec le thème Ao, une police embarquée et 1 px par dp.
func New(max image.Point) *Harness {
	th := theme.Ao()
	th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
	return &Harness{Theme: th, Max: max, PxPerDp: 1}
}

// Frame exécute w dans un nouveau contexte puis publie ses opérations au routeur.
func (h *Harness) Frame(w layout.Widget) layout.Dimensions {
	ops := new(op.Ops)
	gtx := layout.Context{
		Ops:         ops,
		Source:      h.router.Source(),
		Now:         time.Now(),
		Metric:      unit.Metric{PxPerDp: h.PxPerDp, PxPerSp: h.PxPerDp},
		Constraints: layout.Constraints{Max: h.Max},
	}
	dims := w(gtx)
	h.router.Frame(ops)
	return dims
}

// Click met en file un appui puis un relâchement au point p. Les widgets le
// reçoivent au prochain Frame.
func (h *Harness) Click(p image.Point) {
	pos := f32.Pt(float32(p.X), float32(p.Y))
	h.router.Queue(
		pointer.Event{Kind: pointer.Press, Source: pointer.Mouse, Buttons: pointer.ButtonPrimary, Position: pos},
		pointer.Event{Kind: pointer.Release, Source: pointer.Mouse, Position: pos},
	)
}
