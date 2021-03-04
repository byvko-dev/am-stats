package render

import (
	"image/color"

	"github.com/cufee/am-stats/render"
	"github.com/fogleman/gg"
)

// GameModeMap - Map of all known game modes
var GameModeMap map[int]string = make(map[int]string)

// SpecialGameModeTypes - Tournament and training battles
var SpecialGameModeTypes []int = []int{2, 4, 5}

func init() {
	GameModeMap[0] = "Unknown"
	GameModeMap[1] = "Regular Battle"
	GameModeMap[2] = "Training Room"
	GameModeMap[4] = "Tournament Battle"
	GameModeMap[5] = "Tournament Battle"
	GameModeMap[7] = "Rating Battle"
	GameModeMap[8] = "Mad Games"
	GameModeMap[22] = "Realistic Battle"
	GameModeMap[23] = "Uprising"
	GameModeMap[24] = "Gravity Force"
}

type hpBarBlockData struct {
	// HP %
	PercentHP float64
	HPColor   color.RGBA
	HPColorBG color.RGBA
	// General
	Margin  int
	Width   int
	Height  int
	Context *gg.Context
	// Color
	Color color.RGBA
}

// replayBlockData -
type replayBlockData struct {
	//
	TotalTextHeight int
	TotalTextLines  int
	// Block setup
	Color         color.RGBA
	BlockMargin   float64
	BlockTextSize float64
	// Text
	TextLines  []blockTextLine
	TextAlign  int
	TextSize   float64
	TextMargin float64
	TextColor  color.RGBA
	// Colored dot
	DotSize  int
	DotColor color.RGBA
	// Icon
	IconURL         string
	IconSize        int
	IconTextOverlay string
	// General
	Padding int
	Width   int
	Height  int
	Context *gg.Context
}

type blockTextLine struct {
	Text      string
	TextScale float64
	Color     color.RGBA
}

// Blueprint for replay block
func (c *replayBlockData) Defaults() {
	// Icon
	c.IconSize = 50
	// Dot
	c.DotSize = 50
	// Text
	c.TextSize = render.FontSize
	c.TextMargin = float64(render.FrameMargin)
	c.TextColor = render.BigTextColor
}

func resieFont(font float64, coeff int, div int) float64 {
	return float64(int(font) * coeff / div)
}

const battlesTypeSupremacy int = 1
