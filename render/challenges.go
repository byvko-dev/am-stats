package render

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/fogleman/gg"
)

func getPlayerChall(pid int) (blocks []challengeBlock, err error) {
	count := 3

	var block challengeBlock
	block.width = (baseCardWidth - (int(textMargin) * 8)) / 3
	block.height = block.width
	block.height = int(float64(block.width) / 1.5)
	block.isColored = true
	block.hasIcon = false
	block.color = color.RGBA{255, 255, 100, 30}
	block.premiumColor = color.RGBA{255, 255, 0, 30}
	block.shortTxtColor = color.RGBA{255, 255, 255, 255}
	block.longTxtColor = color.RGBA{255, 255, 255, 200}
	block.altTxtColor = color.RGBA{255, 255, 255, 100}
	block.textSize = fontSize

	for i := 0; i <= count; i++ {
		thisBlock := challengeBlock(block)
		thisBlock.status = "complete"
		thisBlock.score = 1.1
		thisBlock.position = i
		thisBlock.prizeTxt = "#1 - $10, #2 - $5"
		thisBlock.shortTxt = "Kill Enemy Tanks"
		thisBlock.longText = "Kill enemy tanks to win and do a lot more of other stuff"
		if i == 1 {
			thisBlock.isPremium = true
			thisBlock.isLocked = true
		}
		if i == 2 {
			thisBlock.isPremium = true
			thisBlock.isLocked = false
		}

		thisBlock, _ = addChallengeCtx(thisBlock)
		blocks = append(blocks, thisBlock)
	}
	return blocks, nil
}

func addChallengeCtx(block challengeBlock) (challengeBlock, error) {
	ctx := gg.NewContext(block.width, block.height)

	// Color is requested
	if block.isColored == true {
		ctx.SetColor(block.color)
		if block.isPremium {
			ctx.SetColor(block.premiumColor)
		}
		ctx.DrawRoundedRectangle(0, 0, float64(block.width), float64(block.height), fontSize)
		ctx.Fill()
	}
	if block.textCoeff == 0 {
		block.textCoeff = 0.6
	}
	topRow := block.textSize + textMargin
	scoreWidth := block.textSize * 3

	if block.isLocked {
		// Locked text
		if err := ctx.LoadFontFace(fontPath, (block.textSize)); err != nil {
			return block, err
		}
		ctx.SetColor(block.shortTxtColor)
		t1 := "Locked"
		t2 := "You need to be a premium member to participate"
		tW, tH := ctx.MeasureString(t1)
		tX := (float64(block.width) - tW) / 2
		tY := (topRow-tH)/2 + tH
		ctx.DrawString(t1, tX, tY)

		// Bottom text
		ctx.SetColor(block.longTxtColor)
		if err := ctx.LoadFontFace(fontPath, (block.textSize * 0.75)); err != nil {
			return block, err
		}
		t2X := textMargin
		t2Y := topRow + textMargin
		ctx.DrawStringWrapped(t2, t2X, t2Y, 0.0, 0.0, (float64(block.width) - textMargin*2), 1.5, gg.AlignLeft)

		// Draw prize string
		if block.prizeTxt != "" {
			ctx.SetColor(block.longTxtColor)
			prize := ("Prize: " + block.prizeTxt)
			pW, _ := ctx.MeasureString(prize)
			pX := (float64(block.width) - textMargin*2 - pW) / 2
			pY := float64(block.height) - textMargin
			ctx.DrawString(prize, pX, pY)
		}

		// return
		block.context = ctx
		return block, nil
	}

	// Left side - Icon, status, score and position
	// Top Row - Icon and short text
	if err := ctx.LoadFontFace(fontPath, (block.textSize)); err != nil {
		return block, err
	}
	ctx.SetColor(block.shortTxtColor)
	tW, tH := ctx.MeasureString(block.shortTxt)
	tX := (float64(block.width) - tW) / 2
	tY := (topRow-tH)/2 + tH
	ctx.DrawString(block.shortTxt, tX, tY)

	// Icon
	if block.hasIcon {
		var iconColor color.RGBA
		switch block.position {
		case 1:
			iconColor = color.RGBA{0, 220, 0, 100}
		case 0:
			iconColor = color.RGBA{200, 200, 200, 100}
		default:
			iconColor = color.RGBA{240, 240, 0, 255}
		}
		ctx.SetColor(iconColor)
		iR := block.textSize / 2.5
		iY := tY - (tY-(iR*2))/2 - iR/2
		iX := iY
		ctx.DrawCircle(iX, iY, iR)
		ctx.Fill()
	}

	// Bottom row - Left side
	// Leaderboard position
	if block.position > 0 {
		ctx.SetColor(block.shortTxtColor)
		if err := ctx.LoadFontFace(fontPath, (block.textSize)); err != nil {
			return block, err
		}
		lbPos := "#" + strconv.Itoa(block.position)
		lbW, _ := ctx.MeasureString(lbPos)
		lbX := (float64(scoreWidth) - lbW) / 2
		lbY := float64(block.height) - textMargin
		ctx.DrawString(lbPos, lbX, lbY)
	}

	// Score alt text
	if err := ctx.LoadFontFace(fontPath, ((block.textSize) * block.textCoeff)); err != nil {
		return block, err
	}
	ctx.SetColor(block.altTxtColor)
	saW, saH := ctx.MeasureString("Score")
	saX := (float64(scoreWidth) - saW) / 2
	saY := topRow + textMargin + saH
	ctx.DrawString("Score", saX, saY)
	ctx.SetColor(block.shortTxtColor)

	// Score
	score := fmt.Sprintf("%.1f", block.score)
	if block.position == 0 {
		score = "-"
	}
	if float64(int(block.score)) == block.score {
		score = fmt.Sprintf("%.0f", block.score)
	}
	scoreCoeff := 1.3 - (float64(len(score)) * 0.1)

	if err := ctx.LoadFontFace(fontPath, (block.textSize * scoreCoeff)); err != nil {
		return block, err
	}
	sW, _ := ctx.MeasureString(score)
	sX := (float64(scoreWidth) - sW) / 2
	sY := saY + (float64(block.height)-textMargin-saY)/2 // Center betweeen position and score alt text
	ctx.DrawString(score, sX, sY)

	// Bottom row - Right side
	rightWidth := float64(block.width) - (scoreWidth)
	// Draw challenge description
	if err := ctx.LoadFontFace(fontPath, (block.textSize * 0.75)); err != nil {
		return block, err
	}
	ctx.SetColor(block.longTxtColor)
	ltX := scoreWidth
	ltY := topRow + textMargin
	ctx.DrawStringWrapped(block.longText, ltX, ltY, 0.0, 0.0, (rightWidth - textMargin), 1.5, gg.AlignLeft)

	// Draw prize string
	if block.prizeTxt != "" {
		ctx.SetColor(block.longTxtColor)
		pX := scoreWidth
		pY := float64(block.height) - textMargin
		ctx.DrawString(block.prizeTxt, pX, pY)
	}

	ctx.SavePNG("test.png")

	// Return new context
	block.context = ctx
	return block, nil
}
