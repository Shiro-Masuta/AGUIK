package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type RadioOption struct {
	Label string
}

// RadioGroup est une sélection exclusive (un seul actif à la fois) au
// style visuel de Checkbox — même socle (toggleCheckedColors,
// checkboxSizeFor), mais un point rempli au lieu d'une coche, et un
// Clickable par option plutôt qu'un widget.Bool (on n'a pas besoin d'un
// état indépendant par item, juste "lequel est actif").
type RadioGroup struct {
	Options  []RadioOption
	Selected int // -1 = aucune sélection
	Variant  ButtonVariant
	Size     ButtonSize
	Disabled bool

	btns []widget.Clickable
}

func (r *RadioGroup) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if r.Selected < -1 || r.Selected >= len(r.Options) {
		r.Selected = -1
	}
	r.btns = glayout.Resize(r.btns, len(r.Options))

	checkedBg, checkedFg := toggleCheckedColors(th, r.Variant)

	children := make([]glayout.Child, len(r.Options))
	for i, opt := range r.Options {
		if !r.Disabled && r.btns[i].Clicked(gtx) {
			r.Selected = i
			interaction.Invalidate(gtx)
		}
		hovered := !r.Disabled && r.btns[i].Hovered()
		pressed := !r.Disabled && r.btns[i].Pressed()
		selected := i == r.Selected

		children[i] = glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return r.layoutOption(gtx, th, opt.Label, selected, hovered, pressed, checkedBg, checkedFg, &r.btns[i])
		})
	}

	return glayout.Column(th.Spacing.SM, layout.Start, children...)(gtx)
}

func (r *RadioGroup) layoutOption(gtx layout.Context, th *theme.Theme, label string, selected, hovered, pressed bool, checkedBg, checkedFg color.NRGBA, btn *widget.Clickable) layout.Dimensions {
	boxPx := gtx.Dp(unit.Dp(checkboxSizeFor(r.Size)))

	bg := th.Colors.Surface
	if selected {
		bg = checkedBg
	}
	if r.Disabled {
		bg = th.Colors.Muted
	}

	frame := NeoFrame(gtx, th, boxPx, boxPx, bg, hovered, pressed)

	if selected {
		drawRadioDot(gtx, frame.BgRect, checkedFg)
	}

	if !r.Disabled {
		interaction.ClickArea(gtx, btn, frame.RestRect)
	}

	boxSize := frame.CanvasDim
	if label == "" {
		return layout.Dimensions{Size: boxSize}
	}

	pad := gtx.Dp(unit.Dp(th.Spacing.SM))
	lblColor := th.Colors.Foreground
	if r.Disabled {
		lblColor = th.Colors.Muted
	}

	lbl := Text(th, label, WithColor(lblColor))

	mgtx := gtx
	mgtx.Constraints.Min = image.Point{}
	macro := op.Record(mgtx.Ops)
	lblDims := lbl(mgtx)
	macro.Stop()

	labelOff := image.Pt(boxSize.X+pad, (boxSize.Y-lblDims.Size.Y)/2)
	labelStack := op.Offset(labelOff).Push(gtx.Ops)
	lbl(gtx)
	labelStack.Pop()

	total := image.Pt(boxSize.X+pad+lblDims.Size.X, boxSize.Y)
	if lblDims.Size.Y > total.Y {
		total.Y = lblDims.Size.Y
	}
	return layout.Dimensions{Size: total}
}

// drawRadioDot peint un cercle plein, centré et rétréci par rapport à rect
// (contrairement à Checkbox où la coche touche presque les bords).
func drawRadioDot(gtx layout.Context, rect image.Rectangle, col color.NRGBA) {
	inset := rect.Dx() / 4
	if inset < 1 {
		inset = 1
	}
	dot := image.Rectangle{
		Min: rect.Min.Add(image.Pt(inset, inset)),
		Max: rect.Max.Sub(image.Pt(inset, inset)),
	}
	e := clip.Ellipse(dot)
	paint.FillShape(gtx.Ops, col, e.Op(gtx.Ops))
}
