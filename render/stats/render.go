package render

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"sort"
	"strconv"
	"sync"

	stats "github.com/cufee/am-stats/dataprep/stats"
	"github.com/cufee/am-stats/render"
	wgapi "github.com/cufee/am-stats/wargamingapi"
	"github.com/cufee/am-stats/winstreak"
	"github.com/fogleman/gg"
)

// ImageFromStats -
func ImageFromStats(data stats.ExportData, sortKey string, tankLimit int, premium bool, verified bool, bgImage image.Image) (finalImage image.Image, err error) {
	var finalCards render.AllCards
	cardsChan := make(chan render.CardData, (2 + len(data.SessionStats.Vehicles)))
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
		headerHeight := 1.0
		header, err := makeStatsHeaderCard(render.PrepNewCard(0, headerHeight), data.PlayerDetails.Name, clanTag, "Random Battles", premium, verified)
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
		allStats, err := makeAllStatsCard(render.PrepNewCard(1, 1.5), data)
		if err != nil {
			log.Println(err)
			return
		}
		cardsChan <- allStats
	}()

	// Sort vehicles
	vehicles := sortTanks(data.SessionStats.Vehicles, sortKey)
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
				tankCard, err = makeDetailedCard(render.PrepNewCard((i+2), 1.0), tank, lastSession)
			} else {
				tankCard, err = makeSlimCard(render.PrepNewCard((i+2), 0.5), tank, lastSession)
			}
			if err != nil {
				log.Println(err)
				return
			}
			cardsChan <- tankCard
		}(tank, i)
	}

	wg.Wait()
	close(cardsChan)

	for c := range cardsChan {
		finalCards.Cards = append(finalCards.Cards, c)
	}
	finalCtx, err := render.AddAllCardsToFrame(finalCards, data.SessionStats.Timestamp, "Session from", bgImage)
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
		ratingBlock.BigIconColor = getRatingColor(session.TankWN8)
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
		ratingBlock.BigIconColor = getRatingColor(session.TankWN8)
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

func makeAllStatsCard(card render.CardData, data stats.ExportData) (render.CardData, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()

	ctx := *card.Context
	if err := ctx.LoadFontFace(render.FontPath, render.FontSize); err != nil {
		return card, err
	}
	ctx.SetColor(color.White)

	// Default Block settings
	blockWidth := card.Context.Width() / 3
	bottomBlockWidth := card.Context.Width() / 4
	availableHeight := (ctx.Height()) / 2
	blockHeight := availableHeight
	var defaultBlock statsBlock
	defaultBlock.TextSize = render.FontSize * 1.5
	defaultBlock.Width = blockWidth
	defaultBlock.Height = blockHeight
	defaultBlock.BigTextColor = render.BigTextColor
	defaultBlock.SmallTextColor = render.SmallTextColor
	defaultBlock.AltTextColor = render.AltTextColor
	// Top Row - 3 Blocks (Battles, WN8, WR)
	badSession := true
	if data.SessionStats.StatsAll.Battles > 0 {
		badSession = false
	}
	if data.PlayerDetails.Stats.All.Battles < 1 {
		return card, fmt.Errorf("this player has no battles")
	}
	// Block 1 - Battles
	battlesBlock := statsBlock(defaultBlock)
	battlesBlock.SmallText = strconv.Itoa(data.PlayerDetails.Stats.All.Battles)
	battlesBlock.BigText = strconv.Itoa(data.SessionStats.BattlesAll)
	battlesBlock.AltText = "Battles"
	battlesBlock, err := addStatsBlockCtx(battlesBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(battlesBlock.Context.Image(), 0, 0)

	// Block 2 - WN8
	ratingBlock := statsBlock(defaultBlock)
	// Icon
	ratingBlock.HasBigIcon = true
	ratingBlock.BigIconColor = getRatingColor(data.SessionStats.SessionRating)
	ratingBlock.HasSmallIcon = true
	ratingBlock.SmallIconColor = getRatingColor(data.PlayerDetails.CareerWN8)
	ratingBlock.Height = blockHeight + int(render.FontSize)
	ratingBlock.TextSize = render.FontSize * 1.75
	careerWN8str := "-"
	if data.PlayerDetails.CareerWN8 > 0 {
		careerWN8str = strconv.Itoa(data.PlayerDetails.CareerWN8)
	}
	ratingBlock.SmallText = careerWN8str
	ratingBlock.BigText = "-"
	if !badSession && data.SessionStats.SessionRating > -1 {
		ratingBlock.BigText = strconv.Itoa(data.SessionStats.SessionRating)
	}
	ratingBlock.AltText = "WN8"
	ratingBlock, err = addStatsBlockCtx(ratingBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(ratingBlock.Context.Image(), blockWidth, 0)
	// Block 3 - WR
	winrateBlock := statsBlock(battlesBlock)
	oldBattles := data.PlayerDetails.Stats.All.Battles - data.SessionStats.StatsAll.Battles
	oldWins := data.PlayerDetails.Stats.All.Wins - data.SessionStats.StatsAll.Wins
	oldWinrate := (float64(oldWins) / float64(oldBattles)) * 100
	winrateAll := ((float64(data.PlayerDetails.Stats.All.Wins) / float64(data.PlayerDetails.Stats.All.Battles)) * 100)
	winrateSession := 0.0
	winrateChange := 0.0
	if !badSession {
		winrateSession = ((float64(data.SessionStats.StatsAll.Wins) / float64(data.SessionStats.StatsAll.Battles)) * 100)
		winrateChange = math.Round((winrateAll-oldWinrate)*100) / 100
	}
	winrateChangeStr := ""
	if winrateChange > 0.00 {
		winrateChangeStr = fmt.Sprintf(" (+%.2f", winrateChange) + "%)"
	}
	if winrateChange < 0.00 {
		winrateChangeStr = fmt.Sprintf(" (%.2f", winrateChange) + "%)"
	}
	winrateAllStr := fmt.Sprintf("%.2f", winrateAll) + "%" + winrateChangeStr
	winrateSessionStr := "-"
	if !badSession {
		winrateSessionStr = fmt.Sprintf("%.2f", winrateSession) + "%"
	}
	winrateBlock.SmallText = winrateAllStr
	winrateBlock.BigText = winrateSessionStr
	winrateBlock.AltText = "Winrate"
	winrateBlock, err = addStatsBlockCtx(winrateBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(winrateBlock.Context.Image(), (blockWidth * 2), 0)

	// Bottom Row - 4 Blocks
	// Block 1 - Avg Damage
	avgDamageBlock := statsBlock(defaultBlock)
	avgDamageBlock.Width = bottomBlockWidth
	avgDamageAll := "-"
	if data.PlayerDetails.Stats.All.Battles > 0 {
		avgDamageAll = strconv.Itoa((data.PlayerDetails.Stats.All.DamageDealt / data.PlayerDetails.Stats.All.Battles))
	}
	avgDamageSession := "-"
	if !badSession {
		avgDamageSession = strconv.Itoa((data.SessionStats.StatsAll.DamageDealt / data.SessionStats.StatsAll.Battles))
	}
	avgDamageBlock.SmallText = avgDamageAll
	avgDamageBlock.BigText = avgDamageSession
	avgDamageBlock.AltText = "Avg. Damage"
	avgDamageBlock, err = addStatsBlockCtx(avgDamageBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(avgDamageBlock.Context.Image(), 0, blockHeight)
	// Block 2 - Damage Ratio
	damageRatioBlock := statsBlock(avgDamageBlock)
	damageRatioAll := "-"
	if data.PlayerDetails.Stats.All.DamageReceived > 0 {
		damageRatioAll = fmt.Sprintf("%.2f", (float64(data.PlayerDetails.Stats.All.DamageDealt) / float64(data.PlayerDetails.Stats.All.DamageReceived)))
	}
	damageRatioSession := "-"
	if !badSession && data.SessionStats.StatsAll.DamageReceived > 1 {
		damageRatioSession = fmt.Sprintf("%.2f", (float64(data.SessionStats.StatsAll.DamageDealt) / float64(data.SessionStats.StatsAll.DamageReceived)))
	}
	damageRatioBlock.SmallText = damageRatioAll
	damageRatioBlock.BigText = damageRatioSession
	damageRatioBlock.AltText = "Damage Ratio"
	damageRatioBlock, err = addStatsBlockCtx(damageRatioBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(damageRatioBlock.Context.Image(), bottomBlockWidth, blockHeight)
	// Block 3 - Destruction Ratio
	destrRatioBlock := statsBlock(avgDamageBlock)
	destrRatioAll := "-"
	if data.PlayerDetails.Stats.All.SurvivedBattles > 0 && data.PlayerDetails.Stats.All.Battles != data.PlayerDetails.Stats.All.SurvivedBattles {
		destrRatioAll = fmt.Sprintf("%.2f", (float64(data.PlayerDetails.Stats.All.Frags) / (float64(data.PlayerDetails.Stats.All.Battles) - float64(data.PlayerDetails.Stats.All.SurvivedBattles))))
	}
	destrRatioSession := "-"
	if !badSession && data.SessionStats.StatsAll.SurvivedBattles > 0 && data.SessionStats.StatsAll.Battles != data.SessionStats.StatsAll.SurvivedBattles {
		destrRatioSession = fmt.Sprintf("%.2f", (float64(data.SessionStats.StatsAll.Frags) / (float64(data.SessionStats.StatsAll.Battles) - float64(data.SessionStats.StatsAll.SurvivedBattles))))
	}
	destrRatioBlock.SmallText = destrRatioAll
	destrRatioBlock.BigText = destrRatioSession
	destrRatioBlock.AltText = "Destruction Ratio"
	destrRatioBlock, err = addStatsBlockCtx(destrRatioBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(destrRatioBlock.Context.Image(), (bottomBlockWidth * 2), blockHeight)
	// Block 4 - Average XP or Win streak
	winStreak, err := winstreak.CheckStreak(data.PlayerDetails.ID, data.PlayerDetails.Stats.All)
	if err != nil {
		log.Print("failed to get a win streak:", err)
	}
	streakBlock := statsBlock(avgDamageBlock)
	streakBlock.BigText = strconv.Itoa(winStreak.Streak)
	streakBlock.SmallText = "-"
	if winStreak.BestStreak > 0 {
		streakBlock.SmallText = strconv.Itoa(winStreak.BestStreak)
	}
	streakBlock.AltText = "Win Streak"
	streakBlock, err = addStatsBlockCtx(streakBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(streakBlock.Context.Image(), (bottomBlockWidth * 3), blockHeight)
	// Draw lines
	ctx.SetColor(render.DecorLinesColor)
	lineX := float64(render.FrameMargin)
	lineY := float64(blockHeight)
	lineHeight := 2.0
	lineWidth := (float64(ctx.Width()) - float64(render.FrameMargin*2) - 80.0) / 2
	ctx.DrawRectangle(lineX, lineY, lineWidth, lineHeight)
	ctx.DrawRectangle((lineX + lineWidth + 80), lineY, lineWidth, lineHeight)
	ctx.Fill()

	// Render image
	card.Image = ctx.Image()
	return card, nil
}

func makeStatsHeaderCard(card render.CardData, playerName, playerClan, battleType string, premium bool, verified bool) (render.CardData, error) {
	ctx := *card.Context
	if err := ctx.LoadFontFace(render.FontPath, render.FontSizeHeader); err != nil {
		return card, err
	}
	playerStr := playerName + " "
	// Calculate text size
	nameStrW, nameStrH := ctx.MeasureString(playerStr)
	clanStrW, _ := ctx.MeasureString(playerClan)
	battleTypeW, battleTypeH := ctx.MeasureString(battleType)
	totalTextW := nameStrW
	if nameStrW < battleTypeW {
		totalTextW = battleTypeW
	}
	totalTextH := nameStrH + render.TextMargin + battleTypeH

	xOffset := ((float64(ctx.Width()) - totalTextW) / 2)
	yOffset := ((float64(ctx.Height()) - totalTextH) / 2)
	// Draw battle type text
	ctx.SetColor(color.RGBA{255, 255, 255, 200})
	btDrawX := ((totalTextW - battleTypeW) / 2) + xOffset
	btDrawY := yOffset + battleTypeH
	ctx.DrawString(battleType, btDrawX, btDrawY)
	// Draw player name and tag text
	ctx.SetColor(color.White)
	if premium {
		ctx.SetColor(render.PremiumColor)
	}
	psDrawX := ((totalTextW - nameStrW - clanStrW) / 2) + xOffset
	psDrawY := totalTextH + yOffset
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
	// Make image
	card.Image = ctx.Image()
	return card, nil
}

// Rating color calculator
func getRatingColor(r int) color.RGBA {
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

// Sorting of vehicles
func sortTanks(vehicles []wgapi.VehicleStats, sortKey string) []wgapi.VehicleStats {
	// Sort based on passed key
	switch sortKey {
	case "+battles":
		sort.Slice(vehicles, func(i, j int) bool {
			return vehicles[i].Battles < vehicles[j].Battles
		})
	case "-battles":
		sort.Slice(vehicles, func(i, j int) bool {
			return vehicles[i].Battles > vehicles[j].Battles
		})
	case "+winrate":
		sort.Slice(vehicles, func(i, j int) bool {
			return (float64(vehicles[i].Wins) / float64(vehicles[i].Battles)) < (float64(vehicles[j].Wins) / float64(vehicles[j].Battles))
		})
	case "-winrate":
		sort.Slice(vehicles, func(i, j int) bool {
			return (float64(vehicles[i].Wins) / float64(vehicles[i].Battles)) > (float64(vehicles[j].Wins) / float64(vehicles[j].Battles))
		})
	case "+wn8":
		sort.Slice(vehicles, func(i, j int) bool {
			return absInt(vehicles[i].TankWN8) < absInt(vehicles[j].TankWN8)
		})
	case "-wn8":
		sort.Slice(vehicles, func(i, j int) bool {
			return absInt(vehicles[i].TankWN8) > absInt(vehicles[j].TankWN8)
		})
	case "+last_battle":
		sort.Slice(vehicles, func(i, j int) bool {
			return absInt(vehicles[i].LastBattleTime) < absInt(vehicles[j].LastBattleTime)
		})
	case "-last_battle":
		sort.Slice(vehicles, func(i, j int) bool {
			return absInt(vehicles[i].LastBattleTime) > absInt(vehicles[j].LastBattleTime)
		})
	case "relevance":
		sort.Slice(vehicles, func(i, j int) bool {
			return (absInt(vehicles[i].TankRawWN8) * vehicles[i].LastBattleTime * vehicles[i].Battles) > (absInt(vehicles[j].TankRawWN8) * vehicles[j].LastBattleTime * vehicles[j].Battles)
		})
	default:
		sort.Slice(vehicles, func(i, j int) bool {
			return vehicles[i].LastBattleTime > vehicles[j].LastBattleTime
		})
	}
	return vehicles
}

// absInt - Absolute value of an integer
func absInt(val int) int {
	if val >= 0 {
		return val
	}
	return -val
}
