package components

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

// Badge est une étiquette non-interactive (statut, tag, compteur...) au
// même cadre néobrutaliste que les autres composants. Contrairement à
// Button/Checkbox/Switch, hovered/pressed valent toujours false — pas de
// widget.Clickable ni de widget.Bool, donc pas de ClickArea non plus : un
// Badge ne réagit jamais au pointeur, il garde en permanence son ombre au
// repos.
type Badge struct {
	Text    string
	Variant ButtonVariant
	Size    ButtonSize
}

// badgePaddingFor renvoie un padding plus compact que paddingFor (Button) :
// un badge est fait pour être petit, pas pour ressembler à un bouton.
func badgePaddingFor(th *theme.Theme, s ButtonSize) float32 {
	switch s {
	case SizeMD:
		return th.Spacing.SM
	case SizeLG:
		return th.Spacing.MD
	default:
		return th.Spacing.XS
	}
}

func (b *Badge) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	bg, fg := variantColors(th, b.Variant)

	padX := gtx.Dp(unit.Dp(badgePaddingFor(th, b.Size)))
	padY := padX / 2

	label := Text(th, b.Text,
		WithColor(fg),
		WithSize(th.Typography.FontSize*0.85),
		WithWeight(font.SemiBold),
	)

	mgtx := gtx
	mgtx.Constraints.Min = image.Point{}
	macro := op.Record(mgtx.Ops)
	labelDims := label(mgtx)
	macro.Stop()

	innerW := labelDims.Size.X + 2*padX
	innerH := labelDims.Size.Y + 2*padY

	// hovered/pressed toujours false : un Badge ne réagit jamais au
	// pointeur, il garde son cadre "au repos" en permanence.
	frame := NeoFrame(gtx, th, innerW, innerH, bg, false, false)

	labelOff := frame.BgRect.Min.Add(image.Pt(padX, padY))
	labelStack := op.Offset(labelOff).Push(gtx.Ops)
	label(gtx)
	labelStack.Pop()

	return layout.Dimensions{Size: frame.CanvasDim}
}
