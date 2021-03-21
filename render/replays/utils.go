package render

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/cufee/am-stats/render"
	wgapi "github.com/cufee/am-stats/wargamingapi"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func getAlignedX(alignment int, fieldW float64, elemW float64) float64 {
	switch alignment {
	case -1:
		return 0
	case 1:
		return (fieldW - elemW)
	default:
		return ((fieldW - elemW) / 2.0)
	}

}

func loadIcon(url string) (img image.Image, err error) {
	// Get image
	response, _ := http.Get(url)
	if response != nil {
		defer response.Body.Close()

		// Decode image
		if img, _, err = image.Decode(response.Body); err != nil {
			return img, err
		}
	}
	return img, err
}

func getTextParams(ctx *gg.Context, block *replayBlockData, size float64, value string) (width float64, height float64, drawX float64) {
	// Return 0
	if value == "" {
		return width, height, drawX
	}

	// Load font and measure text
	if err := ctx.LoadFontFace(render.FontPath, size); err != nil {
		return width, height, drawX
	}
	width, height = ctx.MeasureString(value)
	drawX = getAlignedX(block.TextAlign, float64(block.Width-block.Padding), width) + float64(block.Padding/2)
	block.TotalTextHeight += int(height)
	block.TotalTextLines++
	return width, height, drawX
}

// getField - Get field value by name
func getField(data wgapi.AchievementsFrame, field string) int {
	v := reflect.ValueOf(&data.Achievements)
	f := reflect.Indirect(v).FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == strings.ToLower(field) })
	if f == (reflect.Value{}) {
		return 0
	}
	return int(f.Int())
}

// renderCardBlocks - Render all card blocks
func renderCardBlocks(card *render.CardData) error {
	// Atomic counter
	var wg sync.WaitGroup

	for i, block := range card.Blocks {
		if block.Extra == nil {
			continue
		}

		wg.Add(1)
		go func(block render.Block, i int) {
			defer wg.Done()

			// Fill block Width and Height for legacy code
			if block.ExtraType == "hpbar" {
				blockExtra := block.Extra.(*hpBarBlockData)
				blockExtra.Height = card.Context.Height()
				blockExtra.Width = block.Width + block.Padding
				if blockExtra.Color == (color.RGBA{}) {
					blockExtra.Color = block.Color
				}

				// Render block image
				if err := renderHPBarBlock(block.Extra); err != nil {
					log.Print(err)
					return
				}

				block.Context = blockExtra.Context
			} else {
				blockExtra := block.Extra.(*replayBlockData)
				blockExtra.Width = block.Width + block.Padding
				blockExtra.Height = card.Context.Height()
				blockExtra.Padding = block.Padding
				if blockExtra.Color == (color.RGBA{}) {
					blockExtra.Color = block.Color
				}

				// Render block image
				if err := renderIconBlock(block.Extra); err != nil {
					log.Print(err)
					return
				}

				block.Context = blockExtra.Context
			}

			// Calculate rendering offset
			var offset int = card.FrameMargin
			for _, b := range card.Blocks[:i] {
				offset += b.Width + b.Padding
			}

			// Draw block
			card.Context.DrawImage(block.Context.Image(), offset, 0)
		}(block, i)
	}
	wg.Wait()

	// Render image
	card.Image = card.Context.Image()
	return nil
}

// AddAllCardsToFrame - Render all cards to frame
func renderAllCardsOnFrame(finalCards render.AllCards, header string, bgImage image.Image) (*gg.Context, error) {
	if len(finalCards.Cards) == 0 {
		return nil, fmt.Errorf("no cards to be rendered")
	}

	// Frame height
	maxIndexXMap := make(map[int]int)
	totalStackedCardsMap := make(map[int]int)
	totalCardsHeightMap := make(map[int]int)
	var totalHeaderCardsHeight int
	var totalHeaderCards int
	var totalCardsWidth int
	var maxIndexX int
	var maxWidth int
	for _, card := range finalCards.Cards {
		maxIndexXMap[card.IndexX] = card.IndexX
		if maxIndexXMap[card.IndexX] > maxIndexX {
			maxIndexX = maxIndexXMap[card.IndexX]
		}

		if card.Type == render.CardTypeHeader {
			totalHeaderCards++
			totalHeaderCardsHeight += card.Context.Height()
		} else {
			totalStackedCardsMap[card.IndexX]++
			totalCardsHeightMap[card.IndexX] += card.Context.Height()
		}

		if card.Context.Width() > maxWidth && card.Type != render.CardTypeHeader {
			maxWidth = card.Context.Width()
		}
	}

	// Save max height and cards stacked
	var totalCardsHeight int
	var totalStackedCards int
	for _, t := range totalStackedCardsMap {
		if (t + totalHeaderCards) > totalStackedCards {
			totalStackedCards = t + totalHeaderCards
		}
	}
	for _, h := range totalCardsHeightMap {
		if (h + totalHeaderCardsHeight) > totalCardsHeight {
			totalCardsHeight = h + totalHeaderCardsHeight
		}

	}
	totalCardsHeight += ((totalStackedCards-1)*(render.FrameMargin/4) + (render.FrameMargin * 2))
	totalCardsWidth += (maxIndexX+1)*maxWidth + render.FrameMargin/2

	// Get frame CTX
	ctx, err := prepBgContext(totalCardsHeight, totalCardsWidth, bgImage)
	if err != nil {
		return finalCards.Frame, err
	}
	finalCards.Frame = ctx
	// Sort cards
	sort.Slice(finalCards.Cards, func(i, j int) bool {
		return finalCards.Cards[i].Index < finalCards.Cards[j].Index
	})

	// Render cards
	var lastCardPos map[int]int = make(map[int]int)
	for i := 0; i < len(finalCards.Cards); i++ {
		card := finalCards.Cards[i]

		if lastCardPos[card.IndexX] == 0 {
			lastCardPos[card.IndexX] = render.FrameMargin
		}

		cardMarginH := lastCardPos[card.IndexX] + (render.FrameMargin / 4)
		finalCards.Frame.DrawImage(card.Image, (render.FrameMargin + (render.FrameMargin/2+maxWidth)*card.IndexX + 1), cardMarginH)
		if card.Type == render.CardTypeHeader {
			for n := 0; n <= maxIndexX; n++ {
				lastCardPos[n] = cardMarginH + card.Context.Height()
			}
		} else {
			lastCardPos[card.IndexX] = cardMarginH + card.Context.Height()
		}
	}

	// Draw header
	if err := finalCards.Frame.LoadFontFace(render.FontPath, render.FontSize); err != nil {
		return finalCards.Frame, err
	}
	finalCards.Frame.SetColor(color.RGBA{100, 100, 100, 100})
	headerW, headerH := finalCards.Frame.MeasureString(header)
	headerX := (float64(finalCards.Frame.Width()) - headerW) / 2
	headerY := (float64(render.FrameMargin)-headerH)/2 + headerH
	finalCards.Frame.DrawString(header, headerX, headerY)

	return finalCards.Frame, nil
}

// Prepare a frame background context
func prepBgContext(totalHeight int, width int, bgImage image.Image) (frameCtx *gg.Context, err error) {
	frameWidth := width + (2 * render.FrameMargin)
	frameCtx = gg.NewContext(frameWidth, totalHeight)
	bgImage = imaging.Fill(bgImage, frameCtx.Width(), frameCtx.Height(), imaging.Center, imaging.NearestNeighbor)
	bgImage = imaging.Blur(bgImage, 10.0)
	frameCtx.DrawImage(bgImage, 0, 0)
	return frameCtx, nil
}

// intInSlice - Check if int exuist in slice
func intInSlice(slice []int, i int) bool {
	for _, n := range slice {
		if i == n {
			return true
		}
	}
	return false
}

func masteryToIconURL(mastery int) string {
	switch mastery {
	case 3:
		return "http://glossary-eu-static.gcdn.co/icons/wotb/current/achievements/markOfMasteryI.png"
	case 2:
		return "http://glossary-eu-static.gcdn.co/icons/wotb/current/achievements/markOfMasteryII.png"
	case 1:
		return "http://glossary-eu-static.gcdn.co/icons/wotb/current/achievements/markOfMasteryIII.png"
	case 4:
		return "http://glossary-eu-static.gcdn.co/icons/wotb/current/achievements/markOfMastery.png"
	}
	return ""
}
