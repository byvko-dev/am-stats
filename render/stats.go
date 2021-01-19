package render

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/cufee/am-stats/stats"
	wgapi "github.com/cufee/am-stats/wargamingapi"
	"github.com/cufee/am-stats/winstreak"
	"github.com/fogleman/gg"
)

func makeAllStatsCard(card cardData, data stats.ExportData) (cardData, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()

	ctx := *card.context
	if err := ctx.LoadFontFace(fontPath, fontSize); err != nil {
		return card, err
	}
	ctx.SetColor(color.White)

	// Default Block settings
	blockWidth := card.context.Width() / 3
	bottomBlockWidth := card.context.Width() / 4
	availableHeight := (ctx.Height()) / 2
	blockHeight := availableHeight
	var defaultBlock cardBlock
	defaultBlock.textSize = fontSize * 1.5
	defaultBlock.width = blockWidth
	defaultBlock.height = blockHeight
	defaultBlock.bigTextColor = bigTextColor
	defaultBlock.smallTextColor = smallTextColor
	defaultBlock.altTextColor = altTextColor
	// Top Row - 3 Blocks (Battles, WN8, WR)
	badSession := true
	if data.SessionStats.StatsAll.Battles > 0 {
		badSession = false
	}
	if data.PlayerDetails.Stats.All.Battles < 1 {
		return card, fmt.Errorf("this player has no battles")
	}
	// Block 1 - Battles
	battlesBlock := cardBlock(defaultBlock)
	battlesBlock.smallText = strconv.Itoa(data.PlayerDetails.Stats.All.Battles)
	battlesBlock.bigText = strconv.Itoa(data.SessionStats.BattlesAll)
	battlesBlock.altText = "Battles"
	battlesBlock, err := addBlockCtx(battlesBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(battlesBlock.context.Image(), 0, 0)

	// Block 2 - WN8
	ratingBlock := cardBlock(defaultBlock)
	// Icon
	ratingBlock.hasBigIcon = true
	ratingBlock.bigIconColor = getRatingColor(data.SessionStats.SessionRating)
	ratingBlock.hasSmallIcon = true
	ratingBlock.smallIconColor = getRatingColor(data.PlayerDetails.CareerWN8)
	ratingBlock.height = blockHeight + int(fontSize)
	ratingBlock.textSize = fontSize * 1.75
	careerWN8str := "-"
	if data.PlayerDetails.CareerWN8 > 0 {
		careerWN8str = strconv.Itoa(data.PlayerDetails.CareerWN8)
	}
	ratingBlock.smallText = careerWN8str
	ratingBlock.bigText = "-"
	if !badSession && data.SessionStats.SessionRating > -1 {
		ratingBlock.bigText = strconv.Itoa(data.SessionStats.SessionRating)
	}
	ratingBlock.altText = "WN8"
	ratingBlock, err = addBlockCtx(ratingBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(ratingBlock.context.Image(), blockWidth, 0)
	// Block 3 - WR
	winrateBlock := cardBlock(battlesBlock)
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
	winrateBlock.smallText = winrateAllStr
	winrateBlock.bigText = winrateSessionStr
	winrateBlock.altText = "Winrate"
	winrateBlock, err = addBlockCtx(winrateBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(winrateBlock.context.Image(), (blockWidth * 2), 0)

	// Bottom Row - 4 Blocks
	// Block 1 - Avg Damage
	avgDamageBlock := cardBlock(defaultBlock)
	avgDamageBlock.width = bottomBlockWidth
	avgDamageAll := "-"
	if data.PlayerDetails.Stats.All.Battles > 0 {
		avgDamageAll = strconv.Itoa((data.PlayerDetails.Stats.All.DamageDealt / data.PlayerDetails.Stats.All.Battles))
	}
	avgDamageSession := "-"
	if !badSession {
		avgDamageSession = strconv.Itoa((data.SessionStats.StatsAll.DamageDealt / data.SessionStats.StatsAll.Battles))
	}
	avgDamageBlock.smallText = avgDamageAll
	avgDamageBlock.bigText = avgDamageSession
	avgDamageBlock.altText = "Avg. Damage"
	avgDamageBlock, err = addBlockCtx(avgDamageBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(avgDamageBlock.context.Image(), 0, blockHeight)
	// Block 2 - Damage Ratio
	damageRatioBlock := cardBlock(avgDamageBlock)
	damageRatioAll := "-"
	if data.PlayerDetails.Stats.All.DamageReceived > 0 {
		damageRatioAll = fmt.Sprintf("%.2f", (float64(data.PlayerDetails.Stats.All.DamageDealt) / float64(data.PlayerDetails.Stats.All.DamageReceived)))
	}
	damageRatioSession := "-"
	if !badSession && data.SessionStats.StatsAll.DamageReceived > 1 {
		damageRatioSession = fmt.Sprintf("%.2f", (float64(data.SessionStats.StatsAll.DamageDealt) / float64(data.SessionStats.StatsAll.DamageReceived)))
	}
	damageRatioBlock.smallText = damageRatioAll
	damageRatioBlock.bigText = damageRatioSession
	damageRatioBlock.altText = "Damage Ratio"
	damageRatioBlock, err = addBlockCtx(damageRatioBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(damageRatioBlock.context.Image(), bottomBlockWidth, blockHeight)
	// Block 3 - Destruction Ratio
	destrRatioBlock := cardBlock(avgDamageBlock)
	destrRatioAll := "-"
	if data.PlayerDetails.Stats.All.SurvivedBattles > 0 && data.PlayerDetails.Stats.All.Battles != data.PlayerDetails.Stats.All.SurvivedBattles {
		destrRatioAll = fmt.Sprintf("%.2f", (float64(data.PlayerDetails.Stats.All.Frags) / (float64(data.PlayerDetails.Stats.All.Battles) - float64(data.PlayerDetails.Stats.All.SurvivedBattles))))
	}
	destrRatioSession := "-"
	if !badSession && data.SessionStats.StatsAll.SurvivedBattles > 0 && data.SessionStats.StatsAll.Battles != data.SessionStats.StatsAll.SurvivedBattles {
		destrRatioSession = fmt.Sprintf("%.2f", (float64(data.SessionStats.StatsAll.Frags) / (float64(data.SessionStats.StatsAll.Battles) - float64(data.SessionStats.StatsAll.SurvivedBattles))))
	}
	destrRatioBlock.smallText = destrRatioAll
	destrRatioBlock.bigText = destrRatioSession
	destrRatioBlock.altText = "Destruction Ratio"
	destrRatioBlock, err = addBlockCtx(destrRatioBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(destrRatioBlock.context.Image(), (bottomBlockWidth * 2), blockHeight)
	// Block 4 - Average XP or Win streak
	winStreak, err := winstreak.CheckStreak(data.PlayerDetails.ID, data.PlayerDetails.Stats.All)
	if err != nil {
		log.Print("failed to get a win streak:", err)
	}
	streakBlock := cardBlock(avgDamageBlock)
	streakBlock.bigText = strconv.Itoa(winStreak.Streak)
	streakBlock.smallText = "-"
	if winStreak.BestStreak > 0 {
		streakBlock.smallText = strconv.Itoa(winStreak.BestStreak)
	}
	streakBlock.altText = "Win Streak"
	streakBlock, err = addBlockCtx(streakBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(streakBlock.context.Image(), (bottomBlockWidth * 3), blockHeight)
	// Draw lines
	ctx.SetColor(decorLinesColor)
	lineX := float64(frameMargin)
	lineY := float64(blockHeight)
	lineHeight := 2.0
	lineWidth := (float64(ctx.Width()) - float64(frameMargin*2) - 80.0) / 2
	ctx.DrawRectangle(lineX, lineY, lineWidth, lineHeight)
	ctx.DrawRectangle((lineX + lineWidth + 80), lineY, lineWidth, lineHeight)
	ctx.Fill()

	// Render image
	card.image = ctx.Image()
	return card, nil
}

// Makes a detailed card for a tank
func makeDetailedCard(card cardData, session wgapi.VehicleStats, lastSession wgapi.VehicleStats) (cardData, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()

	ctx := *card.context
	if err := ctx.LoadFontFace(fontPath, (fontSize * 1.25)); err != nil {
		return card, err
	}

	if session.Battles < 1 {
		return card, fmt.Errorf("sessions battles is < 1")
	}

	ctx.SetColor(color.White)
	// Measure tank name
	nameW, nameH := ctx.MeasureString(session.TankName)
	if err := ctx.LoadFontFace(fontPath, (fontSize * 0.75)); err != nil {
		return card, err
	}
	tierW, tierH := ctx.MeasureString(tierToRoman(session.TankTier))
	headerHeigth := int(nameH * 2)
	nameX := (float64(card.context.Width()) - nameW) / 2
	nameY := (float64(headerHeigth)-nameH)/2 + nameH
	tierX := nameX - (fontSize / 2) - tierW
	tierY := (float64(headerHeigth)-tierH)/2 + tierH

	// Draw tank tier
	ctx.DrawString(tierToRoman(session.TankTier), tierX, tierY)
	// Draw tank name
	if err := ctx.LoadFontFace(fontPath, (fontSize * 1.25)); err != nil {
		return card, err
	}
	ctx.DrawString(session.TankName, nameX, nameY)

	blockWidth := card.context.Width() / 4
	availableHeight := int(ctx.Height() - (headerHeigth))
	// Blocks will take 75% of the total card heiht
	blockHeight := availableHeight
	// Default Block
	var defaultBlock cardBlock
	defaultBlock.textSize = fontSize * 1.30
	defaultBlock.textCoeff = 0.75
	defaultBlock.width = blockWidth
	defaultBlock.height = blockHeight
	defaultBlock.bigTextColor = bigTextColor
	defaultBlock.smallTextColor = smallTextColor
	defaultBlock.altTextColor = altTextColor

	// Bottom Row - Avg Damage, Avg XP, Winrate
	// Block 1 - Battles
	battlesBlock := cardBlock(defaultBlock)
	battlesBlock.width = blockWidth
	battlesSession := strconv.Itoa(session.Battles)
	battlesLastSession := "-"
	if lastSession.Battles > 0 {
		battlesLastSession = strconv.Itoa(lastSession.Battles)
	}
	battlesBlock.smallText = battlesLastSession
	battlesBlock.bigText = battlesSession
	battlesBlock.altText = "Battles"
	battlesBlock, err := addBlockCtx(battlesBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(battlesBlock.context.Image(), 0, headerHeigth)
	// Block 2 - Avg Damage
	avgDamageBlock := cardBlock(defaultBlock)
	avgDamageBlock.width = blockWidth
	avgDamageSession := strconv.Itoa((session.DamageDealt / session.Battles))
	avgDamageLastSession := "-"
	avgDamageBlock.hasBigIcon = true
	if lastSession.Battles > 0 {
		avgDamageLastSession = strconv.Itoa((lastSession.DamageDealt / lastSession.Battles))
		if (lastSession.DamageDealt / lastSession.Battles) < (session.DamageDealt / session.Battles) {
			avgDamageBlock.bigArrowDirection = 1
		}
		if (lastSession.DamageDealt / lastSession.Battles) > (session.DamageDealt / session.Battles) {
			avgDamageBlock.bigArrowDirection = -1
		}
	}
	avgDamageBlock.smallText = avgDamageLastSession
	avgDamageBlock.bigText = avgDamageSession
	avgDamageBlock.altText = "Avg. Damage"
	avgDamageBlock, err = addBlockCtx(avgDamageBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(avgDamageBlock.context.Image(), (blockWidth), headerHeigth)
	// Block 1 - Winrate
	winrateBlock := cardBlock(avgDamageBlock)
	winrateSession := ((float64(session.Wins) / float64(session.Battles)) * 100)
	winrateLastSession := "-"
	if lastSession.Battles > 0 {
		winrateLastSession = fmt.Sprintf("%.2f", ((float64(lastSession.Wins)/float64(lastSession.Battles))*100)) + "%"
	}
	winrateBlock.bigText = fmt.Sprintf("%.2f", winrateSession) + "%"
	winrateBlock.smallText = winrateLastSession
	winrateBlock.altText = "Winrate"
	winrateBlock.hasBigIcon = true
	if ((float64(lastSession.Wins) / float64(lastSession.Battles)) * 100) < winrateSession {
		winrateBlock.bigArrowDirection = 1
	}
	if ((float64(lastSession.Wins) / float64(lastSession.Battles)) * 100) > winrateSession {
		winrateBlock.bigArrowDirection = -1
	}
	winrateBlock, err = addBlockCtx(winrateBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(winrateBlock.context.Image(), (blockWidth * 2), headerHeigth)
	// Block 4 - Draw WN8
	ratingBlock := cardBlock(defaultBlock)
	// Icon
	ratingBlock.smallText = "WN8"
	ratingBlock.bigText = "-"
	if session.TankWN8 > -1 {
		ratingBlock.hasBigIcon = true
		ratingBlock.bigIconColor = getRatingColor(session.TankWN8)
		ratingBlock.bigText = strconv.Itoa(session.TankWN8)
	}
	ratingBlock.smallTextColor = altTextColor
	ratingBlock, err = addBlockCtx(ratingBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(ratingBlock.context.Image(), (blockWidth * 3), headerHeigth)

	// Draw lines
	ctx.SetColor(decorLinesColor)
	lineX := float64(frameMargin)
	lineY := float64(headerHeigth)
	lineHeight := 2.0
	lineWidth := (float64(ctx.Width()) - float64(frameMargin*2))
	ctx.DrawRectangle(lineX, lineY, lineWidth, lineHeight)
	ctx.Fill()

	// Render image
	card.image = ctx.Image()
	return card, nil
}

// Makes a slim detailed card for a tank
func makeSlimCard(card cardData, session wgapi.VehicleStats, lastSession wgapi.VehicleStats) (cardData, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()

	ctx := *card.context
	if err := ctx.LoadFontFace(fontPath, (fontSize)); err != nil {
		return card, err
	}

	if session.Battles < 1 {
		return card, fmt.Errorf("sessions battles is < 1")
	}

	ctx.SetColor(color.White)
	tankNameWidth := float64(card.context.Width()) * 0.375
	tankBlockWidth := (float64(card.context.Width()) - tankNameWidth) / 3

	// Default Block
	var defaultBlock cardBlock
	defaultBlock.textSize = fontSize
	defaultBlock.width = int(tankBlockWidth)
	defaultBlock.height = card.context.Height()
	defaultBlock.bigTextColor = bigTextColor
	defaultBlock.smallTextColor = altTextColor

	// Draw tank name
	finalName := ""
	dotsW, _ := ctx.MeasureString("...")

	for _, r := range []rune(session.TankName) {
		w, _ := ctx.MeasureString(finalName)
		if (w + dotsW) > (tankNameWidth - (float64(frameMargin) * 1.5)) {
			finalName = finalName + "..."
			break
		}
		finalName = finalName + string(r)
	}
	_, nameH := ctx.MeasureString(finalName)

	nameY := (float64(card.context.Height()) - ((float64(card.context.Height()) - nameH) / 2))
	ctx.DrawString(finalName, (float64(frameMargin) * 1.5), nameY)

	// Draw tank tier
	if err := ctx.LoadFontFace(fontPath, (fontSize * 0.75)); err != nil {
		return card, err
	}
	tierW, tierH := ctx.MeasureString(tierToRoman(session.TankTier))
	tierX := float64(frameMargin/2) + ((float64(frameMargin) - tierW) / 2)
	tierY := (float64(card.context.Height()) - ((float64(card.context.Height()) - tierH) / 2))
	ctx.DrawString(tierToRoman(session.TankTier), tierX, tierY)

	// 3 Blocks - DMG / WR / WN8
	// Block 3 - Draw WN8
	ratingBlock := cardBlock(defaultBlock)
	// Icon
	ratingBlock.smallText = "WN8"
	ratingBlock.bigText = "-"
	if session.TankWN8 > -1 {
		ratingBlock.hasBigIcon = true
		ratingBlock.bigIconColor = getRatingColor(session.TankWN8)
		ratingBlock.bigText = strconv.Itoa(session.TankWN8)
	}
	ratingBlock, err := addBlockCtx(ratingBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(ratingBlock.context.Image(), int(tankNameWidth+(tankBlockWidth*2)), 0)

	// Block 2 - Winrate
	winrateBlock := cardBlock(defaultBlock)
	winrateSession := ((float64(session.Wins) / float64(session.Battles)) * 100)
	winrateBlock.bigText = fmt.Sprintf("%.1f", winrateSession) + "% (" + strconv.Itoa(session.Battles) + ")"
	winrateBlock.smallText = "Winrate"
	winrateBlock.hasBigIcon = true
	if lastSession.Battles > 0 && ((float64(session.Wins)/float64(session.Battles))*100) > ((float64(lastSession.Wins)/float64(lastSession.Battles))*100) {
		winrateBlock.bigArrowDirection = 1
	}
	if lastSession.Battles > 0 && ((float64(session.Wins)/float64(session.Battles))*100) < ((float64(lastSession.Wins)/float64(lastSession.Battles))*100) {
		winrateBlock.bigArrowDirection = -1
	}
	winrateBlock, err = addBlockCtx(winrateBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(winrateBlock.context.Image(), int(tankNameWidth+(tankBlockWidth*1)), 0)

	// Block 1 - Avg Damage
	avgDamageBlock := cardBlock(defaultBlock)
	avgDamageSession := strconv.Itoa((session.DamageDealt / session.Battles))
	avgDamageBlock.smallText = "Avg. Damage"
	avgDamageBlock.bigText = avgDamageSession
	avgDamageBlock.hasBigIcon = true
	if lastSession.Battles > 0 && (session.DamageDealt/session.Battles) > (lastSession.DamageDealt/lastSession.Battles) {
		avgDamageBlock.bigArrowDirection = 1
	}
	if lastSession.Battles > 0 && (session.DamageDealt/session.Battles) < (lastSession.DamageDealt/lastSession.Battles) {
		avgDamageBlock.bigArrowDirection = -1
	}
	avgDamageBlock, err = addBlockCtx(avgDamageBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(avgDamageBlock.context.Image(), int(tankNameWidth), 0)

	// Render image
	card.image = ctx.Image()
	return card, nil
}
func addBlockCtx(block cardBlock) (cardBlock, error) {
	ctx := gg.NewContext(block.width, block.height)
	// Color is requested
	if block.isColored == true {
		ctx.SetColor(block.color)
		ctx.DrawRectangle(0, 0, float64(block.width), float64(block.height))
		ctx.Fill()
	}
	if block.textCoeff == 0 {
		block.textCoeff = 0.6
	}
	// Calc altText
	var (
		altMargin float64
		aTxtW     float64
		aTxtH     float64
	)
	if block.altText != "" {
		ctx.SetColor(block.altTextColor)
		if err := ctx.LoadFontFace(fontPath, (block.textSize * (block.textCoeff - 0.15))); err != nil {
			return block, err
		}
		aTxtW, aTxtH = ctx.MeasureString(block.altText)
		altMargin = aTxtH
	}
	availHeiht := block.height
	var totalTextHeight float64 = altMargin
	// Calc small text
	if err := ctx.LoadFontFace(fontPath, (block.textSize * block.textCoeff)); err != nil {
		return block, err
	}
	sTxtW, sTxtH := ctx.MeasureString(block.smallText)
	if sTxtH > 0 {
		totalTextHeight += sTxtH
	}
	sX := ((float64(block.width) - sTxtW) / 2.0)
	// Calc Big text
	if err := ctx.LoadFontFace(fontPath, block.textSize); err != nil {
		return block, err
	}
	bTxtW, bTxtH := ctx.MeasureString(block.bigText)
	if bTxtH > 0 {
		totalTextHeight += bTxtH
	}
	bX := ((float64(block.width) - bTxtW) / 2.0)

	// Draw text
	var drawTextMargins float64
	if block.altText != "" {
		drawTextMargins = (float64(availHeiht) - totalTextHeight) / 4
	} else {
		drawTextMargins = (float64(availHeiht) - totalTextHeight) / 3
	}
	// Big text
	ctx.SetColor(block.bigTextColor)
	if err := ctx.LoadFontFace(fontPath, block.textSize); err != nil {
		return block, err
	}
	bY := bTxtH + drawTextMargins
	ctx.DrawString(block.bigText, bX, bY)

	// Small text
	ctx.SetColor(block.smallTextColor)
	if err := ctx.LoadFontFace(fontPath, (block.textSize * block.textCoeff)); err != nil {
		return block, err
	}
	sY := bY + sTxtH + drawTextMargins
	ctx.DrawString(block.smallText, sX, sY)

	if block.altText != "" {
		if err := ctx.LoadFontFace(fontPath, (block.textSize * (block.textCoeff - 0.15))); err != nil {
			return block, err
		}
		ctx.SetColor(block.altTextColor)
		aX := ((float64(block.width) - aTxtW) / 2.0)
		aY := sY + drawTextMargins + aTxtH
		ctx.DrawString(block.altText, aX, aY)
	}

	// Draw icons
	if block.hasBigIcon == true {
		ctx.SetColor(block.bigIconColor)
		if block.bigArrowDirection == 0 {
			iR := 8.0 * (block.textSize / fontSize)
			iX := bX - (iR * 1.5)
			iY := bY - iR - ((bTxtH - (iR * 2)) / 2)
			ctx.DrawCircle(iX, iY, iR)
			ctx.Fill()
		}
		if block.bigArrowDirection == 1 {
			ctx.SetColor(color.RGBA{0, 255, 0, 180})
			iR := 8.0 * (block.textSize / fontSize)
			iX := bX - (iR * 1.5)
			iY := bY - ((bTxtH - (iR)) / 2) - (fontSize / 10)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 0)
			ctx.Fill()
		}
		if block.bigArrowDirection == -1 {
			ctx.SetColor(color.RGBA{255, 0, 0, 180})
			iR := 8.0 * (block.textSize / fontSize)
			iX := bX - (iR * 1.5)
			iY := bY - bTxtH + ((bTxtH - (iR)) / 2) + (fontSize / 10)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 1)
			ctx.Fill()
		}
	}
	if block.hasSmallIcon == true {
		ctx.SetColor(block.smallIconColor)
		if block.smallArrowDirection == 0 {
			iR := 8.0 * 0.75 * (block.textSize / fontSize)
			iX := sX - (iR * 1.5)
			iY := sY - iR - ((sTxtH - (iR * 2)) / 2)
			ctx.DrawCircle(iX, iY, iR)
			ctx.Fill()
		}
		if block.smallArrowDirection == 1 {
			iR := 8.0 * 0.75 * (block.textSize / fontSize)
			iX := sX - (iR * 1.5)
			iY := sY - iR - ((sTxtH - (iR * 2)) / 2)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 0)
			ctx.Fill()
		}
		if block.smallArrowDirection == -1 {
			iR := 8.0 * 0.75 * (block.textSize / fontSize)
			iX := sX - (iR * 1.5)
			iY := sY - iR - ((sTxtH - (iR * 2)) / 2)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 1)
			ctx.Fill()
		}
	}
	block.context = ctx
	return block, nil
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
