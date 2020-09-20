package render

import (
	"image"
	"image/color"
    "github.com/fogleman/gg"
)

// Struct for cards
type cardData struct {
    image   image.Image
    index   int
    context *gg.Context

}
// Card block
type cardBlock struct{
    // Text and color setup
    isColored       bool
    color           color.RGBA
    bigTextColor    color.RGBA
    smallTextColor  color.RGBA
    altTextColor    color.RGBA
    textSize        float64
    textCoeff       float64
    // Icon for WN8 and stats change
    hasBigIcon          bool
    bigArrowDirection   int
    bigIconColor        color.RGBA
    hasSmallIcon        bool
    smallArrowDirection int
    smallIconColor      color.RGBA
    // Text
    bigText         string
    smallText       string
    altText         string
    // Dimensions
    width           int
    height          int
    context         *gg.Context
}
// A slice of all generated cards
type allCards struct {
    cards   []cardData
    frame   *gg.Context
}
