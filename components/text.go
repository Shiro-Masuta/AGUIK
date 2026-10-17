// components/text.go
package components

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

// textStyle regroupe tout ce qui est configurable pour un Text.
// Les valeurs sont préremplies depuis le Theme, puis surchargées par les options.
type textStyle struct {
	size      unit.Sp
	color     color.NRGBA
	font      font.Font
	alignment text.Alignment
	maxLines  int
	truncator string
}

// TextOption modifie une valeur précise de textStyle. Chaque With* ci-dessous
// n'écrase QUE le champ concerné — le reste garde sa valeur par défaut (theme).
type TextOption func(*textStyle)

func WithSize(sp float32) TextOption {
	return func(t *textStyle) { t.size = unit.Sp(sp) }
}

func WithColor(c color.NRGBA) TextOption {
	return func(t *textStyle) { t.color = c }
}

func WithFont(f font.Font) TextOption {
	return func(t *textStyle) { t.font = f }
}

func WithAlignment(a text.Alignment) TextOption {
	return func(t *textStyle) { t.alignment = a }
}

func WithMaxLines(n int) TextOption {
	return func(t *textStyle) { t.maxLines = n }
}

func WithWeight(w font.Weight) TextOption {
	return func(t *textStyle) { t.font.Weight = w }
}

func WithEllipsis() TextOption {
	return func(t *textStyle) {
		if t.maxLines == 0 {
			t.maxLines = 1
		}
		t.truncator = "…"
	}
}

// Text affiche une chaîne stylée. Par défaut, taille et couleur viennent du
// Theme (auto). Chaque option passée surcharge un champ précis :
//
//	components.Text(th, "Salut")                          // 100% auto
//	components.Text(th, "Alerte", components.WithColor(th.Colors.Danger))
//	components.Text(th, "Titre", components.WithSize(24))
func Text(th *theme.Theme, str string, opts ...TextOption) layout.Widget {
	ts := textStyle{
		size:  unit.Sp(th.Typography.FontSize),
		color: th.Colors.Foreground,
		font:  th.Typography.Font,
	}
	for _, opt := range opts {
		opt(&ts)
	}

	return func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = image.Point{} // ne jamais hériter aveuglément d'un Min externe

		macro := op.Record(gtx.Ops)
		paint.ColorOp{Color: ts.color}.Add(gtx.Ops)
		colorCall := macro.Stop()

		l := widget.Label{
			Alignment: ts.alignment,
			MaxLines:  ts.maxLines,
			Truncator: ts.truncator,
		}
		return l.Layout(gtx, th.Shaper, ts.font, ts.size, str, colorCall)
	}
}
