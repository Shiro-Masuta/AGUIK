package interaction

import (
	"image"
	"testing"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/internal/uitest"
)

func TestFocusClearerClearsFocusOnBackgroundClick(t *testing.T) {
	h := uitest.New(image.Pt(200, 100))
	var clearer FocusClearer
	var target widget.Clickable
	draw := func(gtx layout.Context) layout.Dimensions {
		clearer.Layout(gtx)
		stack := op.Offset(image.Pt(0, 0)).Push(gtx.Ops)
		target.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Pt(20, 20)}
		})
		stack.Pop()
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}

	h.Frame(func(gtx layout.Context) layout.Dimensions {
		gtx.Execute(key.FocusCmd{Tag: &target})
		return draw(gtx)
	})

	focused := false
	h.Frame(func(gtx layout.Context) layout.Dimensions {
		focused = gtx.Focused(&target)
		return draw(gtx)
	})
	if !focused {
		t.Fatal("le tag de test devrait avoir le focus avant le clic")
	}

	h.Click(image.Pt(180, 80))
	h.Frame(draw)

	h.Frame(func(gtx layout.Context) layout.Dimensions {
		focused = gtx.Focused(&target)
		return draw(gtx)
	})
	if focused {
		t.Fatal("le clic sur le fond devrait retirer le focus")
	}
}
