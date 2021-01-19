package render

import (
	"fmt"
	"sort"
	"time"

	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"

	"github.com/cufee/am-stats/config"
)

// General settings
var (
	FontPath        = (config.AssetsPath + "/font.ttf")
	FontSizeHeader  = 36.0
	FontSize        = 24.0
	TextMargin      = FontSize / 2
	FrameWidth      = 900
	FrameMargin     = 50
	BaseCardWidth   = FrameWidth - (2 * FrameMargin)
	BaseCardHeigh   = 150
	BaseCardColor   = color.RGBA{30, 30, 30, 204}
	DecorLinesColor = color.RGBA{80, 80, 80, 255}
	BigTextColor    = color.RGBA{255, 255, 255, 255}
	SmallTextColor  = color.RGBA{204, 204, 204, 255}
	AltTextColor    = color.RGBA{100, 100, 100, 255}

	PremiumColor  = color.RGBA{255, 223, 0, 255}
	VerifiedColor = color.RGBA{72, 167, 250, 255}
)

// AddAllCardsToFrame - Render all cards to frame
func AddAllCardsToFrame(finalCards AllCards, timestamp time.Time, timestampText string, bgImage image.Image) (*gg.Context, error) {
	if len(finalCards.Cards) == 0 {
		return nil, fmt.Errorf("no cards to be rendered")
	}

	// Frame height
	var totalCardsHeight int
	for _, card := range finalCards.Cards {
		totalCardsHeight += card.Context.Height()
	}
	totalCardsHeight += ((len(finalCards.Cards)-1)*(FrameMargin/2) + (FrameMargin * 2))
	// Get frame CTX
	ctx, err := prepBgContext(totalCardsHeight, bgImage)
	if err != nil {
		return finalCards.Frame, err
	}
	finalCards.Frame = ctx
	// Sort cards
	sort.Slice(finalCards.Cards, func(i, j int) bool {
		return finalCards.Cards[i].Index < finalCards.Cards[j].Index
	})
	// Render cards
	var lastCardPos int = FrameMargin / 2
	for i := 0; i < len(finalCards.Cards); i++ {
		card := finalCards.Cards[i]
		cardMarginH := lastCardPos + (FrameMargin / 2)
		finalCards.Frame.DrawImage(card.Image, FrameMargin, cardMarginH)
		lastCardPos = cardMarginH + card.Context.Height()
	}

	// Draw timestamp
	if err := finalCards.Frame.LoadFontFace(FontPath, FontSize*0.75); err != nil {
		return finalCards.Frame, err
	}
	finalCards.Frame.SetColor(color.RGBA{100, 100, 100, 100})
	time := timestamp.Format(fmt.Sprintf("%s Jan 2", timestampText))
	timeW, timeH := finalCards.Frame.MeasureString(time)
	timeX := (float64(finalCards.Frame.Width()) - timeW) / 2
	timeY := (float64(FrameMargin)-timeH)/2 + timeH
	finalCards.Frame.DrawString(time, timeX, timeY)

	return finalCards.Frame, nil
}

// Prepare a frame background context
func prepBgContext(totalHeight int, bgImage image.Image) (frameCtx *gg.Context, err error) {
	frameCtx = gg.NewContext(FrameWidth, totalHeight)
	bgImage = imaging.Fill(bgImage, frameCtx.Width(), frameCtx.Height(), imaging.Center, imaging.NearestNeighbor)
	bgImage = imaging.Blur(bgImage, 10.0)
	frameCtx.DrawImage(bgImage, 0, 0)
	return frameCtx, nil
}

// PrepNewCard - Prepare a new cardData struct
func PrepNewCard(index int, heightMod float64) CardData {
	cardHeight := int(float64(BaseCardHeigh) * heightMod)
	cardWidth := BaseCardWidth
	cardCtx := gg.NewContext(cardWidth, cardHeight)
	cardCtx.SetColor(BaseCardColor)
	cardCtx.DrawRoundedRectangle(0, 0, float64(cardWidth), float64(cardHeight), FontSize)
	cardCtx.Fill()
	var card CardData
	card.Context = cardCtx
	card.Index = index
	return card
}
