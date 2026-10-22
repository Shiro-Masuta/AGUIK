package components

import (
	"image"
	"testing"

	"gioui.org/layout"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/Shiro-Masuta/AGUIK/internal/uitest"
)

func testIcon(t *testing.T) *widget.Icon {
	t.Helper()
	ic, err := widget.NewIcon(icons.NavigationArrowUpward)
	if err != nil {
		t.Fatal(err)
	}
	return ic
}

func center(size image.Point) image.Point { return size.Div(2) }

func TestIconButtonClickedLastsOneFrame(t *testing.T) {
	h := uitest.New(image.Pt(300, 300))
	b := IconButton{Icon: testIcon(t)}
	draw := func(gtx layout.Context) layout.Dimensions { return b.Layout(gtx, &h.Theme) }

	dims := h.Frame(draw)
	h.Click(center(dims.Size))

	h.Frame(draw)
	if !b.Clicked() {
		t.Fatal("Clicked = false après un clic")
	}
	h.Frame(draw)
	if b.Clicked() {
		t.Fatal("Clicked doit retomber à false à la frame suivante")
	}
}

func TestIconButtonDisabledIgnoresClicks(t *testing.T) {
	h := uitest.New(image.Pt(300, 300))
	b := IconButton{Icon: testIcon(t), Disabled: true}
	draw := func(gtx layout.Context) layout.Dimensions { return b.Layout(gtx, &h.Theme) }

	dims := h.Frame(draw)
	h.Click(center(dims.Size))
	h.Frame(draw)
	if b.Clicked() {
		t.Fatal("un bouton désactivé ne doit pas réagir")
	}
}

func TestButtonClickAndDisabledState(t *testing.T) {
	h := uitest.New(image.Pt(300, 200))
	b := Button{Text: "Save"}
	draw := func(gtx layout.Context) layout.Dimensions { return b.Layout(gtx, &h.Theme) }
	h.Frame(draw)
	b.Clickable.Click()
	h.Frame(draw)
	if !b.Clicked() {
		t.Fatal("Button.Clicked() should report its programmatic click")
	}

	b.Disabled = true
	b.Clickable.Click()
	h.Frame(draw)
	if b.Clicked() {
		t.Fatal("disabled button must ignore clicks")
	}
}

func TestCheckboxAndSwitchToggleOnPointerClick(t *testing.T) {
	tests := []struct {
		name string
		view func(*testing.T, *uitest.Harness) (layout.Widget, *bool)
	}{
		{
			name: "checkbox",
			view: func(t *testing.T, h *uitest.Harness) (layout.Widget, *bool) {
				c := &Checkbox{}
				return func(gtx layout.Context) layout.Dimensions { return c.Layout(gtx, &h.Theme) }, &c.Bool.Value
			},
		},
		{
			name: "switch",
			view: func(t *testing.T, h *uitest.Harness) (layout.Widget, *bool) {
				s := &Switch{}
				return func(gtx layout.Context) layout.Dimensions { return s.Layout(gtx, &h.Theme) }, &s.Bool.Value
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := uitest.New(image.Pt(300, 200))
			draw, value := tc.view(t, h)
			dims := h.Frame(draw)
			h.Click(center(dims.Size))
			h.Frame(draw)
			if !*value {
				t.Fatal("control did not toggle on click")
			}
		})
	}
}

func TestRadioGroupAndTabsSelection(t *testing.T) {
	t.Run("radio group", func(t *testing.T) {
		h := uitest.New(image.Pt(300, 200))
		r := &RadioGroup{Options: []RadioOption{{Label: "First"}, {Label: "Second"}}, Selected: -1}
		draw := func(gtx layout.Context) layout.Dimensions { return r.Layout(gtx, &h.Theme) }
		h.Frame(draw)
		r.btns[1].Click()
		h.Frame(draw)
		if r.Selected != 1 {
			t.Fatalf("Selected = %d, want 1", r.Selected)
		}
		r.Disabled = true
		r.btns[0].Click()
		h.Frame(draw)
		if r.Selected != 1 {
			t.Fatalf("disabled RadioGroup changed selection to %d", r.Selected)
		}
	})

	t.Run("tabs", func(t *testing.T) {
		h := uitest.New(image.Pt(300, 200))
		tabs := &Tabs{Items: []TabItem{{Label: "First"}, {Label: "Second"}}}
		draw := func(gtx layout.Context) layout.Dimensions { return tabs.Layout(gtx, &h.Theme) }
		h.Frame(draw)
		tabs.btns[1].Click()
		h.Frame(draw)
		if tabs.Active != 1 {
			t.Fatalf("Active = %d, want 1", tabs.Active)
		}
		tabs.Disabled = true
		tabs.btns[0].Click()
		h.Frame(draw)
		if tabs.Active != 1 {
			t.Fatalf("disabled Tabs changed active index to %d", tabs.Active)
		}
	})
}

func TestSidebarClickChangesActive(t *testing.T) {
	h := uitest.New(image.Pt(300, 400))
	th := &h.Theme
	icon := testIcon(t)

	var probe IconButton
	probe.Icon = icon
	btn := h.Frame(func(gtx layout.Context) layout.Dimensions { return probe.Layout(gtx, th) }).Size

	s := Sidebar{Items: []SidebarItem{{Icon: icon}, {Icon: icon}, {Icon: icon}}}
	draw := func(gtx layout.Context) layout.Dimensions { return s.Layout(gtx, th) }
	h.Frame(draw)

	// Origine du contenu : marge du cadre + bordure + padding de la sidebar.
	origin := hoverLiftDp + int(th.Shape.BorderWidth) + int(th.Spacing.SM)
	gap := int(th.Spacing.SM)
	second := image.Pt(origin+btn.X/2, origin+btn.Y+gap+btn.Y/2)

	h.Click(second)
	h.Frame(draw)
	if s.Active != 1 {
		t.Fatalf("Active = %d, want 1", s.Active)
	}
}
