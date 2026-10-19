package table

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/components"
	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

// TableColumn décrit une colonne de Table[T].
type TableColumn[T any] struct {
	Label  string
	Weight float32 // largeur relative ; 0 vaut 1
	Cell   func(item T) string
	Less   func(a, b T) bool // nil : colonne non triable
}

// Table affiche des lignes de type T avec en-tête triable, zébrage et
// réordonnancement optionnel. Sa valeur zéro est utilisable. Elle contient de
// l'état Gio : elle doit être conservée entre les frames.
//
// Le tri ne modifie jamais Rows : il produit un ordre d'affichage. Revenir à
// SortNone restaure donc l'ordre d'origine.
type Table[T any] struct {
	Columns []TableColumn[T]
	Rows    []T

	SortColumn    int
	SortDirection SortDirection

	// Reorderable affiche des boutons haut/bas, actifs seulement hors tri et
	// si MoveUpIcon et MoveDownIcon sont fournies (AGUIK n'a pas de catalogue
	// d'icônes).
	Reorderable  bool
	MoveUpIcon   *widget.Icon
	MoveDownIcon *widget.Icon

	// EmptyText est affiché quand Rows est vide.
	EmptyText string

	// RowActions construit le widget affiché à droite d'une ligne. index est
	// l'indice dans Rows, stable même quand le tri est actif.
	RowActions func(item T, index int) layout.Widget

	headerBtns  []widget.Clickable
	reorderBtns []tableReorderButtons
	ghostBtns   tableReorderButtons // boutons jetables pour mesurer l'en-tête
	list        widget.List

	order      []int // position affichée -> indice dans Rows
	appliedCol int
	appliedDir SortDirection
	stale      bool

	moving           bool
	moveFrom, moveTo int
}

type tableReorderButtons struct {
	up, down components.IconButton
}

// Resort force le recalcul de l'ordre d'affichage. À appeler après avoir
// modifié le contenu de Rows sans en changer la longueur.
func (t *Table[T]) Resort() { t.stale = true }

// Layout dessine la table. Elle attend une hauteur maximale finie.
func (t *Table[T]) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	t.headerBtns = glayout.Resize(t.headerBtns, len(t.Columns))
	t.reorderBtns = glayout.Resize(t.reorderBtns, len(t.Rows))
	t.handleSortClicks(gtx)
	t.rebuildOrder()

	reorder := t.reorderEnabled()
	stripe := glayout.BlendColor(th.Colors.Background, th.Colors.Muted, 0.15)

	t.list.Axis = layout.Vertical
	body := func(gtx layout.Context) layout.Dimensions {
		if len(t.order) == 0 && t.EmptyText != "" {
			return glayout.Inset(th.Spacing.MD, components.Caption(th, t.EmptyText))(gtx)
		}
		return t.list.Layout(gtx, len(t.order), func(gtx layout.Context, pos int) layout.Dimensions {
			bg := th.Colors.Background
			if pos%2 == 1 {
				bg = stripe
			}
			return t.layoutRow(gtx, th, bg, reorder, pos)
		})
	}

	content := glayout.Column(0, layout.Start,
		glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.layoutHeader(gtx, th, reorder)
		}),
		glayout.Flexed(1, body),
	)

	macro := op.Record(gtx.Ops)
	dims := content(gtx)
	call := macro.Stop()

	size := image.Pt(gtx.Constraints.Max.X, dims.Size.Y)
	area := clip.Rect{Max: size}.Push(gtx.Ops)
	call.Add(gtx.Ops)
	area.Pop()

	t.applyMove(gtx)
	return layout.Dimensions{Size: size}
}

func (t *Table[T]) layoutHeader(gtx layout.Context, th *theme.Theme, reorder bool) layout.Dimensions {
	row := t.rowWidget(th,
		func(i int, col TableColumn[T]) layout.Widget { return t.headerCell(th, i, col) },
		t.ghostTrailing(th, reorder),
	)
	dims := plainRow(gtx, th, th.Colors.SurfaceAlt, row)
	return withRule(gtx, dims, th.Colors.Border, gtx.Dp(unit.Dp(th.Shape.BorderWidth)))
}

func (t *Table[T]) layoutRow(gtx layout.Context, th *theme.Theme, bg color.NRGBA, reorder bool, pos int) layout.Dimensions {
	idx := t.order[pos]
	row := t.Rows[idx]

	trailing := make([]glayout.Child, 0, 2)
	for _, w := range t.trailingWidgets(th, reorder, idx, &t.reorderBtns[idx]) {
		trailing = append(trailing, glayout.Rigid(w))
	}

	content := t.rowWidget(th, func(_ int, col TableColumn[T]) layout.Widget {
		return glayout.StretchWidth(components.Text(th, col.Cell(row), components.WithMaxLines(1), components.WithEllipsis()))
	}, trailing)

	dims := plainRow(gtx, th, bg, content)
	hair := glayout.BlendColor(th.Colors.Background, th.Colors.Border, 0.35)
	return withRule(gtx, dims, hair, max(1, gtx.Dp(1)))
}

// rowWidget compose les colonnes (Flexed) suivies des colonnes de fin (Rigid).
// En-tête et lignes passent par ici : ils ne peuvent pas diverger.
func (t *Table[T]) rowWidget(
	th *theme.Theme,
	cell func(i int, col TableColumn[T]) layout.Widget,
	trailing []glayout.Child,
) layout.Widget {
	children := make([]glayout.Child, 0, len(t.Columns)+len(trailing))
	for i, col := range t.Columns {
		children = append(children, glayout.Flexed(columnWeight(col.Weight), cell(i, col)))
	}
	children = append(children, trailing...)
	return glayout.Row(th.Spacing.MD, layout.Middle, children...)
}

func (t *Table[T]) headerCell(th *theme.Theme, i int, col TableColumn[T]) layout.Widget {
	label := col.Label
	if t.SortColumn == i {
		switch t.SortDirection {
		case SortAsc:
			label += " ▲"
		case SortDesc:
			label += " ▼"
		}
	}
	content := components.Text(th, label, components.WithWeight(font.SemiBold), components.WithColor(th.Colors.Muted))

	return func(gtx layout.Context) layout.Dimensions {
		dims := content(gtx)
		// La hitbox couvre toute la colonne, pas seulement le texte.
		size := image.Pt(gtx.Constraints.Max.X, dims.Size.Y)
		if col.Less != nil {
			interaction.ClickArea(gtx, &t.headerBtns[i], image.Rectangle{Max: size})
		}
		return layout.Dimensions{Size: size}
	}
}

// ghostTrailing réserve dans l'en-tête la largeur exacte des colonnes de fin :
// ce sont les mêmes widgets que dans les lignes, mesurés mais jamais dessinés.
func (t *Table[T]) ghostTrailing(th *theme.Theme, reorder bool) []glayout.Child {
	if len(t.order) == 0 {
		return nil
	}
	ws := t.trailingWidgets(th, reorder, t.order[0], &t.ghostBtns)
	children := make([]glayout.Child, len(ws))
	for i, w := range ws {
		children[i] = glayout.Rigid(glayout.Ghost(w))
	}
	return children
}

func (t *Table[T]) trailingWidgets(th *theme.Theme, reorder bool, idx int, btns *tableReorderButtons) []layout.Widget {
	var ws []layout.Widget
	if reorder {
		ws = append(ws, t.reorderWidget(th, idx, btns))
	}
	if t.RowActions != nil {
		ws = append(ws, t.RowActions(t.Rows[idx], idx))
	}
	return ws
}

func (t *Table[T]) reorderEnabled() bool {
	return t.Reorderable && t.MoveUpIcon != nil && t.MoveDownIcon != nil && !t.sortActive()
}

// reorderWidget construit les boutons haut/bas de la ligne idx. Le clic
// enregistre un déplacement appliqué en fin de frame, jamais pendant le
// parcours de la liste.
func (t *Table[T]) reorderWidget(th *theme.Theme, idx int, b *tableReorderButtons) layout.Widget {
	configureMoveButton(&b.up, t.MoveUpIcon, idx == 0)
	configureMoveButton(&b.down, t.MoveDownIcon, idx == len(t.Rows)-1)
	return glayout.Row(th.Spacing.SM, layout.Middle,
		glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
			dims := b.up.Layout(gtx, th)
			if b.up.Clicked {
				t.requestMove(idx, idx-1)
			}
			return dims
		}),
		glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
			dims := b.down.Layout(gtx, th)
			if b.down.Clicked {
				t.requestMove(idx, idx+1)
			}
			return dims
		}),
	)
}

func (t *Table[T]) requestMove(from, to int) {
	t.moving, t.moveFrom, t.moveTo = true, from, to
}

func (t *Table[T]) applyMove(gtx layout.Context) {
	if !t.moving {
		return
	}
	t.moving = false
	n := len(t.Rows)
	if t.moveFrom < 0 || t.moveFrom >= n || t.moveTo < 0 || t.moveTo >= n {
		return
	}
	t.Rows[t.moveFrom], t.Rows[t.moveTo] = t.Rows[t.moveTo], t.Rows[t.moveFrom]
	interaction.Invalidate(gtx)
}

// configureMoveButton met à jour les champs de configuration sans toucher à
// l'état Gio interne du bouton (clic, survol).
func configureMoveButton(b *components.IconButton, icon *widget.Icon, disabled bool) {
	b.Icon = icon
	b.Variant = components.ButtonDefault
	b.Size = components.SizeSM
	b.Disabled = disabled
}
