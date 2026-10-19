package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("AGUIK - showcase2"), app.Size(unit.Dp(1000), unit.Dp(700)))
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := theme.Ao()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	st := newState()
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			st.layout(gtx, &th)
			e.Frame(gtx.Ops)
		}
	}
}
