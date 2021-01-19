package render

import (
	"image"

	"github.com/fogleman/gg"
)

// CardData -
type CardData struct {
	Image   image.Image
	Index   int
	Context *gg.Context
}

// AllCards - A slice of all generated cards
type AllCards struct {
	Cards []CardData
	Frame *gg.Context
}
