package components

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type Alert struct {
	Icon        *widget.Icon
	Title       string
	Message     string
	Variant     ButtonVariant
	Dismissible bool
	MaxWidth    float32

	closeBtn widget.Clickable
	closed   bool
}

func (a *Alert) Closed() bool { return a.closed }

const (
	defaultAlertMaxWidthDp = 280
	alertAccentDp          = 6
	alertIconDp            = 22
	alertCloseDp           = 18
)

// Layout suit la même signature que Dialog/Button : (gtx, th), pas de
// shaper à part — Text() le prend déjà via th.
func (a *Alert) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if a.closed {
		return layout.Dimensions{}
	}

	if a.Dismissible && a.closeBtn.Clicked(gtx) {
		a.closed = true
		interaction.Invalidate(gtx)
		return layout.Dimensions{}
	}

	accent, _ := variantColors(th, a.Variant)
	if a.Variant == ButtonDefault {
		accent = th.Colors.Primary
	}
	cardBg := th.Colors.Surface

	pad := gtx.Dp(unit.Dp(th.Spacing.SM))
	gap := gtx.Dp(unit.Dp(th.Spacing.XS))
	borderW := gtx.Dp(unit.Dp(th.Shape.BorderWidth))
	accentW := gtx.Dp(unit.Dp(alertAccentDp))
	iconPx := gtx.Dp(unit.Dp(alertIconDp))
	closePx := gtx.Dp(unit.Dp(alertCloseDp))

	maxW := a.MaxWidth
	if maxW == 0 {
		maxW = defaultAlertMaxWidthDp
	}
	maxInnerW := gtx.Dp(unit.Dp(maxW)) - 2*borderW

	contentW := maxInnerW - 2*pad - accentW - gap
	if a.Icon != nil {
		contentW -= iconPx + gap
	}
	if a.Dismissible {
		contentW -= closePx + gap
	}
	if contentW < 0 {
		contentW = 0
	}

	mgtx := gtx
	mgtx.Constraints = layout.Constraints{Min: image.Point{}, Max: image.Pt(contentW, gtx.Constraints.Max.Y)}

	var titleW, msgW layout.Widget
	var titleDims, msgDims layout.Dimensions
	if a.Title != "" {
		titleW = Text(th, a.Title, WithWeight(font.Bold))
		m := op.Record(mgtx.Ops)
		titleDims = titleW(mgtx)
		m.Stop()
	}
	if a.Message != "" {
		msgW = Text(th, a.Message, WithSize(th.Typography.FontSize*0.9), WithColor(th.Colors.Muted))
		m := op.Record(mgtx.Ops)
		msgDims = msgW(mgtx)
		m.Stop()
	}

	textW := titleDims.Size.X
	if msgDims.Size.X > textW {
		textW = msgDims.Size.X
	}
	textH := titleDims.Size.Y + msgDims.Size.Y
	if a.Title != "" && a.Message != "" {
		textH += gap
	}

	innerW := 2*pad + accentW + gap + textW
	if a.Icon != nil {
		innerW += iconPx + gap
	}
	if a.Dismissible {
		innerW += closePx + gap
	}
	if innerW > maxInnerW {
		innerW = maxInnerW
	}

	innerH := textH + 2*pad
	if a.Icon != nil && iconPx+2*pad > innerH {
		innerH = iconPx + 2*pad
	}
	if a.Dismissible && closePx+2*pad > innerH {
		innerH = closePx + 2*pad
	}

	frame := NeoFrame(gtx, th, innerW, innerH, cardBg, false, false)

	accentRect := image.Rectangle{
		Min: frame.BgRect.Min,
		Max: image.Pt(frame.BgRect.Min.X+accentW, frame.BgRect.Max.Y),
	}
	paint.FillShape(gtx.Ops, accent, clip.Rect(accentRect).Op())

	cursor := frame.BgRect.Min.Add(image.Pt(accentW+gap, pad))

	if a.Icon != nil {
		iconStack := op.Offset(cursor).Push(gtx.Ops)
		igtx := gtx
		igtx.Constraints.Min.X = iconPx
		a.Icon.Layout(igtx, accent)
		iconStack.Pop()
		cursor = cursor.Add(image.Pt(iconPx+gap, 0))
	}

	textStack := op.Offset(cursor).Push(gtx.Ops)
	ty := 0
	if a.Title != "" {
		off := op.Offset(image.Pt(0, ty)).Push(gtx.Ops)
		titleW(mgtx)
		off.Pop()
		ty += titleDims.Size.Y + gap
	}
	if a.Message != "" {
		off := op.Offset(image.Pt(0, ty)).Push(gtx.Ops)
		msgW(mgtx)
		off.Pop()
	}
	textStack.Pop()

	if a.Dismissible {
		closeOff := image.Pt(frame.BgRect.Max.X-pad-closePx, frame.BgRect.Min.Y+pad)
		interaction.ClickArea(gtx, &a.closeBtn, image.Rectangle{Min: closeOff, Max: closeOff.Add(image.Pt(closePx, closePx))})

		xW := Text(th, "×", WithWeight(font.Bold))
		xmgtx := gtx
		xmgtx.Constraints.Min = image.Point{}
		xm := op.Record(xmgtx.Ops)
		xDims := xW(xmgtx)
		xm.Stop()

		xOff := closeOff.Add(image.Pt((closePx-xDims.Size.X)/2, (closePx-xDims.Size.Y)/2))
		xStack := op.Offset(xOff).Push(gtx.Ops)
		xW(xmgtx)
		xStack.Pop()
	}

	return layout.Dimensions{Size: frame.CanvasDim}
}
