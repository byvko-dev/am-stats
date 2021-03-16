package render

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"
	"sync"

	replays "github.com/cufee/am-stats/dataprep/replays"
	"github.com/cufee/am-stats/render"
	stats "github.com/cufee/am-stats/render/stats"
	"github.com/fogleman/gg"
)

// Render -
func Render(replay replays.ReplaySummary, bgImage image.Image) (image.Image, error) {
	// Init
	var finalCards render.AllCards
	cardsChan := make(chan render.CardData, (2 + len(replay.Details)))
	var wg sync.WaitGroup

	// Configure blueprint for cards
	var slimBlockBP replayBlockData
	slimBlockBP.Defaults()
	slimBlockBP.TextSize = (slimBlockBP.TextSize * 0.90)

	// Max values
	checkCtx := gg.NewContext(1, 1)
	checkBlock := replayBlockData(slimBlockBP)

	var addPlatoon bool
	const hpBarWidth int = 8
	const platoonWidth int = 25
	maxNameLength, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, "Supremacy Points - 000")
	var maxKillsLength float64
	var maxRatingLength float64
	var maxDamageLength float64
	var maxAssistLength float64
	var maxWinrateLength float64
	teamWins := make(map[int]int)
	teamPoints := make(map[int]int)
	teamRating := make(map[int]int)
	teamDamage := make(map[int]int)
	teamBattles := make(map[int]int)

	// Calculate max length for data
	for i, player := range replay.Details {
		// Check if need platoon icon spacing
		if player.SquadIndex > 0 && !addPlatoon && !intInSlice(SpecialGameModeTypes, replay.RoomType) {
			addPlatoon = true
		}

		// Team points
		teamPoints[player.Team] += player.WpPointsEarned
		if player.Team == 1 {
			teamPoints[2] += player.WpPointsStolen
		} else {
			teamPoints[1] += player.WpPointsStolen
		}

		// Team totals
		teamRating[player.Team] += player.TankProfile.TankWN8
		teamWins[player.Team] += player.Profile.Stats.All.Wins
		teamBattles[player.Team] += player.TankProfile.Battles
		teamDamage[player.Team] += player.TankProfile.DamageDealt
		teamBattles[player.Team] += player.Profile.Stats.All.Battles

		// Compile clan tag
		var clanTag string
		if player.ClanTag != "" {
			clanTag = fmt.Sprintf(" [%v]", player.ClanTag)
		}
		// Compile player name and tank
		replay.Details[i].Profile.Name = player.Profile.Name + clanTag

		// Check maximum length
		nW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, player.Profile.Name+clanTag)
		tNW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, player.TankProfile.TankName)
		if nW > maxNameLength { // name
			maxNameLength = nW
		}
		if tNW > maxNameLength { // tank
			maxNameLength = tNW
		}
		kW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, fmt.Sprint(player.EnemiesDestroyed))
		if kW > maxKillsLength { // check kills
			maxKillsLength = kW
		}
		rW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, fmt.Sprint(player.TankProfile.TankWN8))
		if rW > maxRatingLength { // check rating
			maxRatingLength = rW
		}
		dW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, fmt.Sprint(player.DamageMade))
		if dW > maxDamageLength { // check damage
			maxDamageLength = dW
		}
		aW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, fmt.Sprint(player.DamageAssisted))
		if aW > maxAssistLength { // assisted
			maxAssistLength = aW
		}
		wW, _, _ := getTextParams(checkCtx, &checkBlock, slimBlockBP.TextSize, (fmt.Sprintf("%.2f", ((float64(player.Profile.Stats.All.Wins)/float64(player.Profile.Stats.All.Battles))*100)) + "%"))
		if wW > maxWinrateLength { // check winrate
			maxWinrateLength = wW
		}
	}

	// Work on cards in go routines
	for i, player := range replay.Details {
		wg.Add(1)

		go func(player replays.ReplayPlayerData, i int) {
			defer wg.Done()

			// Prep card blocks
			var card render.CardData
			card.FrameMargin = render.FrameMargin / 2
			if player.Team == 2 {
				card.IndexX = 1
			}
			blueprint := slimBlockBP

			// Platoon icon
			if addPlatoon {
				var platoonBlock render.Block
				platoonBlock.Width = platoonWidth
				if player.SquadIndex != 0 {
					// Prep extra block data
					platoonExtra := replayBlockData(blueprint)
					platoonExtra.TextLines = append(platoonExtra.TextLines, blockTextLine{Text: fmt.Sprint(player.SquadIndex), Color: render.AltTextColor})
					platoonExtra.TextSize = blueprint.TextSize
					platoonExtra.TextAlign = -1
					platoonBlock.Extra = &platoonExtra
				}
				card.Blocks = append(card.Blocks, platoonBlock)
			}

			// Add HP bar
			var hpBarBlock render.Block
			hpBarBlock.Width = hpBarWidth
			// Prep extra block data
			var hpBarExtra hpBarBlockData
			hpBarBlock.ExtraType = "hpbar"
			hpBarExtra.HPColorBG = color.RGBA{65, 65, 65, 220}
			if player.KilledBy == 0 {
				hpBarExtra.PercentHP = float64(player.HitpointsLeft) / float64(player.HitpointsLeft+player.DamageReceived)
				hpBarExtra.HPColor = color.RGBA{123, 219, 101, 220}
				hpBarExtra.HPColorBG = color.RGBA{80, 80, 80, 220}
				if hpBarExtra.PercentHP < 0.1 {
					hpBarExtra.PercentHP = 0.1
				}
			}
			if player.Team != 1 {
				hpBarExtra.HPColor = color.RGBA{219, 109, 101, 220}
			}
			hpBarExtra.Margin = int(render.TextMargin)
			hpBarBlock.Extra = &hpBarExtra
			card.Blocks = append(card.Blocks, hpBarBlock)

			// Add player and tank name
			var playerNameBlock render.Block
			playerNameBlock.Width = int(maxNameLength)
			playerNameBlock.Padding = int(render.TextMargin * 3)
			// Prep extra block data
			playerNameExtra := replayBlockData(blueprint)
			var name blockTextLine
			name.Text = player.Profile.Name
			if player.IsProtagonist { // Set color for protagonist
				name.Color = render.ProtagonistColor
			}
			playerNameExtra.TextLines = append(playerNameExtra.TextLines, name)
			playerNameExtra.TextLines = append(playerNameExtra.TextLines, blockTextLine{Text: player.TankProfile.TankName, Color: render.SmallTextColor})
			playerNameExtra.TextAlign = -1
			playerNameExtra.TextSize = blueprint.TextSize
			playerNameBlock.Extra = &playerNameExtra
			card.Blocks = append(card.Blocks, playerNameBlock)

			// Add rating color bar
			var ratingBar render.Block
			ratingBar.Width = hpBarWidth
			// Prep extra block data
			var ratingBarExtra hpBarBlockData
			ratingBar.ExtraType = "hpbar"
			ratingBarExtra.PercentHP = 1
			ratingBarExtra.HPColor = stats.GetRatingColor(player.TankProfile.TankWN8)
			ratingBarExtra.Margin = int(render.TextMargin * 2)
			ratingBar.Extra = &ratingBarExtra
			card.Blocks = append(card.Blocks, ratingBar)

			// Add rating value
			var ratingBlock render.Block
			ratingBlock.Width = int(maxRatingLength)
			ratingBlock.Padding = hpBarWidth * 2
			// Prep extra block data
			ratingBlockExtra := replayBlockData(blueprint)
			rating := "-"
			if player.TankProfile.TankWN8 >= 0 {
				rating = fmt.Sprint(player.TankProfile.TankWN8)
			}
			ratingBlockExtra.TextLines = append(ratingBlockExtra.TextLines, blockTextLine{Text: rating})
			ratingBlock.Extra = &ratingBlockExtra
			card.Blocks = append(card.Blocks, ratingBlock)

			// Add windrate
			var winrateBlock render.Block
			winrateBlock.Width = int(maxWinrateLength)
			winrateBlock.Padding = int(render.TextMargin)
			// Prep extra block data
			winrateBlockExtra := replayBlockData(blueprint)
			winrateBlockExtra.TextLines = append(winrateBlockExtra.TextLines, blockTextLine{Text: (fmt.Sprintf("%.2f", ((float64(player.Profile.Stats.All.Wins)/float64(player.Profile.Stats.All.Battles))*100)) + "%"), Color: render.SmallTextColor})
			winrateBlock.Extra = &winrateBlockExtra
			card.Blocks = append(card.Blocks, winrateBlock)

			// Add damage
			var damageBlock render.Block
			damageBlock.Width = int(maxDamageLength)
			damageBlock.Padding = int(render.TextMargin)
			// Prep extra block data
			damageBlockExtra := replayBlockData(blueprint)
			damageBlockExtra.TextLines = append(damageBlockExtra.TextLines, blockTextLine{Text: fmt.Sprint(player.DamageMade)})
			damageBlock.Extra = &damageBlockExtra
			card.Blocks = append(card.Blocks, damageBlock)

			// Add kills
			var killsBlock render.Block
			killsBlock.Width = int(maxKillsLength)
			killsBlock.Padding = int(render.TextMargin)
			// Prep extra block data
			killsBlockExtra := replayBlockData(blueprint)
			kills := "-"
			if player.EnemiesDestroyed > 0 {
				kills = fmt.Sprint(player.EnemiesDestroyed)
			}
			killsBlockExtra.TextLines = append(killsBlockExtra.TextLines, blockTextLine{Text: kills, Color: render.SmallTextColor})
			killsBlock.Extra = &killsBlockExtra
			card.Blocks = append(card.Blocks, killsBlock)

			// Prep card context
			render.PrepNewCard(&card, 1, 0.5, 0)
			card.Index = i + 10
			if err := renderCardBlocks(&card); err != nil {
				log.Println(err)
				return
			}

			cardsChan <- card
			return
		}(player, i)
	}
	wg.Wait()
	close(cardsChan)

	// Append all cards
	for c := range cardsChan {
		finalCards.Cards = append(finalCards.Cards, c)
	}

	// Add header card
	var card render.CardData
	card.FrameMargin = render.FrameMargin / 2
	card.IndexX = 0
	card.Type = render.CardTypeHeader
	blueprint := slimBlockBP

	// Add Team 1 averages
	var blankBlock render.Block
	blankBlock.Width = int(maxNameLength) + (hpBarWidth * 2) + int(render.TextMargin*3)
	if addPlatoon {
		blankBlock.Width += platoonWidth
	}
	if replay.BattleType == battlesTypeSupremacy {
		// Prep extra block data
		blankBlockExtra := replayBlockData(blueprint)
		blankBlockExtra.TextLines = append(blankBlockExtra.TextLines, blockTextLine{Text: fmt.Sprintf("Supremacy Points - %v", ((teamPoints[1]+49)/50)*50), Color: render.SmallTextColor})
		blankBlockExtra.TextAlign = -1
		blankBlock.Extra = &blankBlockExtra
	}
	card.Blocks = append(card.Blocks, blankBlock)

	// Add rating
	var ratingBlock1 render.Block
	ratingBlock1.Width = int(maxRatingLength) + (hpBarWidth * 2)
	// Prep extra block data
	ratingBlock1Extra := replayBlockData(blueprint)
	ratingBlock1Extra.TextLines = append(ratingBlock1Extra.TextLines, blockTextLine{Text: "WN8", Color: render.AltTextColor})
	// ratingBlock1Extra.TextLines = append(ratingBlock1Extra.TextLines, blockTextLine{Text: fmt.Sprint(teamRating[1] / len(replay.Allies)), Color: render.BigTextColor})
	ratingBlock1.Extra = &ratingBlock1Extra
	card.Blocks = append(card.Blocks, ratingBlock1)

	// Add winrate
	var winrateBlock1 render.Block
	winrateBlock1.Width = int(maxWinrateLength + render.TextMargin)
	// Prep extra block data
	winrateBlock1Extra := replayBlockData(blueprint)
	winrateBlock1Extra.TextLines = append(winrateBlock1Extra.TextLines, blockTextLine{Text: "WR", Color: render.AltTextColor})
	// winrateBlock1Extra.TextLines = append(winrateBlock1Extra.TextLines, blockTextLine{Text: fmt.Sprintf("%.2f", 100*float64(teamWins[1])/float64(teamBattles[1])) + "%", Color: render.SmallTextColor})
	winrateBlock1.Extra = &winrateBlock1Extra
	card.Blocks = append(card.Blocks, winrateBlock1)

	// Add damage
	var damageBlock1 render.Block
	damageBlock1.Width = int(maxDamageLength + render.TextMargin)
	// Prep extra block data
	damageBlock1Extra := replayBlockData(blueprint)
	damageBlock1Extra.TextLines = append(damageBlock1Extra.TextLines, blockTextLine{Text: "DMG", Color: render.AltTextColor})
	// damageBlock1Extra.TextLines = append(damageBlock1Extra.TextLines, blockTextLine{Text: fmt.Sprint(teamDamage[1] / len(replay.Allies)), Color: render.BigTextColor})
	damageBlock1.Extra = &damageBlock1Extra
	card.Blocks = append(card.Blocks, damageBlock1)

	// Add kills
	var killsBlock1 render.Block
	killsBlock1.Width = int(maxKillsLength + render.TextMargin)
	// Prep extra block data
	killsBlock1Extra := replayBlockData(blueprint)
	killsBlock1Extra.TextLines = append(killsBlock1Extra.TextLines, blockTextLine{Text: "K", Color: render.AltTextColor})
	// killsBlock1Extra.TextLines = append(killsBlock1Extra.TextLines, blockTextLine{Text: " ", Color: render.SmallTextColor})
	killsBlock1.Extra = &killsBlock1Extra
	card.Blocks = append(card.Blocks, killsBlock1)

	// Spacing block
	var spacingBlock render.Block
	spacingBlock.Width = int(render.FrameMargin * 3 / 2)
	spacingBlockExtra := replayBlockData(blueprint)
	spacingBlockExtra.TextLines = append(spacingBlockExtra.TextLines, blockTextLine{Text: "|", Color: render.AltTextColor, TextScale: 1.5})
	spacingBlock.Extra = &spacingBlockExtra
	card.Blocks = append(card.Blocks, spacingBlock)

	// Add Team 2 averages
	blankBlock2 := render.Block(blankBlock)
	if replay.BattleType == battlesTypeSupremacy {
		blankBlockExtra2 := replayBlockData(blueprint)
		blankBlockExtra2.TextLines = []blockTextLine{{Text: fmt.Sprintf("Supremacy Points - %v", ((teamPoints[2]+49)/50)*50), Color: render.SmallTextColor}}
		blankBlockExtra2.TextAlign = -1
		blankBlock2.Extra = &blankBlockExtra2
	}
	card.Blocks = append(card.Blocks, blankBlock2)

	// Add rating
	var ratingBlock2 render.Block
	ratingBlock2.Width = int(maxRatingLength) + (hpBarWidth * 2)
	ratingBlock2Extra := replayBlockData(blueprint)
	ratingBlock2Extra.TextLines = append(ratingBlock2Extra.TextLines, blockTextLine{Text: "WN8", Color: render.AltTextColor})
	// ratingBlock2Extra.TextLines = append(ratingBlock2Extra.TextLines, blockTextLine{Text: fmt.Sprint(math.Round(float64(teamRating[2] / len(replay.Enemies)))), Color: render.BigTextColor})
	ratingBlock2.Extra = &ratingBlock2Extra
	card.Blocks = append(card.Blocks, ratingBlock2)

	// Add winrate
	var winrateBlock2 render.Block
	winrateBlock2.Width = int(maxWinrateLength + render.TextMargin)
	// Prep extra block data
	winrateBlock2Extra := replayBlockData(blueprint)
	winrateBlock2Extra.TextLines = append(winrateBlock2Extra.TextLines, blockTextLine{Text: "WR", Color: render.AltTextColor})
	// winrateBlock2Extra.TextLines = append(winrateBlock2Extra.TextLines, blockTextLine{Text: fmt.Sprintf("%.2f", 100*float64(teamWins[2])/float64(teamBattles[2])) + "%", Color: render.SmallTextColor})
	log.Print(teamWins[2], teamBattles[2])
	winrateBlock2.Extra = &winrateBlock2Extra
	card.Blocks = append(card.Blocks, winrateBlock2)

	// Add damage
	var damageBlock2 render.Block
	damageBlock2.Width = int(maxDamageLength + render.TextMargin)
	// Prep extra block data
	damageBlock2Extra := replayBlockData(blueprint)
	damageBlock2Extra.TextLines = append(damageBlock2Extra.TextLines, blockTextLine{Text: "DMG", Color: render.AltTextColor})
	// damageBlock2Extra.TextLines = append(damageBlock2Extra.TextLines, blockTextLine{Text: fmt.Sprint(teamDamage[2] / len(replay.Enemies)), Color: render.BigTextColor})
	damageBlock2.Extra = &damageBlock2Extra
	card.Blocks = append(card.Blocks, damageBlock2)

	// Add kills
	var killsBlock2 render.Block
	killsBlock2.Width = int(maxKillsLength + render.TextMargin)
	// Prep extra block data
	killsBlock2Extra := replayBlockData(blueprint)
	killsBlock2Extra.TextLines = append(killsBlock2Extra.TextLines, blockTextLine{Text: "K", Color: render.AltTextColor})
	// killsBlock2Extra.TextLines = append(killsBlock2Extra.TextLines, blockTextLine{Text: " ", Color: render.SmallTextColor})
	killsBlock2.Extra = &killsBlock2Extra
	card.Blocks = append(card.Blocks, killsBlock2)

	// Prep card context
	render.PrepNewCard(&card, 1, 0.5, 0)
	card.Index = 0
	card.IndexX = 0
	if err := renderCardBlocks(&card); err != nil {
		log.Println(err)
		return nil, err
	}
	finalCards.Cards = append(finalCards.Cards, card)

	// Battle result
	result := "Victory"
	if replay.WinnerTeam != replay.ProtagonistTeam {
		result = "Defeat"
	}

	// Game mode
	gameMode, ok := GameModeMap[replay.RoomType]
	if !ok {
		gameMode = GameModeMap[0]
	}

	header := strings.Join([]string{replay.MapName, gameMode, result}, " - ")
	finalCtx, err := renderAllCardsOnFrame(finalCards, header, bgImage)
	if err != nil {
		return nil, err
	}
	image := finalCtx.Image()

	return image, err
}
