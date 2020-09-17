package render

import (
	"fmt"
	"math"
	"strconv"
	
	"sync"
	"image"
	"image/color"
    "github.com/fogleman/gg"
	"github.com/disintegration/imaging"
	
	"github.com/cufee/am-stats/stats"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)
// General settings
var (  
    fontPath        = "../am-stats/render/assets/font.ttf"            
    fontSizeHeader  = 36.0
    fontSize        = 24.0
    textMargin      = fontSize / 2
    frameWidth      = 900
    frameMargin     = 50
    baseCardWidth   = frameWidth - (2*frameMargin)
    baseCardHeigh   = 150
    baseCardColor   = color.RGBA{0,0,0,100}
    defaultBG       = "../am-stats/render/assets/bg_frame.png"
)
// ImageFromStats - 
func ImageFromStats(data stats.ExportData, sortKey string, tankLimit int) (finalImage image.Image, err error){
	var finalCards allCards
	cardsChan := make(chan cardData, (2 + len(data.SessionStats.Vehicles)))
	var wg sync.WaitGroup
	// Work on cards in go routines
	wg.Add(1)
	go func() {
		defer wg.Done()
		header, err := makeHeaderCard(prepNewCard(0, 1.0), data.PlayerDetails.Name, data.PlayerDetails.ClanTag, "Random Battles")
		if err != nil {
			return
		}
		cardsChan <- header
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		allStats, err := makeAllStatsCard(prepNewCard(1, 1.5), data)
		if err != nil {
			return
		}
		cardsChan <- allStats
	}()
	
	// Render a card per vehicle while under the tankLimit, sort cards by sortKey

	wg.Wait()
	close(cardsChan)
	for c := range cardsChan {
		finalCards.cards = append(finalCards.cards, c)
	}
	finalCtx, err := addAllCardsToFrame(finalCards)
	if err != nil {
		return nil, err
	}
	return finalCtx.Image(), err
}

func addAllCardsToFrame(finalCards allCards) (*gg.Context, error){
	if len(finalCards.cards) == 0 {
		return nil, fmt.Errorf("no cards to be rendered")
	}
	var totalCardsHeight int
    for _, card := range finalCards.cards {
		totalCardsHeight += card.context.Height()
	}
	totalCardsHeight += ((len(finalCards.cards) + 1) * frameMargin)
    finalCards.frame = prepBgContext(totalCardsHeight)
	var lastCardPos int
    for _, card := range finalCards.cards {
        cardMarginH := lastCardPos + frameMargin
		finalCards.frame.DrawImage(card.image, frameMargin, cardMarginH)
		lastCardPos += cardMarginH + card.context.Height()
    }
    return finalCards.frame, nil
}

func makeHeaderCard(card cardData, playerName, playerClan, battleType string) (cardData, error) {
    ctx := *card.context
    if err := ctx.LoadFontFace(fontPath, fontSizeHeader);err != nil {
        return card, err
    }
    playerStr := playerName + " " + playerClan
    // Calculate text size
    playerStrW, playerStrH := ctx.MeasureString(playerStr)
    battleTypeW, battleTypeH := ctx.MeasureString(battleType)
    totalTextW := playerStrW
    if playerStrW < battleTypeW {
		totalTextW = battleTypeW
    }
    totalTextH := playerStrH + textMargin + battleTypeH
    xOffset := ((float64(ctx.Width()) - totalTextW) / 2)
    yOffset := ((float64(ctx.Height()) - totalTextH) / 2)
	// Draw battle type text
	ctx.SetColor(color.RGBA{255,255,255,200})
    btDrawX := ((totalTextW - battleTypeW) / 2) + xOffset
    btDrawY :=  yOffset + battleTypeH
    ctx.DrawString(battleType, btDrawX, btDrawY)
    // Draw player name and tag text
	ctx.SetColor(color.White)
    psDrawX := ((totalTextW - playerStrW) / 2) + xOffset
    psDrawY := totalTextH + yOffset
    ctx.DrawString(playerStr, psDrawX, psDrawY)
    // Make image
    card.image = ctx.Image()
    return card, nil
}

func makeAllStatsCard(card cardData, data stats.ExportData) (cardData, error) {
    ctx := *card.context
    if err := ctx.LoadFontFace(fontPath, fontSize);err != nil {
        return card, err
    }
	ctx.SetColor(color.White)
	// Default Block settings
	blockWidth 			:= card.context.Width() / 3
	bottomBlockWidth 	:= card.context.Width() / 4
	availableHeight 	:= (ctx.Height() - int(fontSize / 2)) / 2
	blockHeight 		:= availableHeight
	var defaultBlock cardBlock
	defaultBlock.textSize 		= fontSize
	defaultBlock.width	  		= blockWidth
	defaultBlock.height			= blockHeight
	defaultBlock.bigTextColor	= color.RGBA{255,255,255,255}
	defaultBlock.smallTextColor	= color.RGBA{255,255,255,200}
	defaultBlock.altTextColor	= color.RGBA{255,255,255,200}
	// Top Row - 3 Blocks (Battles, WN8, WR)
	// Block 1 - Battles
	battlesBlock := cardBlock(defaultBlock)
	battlesBlock.textSize 	= fontSize * 1.25
	battlesBlock.smallText 	= strconv.Itoa(data.PlayerDetails.Stats.All.Battles)
	battlesBlock.bigText 	= strconv.Itoa(data.SessionStats.BattlesAll)
	battlesBlock.altText 	= "Battles"
	battlesBlock, err := addBlockCtx(battlesBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(battlesBlock.context.Image(), 0, 0)

	// Block 2 - WN8
	ratingBlock := cardBlock(defaultBlock)
	// Icon
	ratingBlock.hasBigIcon			= true
	ratingBlock.bigIconColor		= getRatingColor(data.SessionStats.SessionRating)
	ratingBlock.hasSmallIcon		= true
	ratingBlock.smallIconColor		= getRatingColor(data.PlayerDetails.CareerWN8)
	ratingBlock.height				= blockHeight + int(fontSize)
	ratingBlock.textSize 			= fontSize * 1.50
	ratingBlock.smallText 			= strconv.Itoa(data.PlayerDetails.CareerWN8)
	ratingBlock.bigText 			= strconv.Itoa(data.SessionStats.SessionRating)
	ratingBlock.altText 			= "WN8"
	ratingBlock, err = addBlockCtx(ratingBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(ratingBlock.context.Image(), blockWidth, 0)
	// Block 3 - WR
	winrateBlock := cardBlock(battlesBlock)
	oldBattles 				:= data.PlayerDetails.Stats.All.Battles - data.SessionStats.StatsAll.Battles
	oldWins					:= data.PlayerDetails.Stats.All.Wins - data.SessionStats.StatsAll.Wins
	oldWinrate				:= (float64(oldWins) / float64(oldBattles)) * 100
	winrateAll				:= ((float64(data.PlayerDetails.Stats.All.Wins) / float64(data.PlayerDetails.Stats.All.Battles)) * 100)
	winrateSession			:= ((float64(data.SessionStats.StatsAll.Wins) / float64(data.SessionStats.StatsAll.Battles)) * 100)
	winrateChange 			:= math.Round((winrateAll - oldWinrate)*100)/100
	winrateChangeStr := ""
	if winrateChange > 0.00 {
		winrateChangeStr	= fmt.Sprintf(" (+%.2f", winrateChange) + "%)"
	}
	if winrateChange < 0.00 {
		winrateChangeStr	= fmt.Sprintf(" (%.2f", winrateChange) + "%)"
	}
	winrateAllStr			:= fmt.Sprintf("%.2f", winrateAll) + "%" + winrateChangeStr
	winrateSessionStr		:= fmt.Sprintf("%.2f", winrateSession) + "%"
	winrateBlock.smallText 	= winrateAllStr
	winrateBlock.bigText 	= winrateSessionStr
	winrateBlock.altText 	= "Winrate"
	winrateBlock, err = addBlockCtx(winrateBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(winrateBlock.context.Image(), (blockWidth * 2), 0)

	// Bottom Row - 4 Blocks
	// Block 1 - Avg Damage
	avgDamageBlock := cardBlock(defaultBlock)
	avgDamageBlock.textSize 	= fontSize * 1.25
	avgDamageBlock.width 		= bottomBlockWidth
	avgDamageAll				:= strconv.Itoa((data.PlayerDetails.Stats.All.DamageDealt / data.PlayerDetails.Stats.All.Battles))
	avgDamageSession			:= strconv.Itoa((data.SessionStats.StatsAll.DamageDealt / data.SessionStats.StatsAll.Battles))
	avgDamageBlock.smallText 	= avgDamageAll
	avgDamageBlock.bigText 		= avgDamageSession
	avgDamageBlock.altText 		= "Avg. Damage"
	avgDamageBlock, err = addBlockCtx(avgDamageBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(avgDamageBlock.context.Image(), 0, blockHeight)
	// Block 2 - Damage Ratio
	damageRatioBlock := cardBlock(avgDamageBlock)
	damageRatioAll				:= fmt.Sprintf("%.2f", ((float64(data.PlayerDetails.Stats.All.DamageDealt) / float64(data.PlayerDetails.Stats.All.DamageReceived))))
	damageRatioSession			:= fmt.Sprintf("%.2f", ((float64(data.SessionStats.StatsAll.DamageDealt) / float64(data.SessionStats.StatsAll.DamageReceived))))
	damageRatioBlock.smallText 	= damageRatioAll
	damageRatioBlock.bigText 	= damageRatioSession
	damageRatioBlock.altText 	= "Damage Ratio"
	damageRatioBlock, err = addBlockCtx(damageRatioBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(damageRatioBlock.context.Image(), bottomBlockWidth, blockHeight)
	// Block 3 - Destruction Ratio
	destrRatioBlock := cardBlock(avgDamageBlock)
	destrRatioAll				:= fmt.Sprintf("%.2f", ((float64(data.PlayerDetails.Stats.All.Frags) / (float64(data.PlayerDetails.Stats.All.Battles) - float64(data.PlayerDetails.Stats.All.SurvivedBattles)))))
	destrRatioSession			:= fmt.Sprintf("%.2f", ((float64(data.SessionStats.StatsAll.Frags) / (float64(data.SessionStats.StatsAll.Battles) - float64(data.SessionStats.StatsAll.SurvivedBattles)))))
	destrRatioBlock.smallText 	= destrRatioAll
	destrRatioBlock.bigText 	= destrRatioSession
	destrRatioBlock.altText 	= "Destruction Ratio"
	destrRatioBlock, err = addBlockCtx(destrRatioBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(destrRatioBlock.context.Image(), (bottomBlockWidth * 2), blockHeight)
	// Block 4 - Average XP
	avgXPBlock := cardBlock(avgDamageBlock)
	avgXPAll				:= strconv.Itoa((data.PlayerDetails.Stats.All.Xp / data.PlayerDetails.Stats.All.Battles))
	avgXPSession			:= strconv.Itoa((data.SessionStats.StatsAll.Xp / data.SessionStats.StatsAll.Battles))
	avgXPBlock.smallText 	= avgXPAll
	avgXPBlock.bigText 		= avgXPSession
	avgXPBlock.altText 		= "Avg. XP"
	avgXPBlock, err = addBlockCtx(avgXPBlock)
	if err != nil {
		return card, err
	}
	ctx.DrawImage(avgXPBlock.context.Image(), (bottomBlockWidth * 3), blockHeight)

	// Render image
    card.image = ctx.Image()
    return card, nil
}
// Makes a detailed card for a tank
func makeDetailedCard(card cardData, tank wgapi.VehicleStats) (cardData, error) {
    ctx := *card.context
    if err := ctx.LoadFontFace(fontPath, fontSize);err != nil {
        return card, err
    }
	// Default Block settings
	blockWidth 			:= card.context.Width() / 3
	// bottomBlockWidth 	:= card.context.Width() / 4
	availableHeight 	:= (ctx.Height() - int(fontSize / 2)) / 2
	blockHeight 		:= availableHeight
	var defaultBlock cardBlock
	defaultBlock.textSize 		= fontSize
	defaultBlock.width	  		= blockWidth
	defaultBlock.height			= blockHeight
	defaultBlock.bigTextColor	= color.RGBA{255,255,255,255}
	defaultBlock.smallTextColor	= color.RGBA{255,255,255,200}
	defaultBlock.altTextColor	= color.RGBA{255,255,255,200}
	return card, nil
}
func addBlockCtx(block cardBlock) (cardBlock, error){
	ctx := gg.NewContext(block.width, block.height)
	// Color is requested
	if block.isColored == true {
		ctx.SetColor(block.color)
		ctx.DrawRectangle(0,0,float64(block.width),float64(block.height))
		ctx.Fill()
	}
	// Draw altText
	if block.altText != "" {
		ctx.SetColor(block.altTextColor)
		if err := ctx.LoadFontFace(fontPath, (block.textSize * 0.5));err != nil {
			return block, err
		}
		aTxtW, _ := ctx.MeasureString(block.altText)
		sX := ((float64(block.width) - aTxtW) / 2.0)
		sY := float64(block.height) - (block.textSize / 4)
		ctx.DrawString(block.altText, sX, sY)
	}
	// Draw small text
	ctx.SetColor(block.smallTextColor)
	if err := ctx.LoadFontFace(fontPath, (block.textSize * 0.75));err != nil {
        return block, err
    }
	sTxtW, sTxtH := ctx.MeasureString(block.smallText)
	sX := ((float64(block.width) - sTxtW) / 2.0)
	sY := float64(block.height / 2) + sTxtH + (block.textSize / 4)
	ctx.DrawString(block.smallText, sX, sY)
	// Draw Big text
	ctx.SetColor(block.bigTextColor)
	if err := ctx.LoadFontFace(fontPath, block.textSize);err != nil {
        return block, err
    }
	bTxtW, bTxtH := ctx.MeasureString(block.bigText)
	bX := ((float64(block.width) - bTxtW) / 2.0)
	bY := float64(block.height / 2) - (block.textSize / 4)
	ctx.DrawString(block.bigText, bX, bY)
	if block.hasBigIcon == true {
		ctx.SetColor(block.bigIconColor)
		if block.bigArrowDirection == 0 {
			iR := 10.0
			iX := bX - (iR*1.5)
			iY := bY - iR - ((bTxtH - (iR*2)) / 2)
			ctx.DrawCircle(iX, iY, iR)
			ctx.Fill()
		}
		if block.bigArrowDirection == 1 {
			iR := 10.0
			iX := bX - (iR*1.5)
			iY := bY - iR - ((bTxtH - (iR*2)) / 2)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 0)
			ctx.Fill()
		}
		if block.bigArrowDirection == -1 {
			iR := 10.0
			iX := bX - (iR*1.5)
			iY := bY - iR - ((bTxtH - (iR*2)) / 2)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 1)
			ctx.Fill()
		}
	}
	if block.hasSmallIcon == true {
		ctx.SetColor(block.smallIconColor)
		if block.smallArrowDirection == 0 {
			iR := 10.0 * 0.75
			iX := sX - (iR*1.5)
			iY := sY - iR - ((sTxtH - (iR*2)) / 2)
			ctx.DrawCircle(iX, iY, iR)
			ctx.Fill()
		}
		if block.smallArrowDirection == 1 {
			iR := 10.0 * 0.75
			iX := sX - (iR*1.5)
			iY := sY - iR - ((sTxtH - (iR*2)) / 2)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 0)
			ctx.Fill()
		}
		if block.smallArrowDirection == -1 {
			iR := 10.0 * 0.75
			iX := sX - (iR*1.5)
			iY := sY - iR - ((sTxtH - (iR*2)) / 2)
			ctx.DrawRegularPolygon(3, iX, iY, iR, 1)
			ctx.Fill()
		}
	}
	block.context = ctx
	return block, nil
}

// Rating color calculator
func getRatingColor(r int) (color.RGBA) {
	if r > 0 && r < 301 {
		return color.RGBA{0,0,0,0}
	}
	if r > 300 && r < 451 {
		return color.RGBA{251,83,83,180}
	}
	if r > 450 && r < 651 {
		return color.RGBA{255,160,49,180}
	}
	if r > 650 && r < 901 {
		return color.RGBA{255,244,65,180}
	}
	if r > 900 && r < 1201 {
		return color.RGBA{149,245,62,180}
	}
	if r > 1200 && r < 1601 {
		return color.RGBA{103,190,51,180}
	}
	if r > 1600 && r < 2001 {
		return color.RGBA{106,236,255,180}
	}
	if r > 2000 && r < 2451 {
		return color.RGBA{46,174,193,180}
	}
	if r > 2450 && r < 2901 {
		return color.RGBA{208,108,255,180}
	}
	if r > 2900 {
		return color.RGBA{142,65,177,180}
	}
	return color.RGBA{0,0,0,0}
}

// Prepare a frame background context
func prepBgContext(totalHeight int) (*gg.Context) {
    frameCtx := gg.NewContext(frameWidth, totalHeight)
    bgImage, err := gg.LoadImage(defaultBG)
    if err != nil {
        panic(err)
    }
	bgImage = imaging.Fill(bgImage, frameCtx.Width(), frameCtx.Height(), imaging.Center, imaging.NearestNeighbor)
	bgImage = imaging.Blur(bgImage, 4.0)
    frameCtx.DrawImage(bgImage, 0, 0)
    return frameCtx
}
// Prepare a new cardData struct
func prepNewCard(index int, heightMod float64) (cardData) {
    cardHeight := int(float64(baseCardHeigh) * heightMod)
    cardWidth  := baseCardWidth
    cardCtx := gg.NewContext(cardWidth, cardHeight)
    cardCtx.SetColor(baseCardColor)
    cardCtx.DrawRectangle(0, 0, float64(cardWidth), float64(cardHeight))
    cardCtx.Fill()
    var card cardData
    card.context = cardCtx
    card.index   = index
    return card
}