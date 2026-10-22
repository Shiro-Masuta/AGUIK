package overlay

import (
	"image"
	"testing"

	"gioui.org/layout"

	"github.com/Shiro-Masuta/AGUIK/internal/uitest"
)

func TestDialogActionClickReportsIndexAndAutoCloses(t *testing.T) {
	h := uitest.New(image.Pt(400, 300))
	d := &Dialog{
		Open:      true,
		Title:     "Confirm",
		AutoClose: true,
		Actions:   []DialogAction{{Text: "Cancel"}, {Text: "OK"}},
	}
	draw := func(gtx layout.Context) layout.Dimensions { return d.Layout(gtx, &h.Theme) }
	h.Frame(draw)
	d.actionBtns[1].Clickable.Click()
	h.Frame(draw)
	if d.Clicked() != 1 {
		t.Fatalf("Clicked() = %d, want 1", d.Clicked())
	}
	if d.Open {
		t.Fatal("dialog should close after an action when AutoClose is enabled")
	}
	h.Frame(draw)
	if d.Clicked() != -1 {
		t.Fatalf("Clicked() after next frame = %d, want -1", d.Clicked())
	}
}
