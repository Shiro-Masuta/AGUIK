package table

import (
	"slices"
	"testing"
)

func TestNextSort(t *testing.T) {
	tests := []struct {
		name    string
		col     int
		dir     SortDirection
		clicked int
		wantCol int
		wantDir SortDirection
	}{
		{"premier clic", -1, SortNone, 1, 1, SortAsc},
		{"zéro value, colonne 0", 0, SortNone, 0, 0, SortAsc},
		{"croissant vers décroissant", 1, SortAsc, 1, 1, SortDesc},
		{"décroissant vers aucun", 1, SortDesc, 1, -1, SortNone},
		{"autre colonne repart en croissant", 1, SortDesc, 2, 2, SortAsc},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			col, dir := nextSort(tc.col, tc.dir, tc.clicked)
			if col != tc.wantCol || dir != tc.wantDir {
				t.Fatalf("got (%d, %d), want (%d, %d)", col, dir, tc.wantCol, tc.wantDir)
			}
		})
	}
}

func TestRebuildOrderKeepsRows(t *testing.T) {
	tb := Table[int]{
		Columns: []TableColumn[int]{{Less: func(a, b int) bool { return a < b }}},
		Rows:    []int{3, 1, 2},
	}
	original := slices.Clone(tb.Rows)

	check := func(name string, want []int) {
		t.Helper()
		tb.rebuildOrder()
		if !slices.Equal(tb.order, want) {
			t.Fatalf("%s: order = %v, want %v", name, tb.order, want)
		}
		if !slices.Equal(tb.Rows, original) {
			t.Fatalf("%s: Rows modifié : %v", name, tb.Rows)
		}
	}

	check("sans tri", []int{0, 1, 2})

	tb.SortColumn, tb.SortDirection = 0, SortAsc
	check("croissant", []int{1, 2, 0})

	tb.SortDirection = SortDesc
	check("décroissant", []int{0, 2, 1})

	tb.SortColumn, tb.SortDirection = -1, SortNone
	check("retour à l'ordre d'origine", []int{0, 1, 2})
}
