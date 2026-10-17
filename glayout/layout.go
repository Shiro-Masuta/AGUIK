package glayout

import (
	"gioui.org/layout"
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
	return func(gtx layout.Context) layout.Dimensions {
		flexChildren := make([]layout.FlexChild, 0, len(children)*2)
		for i, c := range children {
			c := c
			if c.weight > 0 {
				flexChildren = append(flexChildren, layout.Flexed(c.weight, c.widget))
			} else {
				flexChildren = append(flexChildren, layout.Rigid(c.widget))
			}
			if i < len(children)-1 && spacing > 0 {
				sp := spacing
				flexChildren = append(flexChildren, layout.Rigid(Spacer(sp)))
			}
		}
		return layout.Flex{Axis: layout.Horizontal, Alignment: alignment}.Layout(gtx, flexChildren...)
	}
}

// Column aligne des widgets verticalement avec un espacement constant entre eux.
/*
fonctionnement :

	Column(th.Spacing.MD, layout.Start,
		Rigid(header),
		Flexed(1, list),
		Rigid(footer),
	)

*/
func Column(spacing float32, alignment layout.Alignment, children ...Child) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		flexChildren := make([]layout.FlexChild, 0, len(children)*2)
		for i, c := range children {
			c := c
			if c.weight > 0 {
				flexChildren = append(flexChildren, layout.Flexed(c.weight, c.widget))
			} else {
				flexChildren = append(flexChildren, layout.Rigid(c.widget))
			}
			if i < len(children)-1 && spacing > 0 {
				sp := spacing
				flexChildren = append(flexChildren, layout.Rigid(Spacer(sp)))
			}
		}
		return layout.Flex{Axis: layout.Vertical, Alignment: alignment}.Layout(gtx, flexChildren...)
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
	if columns < 1 {
		columns = 1
	}
	return func(gtx layout.Context) layout.Dimensions {
		rowCount := (len(children) + columns - 1) / columns
		rows := make([]layout.Widget, 0, rowCount)
		for i := 0; i < len(children); i += columns {
			end := i + columns
			if end > len(children) {
				end = len(children)
			}
			rowChildren := children[i:end]
			rows = append(rows, func(gtx layout.Context) layout.Dimensions {
				flexChildren := make([]layout.FlexChild, 0, len(rowChildren)*2)
				for j, child := range rowChildren {
					child := child
					flexChildren = append(flexChildren, layout.Flexed(1, child))
					if j < len(rowChildren)-1 && spacing > 0 {
						sp := spacing
						flexChildren = append(flexChildren, layout.Rigid(Spacer(sp)))
					}
				}
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, flexChildren...)
			})
		}
		return Column(spacing, layout.Start, wrapRigid(rows)...)(gtx)
	}
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

func wrapRigid(ws []layout.Widget) []Child {
	cs := make([]Child, len(ws))
	for i, w := range ws {
		cs[i] = Rigid(w)
	}
	return cs
}
