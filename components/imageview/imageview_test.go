package imageview

import (
	"image"
	"image/color"
	"testing"
)

type sliceBackedImage struct {
	pixels []byte
}

func (sliceBackedImage) ColorModel() color.Model { return color.RGBAModel }

func (sliceBackedImage) Bounds() image.Rectangle { return image.Rect(0, 0, 2, 2) }

func (sliceBackedImage) At(int, int) color.Color { return color.NRGBA{A: 255} }

func TestImageViewFitFillScalesHighlightsAndHitTestsIndependently(t *testing.T) {
	view := ImageView{Fit: FitFill}
	view.computeFit(image.Pt(200, 100), image.Pt(100, 100))
	view.Highlights = []Highlight{{Rect: image.Rect(10, 10, 20, 20)}}

	if view.scaleX != 2 || view.scaleY != 1 {
		t.Fatalf("FitFill scale = (%v, %v), want (2, 1)", view.scaleX, view.scaleY)
	}
	if got, want := view.toScreen(view.Highlights[0].Rect), image.Rect(20, 10, 40, 20); got != want {
		t.Fatalf("toScreen() = %v, want %v", got, want)
	}
	if got := view.HitTest(image.Pt(30, 15)); got != 0 {
		t.Fatalf("HitTest() = %d, want 0", got)
	}
}

func TestImageViewEnsureOpAcceptsNonComparableImage(t *testing.T) {
	view := ImageView{Src: sliceBackedImage{pixels: []byte{1, 2, 3}}}
	view.ensureOp()
	view.ensureOp()
}
