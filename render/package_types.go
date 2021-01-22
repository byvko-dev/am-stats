package render

import (
	"image"

	"github.com/fogleman/gg"
)

// CardData -
type CardData struct {
	Image      image.Image
	Index      int
	Context    *gg.Context
	LastXOffs  int
	BlockWidth float64
}

// AllCards - A slice of all generated cards
type AllCards struct {
	Cards []CardData
	Frame *gg.Context
}
