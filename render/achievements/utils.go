package render

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"

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
		block.TotalTextHeight = float64(block.IconSize)
		block.TotalTextLines++
	}
	_, altTextH, altTextDrwX := getTextParams(ctx, block, (block.BlockTextSize * block.TextCoeff * block.TextCoeff), block.AltText)
	_, smlTextH, smlTextDrwX := getTextParams(ctx, block, (block.BlockTextSize * block.TextCoeff), block.SmallText)
	_, bigTextH, bigTextDrwX := getTextParams(ctx, block, (block.BlockTextSize), block.BigText)

	// Draw text
	var drawTextMargins float64
	var lastY float64
	drawTextMargins = (float64(block.Height) - block.TotalTextHeight) / float64(block.TotalTextLines+1)

	// Icon and Alt text
	if block.IconURL != "" {
		if err := ctx.LoadFontFace(render.FontPath, (block.BlockTextSize * block.TextCoeff * block.BlockTextSize)); err != nil {
			return err
		}

		// Load Icon
		var icon image.Image
		if icon, err = loadIcon(block.IconURL); err != nil {
			return err
		}
		// Resize
		icon = imaging.Fill(icon, block.IconSize, block.IconSize, imaging.Center, imaging.NearestNeighbor)

		// Paste Icon
		drawX := (((block.Width) - (block.IconSize)) / 2.0)
		ctx.DrawImage(icon, drawX, int(lastY+drawTextMargins))
		lastY += (drawTextMargins / 2) + float64(block.IconSize)

		if block.AltText != "" {
			ctx.SetColor(block.AltTextColor)
			lastY := lastY + drawTextMargins + altTextH
			ctx.DrawString(block.AltText, altTextDrwX, lastY)
		}
	}

	// Big text
	if block.BigText != "" && block.IconURL == "" {
		if err := ctx.LoadFontFace(render.FontPath, block.BlockTextSize); err != nil {
			return err
		}
		ctx.SetColor(block.BigTextColor)
		lastY = lastY + bigTextH + drawTextMargins
		ctx.DrawString(block.BigText, bigTextDrwX, lastY)
	}

	// Small text
	if block.SmallText != "" {
		if err := ctx.LoadFontFace(render.FontPath, (block.BlockTextSize * block.TextCoeff)); err != nil {
			return err
		}
		ctx.SetColor(block.SmallTextColor)
		lastY = lastY + smlTextH + drawTextMargins

		ctx.DrawString(block.SmallText, smlTextDrwX, lastY)
	}

	block.Context = ctx
	return nil
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
	drawX = ((float64(block.Width) - width) / 2.0)
	block.TotalTextHeight += height
	block.TotalTextLines++
	return width, height, drawX
}

// Recover from panic
func recoverPanic(err error, funcName string) {
	if r := recover(); r != nil {
		message := fmt.Sprintf("recovered from panic in %s\nlast error: %v", funcName, err)
		log.Print(message)
		debug.PrintStack()
	}
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
