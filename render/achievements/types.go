package render

import (
	"image/color"

	"github.com/fogleman/gg"
)

// cardBlockData -
type cardBlockData struct {
	//
	TotalTextHeight float64
	TotalTextLines  int
	// Block setup
	Color         color.RGBA
	TextSize      float64
	TextCoeff     float64
	BlockTextSize float64
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
