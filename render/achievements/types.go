package render

import (
	"image/color"

	"github.com/cufee/am-stats/render"
	"github.com/fogleman/gg"
)

// cardBlockData -
type cardBlockData struct {
	//
	TotalTextHeight int
	TotalTextLines  int
	// Block setup
	Color         color.RGBA
	TextCoeff     int
	BlockMargin   float64
	BlockTextSize float64
	// Space for player name
	TextAlign         int
	NameMargin        float64
	ClanTagMargin     float64
	SpecialBlockWidth float64
	// Text
	BigTextColor   color.RGBA
	SmallTextColor color.RGBA
	BigText        string
	SmallText      string
	TextSize       float64
	TextMargin     float64
	AltTextSize    float64
	BigTextSize    float64
	SmallTextSize  float64
	// Medal Icon and alt text
	IconSize     int
	IconURL      string
	AltText      string
	AltTextColor color.RGBA
	// Dimensions
	Width   int
	Height  int
	Context *gg.Context
}

func (c *cardBlockData) DefaultSlim() {
	// Blueprint for small blocks
	c.IconSize = 50
	c.TextSize = render.FontSize
	c.BigTextSize = render.FontSize
	c.TextMargin = float64(render.FrameMargin)
	c.AltTextSize = resieFont(c.BigTextSize, c.TextCoeff, 10)
	c.SmallTextSize = resieFont(c.BigTextSize, c.TextCoeff, 10)
	c.BlockTextSize = resieFont(float64(render.FontSize), 125, 100)

	c.BigTextColor = render.BigTextColor
	c.AltTextColor = render.AltTextColor
	c.SmallTextColor = render.SmallTextColor
}

func (c *cardBlockData) ChangeCoeff(coeff int, div int) {
	c.AltTextSize = resieFont(c.BigTextSize, coeff, div)
	c.SmallTextSize = resieFont(c.BigTextSize, coeff, div)
}

func resieFont(font float64, coeff int, div int) float64 {
	return float64(int(font) * coeff / div)
}
