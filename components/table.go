package components

import (
	"image"
	"image/color"
	"sort"

	"gioui.org/font"
	"gioui.org/io/input"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type SortDirection int

const (
	SortNone SortDirection = iota
	SortAsc
	SortDesc
)

type TableColumn[T any] struct {
	Label  string
	Weight float32
	Cell   func(item T) string
	Less   func(a, b T) bool
}

type tableReorderButtons struct {
	up   IconButton
	down IconButton
}

func tableMustIcon(data []byte) *widget.Icon {
	ic, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return ic
}

var (
	moveUpIcon   = tableMustIcon(icons.NavigationArrowUpward)
	moveDownIcon = tableMustIcon(icons.NavigationArrowDownward)
)

type Table[T any] struct {
	Columns []TableColumn[T]
	Rows    []T

	SortColumn    int
	SortDirection SortDirection

	Reorderable bool

	RowActions func(item T, index int) layout.Widget

	headerBtns  []widget.Clickable
	reorderBtns []tableReorderButtons
	list        widget.List
}

func blendColor(base, tint color.NRGBA, ratio float32) color.NRGBA {
	lerp := func(a, b uint8) uint8 {
		return uint8(float32(a) + (float32(b)-float32(a))*ratio)
	}
	return color.NRGBA{
		R: lerp(base.R, tint.R),
		G: lerp(base.G, tint.G),
		B: lerp(base.B, tint.B),
		A: 255,
	}
}

func fillRect(gtx layout.Context, col color.NRGBA, size image.Point) {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// stretchWidth force la dimension REPORTÉE d'un widget à occuper toute la
// largeur disponible (Max.X), même si son contenu réel est plus étroit.
// Row/Flex additionne les tailles réelles de ses enfants, il ne les force
// pas à remplir l'espace qui leur a été alloué — d'où ce wrapper.
func stretchWidth(w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		dims := w(gtx)
		return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, dims.Size.Y)}
	}
}

// ghost layoute w normalement — mêmes dimensions retournées, même code
// exécuté — mais ne peint rien : les opérations enregistrées ne sont
// jamais ajoutées à gtx.Ops (pas de macro.Add()). Aucun tag d'interaction
// n'est donc jamais enregistré, ce qui rend le rendu invisible ET non
// cliquable.
//
// Sert à réserver, dans le header, EXACTEMENT la même largeur que le
// widget réel affiché dans les lignes (actions, réorg.) sans dupliquer sa
// logique de mesure : c'est littéralement le même widget, juste pas
// dessiné. C'est le point qui manquait dans la version précédente — un
// calcul de largeur séparé pouvait diverger du rendu réel.
func ghost(w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		mgtx := gtx
		mgtx.Source = input.Source{}
		macro := op.Record(mgtx.Ops)
		dims := w(mgtx)
		macro.Stop()
		return dims
	}
}

// rowHorizontalMargins retourne les marges gauche/droite (en px) que
// framedRow applique réellement à son contenu. DOIT rester synchronisé
// avec NeoFrame — spécifiquement le cas hovered=false/pressed=true, seul
// utilisé par framedRow : ce cas décale le contenu de baseX vers la
// droite (contentOffX = baseX dans NeoFrame), donc baseX fait partie de
// la marge GAUCHE réelle, pas de la marge droite.
func rowHorizontalMargins(gtx layout.Context, th *theme.Theme) (left, right int) {
	dp := func(v float32) int { return gtx.Dp(unit.Dp(v)) }
	lift := dp(hoverLiftDp)
	baseX := dp(th.Shape.Shadow.OffsetX)
	grow := dp(hoverGrowDp)
	borderW := dp(th.Shape.BorderWidth)
	left = lift + baseX + borderW
	right = borderW + grow
	return left, right
}

// neoFrameOverheadX reproduit le calcul de marge de NeoFrame (padLeft +
// padRight + 2*bordure) pour demander à l'avance un innerW tel que le
// canvas final fasse exactement gtx.Constraints.Max.X de large. Garder
// synchronisé avec NeoFrame si hoverLiftDp/hoverGrowDp changent.
func neoFrameOverheadX(gtx layout.Context, th *theme.Theme) int {
	left, right := rowHorizontalMargins(gtx, th)
	return left + right
}

// framedRow dessine content (mesuré à innerW de large) dans un NeoFrame qui
// occupe toute la largeur disponible. pressed=true donne la variante
// "plate" de NeoFrame (bordure + fond, sans ombre) : une ombre par ligne
// sur toute une liste serait trop chargée visuellement.
func framedRow(gtx layout.Context, th *theme.Theme, bg color.NRGBA, content layout.Widget) layout.Dimensions {
	innerW := gtx.Constraints.Max.X - neoFrameOverheadX(gtx, th)
	if innerW < 0 {
		innerW = 0
	}

	macro := op.Record(gtx.Ops)
	gtx2 := gtx
	gtx2.Constraints = layout.Constraints{
		Min: image.Pt(innerW, 0),
		Max: image.Pt(innerW, gtx.Constraints.Max.Y),
	}
	dims := content(gtx2)
	call := macro.Stop()

	frame := NeoFrame(gtx, th, innerW, dims.Size.Y, bg, false, true)

	off := op.Offset(frame.BgRect.Min).Push(gtx.Ops)
	clipArea := clip.Rect{Max: frame.BgRect.Size()}.Push(gtx.Ops) // rien ne peint hors du cadre
	call.Add(gtx.Ops)
	clipArea.Pop()
	off.Pop()

	return layout.Dimensions{Size: frame.CanvasDim}
}

// reorderWidget construit les deux boutons haut/bas pour la ligne i.
// Extrait à part pour que header (via ghost) et lignes utilisent
// EXACTEMENT le même widget, donc la même largeur en pixels.
func (t *Table[T]) reorderWidget(th *theme.Theme, i int) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		buttons := &t.reorderBtns[i]
		up, down := &buttons.up, &buttons.down
		up.Icon = moveUpIcon
		up.Variant = ButtonDefault
		up.Size = SizeSM
		up.Disabled = i == 0

		down.Icon = moveDownIcon
		down.Variant = ButtonDefault
		down.Size = SizeSM
		down.Disabled = i == len(t.Rows)-1

		return glayout.Row(th.Spacing.SM, layout.Middle,
			glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
				dims := up.Layout(gtx, th)
				if up.Clicked && i > 0 {
					t.Rows[i-1], t.Rows[i] = t.Rows[i], t.Rows[i-1]
				}
				return dims
			}),
			glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
				dims := down.Layout(gtx, th)
				if down.Clicked && i < len(t.Rows)-1 {
					t.Rows[i], t.Rows[i+1] = t.Rows[i+1], t.Rows[i]
				}
				return dims
			}),
		)(gtx)
	}
}

// trailingWidgets retourne, dans l'ordre où elles apparaissent en fin de
// ligne, les colonnes optionnelles réellement actives (réorg. puis
// actions). Header et lignes appellent tous les deux cette méthode :
// impossible qu'ils divergent sur le NOMBRE ou l'ORDRE de ces colonnes.
func (t *Table[T]) trailingWidgets(th *theme.Theme, reorderActive bool, row T, i int) []layout.Widget {
	var ws []layout.Widget
	if reorderActive {
		ws = append(ws, t.reorderWidget(th, i))
	}
	if t.RowActions != nil {
		ws = append(ws, t.RowActions(row, i))
	}
	return ws
}

func (t *Table[T]) syncHeaderButtons() {
	if len(t.headerBtns) == len(t.Columns) {
		return
	}
	btns := make([]widget.Clickable, len(t.Columns))
	copy(btns, t.headerBtns)
	t.headerBtns = btns
}

func (t *Table[T]) syncReorderButtons() {
	if len(t.reorderBtns) == len(t.Rows) {
		return
	}
	buttons := make([]tableReorderButtons, len(t.Rows))
	copy(buttons, t.reorderBtns)
	t.reorderBtns = buttons
}

func (t *Table[T]) handleSortClicks(gtx layout.Context) {
	for i, col := range t.Columns {
		if col.Less == nil || !t.headerBtns[i].Clicked(gtx) {
			continue
		}
		if t.SortColumn != i {
			t.SortColumn = i
			t.SortDirection = SortAsc
			continue
		}
		switch t.SortDirection {
		case SortAsc:
			t.SortDirection = SortDesc
		case SortDesc:
			t.SortColumn = -1
			t.SortDirection = SortNone
		default:
			t.SortDirection = SortAsc
		}
	}
}

func (t *Table[T]) applySort() {
	if t.SortDirection == SortNone || t.SortColumn < 0 || t.SortColumn >= len(t.Columns) {
		return
	}
	less := t.Columns[t.SortColumn].Less
	sort.SliceStable(t.Rows, func(i, j int) bool {
		if t.SortDirection == SortDesc {
			return less(t.Rows[j], t.Rows[i])
		}
		return less(t.Rows[i], t.Rows[j])
	})
}

func (t *Table[T]) headerLabel(th *theme.Theme, i int, col TableColumn[T]) layout.Widget {
	label := col.Label
	switch {
	case t.SortColumn == i && t.SortDirection == SortAsc:
		label += " ▲"
	case t.SortColumn == i && t.SortDirection == SortDesc:
		label += " ▼"
	}
	content := Text(th, label, WithWeight(font.SemiBold), WithColor(th.Colors.Muted))

	if col.Less == nil {
		return content
	}
	btn := &t.headerBtns[i]
	return func(gtx layout.Context) layout.Dimensions {
		dims := content(gtx)
		interaction.ClickArea(gtx, btn, image.Rectangle{Max: dims.Size})
		return dims
	}
}

func columnWeight(w float32) float32 {
	if w == 0 {
		return 1
	}
	return w
}

// headerTrailingRow retourne la ligne (au sens "données") sur laquelle le
// header se base pour fantômer les colonnes de fin — la première ligne
// réelle si elle existe, sinon la valeur zéro de T (aucune ligne : les
// colonnes de fin n'ont de toute façon aucune largeur pertinente à
// respecter, la valeur zéro suffit à faire tourner le code sans paniquer).
func (t *Table[T]) headerTrailingRow() T {
	if len(t.Rows) > 0 {
		return t.Rows[0]
	}
	var zero T
	return zero
}

func (t *Table[T]) buildHeaderRow(th *theme.Theme, reorderActive bool) layout.Widget {
	children := make([]glayout.Child, 0, len(t.Columns)+2)
	for i, col := range t.Columns {
		label := stretchWidth(t.headerLabel(th, i, col))
		children = append(children, glayout.Flexed(columnWeight(col.Weight), label))
	}

	trailing := t.trailingWidgets(th, reorderActive, t.headerTrailingRow(), 0)
	for _, w := range trailing {
		children = append(children, glayout.Rigid(ghost(w)))
	}

	return glayout.Row(th.Spacing.MD, layout.Middle, children...)
}

func (t *Table[T]) buildHeader(th *theme.Theme, reorderActive bool) layout.Widget {
	// Applique EXACTEMENT les mêmes marges gauche/droite que framedRow,
	// pour que les colonnes du header tombent sur les mêmes limites que
	// les cellules des lignes en dessous.
	aligned := func(gtx layout.Context) layout.Dimensions {
		headerRow := t.buildHeaderRow(th, reorderActive)

		left, right := rowHorizontalMargins(gtx, th)
		innerW := gtx.Constraints.Max.X - left - right
		if innerW < 0 {
			innerW = 0
		}

		gtx2 := gtx
		gtx2.Constraints = layout.Constraints{
			Min: image.Pt(innerW, 0),
			Max: image.Pt(innerW, gtx.Constraints.Max.Y),
		}

		macro := op.Record(gtx.Ops)
		dims := headerRow(gtx2)
		call := macro.Stop()

		off := op.Offset(image.Pt(left, 0)).Push(gtx.Ops)
		call.Add(gtx.Ops)
		off.Pop()

		return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, dims.Size.Y)}
	}

	// Padding vertical seulement — n'affecte pas l'alignement horizontal.
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(th.Spacing.SM), Bottom: unit.Dp(th.Spacing.SM)}.Layout(gtx, aligned)
	}
}

func divider(th *theme.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		size := image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(th.Shape.BorderWidth)))
		fillRect(gtx, th.Colors.Border, size)
		return layout.Dimensions{Size: size}
	}
}

func (t *Table[T]) rowBackground(th *theme.Theme, i int) color.NRGBA {
	if i%2 == 1 {
		return blendColor(th.Colors.Background, th.Colors.Muted, 0.15)
	}
	return th.Colors.Background
}

func (t *Table[T]) buildRow(th *theme.Theme, reorderActive bool, i int) layout.Widget {
	row := t.Rows[i]

	cells := func(gtx layout.Context) layout.Dimensions {
		children := make([]glayout.Child, 0, len(t.Columns)+2)
		for _, col := range t.Columns {
			cellText := stretchWidth(Text(th, col.Cell(row), WithMaxLines(1), WithEllipsis()))
			children = append(children, glayout.Flexed(columnWeight(col.Weight), cellText))
		}

		for _, w := range t.trailingWidgets(th, reorderActive, row, i) {
			children = append(children, glayout.Rigid(w))
		}

		return glayout.Row(th.Spacing.MD, layout.Middle, children...)(gtx)
	}

	bg := t.rowBackground(th, i)
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(th.Spacing.SM)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return framedRow(gtx, th, bg, cells)
		})
	}
}

func (t *Table[T]) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	t.syncHeaderButtons()
	t.syncReorderButtons()
	t.handleSortClicks(gtx)
	t.applySort()

	reorderActive := t.Reorderable && t.SortDirection == SortNone

	header := t.buildHeader(th, reorderActive)

	t.list.Axis = layout.Vertical
	body := func(gtx layout.Context) layout.Dimensions {
		return t.list.Layout(gtx, len(t.Rows), func(gtx layout.Context, i int) layout.Dimensions {
			return t.buildRow(th, reorderActive, i)(gtx)
		})
	}

	content := glayout.Column(0, layout.Start,
		glayout.Rigid(header),
		glayout.Rigid(divider(th)),
		glayout.Flexed(1, body),
	)

	macro := op.Record(gtx.Ops)
	dims := content(gtx)
	call := macro.Stop()

	clipSize := image.Pt(gtx.Constraints.Max.X, dims.Size.Y)
	clipArea := clip.Rect{Max: clipSize}.Push(gtx.Ops)
	call.Add(gtx.Ops)
	clipArea.Pop()

	return layout.Dimensions{Size: clipSize}
}
