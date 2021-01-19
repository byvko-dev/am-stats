package render

import (
	"image/color"

	"github.com/fogleman/gg"
)

// ChallengeBlock -
type challengeBlock struct {
	// Text and color setup
	isPremium     bool
	isLocked      bool
	isColored     bool
	hasIcon       bool
	color         color.RGBA
	premiumColor  color.RGBA
	shortTxtColor color.RGBA
	longTxtColor  color.RGBA
	altTxtColor   color.RGBA
	textSize      float64
	textCoeff     float64
	// Icon and score
	status   string
	score    float64
	position int
	// Text
	shortTxt string
	longText string
	prizeTxt string
	// Dimensions
	width   int
	height  int
	context *gg.Context
}
