package render

import (
	"image"
	"image/color"

	"github.com/cufee/am-stats/render"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func renderHPBarBlock(blockInterface interface{}) (err error) {
	block := blockInterface.(*hpBarBlockData)
	ctx := gg.NewContext(block.Width, block.Height)
	defer func() { block.Context = ctx }()

	// Color is requested
	if block.Color != (color.RGBA{}) {
		ctx.SetColor(block.Color)
		ctx.DrawRectangle(0, 0, float64(block.Width), float64(block.Height))
		ctx.Fill()
	}

	// Draw background
	barHeight := block.Height - block.Margin*2
	ctx.SetColor(block.HPColorBG)
	ctx.DrawRoundedRectangle(0, float64(block.Margin), float64(block.Width), float64(barHeight), float64(block.Width)/2.1)
	ctx.Fill()

	if block.PercentHP > 0 {
		// Draw hp
		hpMargin := (1 - block.PercentHP) * float64(barHeight)
		ctx.SetColor(block.HPColor)
		ctx.DrawRoundedRectangle(0, float64(block.Margin)+hpMargin, float64(block.Width), float64(barHeight)*block.PercentHP, float64(block.Width)/2.1)
		ctx.Fill()
	}
	return err
}

func renderIconBlock(blockInterface interface{}) (err error) {
	block := blockInterface.(*replayBlockData)

	ctx := gg.NewContext(block.Width, block.Height)
	defer func() { block.Context = ctx }()

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
	} else {
		// Get total text height and lines
		for _, line := range block.TextLines {
			textSize := block.TextSize
			if line.TextScale > 0 {
				textSize = line.TextScale * textSize
			}
			getTextParams(ctx, block, textSize, line.Text)
		}
	}

	// Draw text
	var lastY float64
	var drawTextMargins float64
	drawTextMargins = float64(((block.Height) - block.TotalTextHeight) / (block.TotalTextLines + 1))

	// Icon and Alt text
	if block.IconURL != "" {
		lastY = float64(block.Height-block.IconSize) / 2

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
		return nil
	}

	// Text
	if len(block.TextLines) > 0 && block.IconURL == "" {
		// Load font
		if err := ctx.LoadFontFace(render.FontPath, block.TextSize); err != nil {
			return err
		}

		// Draw each line
		for _, line := range block.TextLines {
			textSize := block.TextSize
			if line.TextScale > 0 {
				textSize = line.TextScale * textSize
			}
			_, textH, drawX := getTextParams(ctx, block, textSize, line.Text)

			if line.Color != (color.RGBA{}) {
				ctx.SetColor(line.Color)
			} else {
				ctx.SetColor(block.TextColor)
			}
			lastY = lastY + textH + drawTextMargins
			ctx.DrawString(line.Text, drawX, lastY)
		}
	}
	return err
}
