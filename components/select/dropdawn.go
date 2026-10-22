package selectwidget

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/components"
	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type selectRow struct {
	call op.CallOp
	dims layout.Dimensions
}

// layoutDropdown dessine un panneau UNIQUE (un seul cadre, pas un par
// option) contenant une liste plate : chaque ligne n'est qu'un
// rectangle de surbrillance au survol + texte, pas un bouton individuel.
// Positionné sous le trigger si la place suffit, sinon au-dessus — la
// place disponible est approximée depuis gtx.Constraints.Max.Y (l'espace
// RESTANT dans la Column englobante), pas la taille réelle de la fenêtre.
func (s *Select) layoutDropdown(gtx layout.Context, th *theme.Theme, triggerSize image.Point) {
	pad := gtx.Dp(unit.Dp(th.Spacing.SM))
	rowPadY := gtx.Dp(unit.Dp(th.Spacing.XS))

	listW := triggerSize.X
	rowW := max(0, listW-2*pad)
	rows, rowsH := s.measureDropdownRows(gtx, th, rowW, rowPadY)
	listH := rowsH + 2*pad
	spaceBelow := gtx.Constraints.Max.Y - triggerSize.Y
	openAbove := listH > spaceBelow && triggerSize.Y < gtx.Constraints.Max.Y
	s.drawDropdownPanel(gtx, th, rows, listW, listH, pad, triggerSize, openAbove)
}

func (s *Select) measureDropdownRows(gtx layout.Context, th *theme.Theme, rowW, rowPadY int) ([]selectRow, int) {
	rows := make([]selectRow, len(s.Options))
	rowsH := 0
	for i, opt := range s.Options {
		clicked := s.optionBtns[i].Clicked(gtx)
		hovered := s.optionBtns[i].Hovered()

		rgtx := gtx
		rgtx.Constraints = layout.Constraints{Min: image.Point{}, Max: image.Pt(rowW, gtx.Constraints.Max.Y)}

		m := op.Record(gtx.Ops)

		rowDims := drawSelectRow(rgtx, th, opt.Label, hovered, rowPadY)
		call := m.Stop()

		rows[i] = selectRow{call: call, dims: rowDims}
		rowsH += rowDims.Size.Y

		if clicked {
			s.Selected = i
			s.open = false
			interaction.Invalidate(gtx)
		}
	}
	return rows, rowsH
}

func (s *Select) drawDropdownPanel(gtx layout.Context, th *theme.Theme, rows []selectRow, listW, listH, pad int, triggerSize image.Point, openAbove bool) {
	panelMacro := op.Record(gtx.Ops)
	frame := components.NeoFrame(gtx, th, listW-2*gtx.Dp(unit.Dp(th.Shape.BorderWidth)), listH, th.Colors.Surface, false, false)

	y := frame.BgRect.Min.Y + pad
	for i, r := range rows {
		rowOff := op.Offset(image.Pt(frame.BgRect.Min.X+pad, y)).Push(gtx.Ops)
		r.call.Add(gtx.Ops)
		rowOff.Pop()

		// Zone de clic sur toute la largeur de la ligne, indépendamment
		// du texte — enregistrée ici, sur les coordonnées réelles.
		rowRect := image.Rectangle{
			Min: image.Pt(frame.BgRect.Min.X+pad, y),
			Max: image.Pt(frame.BgRect.Max.X-pad, y+r.dims.Size.Y),
		}
		interaction.ClickArea(gtx, &s.optionBtns[i], rowRect)

		y += r.dims.Size.Y
	}
	panelCall := panelMacro.Stop()

	posY := triggerSize.Y
	if openAbove {
		posY = -listH
	}
	posStack := op.Offset(image.Pt(0, posY)).Push(gtx.Ops)
	op.Defer(gtx.Ops, panelCall)
	posStack.Pop()
}
