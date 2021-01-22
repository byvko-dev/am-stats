package render

import (
	"fmt"
	"image"
	"log"
	"sync"

	dataprep "github.com/cufee/am-stats/dataprep/achievements"
	dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	mongodbapi "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	dbGloss "github.com/cufee/am-stats/mongodbapi/v1/glossary"
	"github.com/cufee/am-stats/render"
	"github.com/fogleman/gg"
)

// PlayerAchievementsLbImage -
func PlayerAchievementsLbImage(data []dbAch.AchievementsPlayerData, checkData dataprep.AchievementsPIDPos, bgImage image.Image, medals []mongodbapi.MedalWeight) (finalImage image.Image, err error) {
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
	var maxPositionWidth float64
	checkCtx := gg.NewContext(1, 1)
	checkBlock := cardBlockData(slimBlockBP)
	maxScoreWidth, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.SmallTextSize, "Score")

	// Prep block blueprints
	for i, player := range data {
		// Fix player clan tag
		if player.ClanTag != "" {
			player.ClanTag = fmt.Sprintf("[%s]", player.ClanTag)
			data[i].ClanTag = player.ClanTag
		}

		// Get text size
		cW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, player.ClanTag)
		nW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, player.Nickname)
		pW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, fmt.Sprintf("#%v", i+1))
		sW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.BigTextSize, fmt.Sprint(player.Score))
		if cW > maxClanTagWidth { // Check clan tag width
			maxClanTagWidth = cW
		}
		if nW > maxNameWidth { // Check name width
			maxNameWidth = nW
		}
		if pW > maxPositionWidth { // Check position width
			maxPositionWidth = pW
		}
		if sW > maxScoreWidth { // Check score width
			maxScoreWidth = sW
		}
	}

	// Work on cards in go routines
	for i, player := range data {
		wg.Add(1)

		go func(player dbAch.AchievementsPlayerData, i int) {
			defer wg.Done()

			// Pre card blocks
			var card render.CardData
			card.FrameMargin = render.FrameMargin / 2
			blueprint := slimBlockBP
			if player.PID == checkData.PID {
				blueprint = cardBlockData(slimBlockBP)
				blueprint.TextColor = render.ProtagonistColor
			}

			// Position
			var posBlock render.Block
			posBlock.Width = int(maxPositionWidth + blueprint.TextMargin/2)
			// Prep extra block data
			posExtra := cardBlockData(blueprint)
			posExtra.BigText = fmt.Sprintf("#%v", i+1)
			if player.PID == checkData.PID {
				posExtra.BigText = fmt.Sprintf("#%v", checkData.Position)
			}
			posExtra.BigTextSize = blueprint.TextSize
			posExtra.BigTextColor = blueprint.SmallTextColor
			posExtra.TextAlign = -1
			posBlock.Extra = &posExtra
			card.Blocks = append(card.Blocks, posBlock)

			// Name
			var nameBlock render.Block
			nameBlock.Width = int(maxNameWidth + blueprint.TextMargin/2)
			if player.PID == checkData.PID {
				nameBlock.TextColor = render.ProtagonistColor
			}
			// Prep extra block data
			nameExtra := cardBlockData(blueprint)
			nameExtra.BigText = player.Nickname
			nameExtra.BigTextSize = blueprint.TextSize
			nameExtra.BigTextColor = blueprint.TextColor
			nameExtra.TextAlign = -1
			nameBlock.Extra = &nameExtra
			card.Blocks = append(card.Blocks, nameBlock)

			// ClanTag
			var tagBlock render.Block
			tagBlock.Width = int(maxClanTagWidth + blueprint.TextMargin/2)
			// Prep extra block data
			clanExtra := cardBlockData(blueprint)
			clanExtra.BigText = player.ClanTag
			clanExtra.BigTextSize = blueprint.TextSize
			clanExtra.BigTextColor = blueprint.SmallTextColor
			clanExtra.TextAlign = -1
			tagBlock.Extra = &clanExtra
			card.Blocks = append(card.Blocks, tagBlock)

			// Score
			var scoreBlock render.Block
			scoreBlock.Width = int(maxScoreWidth + blueprint.TextMargin/2)
			// Prep extra block data
			scoreExtra := cardBlockData(blueprint)
			scoreExtra.BigText = fmt.Sprint(player.Score)
			scoreExtra.SmallText = "Score"
			scoreBlock.Extra = &scoreExtra
			card.Blocks = append(card.Blocks, scoreBlock)

			// Fill medal scores and blocks
			for _, m := range medals {
				var medalBlock render.Block
				medalBlock.Width = int(blueprint.IconSize) * 3 / 2
				medalExtra := cardBlockData(blueprint)
				// Prep extra block data
				medalExtra.AltText = fmt.Sprint(getField(player.Data, m.Name))
				medalExtra.AltTextColor = blueprint.SmallTextColor
				medalExtra.IconURL = m.IconURL
				medalExtra.TextAlign = 1
				medalBlock.Extra = &medalExtra
				card.Blocks = append(card.Blocks, medalBlock)
			}

			// Prep card context
			render.PrepNewCard(&card, 1, 0.5, 0)
			card.Index = i
			if err := renderCardBlocks(&card, &player, i, player.Medals); err != nil {
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
func renderCardBlocks(card *render.CardData, data *dbAch.AchievementsPlayerData, position int, medals []mongodbapi.MedalWeight) error {
	if data.Score < 1 {
		return fmt.Errorf("score is 0")
	}

	// Atomic counter
	var wg sync.WaitGroup

	for i, block := range card.Blocks {
		if block.Extra == nil {
			continue
		}

		wg.Add(1)
		go func(block render.Block, i int) {
			defer wg.Done()

			// Fill block Width and Height for legacy code
			blockExtra := block.Extra.(*cardBlockData)
			blockExtra.Height = card.Context.Height()
			blockExtra.Width = block.Width

			// Render block image
			if err := renderBlock(blockExtra); err != nil {
				log.Print(err)
				return
			}

			// Calculate rendering offset
			var offset int = card.FrameMargin
			for _, b := range card.Blocks[:i] {
				offset += b.Width
			}

			// Draw block
			card.Context.DrawImage(blockExtra.Context.Image(), offset, 0)
		}(block, i)
	}
	wg.Wait()

	// Render image
	card.Image = card.Context.Image()
	return nil
}
