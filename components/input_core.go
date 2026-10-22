package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/widget/material"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

// inputDims regroupe la taille de contenu voulue (hors padding/bordure) —
// calculée différemment selon Input (une ligne) ou TextArea (N lignes).
type inputDims struct {
	width, height int
}

func constrainInputDims(dims inputDims, pad int, max, overhead image.Point) inputDims {
	maxWidth := max.X - overhead.X - 2*pad
	maxHeight := max.Y - overhead.Y - 2*pad
	if maxWidth < 0 {
		maxWidth = 0
	}
	if maxHeight < 0 {
		maxHeight = 0
	}
	if dims.width > maxWidth {
		dims.width = maxWidth
	}
	if dims.height > maxHeight {
		dims.height = maxHeight
	}
	if dims.width < 0 {
		dims.width = 0
	}
	if dims.height < 0 {
		dims.height = 0
	}
	return dims
}

// layoutInputFrame dessine le cadre néobrutaliste autour d'un
// material.EditorStyle déjà configuré, et le positionne dedans. Partagé
// par Input et TextArea — seule la manière dont dims est calculée diffère.
func layoutInputFrame(gtx layout.Context, th *theme.Theme, ed material.EditorStyle, dims inputDims, pad int, bg color.NRGBA, focused bool) layout.Dimensions {
	dims = constrainInputDims(dims, pad, gtx.Constraints.Max, NeoFrameOverhead(gtx, th))
	innerW := dims.width + 2*pad
	innerH := dims.height + 2*pad

	clipArea := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	frame := NeoFrame(gtx, th, innerW, innerH, bg, false, focused)

	off := frame.BgRect.Min.Add(image.Pt(pad, pad))
	stack := op.Offset(off).Push(gtx.Ops)
	tgtx := gtx
	tgtx.Constraints = layout.Constraints{
		Min: image.Pt(dims.width, dims.height),
		Max: image.Pt(dims.width, dims.height),
	}
	ed.Layout(tgtx)
	stack.Pop()
	clipArea.Pop()

	return layout.Dimensions{Size: gtx.Constraints.Constrain(frame.CanvasDim)}
}
