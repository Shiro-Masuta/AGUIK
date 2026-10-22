package table

import (
	"image"
	"slices"
	"strconv"
	"testing"

	"gioui.org/layout"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/Shiro-Masuta/AGUIK/internal/uitest"
)

func newIntTable() *Table[int] {
	return &Table[int]{
		Columns: []TableColumn[int]{{
			Label: "N",
			Cell:  strconv.Itoa,
			Less:  func(a, b int) bool { return a < b },
		}},
		Rows:       []int{3, 1, 2},
		SortColumn: -1,
	}
}

func TestHeaderClickCyclesSort(t *testing.T) {
	h := uitest.New(image.Pt(500, 300))
	th := &h.Theme
	tb := newIntTable()
	draw := func(gtx layout.Context) layout.Dimensions { return tb.Layout(gtx, th) }

	// Dans la première cellule d'en-tête : après le padding de plainRow.
	header := image.Pt(int(th.Spacing.MD)+4, int(th.Spacing.SM)+2)
	h.Frame(draw)

	steps := []struct {
		dir   SortDirection
		order []int
	}{
		{SortAsc, []int{1, 2, 0}},
		{SortDesc, []int{0, 2, 1}},
		{SortNone, []int{0, 1, 2}},
	}
	for _, st := range steps {
		h.Click(header)
		h.Frame(draw)
		if tb.SortDirection != st.dir {
			t.Fatalf("SortDirection = %d, want %d", tb.SortDirection, st.dir)
		}
		if !slices.Equal(tb.order, st.order) {
			t.Fatalf("order = %v, want %v", tb.order, st.order)
		}
	}
	if !slices.Equal(tb.Rows, []int{3, 1, 2}) {
		t.Fatalf("Rows modifié par le tri : %v", tb.Rows)
	}
}

func TestApplyMoveSwapsRows(t *testing.T) {
	h := uitest.New(image.Pt(500, 300))
	tb := newIntTable()
	apply := func() {
		h.Frame(func(gtx layout.Context) layout.Dimensions {
			tb.applyMove(gtx)
			return layout.Dimensions{}
		})
	}

	tb.requestMove(0, 1)
	apply()
	if want := []int{1, 3, 2}; !slices.Equal(tb.Rows, want) {
		t.Fatalf("Rows = %v, want %v", tb.Rows, want)
	}

	for _, mv := range [][2]int{{0, -1}, {2, 3}, {-1, 0}} {
		tb.requestMove(mv[0], mv[1])
		apply()
		if want := []int{1, 3, 2}; !slices.Equal(tb.Rows, want) {
			t.Fatalf("déplacement hors limites %v appliqué : %v", mv, tb.Rows)
		}
	}
}

func TestReorderButtonClickMovesRow(t *testing.T) {
	h := uitest.New(image.Pt(500, 300))
	up, err := widget.NewIcon(icons.NavigationArrowUpward)
	if err != nil {
		t.Fatal(err)
	}
	down, err := widget.NewIcon(icons.NavigationArrowDownward)
	if err != nil {
		t.Fatal(err)
	}
	tb := newIntTable()
	tb.SortColumn = -1
	tb.Reorderable = true
	tb.MoveUpIcon = up
	tb.MoveDownIcon = down
	draw := func(gtx layout.Context) layout.Dimensions { return tb.Layout(gtx, &h.Theme) }

	h.Frame(draw)
	tb.reorderBtns[0].down.Clickable.Click()
	h.Frame(draw)
	if want := []int{1, 3, 2}; !slices.Equal(tb.Rows, want) {
		t.Fatalf("Rows after row's down click = %v, want %v", tb.Rows, want)
	}
}
