package selectwidget

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/components"
	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

const defaultSelectWidthDp = 200

type SelectOption struct {
	Label string
}

type Select struct {
	Options     []SelectOption
	Selected    int
	Placeholder string
	Width       float32
	Size        components.ButtonSize
	Disabled    bool

	open       bool
	trigger    widget.Clickable
	optionBtns []widget.Clickable
}

func (s *Select) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if s.Selected < -1 || s.Selected >= len(s.Options) {
		s.Selected = -1
	}
	s.optionBtns = glayout.Resize(s.optionBtns, len(s.Options))

	if !s.Disabled && s.trigger.Clicked(gtx) {
		s.open = !s.open
		interaction.Invalidate(gtx)
	}

	label := s.selectedLabel()
	frame := s.layoutTrigger(gtx, th, label)

	if !s.Disabled {
		interaction.ClickArea(gtx, &s.trigger, frame.RestRect)
	}
	if s.open && len(s.Options) > 0 {
		s.layoutDropdown(gtx, th, frame.CanvasDim)
	}

	return layout.Dimensions{Size: gtx.Constraints.Constrain(frame.CanvasDim)}
}

func (s *Select) selectedLabel() string {
	if s.Selected >= 0 {
		return s.Options[s.Selected].Label
	}
	return s.Placeholder
}
