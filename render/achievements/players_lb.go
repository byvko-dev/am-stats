package render

import (
	"fmt"
	"image"
	"log"
	"sync"

	dataprep "github.com/cufee/am-stats/dataprep/achievements"
	dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	dbGloss "github.com/cufee/am-stats/mongodbapi/v1/glossary"
	"github.com/cufee/am-stats/render"
	"github.com/fogleman/gg"
)

// PlayerAchievementsLbImage -
func PlayerAchievementsLbImage(data []dbAch.AchievementsPlayerData, bgImage image.Image, medals []dataprep.MedalWeight) (finalImage image.Image, err error) {
	// Get icon URLs
	for i, m := range medals {
		m.IconURL, err = dbGloss.GetAchievementIcon(m.Name)
		medals[i] = m
	}

	// Init
	var finalCards render.AllCards
	cardsChan := make(chan render.CardData, (2 + len(data)))
	var wg sync.WaitGroup

	// Configure blueprint for cards
	var slimBlockBP cardBlockData
	slimBlockBP.DefaultSlim()
	slimBlockBP.ChangeCoeff(8, 10)
	slimBlockBP.TextSize = (slimBlockBP.TextSize * 12 / 10)
	slimBlockBP.BigTextSize = (slimBlockBP.BigTextSize * 3 / 2)

	// Get longest name
	var maxNameWidth float64
	var maxClanTagWidth float64
	checkCtx := gg.NewContext(1, 1)
	checkBlock := cardBlockData(slimBlockBP)
	maxScoreWidth, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.SmallTextSize, "Score")

	for i, player := range data {
		// Fix player clan tag
		if player.ClanTag != "" {
			player.ClanTag = fmt.Sprintf("[%s] ", player.ClanTag)
			data[i].ClanTag = player.ClanTag
		}

		// Get text size
		cW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, player.ClanTag)
		nW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, player.Nickname)
		sW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.BigTextSize, fmt.Sprint(player.Score))
		if cW > maxClanTagWidth { // Check clan tag width
			maxClanTagWidth = cW
		}
		if nW > maxNameWidth { // Check name width
			maxNameWidth = nW
		}
		if sW > maxScoreWidth { // Check score width
			maxScoreWidth = sW
		}
	}
	slimBlockBP.NameMargin = maxNameWidth
	slimBlockBP.ClanTagMargin = maxClanTagWidth
	slimBlockBP.SpecialBlockWidth = (maxScoreWidth + slimBlockBP.TextMargin)

	// Calculate required card width
	cardWidth := (0 +
		int(maxNameWidth) + // Name width
		int(slimBlockBP.ClanTagMargin) + // Clan tag width
		(2 * int(slimBlockBP.TextMargin)) + // Card padding
		int(slimBlockBP.SpecialBlockWidth) + // Score block width
		((len(medals)) * (int(slimBlockBP.IconSize) * 3 / 2)) + 0) // Medal blocks width

	// Work on cards in go routines
	for i, player := range data {
		wg.Add(1)

		go func(player dbAch.AchievementsPlayerData, i int) {
			defer wg.Done()

			// Prep card context
			card, err := makeSlimPlayerCard(render.PrepNewCard(1, 0.5, cardWidth), slimBlockBP, player, i, medals)
			card.Index = i
			if err != nil {
				log.Println(err)
				return
			}

			cardsChan <- card
			return
		}(player, i)
	}
	wg.Wait()
	close(cardsChan)

	for c := range cardsChan {
		finalCards.Cards = append(finalCards.Cards, c)
	}

	finalCtx, err := render.AddAllCardsToFrame(finalCards, "Achievements Leaderboard", bgImage)
	if err != nil {
		return nil, err
	}
	return finalCtx.Image(), err
}

// Make large player card
func makeLargePlayerCard(card render.CardData, blueprint cardBlockData, player dbAch.AchievementsPlayerData, position int, medals []dataprep.MedalWeight) (_ render.CardData, err error) {
	return makeSlimPlayerCard(card, blueprint, player, position, medals)
}

// Make slim player card
func makeSlimPlayerCard(card render.CardData, blueprint cardBlockData, player dbAch.AchievementsPlayerData, position int, medals []dataprep.MedalWeight) (_ render.CardData, err error) {
	ctx := *card.Context
	if err := ctx.LoadFontFace(render.FontPath, (blueprint.TextSize)); err != nil {
		return card, err
	}

	if player.Score < 1 {
		return card, fmt.Errorf("player score is 0")
	}

	ctx.SetColor(render.BigTextColor)
	blocksCount := float64(len(medals))
	blockWidth := (float64(card.Context.Width()) - (blueprint.NameMargin + blueprint.ClanTagMargin + blueprint.SpecialBlockWidth + (blueprint.TextMargin * 2))) / blocksCount

	blueprint.Width = int(blockWidth)
	blueprint.Height = card.Context.Height()

	var lastXOffs int = int(blueprint.TextMargin)

	// Name Block
	playerName := player.Nickname
	nameBlock := cardBlockData(blueprint)
	nameBlock.BigText = playerName
	nameBlock.Width = int(blueprint.NameMargin)
	nameBlock.BigTextSize = blueprint.TextSize
	nameBlock.TextAlign = -1

	if err = renderBlock(&nameBlock); err != nil {
		return card, err
	}
	ctx.DrawImage(nameBlock.Context.Image(), lastXOffs, 0)
	lastXOffs += nameBlock.Width

	// Clan tag Block
	clanBlock := cardBlockData(blueprint)
	clanBlock.BigText = player.ClanTag
	clanBlock.Width = int(blueprint.ClanTagMargin)
	clanBlock.BigTextSize = blueprint.TextSize
	clanBlock.BigTextColor = blueprint.SmallTextColor

	if err = renderBlock(&clanBlock); err != nil {
		return card, err
	}
	ctx.DrawImage(clanBlock.Context.Image(), lastXOffs, 0)
	lastXOffs += clanBlock.Width

	// Score Block
	scoreBlock := cardBlockData(blueprint)
	scoreBlock.BigText = fmt.Sprint(player.Score)
	scoreBlock.Width = int(blueprint.SpecialBlockWidth)
	scoreBlock.SmallText = "Score"

	if err = renderBlock(&scoreBlock); err != nil {
		return card, err
	}
	ctx.DrawImage(scoreBlock.Context.Image(), lastXOffs, 0)
	lastXOffs += scoreBlock.Width

	//  Medal Blocks
	for _, m := range medals {
		medalBlock := cardBlockData(blueprint)
		medalBlock.AltText = fmt.Sprint(getField(player.Data, m.Name))
		medalBlock.AltTextColor = blueprint.SmallTextColor
		medalBlock.IconURL = m.IconURL

		if err = renderBlock(&medalBlock); err != nil {
			return card, err
		}
		ctx.DrawImage(medalBlock.Context.Image(), lastXOffs, 0)
		lastXOffs += int(blockWidth)
	}

	// Render image
	card.Image = ctx.Image()
	return card, nil
}
