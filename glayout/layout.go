package glayout

import (
	"image"
	"slices"

	"gioui.org/io/input"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

type Child struct {
	widget layout.Widget
	weight float32
}

func Rigid(w layout.Widget) Child                  { return Child{widget: w} }
func Flexed(weight float32, w layout.Widget) Child { return Child{widget: w, weight: weight} }

// Inset applique un padding uniforme basé sur un token de Spacing.
func Inset(space float32, w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(space)).Layout(gtx, w)
	}
}

// Row aligne des widgets horizontalement avec un espacement constant entre eux.
func Row(spacing float32, alignment layout.Alignment, children ...Child) layout.Widget {
	return flex(layout.Horizontal, spacing, alignment, children)
}

// Column aligne des widgets verticalement avec un espacement constant entre eux.
//
//	Column(th.Spacing.MD, layout.Start,
//		Rigid(header),
//		Flexed(1, list),
//		Rigid(footer),
//	)
func Column(spacing float32, alignment layout.Alignment, children ...Child) layout.Widget {
	return flex(layout.Vertical, spacing, alignment, children)
}

func flex(axis layout.Axis, spacing float32, alignment layout.Alignment, children []Child) layout.Widget {
	spacer := Spacer(spacing)
	return func(gtx layout.Context) layout.Dimensions {
		flexChildren := make([]layout.FlexChild, 0, 2*len(children))
		for i, c := range children {
			if c.weight > 0 {
				flexChildren = append(flexChildren, layout.Flexed(c.weight, c.widget))
			} else {
				flexChildren = append(flexChildren, layout.Rigid(c.widget))
			}
			if spacing > 0 && i < len(children)-1 {
				flexChildren = append(flexChildren, layout.Rigid(spacer))
			}
		}
		return layout.Flex{Axis: axis, Alignment: alignment}.Layout(gtx, flexChildren...)
	}
}

// Spacer crée un espace vide de taille fixe, dans l'axe du parent (Row ou Column).
func Spacer(space float32) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Width: unit.Dp(space), Height: unit.Dp(space)}.Layout(gtx)
	}
}

// Grid organise des widgets en grille à colonnes fixes, chaque colonne se
// partageant la largeur disponible à parts égales (layout.Flexed(1)). Les
// enfants sont répartis ligne par ligne, dans l'ordre — la dernière ligne
// peut être incomplète si len(children) n'est pas un multiple de columns
// (les colonnes manquantes restent simplement vides).
func Grid(spacing float32, columns int, children ...layout.Widget) layout.Widget {
	columns = max(columns, 1)
	rows := make([]Child, 0, (len(children)+columns-1)/columns)
	for chunk := range slices.Chunk(children, columns) {
		cells := make([]Child, len(chunk))
		for i, w := range chunk {
			cells[i] = Flexed(1, w)
		}
		rows = append(rows, Rigid(Row(spacing, layout.Start, cells...)))
	}
	return Column(spacing, layout.Start, rows...)
}

// Stack superpose des widgets les uns sur les autres, chacun occupant tout
// l'espace disponible (layout.Stack + layout.Expanded, sans le répéter à
// chaque appel). Le dernier enfant est dessiné en dernier, donc au-dessus —
// c'est le moyen normal de composer un overlay (Dialog...) par-dessus le
// reste de l'UI : PAS comme un enfant de plus dans un Column (qui ne lui
// donnerait qu'un espace résiduel), mais comme un widget frère du contenu
// principal, à l'intérieur d'un Stack commun.
func Stack(children ...layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		stackChildren := make([]layout.StackChild, len(children))
		for i, child := range children {
			stackChildren[i] = layout.Expanded(child)
		}
		return layout.Stack{}.Layout(gtx, stackChildren...)
	}
}

// StretchWidth force la largeur reportée de w à gtx.Constraints.Max.X.
// Flex additionne les tailles réelles de ses enfants sans les étirer.
func StretchWidth(w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		dims := w(gtx)
		return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, dims.Size.Y)}
	}
}

// Ghost layoute w pour obtenir ses dimensions, mais ne rejoue jamais ses
// opérations : rien n'est peint ni cliquable.
func Ghost(w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		gtx.Source = input.Source{}
		macro := op.Record(gtx.Ops)
		dims := w(gtx)
		macro.Stop()
		return dims
	}
}
