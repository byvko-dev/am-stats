package render

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"image"
	"image/color"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"

	"github.com/cufee/am-stats/config"
	"github.com/cufee/am-stats/stats"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// General settings
var (
	fontPath        = (config.AssetsPath + "/font.ttf")
	fontSizeHeader  = 36.0
	fontSize        = 24.0
	textMargin      = fontSize / 2
	frameWidth      = 900
	frameMargin     = 50
	baseCardWidth   = frameWidth - (2 * frameMargin)
	baseCardHeigh   = 150
	baseCardColor   = color.RGBA{30, 30, 30, 204}
	decorLinesColor = color.RGBA{80, 80, 80, 255}
	bigTextColor    = color.RGBA{255, 255, 255, 255}
	smallTextColor  = color.RGBA{204, 204, 204, 255}
	altTextColor    = color.RGBA{100, 100, 100, 255}

	premiumColor  = color.RGBA{255, 223, 0, 255}
	verifiedColor = color.RGBA{72, 167, 250, 255}
)

// ImageFromStats -
func ImageFromStats(data stats.ExportData, sortKey string, tankLimit int, premium bool, verified bool, bgImage image.Image) (finalImage image.Image, err error) {
	var finalCards allCards
	cardsChan := make(chan cardData, (2 + len(data.SessionStats.Vehicles)))
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
		header, err := makeHeaderCard(prepNewCard(0, headerHeight), data.PlayerDetails.Name, clanTag, "Random Battles", premium, verified)
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
		allStats, err := makeAllStatsCard(prepNewCard(1, 1.5), data)
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
			var tankCard cardData
			if i < 3 {
				tankCard, err = makeDetailedCard(prepNewCard((i+2), 1.0), tank, lastSession)
			} else {
				tankCard, err = makeSlimCard(prepNewCard((i+2), 0.5), tank, lastSession)
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
		finalCards.cards = append(finalCards.cards, c)
	}
	finalCtx, err := addAllCardsToFrame(finalCards, data.SessionStats.Timestamp, bgImage)
	if err != nil {
		return nil, err
	}
	return finalCtx.Image(), err
}

func addAllCardsToFrame(finalCards allCards, timestamp time.Time, bgImage image.Image) (*gg.Context, error) {
	if len(finalCards.cards) == 0 {
		return nil, fmt.Errorf("no cards to be rendered")
	}

	// Frame height
	var totalCardsHeight int
	for _, card := range finalCards.cards {
		totalCardsHeight += card.context.Height()
	}
	totalCardsHeight += ((len(finalCards.cards)-1)*(frameMargin/2) + (frameMargin * 2))
	// Get frame CTX
	ctx, err := prepBgContext(totalCardsHeight, bgImage)
	if err != nil {
		return finalCards.frame, err
	}
	finalCards.frame = ctx
	// Sort cards
	sort.Slice(finalCards.cards, func(i, j int) bool {
		return finalCards.cards[i].index < finalCards.cards[j].index
	})
	// Render cards
	var lastCardPos int = frameMargin / 2
	for i := 0; i < len(finalCards.cards); i++ {
		card := finalCards.cards[i]
		cardMarginH := lastCardPos + (frameMargin / 2)
		finalCards.frame.DrawImage(card.image, frameMargin, cardMarginH)
		lastCardPos = cardMarginH + card.context.Height()
	}

	// Draw timestamp
	if err := finalCards.frame.LoadFontFace(fontPath, fontSize*0.75); err != nil {
		return finalCards.frame, err
	}
	finalCards.frame.SetColor(color.RGBA{100, 100, 100, 100})
	time := timestamp.Format("Session from Jan 2")
	timeW, timeH := finalCards.frame.MeasureString(time)
	timeX := (float64(finalCards.frame.Width()) - timeW) / 2
	timeY := (float64(frameMargin)-timeH)/2 + timeH
	finalCards.frame.DrawString(time, timeX, timeY)

	return finalCards.frame, nil
}

func makeHeaderCard(card cardData, playerName, playerClan, battleType string, premium bool, verified bool) (cardData, error) {
	ctx := *card.context
	if err := ctx.LoadFontFace(fontPath, fontSizeHeader); err != nil {
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
	totalTextH := nameStrH + textMargin + battleTypeH

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
		ctx.SetColor(premiumColor)
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
		radius := (fontSizeHeader / 3)
		ctx.SetColor(verifiedColor)
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
	card.image = ctx.Image()
	return card, nil
}

// Prepare a frame background context
func prepBgContext(totalHeight int, bgImage image.Image) (frameCtx *gg.Context, err error) {
	frameCtx = gg.NewContext(frameWidth, totalHeight)
	bgImage = imaging.Fill(bgImage, frameCtx.Width(), frameCtx.Height(), imaging.Center, imaging.NearestNeighbor)
	bgImage = imaging.Blur(bgImage, 10.0)
	frameCtx.DrawImage(bgImage, 0, 0)
	return frameCtx, nil
}

// Prepare a new cardData struct
func prepNewCard(index int, heightMod float64) cardData {
	cardHeight := int(float64(baseCardHeigh) * heightMod)
	cardWidth := baseCardWidth
	cardCtx := gg.NewContext(cardWidth, cardHeight)
	cardCtx.SetColor(baseCardColor)
	cardCtx.DrawRoundedRectangle(0, 0, float64(cardWidth), float64(cardHeight), fontSize)
	cardCtx.Fill()
	var card cardData
	card.context = cardCtx
	card.index = index
	return card
}
