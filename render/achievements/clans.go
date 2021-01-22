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

// ClanAchievementsLbImage -
func ClanAchievementsLbImage(data []dbAch.ClanAchievements, bgImage image.Image, medals []mongodbapi.MedalWeight) (finalImage image.Image, err error) {
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

	for i, clan := range data {
		// Fix player clan tag
		if clan.ClanTag != "" {
			clan.ClanTag = fmt.Sprintf("[%s]", clan.ClanTag)
			data[i].ClanTag = clan.ClanTag
		}

		// Get text size
		cW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, clan.ClanTag)
		pW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, fmt.Sprintf("#%v", i+1))
		sW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.BigTextSize, fmt.Sprint(clan.Score))
		if cW > maxClanTagWidth { // Check clan tag width
			maxClanTagWidth = cW
		}
		if pW > maxNameWidth { // Check clan tag width
			maxNameWidth = cW
		}
		if sW > maxScoreWidth { // Check score width
			maxScoreWidth = sW
		}

		// Fill medal scores
		for _, m := range medals {
			clanMedal := mongodbapi.MedalWeight(m)
			clanMedal.Score = getField(clan.Data, m.Name)
			data[i].Medals = append(data[i].Medals, clanMedal)
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
		((len(medals)) * (int(slimBlockBP.IconSize) * 3 / 2)) + 0) + // Medal blocks width
		int(slimBlockBP.SpecialBlockWidth) // Players counter block

	// Work on cards in go routines
	for i, clan := range data {
		wg.Add(1)

		go func(clan dbAch.ClanAchievements, i int) {
			defer wg.Done()

			// Prep card context
			card := render.PrepNewCard(1, 0.5, cardWidth)
			card.Index = i
			if err := makeClanSlimCard(&card, slimBlockBP, &clan, i, clan.Medals); err != nil {
				log.Println(err)
				return
			}

			cardsChan <- card
			return
		}(clan, i)
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
func makeClanSlimCard(card *render.CardData, blueprint cardBlockData, data *dbAch.ClanAchievements, position int, medals []mongodbapi.MedalWeight) error {
	if err := card.Context.LoadFontFace(render.FontPath, (blueprint.TextSize)); err != nil {
		return err
	}

	if data.Score < 1 {
		return fmt.Errorf("score is 0")
	}

	card.Context.SetColor(render.BigTextColor)
	blocksCount := float64(len(medals))
	card.BlockWidth = (float64(card.Context.Width()) - (blueprint.NameMargin + blueprint.ClanTagMargin + (2 * blueprint.SpecialBlockWidth) + (blueprint.TextMargin * 2))) / blocksCount

	blueprint.Width = int(card.BlockWidth)
	blueprint.Height = card.Context.Height()

	card.LastXOffs = int(blueprint.TextMargin)

	// Position Block
	posBlock := cardBlockData(blueprint)
	posBlock.BigText = fmt.Sprintf("#%v", position+1)
	posBlock.Width = int(blueprint.NameMargin)
	posBlock.BigTextSize = blueprint.TextSize
	posBlock.BigTextColor = blueprint.SmallTextColor

	if err := renderBlock(&posBlock); err != nil {
		return err
	}
	card.Context.DrawImage(posBlock.Context.Image(), card.LastXOffs, 0)
	card.LastXOffs += posBlock.Width

	// Clan tag Block
	clanBlock := cardBlockData(blueprint)
	clanBlock.BigText = data.ClanTag
	clanBlock.Width = int(blueprint.ClanTagMargin)
	clanBlock.BigTextSize = blueprint.TextSize

	if err := renderBlock(&clanBlock); err != nil {
		return err
	}
	card.Context.DrawImage(clanBlock.Context.Image(), card.LastXOffs, 0)
	card.LastXOffs += clanBlock.Width

	// Members Block
	membersBlock := cardBlockData(blueprint)
	membersBlock.BigText = fmt.Sprintf("%v", data.Members)
	membersBlock.SmallText = "Players"
	membersBlock.Width = int(blueprint.SpecialBlockWidth)
	membersBlock.BigTextSize = blueprint.TextSize
	membersBlock.BigTextColor = blueprint.SmallTextColor

	if err := renderBlock(&membersBlock); err != nil {
		return err
	}
	card.Context.DrawImage(membersBlock.Context.Image(), card.LastXOffs, 0)
	card.LastXOffs += membersBlock.Width

	// Add score
	err := addScoreAndMedals(card, blueprint, data.Score, position, medals)
	if err != nil {
		return err
	}

	// Render image
	card.Image = card.Context.Image()
	return err
}
