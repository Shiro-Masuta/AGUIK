package components

import (
	"fmt"
	"image"
	"testing"

	"gioui.org/layout"

	"github.com/Shiro-Masuta/AGUIK/internal/uitest"
)

func TestNeoFrameCanvasMatchesSharedOverhead(t *testing.T) {
	states := []struct {
		name             string
		hovered, pressed bool
	}{
		{"repos", false, false},
		{"survol", true, false},
		{"appui", false, true},
	}

	for _, scale := range []float32{1, 1.75, 2.5} {
		h := uitest.New(image.Pt(400, 400))
		h.PxPerDp = scale
		th := &h.Theme

		for _, st := range states {
			t.Run(fmt.Sprintf("x%v/%s", scale, st.name), func(t *testing.T) {
				h.Frame(func(gtx layout.Context) layout.Dimensions {
					f := NeoFrame(gtx, th, 100, 40, th.Colors.Surface, st.hovered, st.pressed)
					want := image.Pt(100, 40).Add(NeoFrameOverhead(gtx, th))
					if f.CanvasDim != want {
						t.Errorf("CanvasDim = %v, want %v", f.CanvasDim, want)
					}
					return layout.Dimensions{}
				})
			})
		}
	}
}
