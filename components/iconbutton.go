package components

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type IconButton struct {
	Clickable widget.Clickable
	Icon      *widget.Icon
	Variant   ButtonVariant
	Size      ButtonSize
	Disabled  bool
	clicked   bool

	TooltipText      string
	TooltipPlacement TooltipPlacement
}

func (b *IconButton) Clicked() bool { return b.clicked }

func (b *IconButton) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	clicked := false
	if !b.Disabled {
		clicked = b.Clickable.Clicked(gtx)
		if clicked {
			interaction.ClearFocus(gtx)
			interaction.Invalidate(gtx)
		}
	}
	b.clicked = clicked

	bg, fg := variantColors(th, b.Variant)
	if b.Disabled {
		bg, fg = th.Colors.Muted, th.Colors.Foreground
	}
	hovered := !b.Disabled && b.Clickable.Hovered()
	pressed := !b.Disabled && b.Clickable.Pressed()

	content := Icon(th, b.Icon, WithIconSize(iconSizeFor(b.Size)), WithIconColor(fg))

	dims := layoutButtonFrame(gtx, th, content, paddingFor(th, b.Size), bg, hovered, pressed, b.Disabled, &b.Clickable)

	if hovered && b.TooltipText != "" {
		LayoutTooltipBubble(gtx, th, b.TooltipText, b.TooltipPlacement, dims.Size)
	}

	return dims
}

func iconSizeFor(s ButtonSize) float32 {
	switch s {
	case SizeMD:
		return 20
	case SizeLG:
		return 24
	default:
		return 16
	}
}
