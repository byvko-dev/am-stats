package render

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"sort"

	db "github.com/cufee/am-stats/mongodbapi/v1/players"
	"github.com/cufee/am-stats/render"
	"github.com/disintegration/imaging"
)

func makeStatsPlusHeaderCard(card render.CardData, playerName, playerClan string, premium, verified bool, pins []db.UserPin) (render.CardData, error) {
	// Prep pins
	sort.Slice(pins, func(i, j int) bool { return pins[i].Weight > pins[j].Weight })
	if len(pins) > 7 {
		pins = pins[:7]
	}

	ctx := *card.Context
	if err := ctx.LoadFontFace(render.FontPath, render.FontSizeHeader); err != nil {
		return card, err
	}
	if playerClan != "" {
		playerName += " "
	}
	// Calculate text size
	nameStrW, nameStrH := ctx.MeasureString(playerName)
	clanStrW, _ := ctx.MeasureString(playerClan)
	totalTextW := nameStrW + clanStrW

	xOffset := ((float64(ctx.Width()) - totalTextW) / 2)
	yOffset := float64(render.FrameMargin / 2)
	// Draw player name and tag text
	ctx.SetColor(color.White)
	if premium {
		ctx.SetColor(render.PremiumColor)
	}
	psDrawX := ((totalTextW - nameStrW - clanStrW) / 2) + xOffset
	psDrawY := nameStrH + yOffset
	// Draw name
	ctx.DrawString(playerName, psDrawX, psDrawY)
	// Draw tag
	ctx.SetColor(color.White)
	ctx.DrawString(playerClan, (psDrawX + nameStrW), psDrawY)
	// Draw verified icon
	if verified {
		// Draw icon
		radius := (render.FontSizeHeader / 3)
		ctx.SetColor(render.VerifiedColor)
		iX := psDrawX - radius*2
		iY := psDrawY - (nameStrH / 2)
		ctx.DrawCircle(iX, iY, radius)
		ctx.Fill()

		// Draw checkmark
		ctx.SetColor(color.White)
		ctx.SetLineWidth(radius / 3)
		lineX1 := iX - radius/4 - radius/8
		lineY1 := iY - radius/4 + radius/4
		lineX2 := iX - radius/8
		lineY2 := iY + radius/4
		lineX3 := iX + radius/2 - radius/8
		lineY3 := iY - radius/2 + radius/4
		ctx.DrawLine(lineX1, lineY1, lineX2, lineY2)
		ctx.DrawLine(lineX2, lineY2, lineX3, lineY3)
		ctx.Stroke()
	}

	// Draw lines
	ctx.SetColor(render.DecorLinesColor)
	lineX := float64(render.FrameMargin)
	lineY := psDrawY + nameStrH - float64(render.FrameMargin/4)
	lineHeight := 2.0
	lineWidth := (float64(ctx.Width()) - float64(render.FrameMargin*2))
	ctx.DrawRectangle(lineX, lineY, lineWidth, lineHeight)
	ctx.Fill()

	// Draw label
	if err := ctx.LoadFontFace(render.FontPath, render.FontSizeHeader/2); err != nil {
		return card, err
	}
	ctx.SetColor(render.AltTextColor)
	labelStrW, labelStrH := ctx.MeasureString("Pin Collection")
	var labelX float64 = (float64(ctx.Width()) - labelStrW) / 2
	var labelY float64 = lineY + labelStrH + labelStrH/1.25
	ctx.DrawString("Pin Collection", labelX, labelY)

	// Draw pins
	var iconMarginX int = render.FrameMargin / 2
	var pinSize int = (ctx.Height() - int(labelY) - iconMarginX)
	var lastX int = (ctx.Width() - (len(pins)*pinSize + (len(pins)-1)*(iconMarginX))) / 2
	var iconsY int = int(labelY) + (ctx.Height()-pinSize-int(labelY))/2
	for i, pin := range pins {
		var drawX int
		if i == 0 {
			drawX = lastX
		} else {
			drawX = lastX + pinSize + iconMarginX
		}
		lastX = drawX
		if pin.Label != "" {
			if err := ctx.LoadFontFace(render.FontPath, render.FontSizeHeader/2); err != nil {
				return card, err
			}
			ctx.SetColor(render.BigTextColor)
			// Calculate text size
			labelStrW, labelStrH := ctx.MeasureString(pin.Label)
			pin.Size = pinSize - render.FrameMargin/6 - int(labelStrH)
			// Draw text
			labelX := float64((drawX + ((pinSize - pin.Size) / 2))) + (float64(pin.Size)-labelStrW)/2
			labelY := float64(iconsY+pin.Size+render.FrameMargin/6) + labelStrH
			ctx.DrawString(pin.Label, labelX, labelY)
		} else {
			pin.Size = pinSize
		}

		// Load Icon
		var err error
		var icon image.Image
		if icon, err = render.LoadIcon(pin.URL); err != nil {
			return card, err
		}
		if icon == nil {
			log.Print("nil image")
			return card, fmt.Errorf("image was nil")
		}
		// Resize
		icon = imaging.Fill(icon, pin.Size, pin.Size, imaging.Center, imaging.Box)
		// Paste icon
		ctx.DrawImage(icon, int(drawX+((pinSize-pin.Size)/2)), iconsY)
	}

	// Make image
	card.Image = ctx.Image()
	return card, nil
}
