package components

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

// IconText affiche une icône suivie d'un texte, alignés et espacés selon
// le theme. Composition de Icon + Text, aucune logique propre.
func IconText(th *theme.Theme, icon *widget.Icon, str string, iconOpts []IconOption, textOpts []TextOption) layout.Widget {
	return glayout.Row(th.Spacing.SM, layout.Middle,
		glayout.Rigid(Icon(th, icon, iconOpts...)),
		glayout.Rigid(Text(th, str, textOpts...)),
	)
}
