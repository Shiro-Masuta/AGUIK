package codeeditor

import (
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/components"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

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
			return components.Text(th, gutterNumbers(lineCount),
				components.WithColor(th.Colors.Muted),
				components.WithSize(textSize),
				components.WithFont(fnt),
				components.WithAlignment(text.End),
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
