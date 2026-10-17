package interaction

import (
	"gioui.org/io/key"
	"gioui.org/layout"
)

func ClearFocus(gtx layout.Context) {
	gtx.Execute(key.FocusCmd{Tag: nil})
}
