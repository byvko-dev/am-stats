package render

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"

	mongodbapi "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	"github.com/cufee/am-stats/render"
	wgapi "github.com/cufee/am-stats/wargamingapi"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func renderBlock(block *cardBlockData) (err error) {
	ctx := gg.NewContext(block.Width, block.Height)

	// Color is requested
	if block.Color != (color.RGBA{}) {
		ctx.SetColor(block.Color)
		ctx.DrawRectangle(0, 0, float64(block.Width), float64(block.Height))
		ctx.Fill()
	}

	// Margins
	if block.IconURL != "" {
		block.TotalTextHeight = block.IconSize
		block.TotalTextLines++
	}
	altTextW, altTextH, _ := getTextParams(ctx, block, (block.AltTextSize), block.AltText)
	_, smlTextH, smlTextDrwX := getTextParams(ctx, block, (block.SmallTextSize), block.SmallText)
	_, bigTextH, bigTextDrwX := getTextParams(ctx, block, (block.BigTextSize), block.BigText)

	// Draw text
	var lastY float64
	var drawTextMargins float64
	drawTextMargins = float64(((block.Height) - block.TotalTextHeight) / (block.TotalTextLines + 1))

	// Icon and Alt text
	if block.IconURL != "" {
		// Load Icon
		var icon image.Image
		if icon, err = loadIcon(block.IconURL); err != nil {
			return err
		}

		// Resize
		icon = imaging.Fill(icon, block.IconSize, block.IconSize, imaging.Center, imaging.Box)

		// Paste Icon
		IcondrawX := getAlignedX(block.TextAlign, float64(block.Width), float64(block.IconSize))
		ctx.DrawImage(icon, int(IcondrawX), int(lastY))
		lastY += (drawTextMargins / 2) + float64(block.IconSize)

		if block.AltText != "" {
			if err := ctx.LoadFontFace(render.FontPath, block.AltTextSize); err != nil {
				return err
			}
			ctx.SetColor(block.AltTextColor)
			lastY := lastY + altTextH
			drawX := getAlignedX(0, float64(block.IconSize), altTextW) + IcondrawX
			ctx.DrawString(block.AltText, drawX, lastY)
		}
	}

	// Big text
	if block.BigText != "" && block.IconURL == "" {
		if err := ctx.LoadFontFace(render.FontPath, block.BigTextSize); err != nil {
			return err
		}
		ctx.SetColor(block.BigTextColor)
		lastY = lastY + bigTextH + drawTextMargins
		ctx.DrawString(block.BigText, bigTextDrwX, lastY)
	}

	// Small text
	if block.SmallText != "" {
		if err := ctx.LoadFontFace(render.FontPath, (block.SmallTextSize)); err != nil {
			return err
		}
		ctx.SetColor(block.SmallTextColor)
		lastY = lastY + smlTextH + drawTextMargins

		ctx.DrawString(block.SmallText, smlTextDrwX, lastY)
	}

	block.Context = ctx
	return err
}

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

func getTextParams(ctx *gg.Context, block *cardBlockData, size float64, value string) (width float64, height float64, drawX float64) {
	// Return 0
	if value == "" {
		return width, height, drawX
	}

	// Load font and measure text
	if err := ctx.LoadFontFace(render.FontPath, size); err != nil {
		return width, height, drawX
	}
	width, height = ctx.MeasureString(value)
	drawX = getAlignedX(block.TextAlign, float64(block.Width), width)
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

func addScoreAndMedals(card *render.CardData, blueprint cardBlockData, score int, position int, medals []mongodbapi.MedalWeight) error {
	// Score Block
	scoreBlock := cardBlockData(blueprint)
	scoreBlock.BigText = fmt.Sprint(score)
	scoreBlock.Width = int(blueprint.SpecialBlockWidth)
	scoreBlock.SmallText = "Score"

	if err := renderBlock(&scoreBlock); err != nil {
		return err
	}
	card.Context.DrawImage(scoreBlock.Context.Image(), card.LastXOffs, 0)
	card.LastXOffs += scoreBlock.Width

	//  Medal Blocks
	for _, m := range medals {
		medalBlock := cardBlockData(blueprint)
		medalBlock.AltText = fmt.Sprint(m.Score)
		medalBlock.AltTextColor = blueprint.SmallTextColor
		medalBlock.IconURL = m.IconURL
		medalBlock.TextAlign = 1

		if err := renderBlock(&medalBlock); err != nil {
			return err
		}
		card.Context.DrawImage(medalBlock.Context.Image(), card.LastXOffs, 0)
		card.LastXOffs += int(card.BlockWidth)
	}

	// Render image
	card.Image = card.Context.Image()
	return nil
}

// renderCardBlocks - Render all card blocks
func renderCardBlocks(card *render.CardData, position int, medals []mongodbapi.MedalWeight) error {
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
			blockExtra := block.Extra.(*cardBlockData)
			blockExtra.Height = card.Context.Height()
			blockExtra.Width = block.Width

			// Render block image
			if err := renderBlock(blockExtra); err != nil {
				log.Print(err)
				return
			}

			// Calculate rendering offset
			var offset int = card.FrameMargin
			for _, b := range card.Blocks[:i] {
				offset += b.Width
			}

			// Draw block
			card.Context.DrawImage(blockExtra.Context.Image(), offset, 0)
		}(block, i)
	}
	wg.Wait()

	// Render image
	card.Image = card.Context.Image()
	return nil
}
