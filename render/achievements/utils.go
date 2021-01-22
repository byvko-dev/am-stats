package render

import (
	"image"
	"image/color"
	"net/http"
	"reflect"
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
		block.TotalTextHeight = block.IconSize
		block.TotalTextLines++
	}
	_, altTextH, altTextDrwX := getTextParams(ctx, block, (block.AltTextSize), block.AltText)
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
		drawX := (((block.Width) - (block.IconSize)) / 2.0)
		ctx.DrawImage(icon, drawX, int(lastY))
		lastY += (drawTextMargins / 2) + float64(block.IconSize)

		if block.AltText != "" {
			if err := ctx.LoadFontFace(render.FontPath, block.AltTextSize); err != nil {
				return err
			}
			ctx.SetColor(block.AltTextColor)
			lastY := lastY + altTextH
			ctx.DrawString(block.AltText, altTextDrwX, lastY)
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
