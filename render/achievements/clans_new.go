package render

import (
	"fmt"
	"image"
	"log"
	"sync"
	"time"

	dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	mongodbapi "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	dbGloss "github.com/cufee/am-stats/mongodbapi/v1/glossary"
	"github.com/cufee/am-stats/render"
	"github.com/fogleman/gg"
)

// ClansAchievementsLbImage -
func ClansAchievementsLbImage(data []dbAch.ClanAchievements, checkData dbAch.ClanAchievements, bgImage image.Image, medals []mongodbapi.MedalWeight) (finalImage image.Image, err error) {
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
	var maxTimestamp time.Time
	var maxClanTagWidth float64
	var maxPositionWidth float64
	checkCtx := gg.NewContext(1, 1)
	checkBlock := cardBlockData(slimBlockBP)
	maxScoreWidth, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.SmallTextSize, "Score")
	maxClanPlayersWidth, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.SmallTextSize, "Players")

	// Prep block blueprints
	for i, clan := range data {
		// Fix player clan tag
		clan.ClanTag = fmt.Sprintf("[%s]", clan.ClanTag)
		position := i + 1
		if clan.ClanID == checkData.ClanID {
			position = checkData.Position
		}
		data[i].Position = position

		// Get text size
		cW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, clan.ClanTag)
		pW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, fmt.Sprintf("#%v", position))
		sW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.BigTextSize, fmt.Sprint(clan.Score))
		mW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.BigTextSize, fmt.Sprint(clan.Members))
		if cW > maxClanTagWidth { // Check clan tag width
			maxClanTagWidth = cW
		}
		if pW > maxPositionWidth { // Check position width
			maxPositionWidth = pW
		}
		if sW > maxScoreWidth { // Check score width
			maxScoreWidth = sW
		}
		if mW > maxClanPlayersWidth { // Check score width
			maxClanPlayersWidth = mW
		}

		// Get max timestamp
		if clan.Timestamp.After(maxTimestamp) {
			maxTimestamp = clan.Timestamp
		}
	}

	// Work on cards in go routines
	for i, clan := range data {
		wg.Add(1)

		go func(clan dbAch.ClanAchievements, i int) {
			defer wg.Done()

			// Pre card blocks
			var card render.CardData
			card.FrameMargin = render.FrameMargin / 2
			blueprint := slimBlockBP
			// Position
			var posBlock render.Block
			posBlock.Width = int(maxPositionWidth + blueprint.TextMargin/2)
			// Prep extra block data
			posExtra := cardBlockData(blueprint)
			posExtra.BigText = fmt.Sprintf("#%v", clan.Position)
			posExtra.BigTextSize = blueprint.TextSize
			posExtra.BigTextColor = blueprint.SmallTextColor
			posExtra.TextAlign = -1
			posBlock.Extra = &posExtra
			card.Blocks = append(card.Blocks, posBlock)

			// ClanTag
			var tagBlock render.Block
			tagBlock.Width = int(maxClanTagWidth + blueprint.TextMargin/2)
			// Prep extra block data
			clanExtra := cardBlockData(blueprint)
			if clan.ClanID == checkData.ClanID {
				clanExtra.BigTextColor = render.ProtagonistColor
			}
			clanExtra.BigText = fmt.Sprintf("[%s]", clan.ClanTag)
			clanExtra.BigTextSize = blueprint.TextSize
			tagBlock.Extra = &clanExtra
			card.Blocks = append(card.Blocks, tagBlock)

			// Clan players
			var membersBlock render.Block
			membersBlock.Width = int(maxClanPlayersWidth + blueprint.TextMargin/2)
			// Prep extra block data
			membersExtra := cardBlockData(blueprint)
			membersExtra.BigText = fmt.Sprintf("%v", clan.Members)
			membersExtra.SmallText = "Players"
			membersExtra.BigTextSize = blueprint.TextSize
			membersExtra.BigTextColor = blueprint.SmallTextColor
			membersBlock.Extra = &membersExtra
			card.Blocks = append(card.Blocks, membersBlock)

			// Score
			var scoreBlock render.Block
			scoreBlock.Width = int(maxScoreWidth + blueprint.TextMargin/2)
			// Prep extra block data
			scoreExtra := cardBlockData(blueprint)
			scoreExtra.BigText = fmt.Sprint(clan.Score)
			scoreExtra.SmallText = "Score"
			scoreBlock.Extra = &scoreExtra
			card.Blocks = append(card.Blocks, scoreBlock)

			// Fill medal scores and blocks
			for _, m := range medals {
				var medalBlock render.Block
				medalBlock.Width = int(blueprint.IconSize)
				medalExtra := cardBlockData(blueprint)
				// Prep extra block data
				medalExtra.AltText = fmt.Sprint(getField(clan.Data, m.Name))
				medalExtra.AltTextColor = blueprint.SmallTextColor
				medalExtra.IconURL = m.IconURL
				medalExtra.TextAlign = 1
				medalBlock.Extra = &medalExtra
				card.Blocks = append(card.Blocks, medalBlock)
			}

			// Prep card context
			render.PrepNewCard(&card, 1, 0.5, 0)
			card.Index = i
			if err := renderCardBlocks(&card, i, clan.Medals); err != nil {
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

	header := fmt.Sprintf("Achievements Leaderboard | Updated %v min ago", int(time.Now().Sub(maxTimestamp).Minutes()))
	finalCtx, err := render.AddAllCardsToFrame(finalCards, header, bgImage)
	if err != nil {
		return nil, err
	}
	return finalCtx.Image(), err
}
