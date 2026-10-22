package glayout

import (
	"image"
	"image/color"
	"testing"

	"gioui.org/layout"

	"github.com/Shiro-Masuta/AGUIK/internal/uitest"
)

func fixed(size image.Point) layout.Widget {
	return func(layout.Context) layout.Dimensions { return layout.Dimensions{Size: size} }
}

func TestRowAndColumnDimensions(t *testing.T) {
	h := uitest.New(image.Pt(200, 120))

	row := Row(5, layout.Start, Rigid(fixed(image.Pt(30, 10))), Rigid(fixed(image.Pt(50, 20))))
	if got, want := h.Frame(row).Size, image.Pt(85, 20); got != want {
		t.Fatalf("Row size = %v, want %v", got, want)
	}

	column := Column(4, layout.Start, Rigid(fixed(image.Pt(30, 10))), Rigid(fixed(image.Pt(50, 20))))
	if got, want := h.Frame(column).Size, image.Pt(50, 34); got != want {
		t.Fatalf("Column size = %v, want %v", got, want)
	}
}

func TestFlexedChildReceivesRemainingSpace(t *testing.T) {
	h := uitest.New(image.Pt(200, 80))
	w := Row(0, layout.Start,
		Rigid(fixed(image.Pt(40, 10))),
		Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
	)

	if got, want := h.Frame(w).Size, image.Pt(200, 80); got != want {
		t.Fatalf("row size = %v, want %v", got, want)
	}
}

func TestGridHandlesIncompleteLastRow(t *testing.T) {
	h := uitest.New(image.Pt(120, 100))
	var laidOut int
	cell := func(gtx layout.Context) layout.Dimensions {
		laidOut++
		return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, 10)}
	}

	got := h.Frame(Grid(5, 2, cell, cell, cell))
	if want := image.Pt(120, 25); got.Size != want {
		t.Fatalf("Grid size = %v, want %v", got.Size, want)
	}
	if laidOut != 3 {
		t.Fatalf("laid out %d cells, want 3", laidOut)
	}
}

func TestInsetStackAndStretchWidth(t *testing.T) {
	h := uitest.New(image.Pt(120, 80))
	inset := Inset(10, func(gtx layout.Context) layout.Dimensions {
		if got, want := gtx.Constraints.Max, image.Pt(100, 60); got != want {
			t.Errorf("inset child max = %v, want %v", got, want)
		}
		return layout.Dimensions{Size: gtx.Constraints.Max}
	})
	if got, want := h.Frame(inset).Size, image.Pt(120, 80); got != want {
		t.Fatalf("Inset size = %v, want %v", got, want)
	}

	expanded := func(gtx layout.Context) layout.Dimensions {
		if got, want := gtx.Constraints.Min, image.Pt(20, 10); got != want {
			t.Errorf("expanded min = %v, want %v", got, want)
		}
		if got, want := gtx.Constraints.Max, image.Pt(120, 80); got != want {
			t.Errorf("expanded max = %v, want %v", got, want)
		}
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
	stack := Stack(fixed(image.Pt(20, 10)), expanded)
	if got, want := h.Frame(stack).Size, image.Pt(20, 10); got != want {
		t.Fatalf("Stack size = %v, want %v", got, want)
	}

	stretch := StretchWidth(fixed(image.Pt(20, 7)))
	if got, want := h.Frame(stretch).Size, image.Pt(120, 7); got != want {
		t.Fatalf("StretchWidth size = %v, want %v", got, want)
	}
}

func TestGhostMeasuresWithoutPaintingOrRegisteringEvents(t *testing.T) {
	h := uitest.New(image.Pt(100, 60))
	called := false
	w := Ghost(func(gtx layout.Context) layout.Dimensions {
		called = true
		return layout.Dimensions{Size: image.Pt(35, 12)}
	})

	if got, want := h.Frame(w).Size, image.Pt(35, 12); got != want {
		t.Fatalf("Ghost size = %v, want %v", got, want)
	}
	if !called {
		t.Fatal("Ghost did not run the widget for measurement")
	}
}

func TestResizePreservesElements(t *testing.T) {
	original := []int{3, 5, 8}
	shrunk := Resize(original, 2)
	if len(shrunk) != 2 || shrunk[1] != 5 {
		t.Fatalf("shrunk = %v, want [3 5]", shrunk)
	}
	grown := Resize(shrunk, 4)
	if len(grown) != 4 || grown[0] != 3 || grown[1] != 5 || grown[2] != 0 {
		t.Fatalf("grown = %v, want [3 5 0 0]", grown)
	}
}

func TestColorUtilities(t *testing.T) {
	base := color.NRGBA{R: 0, G: 20, B: 40, A: 17}
	tint := color.NRGBA{R: 100, G: 120, B: 140, A: 240}
	if got, want := BlendColor(base, tint, 0.5), (color.NRGBA{R: 50, G: 70, B: 90, A: 255}); got != want {
		t.Fatalf("BlendColor = %#v, want %#v", got, want)
	}
	if got, want := WithAlpha(base, 128), (color.NRGBA{R: 0, G: 20, B: 40, A: 128}); got != want {
		t.Fatalf("WithAlpha = %#v, want %#v", got, want)
	}
}
