package theme

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/text"
)

type Theme struct {
	Colors     Colors
	Shape      Shape
	Spacing    Spacing
	Typography Typography

	Shaper *text.Shaper
}

type Colors struct {
	Background color.NRGBA
	Surface    color.NRGBA
	SurfaceAlt color.NRGBA

	Foreground color.NRGBA
	Muted      color.NRGBA
	Border     color.NRGBA

	Primary   color.NRGBA
	Secondary color.NRGBA
	Danger    color.NRGBA
	Success   color.NRGBA
	Warning   color.NRGBA
}

type Shape struct {
	Radius      float32
	BorderWidth float32
	Shadow      Shadow
}

type Shadow struct {
	Color   color.NRGBA
	OffsetX float32
	OffsetY float32
	Blur    float32
	Spread  float32
}

type Spacing struct {
	XS float32
	SM float32
	MD float32
	LG float32
	XL float32
}

type Typography struct {
	FontSize float32
	Font     font.Font
}
