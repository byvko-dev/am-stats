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
	TextSize      int
	TextCoeff     int
	BlockTextSize int
	// Space for player name
	NameMarginCoef float64
	// Text
	BigTextColor   color.RGBA
	SmallTextColor color.RGBA
	BigText        string
	SmallText      string
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
	c.TextCoeff = 6
	c.NameMarginCoef = 0.5
	c.BlockTextSize = int(render.FontSize) * 4 / 3
	c.TextSize = int(render.FontSize)
	c.BigTextColor = render.BigTextColor
	c.AltTextColor = render.AltTextColor
	c.SmallTextColor = render.SmallTextColor
}
