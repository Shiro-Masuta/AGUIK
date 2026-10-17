package components

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

// CodeEditor est un champ de texte multiligne éditable avec une gouttière de
// numéros de ligne, dans le style d'un éditeur de code. Comme Input et
// TextArea, il encapsule un widget.Editor Gio et doit être conservé dans une
// structure d'état persistante entre les frames — ne jamais le recréer dans
// la fonction de rendu.
//
// L'écriture, le déplacement du curseur, la sélection (surlignée) et le
// copier-coller sont ceux, natifs, de widget.Editor : CodeEditor ne fait que
// l'habiller et ajouter la gouttière.
//
// Limite connue : les numéros de ligne correspondent aux lignes logiques
// (séparées par '\n'), pas aux lignes visuelles. Si une ligne très longue
// retourne à la ligne automatiquement (l'éditeur ne fait pas défiler
// horizontalement), la numérotation peut se décaler légèrement en dessous
// de cette ligne. Pour du code raisonnablement large, ce n'est pas gênant en
// pratique ; on peut ajouter un défilement horizontal séparé plus tard si
// besoin.
//
// Exemple :
//
//	editor := components.CodeEditor{
//	    Placeholder: "// tapez votre code ici",
//	    Height:      360,
//	}
//	// dans le rendu :
//	return state.editor.Layout(gtx, th)
type CodeEditor struct {
	// Editor est l'état Gio du texte. Accessible directement :
	// editor.Editor.Text(), editor.Editor.SetText(...).
	Editor widget.Editor

	// list assure le défilement vertical de la gouttière et de l'éditeur
	// ensemble, comme un seul bloc.
	list widget.List

	Placeholder string

	// Width et Height sont en dp. Width à 0 prend toute la largeur
	// disponible. Height vaut 320 dp par défaut si non renseigné.
	Width  float32
	Height float32

	Disabled bool

	// HideLineNumbers masque la gouttière si true.
	HideLineNumbers bool

	// GutterMinDigits force une largeur minimale de gouttière, en nombre
	// de chiffres. 3 par défaut.
	GutterMinDigits int

	// Font permet d'utiliser une police à chasse fixe (recommandé pour du
	// code). Si vide, la police du thème (th.Typography.Font) est utilisée.
	Font font.Font

	// TextSize surcharge la taille de police en sp. Si 0, utilise
	// th.Typography.FontSize.
	TextSize float32
}

// Text retourne le contenu courant de l'éditeur.
func (c *CodeEditor) Text() string { return c.Editor.Text() }

// SetText remplace le contenu de l'éditeur.
func (c *CodeEditor) SetText(s string) { c.Editor.SetText(s) }

func (c *CodeEditor) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	height := c.Height
	if height <= 0 {
		height = 320
	}

	width := gtx.Constraints.Max.X
	if c.Width > 0 {
		width = gtx.Dp(unit.Dp(c.Width))
	}
	boxHeight := gtx.Dp(unit.Dp(height))
	size := gtx.Constraints.Constrain(image.Pt(width, boxHeight))

	gtx.Constraints.Min = size
	gtx.Constraints.Max = size

	bg := th.Colors.Surface
	if c.Disabled {
		bg = th.Colors.SurfaceAlt
	}

	return drawFrame(gtx, th, bg, th.Colors.Border, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(th.Spacing.SM)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return c.layoutScroll(gtx, th)
		})
	})
}

// layoutScroll place la gouttière et l'éditeur côte à côte, sans hauteur
// bornée (donc sans scroll interne propre à l'éditeur), puis fait défiler
// ce bloc entier via une unique liste verticale native. Gouttière et texte
// bougent donc toujours ensemble, car ils appartiennent à la même passe de
// layout.
func (c *CodeEditor) layoutScroll(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	c.list.Axis = layout.Vertical
	return c.list.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
		return c.layoutContent(gtx, th)
	})
}

func (c *CodeEditor) layoutContent(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	c.Editor.ReadOnly = c.Disabled
	c.Editor.SingleLine = false
	c.Editor.Submit = false

	textSize := c.TextSize
	if textSize <= 0 {
		textSize = th.Typography.FontSize
	}
	fnt := c.Font
	if fnt == (font.Font{}) {
		fnt = th.Typography.Font
	}

	// Thème "material" minimal, requis par material.Editor, qui délègue à
	// AGUIK la couleur, la bordure et le fond réels.
	mth := material.NewTheme()
	mth.Shaper = th.Shaper

	children := []glayout.Child{}

	if !c.HideLineNumbers {
		lineCount := strings.Count(c.Editor.Text(), "\n") + 1
		children = append(children,
			glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return c.layoutGutter(gtx, th, fnt, textSize, lineCount)
			}),
			glayout.Rigid(glayout.Spacer(th.Spacing.SM)),
		)
	}

	children = append(children, glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
		ed := material.Editor(mth, &c.Editor, c.Placeholder)
		ed.Font = fnt
		ed.TextSize = unit.Sp(textSize)
		ed.Color = th.Colors.Foreground
		ed.HintColor = th.Colors.Muted
		// Couleur de sélection bien visible : c'est ce qui "montre" le
		// code sélectionné.
		ed.SelectionColor = withAlpha(th.Colors.Primary, 110)
		return ed.Layout(gtx)
	}))

	return glayout.Row(0, layout.Start, children...)(gtx)
}

// layoutGutter affiche les numéros de ligne alignés à droite, comme un seul
// bloc de texte multiligne (pas un second éditeur).
func (c *CodeEditor) layoutGutter(gtx layout.Context, th *theme.Theme, fnt font.Font, textSize float32, lineCount int) layout.Dimensions {
	minDigits := c.GutterMinDigits
	if minDigits < 1 {
		minDigits = 3
	}
	if d := len(strconv.Itoa(lineCount)); d > minDigits {
		minDigits = d
	}

	// Largeur approximative : correcte pour une police à chasse fixe (ou
	// proche). Ajustable en changeant GutterMinDigits si besoin.
	gw := float32(minDigits)*textSize*0.62 + th.Spacing.SM*2
	gutterWidthPx := gtx.Dp(unit.Dp(gw))

	return drawFrame(gtx, th, th.Colors.SurfaceAlt, th.Colors.Border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gutterWidthPx
		gtx.Constraints.Max.X = gutterWidthPx
		return layout.UniformInset(unit.Dp(th.Spacing.XS)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return Text(th, gutterNumbers(lineCount),
				WithColor(th.Colors.Muted),
				WithSize(textSize),
				WithFont(fnt),
				WithAlignment(text.End),
			)(gtx)
		})
	})
}

func gutterNumbers(lineCount int) string {
	var b strings.Builder
	for i := 1; i <= lineCount; i++ {
		b.WriteString(strconv.Itoa(i))
		if i != lineCount {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// drawFrame dessine un cadre simple (fond + bordure arrondie) autour du
// contenu, en suivant le pattern "mesurer puis dessiner" du reste d'AGUIK :
// on enregistre les opérations du contenu, on connaît alors ses dimensions
// réelles, on dessine le fond à cette taille, puis on rejoue le contenu.
func drawFrame(gtx layout.Context, th *theme.Theme, bg, border color.NRGBA, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	rr := gtx.Dp(unit.Dp(th.Shape.Radius))
	rect := clip.RRect{
		Rect: image.Rectangle{Max: dims.Size},
		NE:   rr, NW: rr, SE: rr, SW: rr,
	}

	paint.FillShape(gtx.Ops, bg, rect.Op(gtx.Ops))

	bw := gtx.Dp(unit.Dp(th.Shape.BorderWidth))
	if bw > 0 {
		stroke := clip.Stroke{Path: rect.Path(gtx.Ops), Width: float32(bw)}
		paint.FillShape(gtx.Ops, border, stroke.Op())
	}

	call.Add(gtx.Ops)
	return dims
}

func withAlpha(c color.NRGBA, a uint8) color.NRGBA {
	c.A = a
	return c
}
