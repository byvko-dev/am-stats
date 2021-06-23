package render

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"strconv"
	"sync"

	dataprep "github.com/cufee/am-stats/dataprep/stats"
	stats "github.com/cufee/am-stats/dataprep/stats"
	"github.com/cufee/am-stats/render"
	wgapi "github.com/cufee/am-stats/wargamingapi"
	"github.com/fogleman/gg"
)

// ImageFromStats -
func ImageFromStats(data stats.ExportData, sortKey string, tankLimit int, premium bool, verified bool, bgImage image.Image) (finalImage image.Image, err error) {
	// // Calculate card width
	// checkCtx := gg.NewContext(1, 1)
	// // Measure player name and clan
	// if err := checkCtx.LoadFontFace(render.FontPath, render.FontSizeHeader); err != nil {
	// 	return nil, err
	// }
	// playerNameW, _ := checkCtx.MeasureString(data.PlayerDetails.Name + " " + data.PlayerDetails.ClanTag)
	// if err := checkCtx.LoadFontFace(render.FontPath, render.FontSizeHeader); err != nil {
	// 	return nil, err
	// }
	// maxCardWidth := playerNameW + render.FontSizeHeader

	var finalCards render.AllCards
	cardsChan := make(chan render.CardData, (3 + len(data.SessionStats.Vehicles)))
	var wg sync.WaitGroup
	// Work on cards in go routines
	wg.Add(1)
	// Header card
	go func() {
		defer wg.Done()
		// Compile Clan tag
		clanTag := ""
		if data.PlayerDetails.ClanTag != "" {
			clanTag = "[" + data.PlayerDetails.ClanTag + "]"
		}

		// Make Header card
		headerHeight := 0.5
		var header render.CardData
		render.PrepNewCard(&header, 0, headerHeight, 0)
		header, err := makeStatsHeaderCard(header, data.PlayerDetails.Name, clanTag, premium, verified)
		if err != nil {
			log.Println(err)
			return
		}
		cardsChan <- header
	}()

	// All stats card
	wg.Add(1)
	go func() {
		defer wg.Done()
		if data.SessionStats.StatsAll.Battles > 0 {
			var randomStats render.CardData
			render.PrepNewCard(&randomStats, 1, 1.25, 0)
			randomStats, err := makeSessionStatsCard(randomStats, "Random Battles", data.SessionStats.StatsAll, data.LastSession.StatsAll, data.SessionStats.SessionRating, data.PlayerDetails.CareerWN8)
			if err != nil {
				log.Println(err)
			} else {
				cardsChan <- randomStats
			}
		}
		if data.SessionStats.StatsRating.Battles > 0 {
			var ratingStats render.CardData
			render.PrepNewCard(&ratingStats, 1, 1.25, 0)
			ratingStats, err := makeSessionStatsCard(ratingStats, "Rating Battles", data.SessionStats.StatsRating, data.LastSession.StatsRating, -1, -1)
			if err != nil {
				log.Println(err)
			} else {
				cardsChan <- ratingStats
			}
		}
		return
	}()

	// Sort vehicles
	vehicles := dataprep.SortTanks(data.SessionStats.Vehicles, sortKey)
	// Create cards for each vehicle in routines
	for i, tank := range vehicles {
		if i == tankLimit {
			break
		}
		wg.Add(1)
		go func(tank wgapi.VehicleStats, i int) {
			defer wg.Done()
			lastSession := data.LastSession.Vehicles[strconv.Itoa(tank.TankID)]
			var tankCard render.CardData
			if i < 3 {
				render.PrepNewCard(&tankCard, (i + 2), 1.0, 0)
				tankCard, err = makeDetailedCard(tankCard, tank, lastSession)
			} else {
				render.PrepNewCard(&tankCard, (i + 2), 0.5, 0)
				tankCard, err = makeSlimCard(tankCard, tank, lastSession)
			}
			if err != nil {
				log.Println(err)
				return
			}
			cardsChan <- tankCard
			return
		}(tank, i)
	}
	wg.Wait()
	close(cardsChan)

	for c := range cardsChan {
		finalCards.Cards = append(finalCards.Cards, c)
	}

	header := data.SessionStats.Timestamp.Format(fmt.Sprintf("%s Jan 2", "Session from"))
	finalCtx, err := render.AddAllCardsToFrame(finalCards, header, data.PlayerDetails.Realm, bgImage)
	if err != nil {
		return nil, err
	}
	return finalCtx.Image(), err
}

// Makes a detailed card for a tank
func makeDetailedCard(card render.CardData, session wgapi.VehicleStats, lastSession wgapi.VehicleStats) (render.CardData, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()

	ctx := *card.Context
	if err := ctx.LoadFontFace(render.FontPath, (render.FontSize * 1.25)); err != nil {
		return card, err
	}

	if session.Battles < 1 {
		return card, fmt.Errorf("sessions battles is < 1")
	}

	ctx.SetColor(color.White)
	// Measure tank name
	nameW, nameH := ctx.MeasureString(session.TankName)
	if err := ctx.LoadFontFace(render.FontPath, (render.FontSize * 0.75)); err != nil {
		return card, err
	}
	tierW, tierH := ctx.MeasureString(tierToRoman(session.TankTier))
	headerHeigth := int(nameH * 2)
	nameX := (float64(card.Context.Width()) - nameW) / 2
	nameY := (float64(headerHeigth)-nameH)/2 + nameH
	tierX := nameX - (render.FontSize / 2) - tierW
	tierY := (float64(headerHeigth)-tierH)/2 + tierH

	// Draw tank tier
	ctx.DrawString(tierToRoman(session.TankTier), tierX, tierY)
	// Draw tank name
	if err := ctx.LoadFontFace(render.FontPath, (render.FontSize * 1.25)); err != nil {
		return card, err
	}
	ctx.DrawString(session.TankName, nameX, nameY)

	blockWidth := card.Context.Width() / 4
	availableHeight := int(ctx.Height() - (headerHeigth))
	// Blocks will take 75% of the total card heiht
	blockHeight := availableHeight
	// Default Block
	var defaultBlock statsBlock
	defaultBlock.TextSize = render.FontSize * 1.30
	defaultBlock.TextCoeff = 0.75
	defaultBlock.Width = blockWidth
	defaultBlock.Height = blockHeight
	defaultBlock.BigTextColor = render.BigTextColor
	defaultBlock.SmallTextColor = render.SmallTextColor
	defaultBlock.AltTextColor = render.AltTextColor

	// Bottom Row - Avg Damage, Avg XP, Winrate
	// Block 1 - Battles
	battlesBlock := statsBlock(defaultBlock)
	battlesBlock.Width = blockWidth
	battlesSession := strconv.Itoa(session.Battles)
	battlesLastSession := "-"
	if lastSession.Battles > 0 {
		battlesLastSession = strconv.Itoa(lastSession.Battles)
	}
	battlesBlock.SmallText = battlesLastSession
	battlesBlock.BigText = battlesSession
	battlesBlock.AltText = "Battles"
	battlesBlock, err := addStatsBlockCtx(battlesBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(battlesBlock.Context.Image(), 0, headerHeigth)
	// Block 2 - Avg Damage
	avgDamageBlock := statsBlock(defaultBlock)
	avgDamageBlock.Width = blockWidth
	avgDamageSession := strconv.Itoa((session.DamageDealt / session.Battles))
	avgDamageLastSession := "-"
	avgDamageBlock.HasBigIcon = true
	if lastSession.Battles > 0 {
		avgDamageLastSession = strconv.Itoa((lastSession.DamageDealt / lastSession.Battles))
		if (lastSession.DamageDealt / lastSession.Battles) < (session.DamageDealt / session.Battles) {
			avgDamageBlock.BigArrowDirection = 1
		}
		if (lastSession.DamageDealt / lastSession.Battles) > (session.DamageDealt / session.Battles) {
			avgDamageBlock.BigArrowDirection = -1
		}
	}
	avgDamageBlock.SmallText = avgDamageLastSession
	avgDamageBlock.BigText = avgDamageSession
	avgDamageBlock.AltText = "Avg. Damage"
	avgDamageBlock, err = addStatsBlockCtx(avgDamageBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(avgDamageBlock.Context.Image(), (blockWidth), headerHeigth)
	// Block 1 - Winrate
	winrateBlock := statsBlock(avgDamageBlock)
	winrateSession := ((float64(session.Wins) / float64(session.Battles)) * 100)
	winrateLastSession := "-"
	if lastSession.Battles > 0 {
		winrateLastSession = fmt.Sprintf("%.2f", ((float64(lastSession.Wins)/float64(lastSession.Battles))*100)) + "%"
	}
	winrateBlock.BigText = fmt.Sprintf("%.2f", winrateSession) + "%"
	winrateBlock.SmallText = winrateLastSession
	winrateBlock.AltText = "Winrate"
	winrateBlock.HasBigIcon = true
	if ((float64(lastSession.Wins) / float64(lastSession.Battles)) * 100) < winrateSession {
		winrateBlock.BigArrowDirection = 1
	}
	if ((float64(lastSession.Wins) / float64(lastSession.Battles)) * 100) > winrateSession {
		winrateBlock.BigArrowDirection = -1
	}
	winrateBlock, err = addStatsBlockCtx(winrateBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(winrateBlock.Context.Image(), (blockWidth * 2), headerHeigth)
	// Block 4 - Draw WN8
	ratingBlock := statsBlock(defaultBlock)
	// Icon
	ratingBlock.SmallText = "WN8"
	ratingBlock.BigText = "-"
	if session.TankWN8 > -1 {
		ratingBlock.HasBigIcon = true
		ratingBlock.BigIconColor = GetRatingColor(session.TankWN8)
		ratingBlock.BigText = strconv.Itoa(session.TankWN8)
	}
	ratingBlock.SmallTextColor = render.AltTextColor
	ratingBlock, err = addStatsBlockCtx(ratingBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(ratingBlock.Context.Image(), (blockWidth * 3), headerHeigth)

	// Draw lines
	ctx.SetColor(render.DecorLinesColor)
	lineX := float64(render.FrameMargin)
	lineY := float64(headerHeigth)
	lineHeight := 2.0
	lineWidth := (float64(ctx.Width()) - float64(render.FrameMargin*2))
	ctx.DrawRectangle(lineX, lineY, lineWidth, lineHeight)
	ctx.Fill()

	// Render image
	card.Image = ctx.Image()
	return card, nil
}

// Makes a slim detailed card for a tank
func makeSlimCard(card render.CardData, session wgapi.VehicleStats, lastSession wgapi.VehicleStats) (render.CardData, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()

	ctx := *card.Context
	if err := ctx.LoadFontFace(render.FontPath, (render.FontSize)); err != nil {
		return card, err
	}

	if session.Battles < 1 {
		return card, fmt.Errorf("sessions battles is < 1")
	}

	ctx.SetColor(color.White)
	tankNameWidth := float64(card.Context.Width()) * 0.375
	tankBlockWidth := (float64(card.Context.Width()) - tankNameWidth) / 3

	// Default Block
	var defaultBlock statsBlock
	defaultBlock.TextSize = render.FontSize
	defaultBlock.Width = int(tankBlockWidth)
	defaultBlock.Height = card.Context.Height()
	defaultBlock.BigTextColor = render.BigTextColor
	defaultBlock.SmallTextColor = render.AltTextColor

	// Draw tank name
	finalName := ""
	dotsW, _ := ctx.MeasureString("...")

	for _, r := range []rune(session.TankName) {
		w, _ := ctx.MeasureString(finalName)
		if (w + dotsW) > (tankNameWidth - (float64(render.FrameMargin) * 1.5)) {
			finalName = finalName + "..."
			break
		}
		finalName = finalName + string(r)
	}
	_, nameH := ctx.MeasureString(finalName)

	nameY := (float64(card.Context.Height()) - ((float64(card.Context.Height()) - nameH) / 2))
	ctx.DrawString(finalName, (float64(render.FrameMargin) * 1.5), nameY)

	// Draw tank tier
	if err := ctx.LoadFontFace(render.FontPath, (render.FontSize * 0.75)); err != nil {
		return card, err
	}
	tierW, tierH := ctx.MeasureString(tierToRoman(session.TankTier))
	tierX := float64(render.FrameMargin/2) + ((float64(render.FrameMargin) - tierW) / 2)
	tierY := (float64(card.Context.Height()) - ((float64(card.Context.Height()) - tierH) / 2))
	ctx.DrawString(tierToRoman(session.TankTier), tierX, tierY)

	// 3 Blocks - DMG / WR / WN8
	// Block 3 - Draw WN8
	ratingBlock := statsBlock(defaultBlock)
	// Icon
	ratingBlock.SmallText = "WN8"
	ratingBlock.BigText = "-"
	if session.TankWN8 > -1 {
		ratingBlock.HasBigIcon = true
		ratingBlock.BigIconColor = GetRatingColor(session.TankWN8)
		ratingBlock.BigText = strconv.Itoa(session.TankWN8)
	}
	ratingBlock, err := addStatsBlockCtx(ratingBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(ratingBlock.Context.Image(), int(tankNameWidth+(tankBlockWidth*2)), 0)

	// Block 2 - Winrate
	winrateBlock := statsBlock(defaultBlock)
	winrateSession := ((float64(session.Wins) / float64(session.Battles)) * 100)
	winrateBlock.BigText = fmt.Sprintf("%.1f", winrateSession) + "% (" + strconv.Itoa(session.Battles) + ")"
	winrateBlock.SmallText = "Winrate"
	winrateBlock.HasBigIcon = true
	if lastSession.Battles > 0 && ((float64(session.Wins)/float64(session.Battles))*100) > ((float64(lastSession.Wins)/float64(lastSession.Battles))*100) {
		winrateBlock.BigArrowDirection = 1
	}
	if lastSession.Battles > 0 && ((float64(session.Wins)/float64(session.Battles))*100) < ((float64(lastSession.Wins)/float64(lastSession.Battles))*100) {
		winrateBlock.BigArrowDirection = -1
	}
	winrateBlock, err = addStatsBlockCtx(winrateBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(winrateBlock.Context.Image(), int(tankNameWidth+(tankBlockWidth*1)), 0)

	// Block 1 - Avg Damage
	avgDamageBlock := statsBlock(defaultBlock)
	avgDamageSession := strconv.Itoa((session.DamageDealt / session.Battles))
	avgDamageBlock.SmallText = "Avg. Damage"
	avgDamageBlock.BigText = avgDamageSession
	avgDamageBlock.HasBigIcon = true
	if lastSession.Battles > 0 && (session.DamageDealt/session.Battles) > (lastSession.DamageDealt/lastSession.Battles) {
		avgDamageBlock.BigArrowDirection = 1
	}
	if lastSession.Battles > 0 && (session.DamageDealt/session.Battles) < (lastSession.DamageDealt/lastSession.Battles) {
		avgDamageBlock.BigArrowDirection = -1
	}
	avgDamageBlock, err = addStatsBlockCtx(avgDamageBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(avgDamageBlock.Context.Image(), int(tankNameWidth), 0)

	// Render image
	card.Image = ctx.Image()
	return card, nil
}

func addStatsBlockCtx(block statsBlock) (statsBlock, error) {
	ctx := gg.NewContext(block.Width, block.Height)
	// Color is requested
	if block.IsColored == true {
		ctx.SetColor(block.Color)
		ctx.DrawRectangle(0, 0, float64(block.Width), float64(block.Height))
		ctx.Fill()
	}
	if block.TextCoeff == 0 {
		block.TextCoeff = 0.6
	}
	// Calc altText
	var (
		altMargin float64
		aTxtW     float64
		aTxtH     float64
	)
	if block.AltText != "" {
		ctx.SetColor(block.AltTextColor)
		if err := ctx.LoadFontFace(render.FontPath, (block.TextSize * (block.TextCoeff - 0.15))); err != nil {
			return block, err
		}
		aTxtW, aTxtH = ctx.MeasureString(block.AltText)
		altMargin = aTxtH
	}
	availHeiht := block.Height
	var totalTextHeight float64 = altMargin
	// Calc small text
	if err := ctx.LoadFontFace(render.FontPath, (block.TextSize * block.TextCoeff)); err != nil {
		return block, err
	}
	sTxtW, sTxtH := ctx.MeasureString(block.SmallText)
	if sTxtH > 0 {
		totalTextHeight += sTxtH
	}
	sX := ((float64(block.Width) - sTxtW) / 2.0)
	// Calc Big text
	if err := ctx.LoadFontFace(render.FontPath, block.TextSize); err != nil {
		return block, err
	}
	bTxtW, bTxtH := ctx.MeasureString(block.BigText)
	if bTxtH > 0 {
		totalTextHeight += bTxtH
	}
	bX := ((float64(block.Width) - bTxtW) / 2.0)

	// Draw text
	var drawTextMargins float64
	if block.AltText != "" {
		drawTextMargins = (float64(availHeiht) - totalTextHeight) / 4
	} else {
		drawTextMargins = (float64(availHeiht) - totalTextHeight) / 3
	}
	// Big text
	ctx.SetColor(block.BigTextColor)
	if err := ctx.LoadFontFace(render.FontPath, block.TextSize); err != nil {
		return block, err
	}
	bY := bTxtH + drawTextMargins
	ctx.DrawString(block.BigText, bX, bY)

	// Small text
	ctx.SetColor(block.SmallTextColor)
	if err := ctx.LoadFontFace(render.FontPath, (block.TextSize * block.TextCoeff)); err != nil {
		return block, err
	}
	sY := bY + sTxtH + drawTextMargins
	ctx.DrawString(block.SmallText, sX, sY)

	if block.AltText != "" {
		if err := ctx.LoadFontFace(render.FontPath, (block.TextSize * (block.TextCoeff - 0.15))); err != nil {
			return block, err
		}
		ctx.SetColor(block.AltTextColor)
		aX := ((float64(block.Width) - aTxtW) / 2.0)
		aY := sY + drawTextMargins + aTxtH
		ctx.DrawString(block.AltText, aX, aY)
	}

	// Draw icons
	if block.HasBigIcon == true {
		ctx.SetColor(block.BigIconColor)
		if block.BigArrowDirection == 0 {
			iR := 8.0 * (block.TextSize / render.FontSize)
			iX := bX - (iR * 1.5)
			iY := bY - iR - ((bTxtH - (iR * 2)) / 2)
			ctx.DrawCircle(iX, iY, iR)
			ctx.Fill()
		}
		if block.BigArrowDirection == 1 {
			ctx.SetColor(color.RGBA{0, 255, 0, 180})
			iR := 8.0 * (block.TextSize / render.FontSize)
			iX := bX - (iR * 1.5)
			iY := bY - ((bTxtH - (iR)) / 2) - (render.FontSize / 10)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 0)
			ctx.Fill()
		}
		if block.BigArrowDirection == -1 {
			ctx.SetColor(color.RGBA{255, 0, 0, 180})
			iR := 8.0 * (block.TextSize / render.FontSize)
			iX := bX - (iR * 1.5)
			iY := bY - bTxtH + ((bTxtH - (iR)) / 2) + (render.FontSize / 10)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 1)
			ctx.Fill()
		}
	}
	if block.HasSmallIcon == true {
		ctx.SetColor(block.SmallIconColor)
		if block.SmallArrowDirection == 0 {
			iR := 8.0 * 0.75 * (block.TextSize / render.FontSize)
			iX := sX - (iR * 1.5)
			iY := sY - iR - ((sTxtH - (iR * 2)) / 2)
			ctx.DrawCircle(iX, iY, iR)
			ctx.Fill()
		}
		if block.SmallArrowDirection == 1 {
			iR := 8.0 * 0.75 * (block.TextSize / render.FontSize)
			iX := sX - (iR * 1.5)
			iY := sY - iR - ((sTxtH - (iR * 2)) / 2)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 0)
			ctx.Fill()
		}
		if block.SmallArrowDirection == -1 {
			iR := 8.0 * 0.75 * (block.TextSize / render.FontSize)
			iX := sX - (iR * 1.5)
			iY := sY - iR - ((sTxtH - (iR * 2)) / 2)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 1)
			ctx.Fill()
		}
	}
	block.Context = ctx
	return block, nil
}

func makeSessionStatsCard(card render.CardData, header string, currentSession, compareSession wgapi.StatsFrame, currentRating, compareRating int) (render.CardData, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()

	ctx := *card.Context
	if err := ctx.LoadFontFace(render.FontPath, (render.FontSize * 1.25)); err != nil {
		return card, err
	}

	if currentSession.Battles < 1 {
		return card, fmt.Errorf("sessions battles is < 1")
	}

	ctx.SetColor(color.White)
	// Measure header
	headerW, headerH := ctx.MeasureString(header)
	if err := ctx.LoadFontFace(render.FontPath, (render.FontSize * 0.75)); err != nil {
		return card, err
	}
	headerHeigth := int(headerH * 2)
	nameX := (float64(card.Context.Width()) - headerW) / 2
	nameY := (float64(headerHeigth)-headerH)/2 + headerH

	// Draw header
	if err := ctx.LoadFontFace(render.FontPath, (render.FontSize * 1.25)); err != nil {
		return card, err
	}
	ctx.DrawString(header, nameX, nameY)

	blocksCnt := 4
	blockWidth := card.Context.Width() / blocksCnt
	availableHeight := int(ctx.Height() - (headerHeigth))
	// Blocks will take 75% of the total card height
	blockHeight := availableHeight
	// Default Block
	var defaultBlock statsBlock
	defaultBlock.TextSize = render.FontSize * 1.65
	defaultBlock.TextCoeff = 0.75
	defaultBlock.Width = blockWidth
	defaultBlock.Height = blockHeight
	defaultBlock.BigTextColor = render.BigTextColor
	defaultBlock.SmallTextColor = render.SmallTextColor
	defaultBlock.AltTextColor = render.AltTextColor

	// Bottom Row - Avg Damage, Avg XP, Winrate
	// Block 1 - Battles
	battlesBlock := statsBlock(defaultBlock)
	battlesBlock.Width = blockWidth
	battlesSession := strconv.Itoa(currentSession.Battles)
	battlesLastSession := "-"
	if compareSession.Battles > 0 {
		battlesLastSession = strconv.Itoa(compareSession.Battles)
	}
	battlesBlock.SmallText = battlesLastSession
	battlesBlock.BigText = battlesSession
	battlesBlock.AltText = "Battles"
	battlesBlock, err := addStatsBlockCtx(battlesBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(battlesBlock.Context.Image(), 0, headerHeigth)
	// Block 2 - Avg Damage
	avgDamageBlock := statsBlock(defaultBlock)
	avgDamageBlock.Width = blockWidth
	avgDamageSession := strconv.Itoa((currentSession.DamageDealt / currentSession.Battles))
	avgDamageLastSession := "-"
	if compareSession.Battles > 0 {
		avgDamageLastSession = strconv.Itoa((compareSession.DamageDealt / compareSession.Battles))
	}
	avgDamageBlock.SmallText = avgDamageLastSession
	avgDamageBlock.BigText = avgDamageSession
	avgDamageBlock.AltText = "Avg. Damage"
	avgDamageBlock, err = addStatsBlockCtx(avgDamageBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(avgDamageBlock.Context.Image(), (blockWidth), headerHeigth)
	// Block 1 - Winrate
	winrateBlock := statsBlock(avgDamageBlock)
	winrateSession := ((float64(currentSession.Wins) / float64(currentSession.Battles)) * 100)
	winrateLastSession := "-"
	if compareSession.Battles > 0 {
		winrateLastSession = fmt.Sprintf("%.2f", ((float64(compareSession.Wins)/float64(compareSession.Battles))*100)) + "%"
	}
	winrateBlock.BigText = fmt.Sprintf("%.2f", winrateSession) + "%"
	winrateBlock.SmallText = winrateLastSession
	winrateBlock.AltText = "Winrate"
	winrateBlock, err = addStatsBlockCtx(winrateBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(winrateBlock.Context.Image(), (blockWidth * 2), headerHeigth)
	// Block 4 - Draw WN8
	if currentRating > -1 || compareRating > -1 {
		ratingBlock := statsBlock(defaultBlock)
		// Icon
		ratingBlock.BigText = "-"
		if currentRating > -1 {
			ratingBlock.HasBigIcon = true
			ratingBlock.BigIconColor = GetRatingColor(currentRating)
			ratingBlock.BigText = strconv.Itoa(currentRating)
		}
		ratingBlock.SmallText = "-"
		if compareRating > -1 {
			ratingBlock.HasSmallIcon = true
			ratingBlock.SmallIconColor = GetRatingColor(compareRating)
			ratingBlock.SmallText = strconv.Itoa(compareRating)
		}
		ratingBlock.AltText = "WN8"
		ratingBlock, err = addStatsBlockCtx(ratingBlock)
		if err != nil {
			return card, err
		}
		ctx.DrawImage(ratingBlock.Context.Image(), (blockWidth * 3), headerHeigth)
	} else {
		// Accuracy Block to replace WN8
		avgAccuracyBlock := statsBlock(defaultBlock)
		avgAccuracyBlock.Width = blockWidth
		avgAccuracyBlock.BigText = "-"
		if currentSession.Shots > 0 {
			avgAccuracyBlock.BigText = fmt.Sprintf("%.2f", ((float64(currentSession.Hits)/float64(currentSession.Shots))*100)) + "%"
		}
		avgAccuracyBlock.SmallText = "-"
		if compareSession.Shots > 0 {
			avgAccuracyBlock.SmallText = fmt.Sprintf("%.2f", ((float64(compareSession.Hits)/float64(compareSession.Shots))*100)) + "%"
		}
		avgAccuracyBlock.AltText = "Accuracy"
		avgAccuracyBlock, err = addStatsBlockCtx(avgAccuracyBlock)
		if err != nil {
			return card, err
		}
		ctx.DrawImage(avgAccuracyBlock.Context.Image(), (blockWidth * 3), headerHeigth)
	}
	// Draw lines
	ctx.SetColor(render.DecorLinesColor)
	lineX := float64(render.FrameMargin)
	lineY := float64(headerHeigth)
	lineHeight := 2.0
	lineWidth := (float64(ctx.Width()) - float64(render.FrameMargin*2))
	ctx.DrawRectangle(lineX, lineY, lineWidth, lineHeight)
	ctx.Fill()

	// Render image
	card.Image = ctx.Image()
	return card, nil
}

func makeStatsHeaderCard(card render.CardData, playerName, playerClan string, premium, verified bool) (render.CardData, error) {
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
	yOffset := ((float64(ctx.Height()) - nameStrH) / 2)
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
	// Draw server icon

	// Make image
	card.Image = ctx.Image()
	return card, nil
}

// GetRatingColor - Rating color calculator
func GetRatingColor(r int) color.RGBA {
	if r > 0 && r < 301 {
		return color.RGBA{255, 0, 0, 180}
	}
	if r > 300 && r < 451 {
		return color.RGBA{251, 83, 83, 180}
	}
	if r > 450 && r < 651 {
		return color.RGBA{255, 160, 49, 180}
	}
	if r > 650 && r < 901 {
		return color.RGBA{255, 244, 65, 180}
	}
	if r > 900 && r < 1201 {
		return color.RGBA{149, 245, 62, 180}
	}
	if r > 1200 && r < 1601 {
		return color.RGBA{103, 190, 51, 180}
	}
	if r > 1600 && r < 2001 {
		return color.RGBA{106, 236, 255, 180}
	}
	if r > 2000 && r < 2451 {
		return color.RGBA{46, 174, 193, 180}
	}
	if r > 2450 && r < 2901 {
		return color.RGBA{208, 108, 255, 180}
	}
	if r > 2900 {
		return color.RGBA{142, 65, 177, 180}
	}
	return color.RGBA{0, 0, 0, 0}
}

// Tank tier to roman numeral
func tierToRoman(t int) string {
	switch t {
	case 1:
		return "I"
	case 2:
		return "II"
	case 3:
		return "III"
	case 4:
		return "IV"
	case 5:
		return "V"
	case 6:
		return "VI"
	case 7:
		return "VII"
	case 8:
		return "VIII"
	case 9:
		return "IX"
	case 10:
		return "X"
	default:
		return ""
	}
}
