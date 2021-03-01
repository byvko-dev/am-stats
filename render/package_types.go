package render

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

// CardData -
type CardData struct {
	Image       image.Image
	IndexX      int
	Index       int
	Context     *gg.Context
	LastXOffs   int
	FrameMargin int
	BlockWidth  float64
	Type        string
	Blocks      []Block
}

// CardTypeHeader - Header card type
const CardTypeHeader string = "headerCard"

// AllCards - A slice of all generated cards
type AllCards struct {
	Cards []CardData
	Frame *gg.Context
}

// Block -
type Block struct {
	// Additional data
	Extra     interface{}
	ExtraType string
	// Block setup
	Color color.RGBA
	// Text
	TextColor      color.RGBA
	SmallTextColor color.RGBA
	TextSize       float64
	TextMargin     float64
	TextAlign      int
	// Dimensions
	Padding int
	Width   int
	Height  int
	Context *gg.Context
}
