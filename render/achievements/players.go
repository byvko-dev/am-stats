package render

import (
	"fmt"
	"image"
	"log"
	"sync"

	dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	mongodbapi "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	dbGloss "github.com/cufee/am-stats/mongodbapi/v1/glossary"
	"github.com/cufee/am-stats/render"
	"github.com/fogleman/gg"
)

// PlayerAchievementsLbImage -
func PlayerAchievementsLbImage(data []dbAch.AchievementsPlayerData, bgImage image.Image, medals []mongodbapi.MedalWeight) (finalImage image.Image, err error) {
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
			player.ClanTag = fmt.Sprintf("[%s]", player.ClanTag)
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

		// Fill medal scores
		for _, m := range medals {
			playerMedal := mongodbapi.MedalWeight(m)
			playerMedal.Score = getField(player.Data, m.Name)
			data[i].Medals = append(data[i].Medals, playerMedal)
		}
	}
	slimBlockBP.NameMargin = maxNameWidth
	slimBlockBP.ClanTagMargin = maxClanTagWidth + slimBlockBP.TextMargin
	slimBlockBP.SpecialBlockWidth = maxScoreWidth + slimBlockBP.TextMargin

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
			card := render.PrepNewCard(1, 0.5, cardWidth)
			card.Index = i
			if err := makePlayerSlimCard(&card, slimBlockBP, &player, i, player.Medals); err != nil {
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

// Make slim player card
func makePlayerSlimCard(card *render.CardData, blueprint cardBlockData, data *dbAch.AchievementsPlayerData, position int, medals []mongodbapi.MedalWeight) error {
	if err := card.Context.LoadFontFace(render.FontPath, (blueprint.TextSize)); err != nil {
		return err
	}

	if data.Score < 1 {
		return fmt.Errorf("score is 0")
	}

	card.Context.SetColor(render.BigTextColor)
	blocksCount := float64(len(medals))
	card.BlockWidth = (float64(card.Context.Width()) - (blueprint.NameMargin + blueprint.ClanTagMargin + blueprint.SpecialBlockWidth + (blueprint.TextMargin * 2))) / blocksCount

	blueprint.Width = int(card.BlockWidth)
	blueprint.Height = card.Context.Height()

	card.LastXOffs = int(blueprint.TextMargin)

	// Name Block
	playerName := data.Nickname
	nameBlock := cardBlockData(blueprint)
	nameBlock.BigText = playerName
	nameBlock.Width = int(blueprint.NameMargin)
	nameBlock.BigTextSize = blueprint.TextSize
	nameBlock.TextAlign = -1

	if err := renderBlock(&nameBlock); err != nil {
		return err
	}
	card.Context.DrawImage(nameBlock.Context.Image(), card.LastXOffs, 0)
	card.LastXOffs += nameBlock.Width

	// Clan tag Block
	clanBlock := cardBlockData(blueprint)
	clanBlock.BigText = data.ClanTag
	clanBlock.Width = int(blueprint.ClanTagMargin)
	clanBlock.BigTextSize = blueprint.TextSize
	clanBlock.BigTextColor = blueprint.SmallTextColor

	if err := renderBlock(&clanBlock); err != nil {
		return err
	}
	card.Context.DrawImage(clanBlock.Context.Image(), card.LastXOffs, 0)
	card.LastXOffs += clanBlock.Width

	// Add score
	err := addScoreAndMedals(card, blueprint, data.Score, position, medals)
	if err != nil {
		return err
	}

	// Render image
	card.Image = card.Context.Image()
	return err
}
