package render

import (
	"fmt"
	"sort"

	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"

	"github.com/cufee/am-stats/config"
)

// General settings
var (
	FontPath         = (config.AssetsPath + "/font.ttf")
	FontSizeHeader   = 36.0
	FontSize         = 24.0
	TextMargin       = FontSize / 2
	FrameWidth       = 900
	FrameMargin      = 50
	BaseCardWidth    = FrameWidth - (2 * FrameMargin)
	BaseCardHeigh    = 150
	BaseCardColor    = color.RGBA{30, 30, 30, 204}
	DecorLinesColor  = color.RGBA{80, 80, 80, 255}
	BigTextColor     = color.RGBA{255, 255, 255, 255}
	SmallTextColor   = color.RGBA{204, 204, 204, 255}
	AltTextColor     = color.RGBA{100, 100, 100, 255}
	ProtagonistColor = color.RGBA{255, 165, 0, 255}

	PremiumColor  = color.RGBA{255, 223, 0, 255}
	VerifiedColor = color.RGBA{72, 167, 250, 255}

	// DEBUG

	DebugColorRed   = color.RGBA{255, 0, 0, 255}
	DebugColorPink  = color.RGBA{255, 192, 203, 255}
	DebugColorGreen = color.RGBA{20, 160, 20, 255}
	BebugIconURL    = "https://images.vexels.com/media/users/3/141120/isolated/preview/a5ff757d7423e6c757795e7b60183180-rocket-round-icon-by-vexels.png"
)

// AddAllCardsToFrame - Render all cards to frame
func AddAllCardsToFrame(finalCards AllCards, header string, bgImage image.Image) (*gg.Context, error) {
	if len(finalCards.Cards) == 0 {
		return nil, fmt.Errorf("no cards to be rendered")
	}

	// Frame height
	var totalCardsHeight int
	var maxWidth int
	for _, card := range finalCards.Cards {
		totalCardsHeight += card.Context.Height()
		if card.Context.Width() > maxWidth {
			maxWidth = card.Context.Width()
		}
	}
	totalCardsHeight += ((len(finalCards.Cards)-1)*(FrameMargin/4) + (FrameMargin * 2))
	// Get frame CTX
	ctx, err := prepBgContext(totalCardsHeight, maxWidth, bgImage)
	if err != nil {
		return finalCards.Frame, err
	}
	finalCards.Frame = ctx
	// Sort cards
	sort.Slice(finalCards.Cards, func(i, j int) bool {
		return finalCards.Cards[i].Index < finalCards.Cards[j].Index
	})
	// Render cards
	var lastCardPos int = FrameMargin * 3 / 4
	for i := 0; i < len(finalCards.Cards); i++ {
		card := finalCards.Cards[i]
		cardMarginH := lastCardPos + (FrameMargin / 4)
		finalCards.Frame.DrawImage(card.Image, FrameMargin, cardMarginH)
		lastCardPos = cardMarginH + card.Context.Height()
	}

	// Draw timestamp
	if err := finalCards.Frame.LoadFontFace(FontPath, FontSize*0.75); err != nil {
		return finalCards.Frame, err
	}
	finalCards.Frame.SetColor(color.RGBA{100, 100, 100, 100})
	headerW, headerH := finalCards.Frame.MeasureString(header)
	headerX := (float64(finalCards.Frame.Width()) - headerW) / 2
	headerY := (float64(FrameMargin)-headerH)/2 + headerH
	finalCards.Frame.DrawString(header, headerX, headerY)

	return finalCards.Frame, nil
}

// Prepare a frame background context
func prepBgContext(totalHeight int, width int, bgImage image.Image) (frameCtx *gg.Context, err error) {
	frameWidth := width + (2 * FrameMargin)
	frameCtx = gg.NewContext(frameWidth, totalHeight)
	bgImage = imaging.Fill(bgImage, frameCtx.Width(), frameCtx.Height(), imaging.Center, imaging.NearestNeighbor)
	bgImage = imaging.Blur(bgImage, 10.0)
	frameCtx.DrawImage(bgImage, 0, 0)
	return frameCtx, nil
}

// PrepNewCard - Prepare a new cardData struct
func PrepNewCard(card *CardData, index int, heightMod float64, width int) {
	if len(card.Blocks) > 0 {
		for _, b := range card.Blocks {
			width += b.Width + b.Padding
		}
		width += card.FrameMargin * 2
	}
	if width == 0 {
		width = BaseCardWidth
	}
	cardHeight := int(float64(BaseCardHeigh) * heightMod)
	cardWidth := width
	cardCtx := gg.NewContext(width, cardHeight)
	cardCtx.SetColor(BaseCardColor)
	cardCtx.DrawRoundedRectangle(0, 0, float64(cardWidth), float64(cardHeight), FontSize)
	cardCtx.Fill()
	card.Context = cardCtx
	card.Index = index
}
