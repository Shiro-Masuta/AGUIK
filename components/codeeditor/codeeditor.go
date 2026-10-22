package codeeditor

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

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
//	editor := codeeditor.CodeEditor{
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
