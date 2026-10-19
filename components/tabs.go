package components

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type TabItem struct {
	Label string
}

type Tabs struct {
	Items    []TabItem
	Active   int
	Variant  ButtonVariant
	Size     ButtonSize
	Disabled bool

	btns []widget.Clickable
}

func (t *Tabs) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if t.Active < 0 || t.Active >= len(t.Items) {
		t.Active = 0
	}
	t.btns = glayout.Resize(t.btns, len(t.Items))

	activeBg, activeFg := variantColors(th, t.Variant)
	if t.Variant == ButtonDefault {
		activeBg, activeFg = th.Colors.Primary, th.Colors.Foreground
	}

	for i := range t.Items {
		if !t.Disabled && t.btns[i].Clicked(gtx) {
			t.Active = i
		}
	}

	children := make([]glayout.Child, len(t.Items))
	for i, item := range t.Items {
		bg, fg := th.Colors.Surface, th.Colors.Foreground
		if i == t.Active {
			bg, fg = activeBg, activeFg
		}
		if t.Disabled {
			bg, fg = th.Colors.Muted, th.Colors.Foreground
		}

		hovered := !t.Disabled && t.btns[i].Hovered()
		pressed := !t.Disabled && t.btns[i].Pressed()
		content := Text(th, item.Label, WithColor(fg), WithWeight(font.SemiBold))

		children[i] = glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutButtonFrame(gtx, th, content, paddingFor(th, t.Size), bg, hovered, pressed, t.Disabled, &t.btns[i])
		})
	}

	return glayout.Row(th.Spacing.XS, layout.Middle, children...)(gtx)
}
