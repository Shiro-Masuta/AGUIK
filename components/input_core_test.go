package components

import (
	"image"
	"testing"
)

func TestConstrainInputDimsRespectsAvailableSpace(t *testing.T) {
	got := constrainInputDims(
		inputDims{width: 240, height: 80},
		5,
		image.Pt(180, 60),
		image.Pt(20, 20),
	)
	if want := (inputDims{width: 150, height: 30}); got != want {
		t.Fatalf("constrainInputDims() = %+v, want %+v", got, want)
	}
}

func TestConstrainInputDimsHandlesParentSmallerThanFrame(t *testing.T) {
	got := constrainInputDims(
		inputDims{width: 100, height: 40},
		5,
		image.Pt(10, 10),
		image.Pt(20, 20),
	)
	if want := (inputDims{}); got != want {
		t.Fatalf("constrainInputDims() = %+v, want %+v", got, want)
	}
}
