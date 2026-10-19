package components

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/interaction"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type SelectOption struct {
	Label string
}

type Select struct {
	Options     []SelectOption
	Selected    int
	Placeholder string
	Width       float32
	Size        ButtonSize
	Disabled    bool

	open       bool
	trigger    widget.Clickable
	optionBtns []widget.Clickable
}

func (s *Select) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if s.Selected < -1 || s.Selected >= len(s.Options) {
		s.Selected = -1
	}
	s.optionBtns = glayout.Resize(s.optionBtns, len(s.Options))

	if !s.Disabled && s.trigger.Clicked(gtx) {
		s.open = !s.open
		interaction.Invalidate(gtx)
	}

	label := s.Placeholder
	if s.Selected >= 0 {
		label = s.Options[s.Selected].Label
	}

	bg, fg := th.Colors.Surface, th.Colors.Foreground
	if s.Disabled {
		bg, fg = th.Colors.Muted, th.Colors.Foreground
	}
	hovered := !s.Disabled && s.trigger.Hovered()
	// Le style "actif" reflète simplement l'ouverture du dropdown — pas
	// de focus clavier à suivre ici : contrairement à widget.Editor
	// (Input), widget.Clickable ne gère que les clics pointeur, jamais
	// le focus. gtx.Focused ne se déclenche jamais pour lui, ce qui
	// causait l'ouverture/fermeture instantanée qu'on a corrigée.
	active := !s.Disabled && s.open

	pad := gtx.Dp(unit.Dp(paddingFor(th, s.Size)))
	widthDp := s.Width
	if widthDp == 0 {
		widthDp = defaultInputWidthDp
	}
	fieldW := gtx.Dp(unit.Dp(widthDp))
	maxFieldW := gtx.Constraints.Max.X - neoFrameOverheadX(gtx, th)
	if maxFieldW < 0 {
		maxFieldW = 0
	}
	if fieldW > maxFieldW {
		fieldW = maxFieldW
	}

	content := Text(th, label, WithColor(fg))

	mgtx := gtx
	mgtx.Constraints.Min = image.Point{}
	mgtx.Constraints.Max.X = fieldW - 2*pad
	if mgtx.Constraints.Max.X < 0 {
		mgtx.Constraints.Max.X = 0
	}
	macro := op.Record(mgtx.Ops)
	contentDims := content(mgtx)
	macro.Stop()

	innerW := fieldW
	innerH := contentDims.Size.Y + 2*pad
	frame := NeoFrame(gtx, th, innerW, innerH, bg, hovered, active)

	off := frame.BgRect.Min.Add(image.Pt(pad, pad))
	stack := op.Offset(off).Push(gtx.Ops)
	content(gtx)
	stack.Pop()

	if !s.Disabled {
		interaction.ClickArea(gtx, &s.trigger, frame.RestRect)
	}

	if s.open && len(s.Options) > 0 {
		s.layoutDropdown(gtx, th, frame.CanvasDim)
	}

	return layout.Dimensions{Size: gtx.Constraints.Constrain(frame.CanvasDim)}
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

	type recordedRow struct {
		call op.CallOp
		dims layout.Dimensions
	}

	rows := make([]recordedRow, len(s.Options))
	rowsH := 0
	listW := triggerSize.X

	for i, opt := range s.Options {
		clicked := s.optionBtns[i].Clicked(gtx)
		hovered := s.optionBtns[i].Hovered()

		rowW := listW - 2*pad
		rgtx := gtx
		rgtx.Constraints = layout.Constraints{Min: image.Point{}, Max: image.Pt(rowW, gtx.Constraints.Max.Y)}

		m := op.Record(gtx.Ops)

		rowDims := drawSelectRow(rgtx, th, opt.Label, hovered, rowPadY)
		call := m.Stop()

		rows[i] = recordedRow{call: call, dims: rowDims}
		rowsH += rowDims.Size.Y

		if clicked {
			s.Selected = i
			s.open = false
			interaction.Invalidate(gtx)
		}
	}

	listH := rowsH + 2*pad

	// Approximation : espace vertical restant sous le trigger, déduit de
	// gtx.Constraints.Max.Y (voir la limite documentée plus haut).
	spaceBelow := gtx.Constraints.Max.Y - triggerSize.Y
	openAbove := listH > spaceBelow && triggerSize.Y < gtx.Constraints.Max.Y

	panelMacro := op.Record(gtx.Ops)
	frame := NeoFrame(gtx, th, listW-2*gtx.Dp(unit.Dp(th.Shape.BorderWidth)), listH, th.Colors.Surface, false, false)

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

// drawSelectRow peint une ligne de liste plate : fond de surbrillance
// optionnel (hover) + texte. Le texte est mesuré AVANT de peindre le fond
// (sinon le rectangle de surbrillance serait peint à hauteur 0, puisque
// sa taille dépendrait d'un texte pas encore layouté).
func drawSelectRow(gtx layout.Context, th *theme.Theme, label string, hovered bool, padY int) layout.Dimensions {
	content := Text(th, label, WithColor(th.Colors.Foreground))

	mgtx := gtx
	mgtx.Constraints.Min = image.Point{}
	macro := op.Record(mgtx.Ops)
	textDims := content(mgtx)
	macro.Stop() // dry-run : jeté, on garde juste la taille

	rowH := textDims.Size.Y + 2*padY

	if hovered {
		paint.FillShape(gtx.Ops, th.Colors.SurfaceAlt, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, rowH)}.Op())
	}

	off := op.Offset(image.Pt(0, padY)).Push(gtx.Ops)
	content(gtx)
	off.Pop()

	return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, rowH)}
}
