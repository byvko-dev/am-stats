package render

import (
	"image/color"

	"github.com/fogleman/gg"
)

// statsCardBlock -
type statsBlock struct {
	// Text and color setup
	IsColored      bool
	Color          color.RGBA
	BigTextColor   color.RGBA
	SmallTextColor color.RGBA
	AltTextColor   color.RGBA
	TextSize       float64
	TextCoeff      float64
	// Icon for WN8 and stats change
	HasBigIcon          bool
	BigArrowDirection   int
	BigIconColor        color.RGBA
	HasSmallIcon        bool
	SmallArrowDirection int
	SmallIconColor      color.RGBA
	// Text
	BigText   string
	SmallText string
	AltText   string
	// Dimensions
	Width   int
	Height  int
	Context *gg.Context
}
