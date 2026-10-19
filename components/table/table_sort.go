package table

import (
	"sort"

	"gioui.org/layout"
)

type SortDirection int

const (
	SortNone SortDirection = iota
	SortAsc
	SortDesc
)

// nextSort calcule l'état de tri après un clic sur l'en-tête de la colonne
// clicked. Cycle : aucun -> croissant -> décroissant -> aucun.
func nextSort(col int, dir SortDirection, clicked int) (int, SortDirection) {
	switch {
	case clicked != col || dir == SortNone:
		return clicked, SortAsc
	case dir == SortAsc:
		return clicked, SortDesc
	default:
		return -1, SortNone
	}
}

func (t *Table[T]) handleSortClicks(gtx layout.Context) {
	for i, col := range t.Columns {
		if col.Less != nil && t.headerBtns[i].Clicked(gtx) {
			t.SortColumn, t.SortDirection = nextSort(t.SortColumn, t.SortDirection, i)
		}
	}
}

func (t *Table[T]) sortActive() bool {
	return t.SortDirection != SortNone &&
		t.SortColumn >= 0 && t.SortColumn < len(t.Columns) &&
		t.Columns[t.SortColumn].Less != nil
}

// rebuildOrder recalcule t.order uniquement si le tri, la longueur de Rows ou
// un appel à Resort l'exigent. Rows n'est jamais modifié.
func (t *Table[T]) rebuildOrder() {
	n := len(t.Rows)
	if !t.stale && len(t.order) == n && t.appliedCol == t.SortColumn && t.appliedDir == t.SortDirection {
		return
	}

	t.order = t.order[:0]
	for i := range n {
		t.order = append(t.order, i)
	}

	if t.sortActive() {
		less := t.Columns[t.SortColumn].Less
		desc := t.SortDirection == SortDesc
		sort.SliceStable(t.order, func(a, b int) bool {
			x, y := t.Rows[t.order[a]], t.Rows[t.order[b]]
			if desc {
				return less(y, x)
			}
			return less(x, y)
		})
	}

	t.appliedCol, t.appliedDir, t.stale = t.SortColumn, t.SortDirection, false
}
