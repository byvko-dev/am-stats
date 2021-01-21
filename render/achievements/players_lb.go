package render

import (
	"fmt"
	"image"
	"log"
	"sync"

	"github.com/cufee/am-stats/config"
	dataprep "github.com/cufee/am-stats/dataprep/achievements"
	dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	dbGloss "github.com/cufee/am-stats/mongodbapi/v1/glossary"
	"github.com/cufee/am-stats/render"
	"github.com/fogleman/gg"
)

var slimBlockBP cardBlockData
var largeBlockBP cardBlockData

func init() {
	// Blueprint for small blocks
	slimBlockBP.IconSize = 50
	slimBlockBP.TextCoeff = 0.6
	slimBlockBP.NameMarginCoef = 0.5
	slimBlockBP.BlockTextSize = render.FontSize * 1.25
	slimBlockBP.TextSize = render.FontSize
	slimBlockBP.BigTextColor = render.BigTextColor
	slimBlockBP.AltTextColor = render.AltTextColor
	slimBlockBP.SmallTextColor = render.SmallTextColor

	// Blueprint for large blocks
	largeBlockBP = cardBlockData(slimBlockBP)
	largeBlockBP.BlockTextSize = render.FontSize * 1.75
	largeBlockBP.TextCoeff = 0.45
	largeBlockBP.IconSize = 75
}

// PlayerAchievementsLbImage -
func PlayerAchievementsLbImage(data []dbAch.AchievementsPlayerData, medals []dataprep.MedalWeight) (finalImage image.Image, err error) {
	defer recoverPanic(err, "PlayerAchievementsLbImage")

	// Get icon URLs
	for i, m := range medals {
		m.IconURL, err = dbGloss.GetAchievementIcon(m.Name)
		medals[i] = m
	}

	// Init
	var finalCards render.AllCards
	cardsChan := make(chan render.CardData, (2 + len(data)))
	var wg sync.WaitGroup

	// Work on cards in go routines
	for i, player := range data {
		wg.Add(1)

		go func(player dbAch.AchievementsPlayerData, i int) {
			defer wg.Done()

			// Prep card context
			card, err := makeSlimPlayerCard(render.PrepNewCard(1, 0.5), slimBlockBP, player, i, medals)
			card.Index = i
			if err != nil {
				log.Println(err)
				return
			}

			cardsChan <- card
		}(player, i)
	}
	wg.Wait()
	close(cardsChan)

	for c := range cardsChan {
		finalCards.Cards = append(finalCards.Cards, c)
	}

	// Get BG
	var bgImage image.Image
	bgImage, err = gg.LoadImage(config.AssetsPath + config.DefaultBG)
	if err != nil {
		return nil, err
	}

	finalCtx, err := render.AddAllCardsToFrame(finalCards, "Achievements Leaderboard", bgImage)
	if err != nil {
		return nil, err
	}
	return finalCtx.Image(), err
}

// Make large player card
func makeLargePlayerCard(card render.CardData, blueprint cardBlockData, player dbAch.AchievementsPlayerData, position int, medals []dataprep.MedalWeight) (_ render.CardData, err error) {
	defer recoverPanic(err, "makeLargePlayerCard")
	return makeSlimPlayerCard(card, blueprint, player, position, medals)
}

// Make slim player card
func makeSlimPlayerCard(card render.CardData, blueprint cardBlockData, player dbAch.AchievementsPlayerData, position int, medals []dataprep.MedalWeight) (_ render.CardData, err error) {
	defer recoverPanic(err, "makeSlimPlayerCard")

	ctx := *card.Context
	if err := ctx.LoadFontFace(render.FontPath, (render.FontSize)); err != nil {
		return card, err
	}

	if player.Score < 1 {
		return card, fmt.Errorf("player score is 0")
	}

	ctx.SetColor(render.BigTextColor)
	playerNameWidth := float64(card.Context.Width()) * blueprint.NameMarginCoef
	blocksCount := float64(len(medals) + 1)
	blockWidth := (float64(card.Context.Width()) - playerNameWidth) / blocksCount

	blueprint.Width = int(blockWidth)
	blueprint.Height = card.Context.Height()

	// Draw name
	playerName := player.Nickname
	clanTag := ""
	if player.ClanTag != "" {
		clanTag = fmt.Sprintf("[%s] ", player.ClanTag)
	}
	_, nameH := ctx.MeasureString(playerName + clanTag)
	tagW, _ := ctx.MeasureString(clanTag)

	nameY := (float64(card.Context.Height()) - ((float64(card.Context.Height()) - nameH) / 2))
	ctx.DrawString(playerName, (float64(render.FrameMargin) + tagW), nameY)
	ctx.SetColor(render.SmallTextColor)
	ctx.DrawString(clanTag, (float64(render.FrameMargin)), nameY)

	// Score Block
	scoreBlock := cardBlockData(blueprint)
	scoreBlock.BigText = fmt.Sprint(player.Score)
	scoreBlock.SmallText = "Score"

	// scoreBlock.Color = render.DebugColorRed

	if err = renderBlock(&scoreBlock); err != nil {
		return card, err
	}
	ctx.DrawImage(scoreBlock.Context.Image(), int(playerNameWidth+(blockWidth*0)), 0)

	//  Medal Blocks
	for i, m := range medals {
		medalBlock := cardBlockData(blueprint)
		medalBlock.SmallText = fmt.Sprint(getField(player.Data, m.Name))
		medalBlock.IconURL = m.IconURL

		// medalBlock.Color = render.DebugColorPink

		if err = renderBlock(&medalBlock); err != nil {
			return card, err
		}
		ctx.DrawImage(medalBlock.Context.Image(), int(playerNameWidth+(blockWidth*float64(i+1))), 0)
	}

	// Render image
	card.Image = ctx.Image()
	return card, nil
}
