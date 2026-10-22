package selectwidget

import (
	"image"
	"testing"

	"gioui.org/layout"

	"github.com/Shiro-Masuta/AGUIK/internal/uitest"
)

func TestSelectOpensAndSelectsOption(t *testing.T) {
	h := uitest.New(image.Pt(400, 400))
	th := &h.Theme
	s := Select{
		Options:  []SelectOption{{Label: "A"}, {Label: "B"}, {Label: "C"}},
		Selected: -1,
		Width:    200,
	}
	draw := func(gtx layout.Context) layout.Dimensions { return s.Layout(gtx, th) }

	trigger := h.Frame(draw).Size
	h.Click(trigger.Div(2))
	h.Frame(draw)
	if !s.open {
		t.Fatal("le clic sur le trigger doit ouvrir la liste")
	}

	var rowH int
	uitest.New(h.Max).Frame(func(gtx layout.Context) layout.Dimensions {
		dims := drawSelectRow(gtx, th, "B", false, int(th.Spacing.XS))
		rowH = dims.Size.Y
		return dims
	})

	const option = 1
	panelTop := trigger.Y + 2 + int(th.Shape.BorderWidth) + int(th.Spacing.SM)
	h.Click(image.Pt(trigger.X/2, panelTop+option*rowH+rowH/2))
	h.Frame(draw)

	if s.Selected != option {
		t.Fatalf("Selected = %d, want %d", s.Selected, option)
	}
	if s.open {
		t.Fatal("le menu doit se fermer après la sélection")
	}
}

func TestSelectClampsToParentWidth(t *testing.T) {
	h := uitest.New(image.Pt(90, 160))
	s := Select{Placeholder: "Select", Width: 240}
	dims := h.Frame(func(gtx layout.Context) layout.Dimensions {
		return s.Layout(gtx, &h.Theme)
	})
	if dims.Size.X > h.Max.X {
		t.Fatalf("Select width = %d, exceeds parent width %d", dims.Size.X, h.Max.X)
	}
}
