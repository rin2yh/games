package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type jengaPhase int

const (
	jengaIdle jengaPhase = iota
	jengaRolling
	jengaJudgement
)

type Jenga struct {
	phase jengaPhase
	dice  *Dice
}

func NewJenga() *Jenga {
	return &Jenga{
		phase: jengaIdle,
		dice:  NewDice(),
	}
}

func (j *Jenga) Reset() {
	j.phase = jengaIdle
	j.dice.Reset()
}

func (j *Jenga) Update() {
	if j.phase != jengaIdle {
		j.dice.Update()
	}
	switch j.phase {
	case jengaIdle:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			j.dice.Roll()
			j.phase = jengaRolling
		}
	case jengaRolling:
		if j.dice.State() == stateResult {
			j.phase = jengaJudgement
		}
	case jengaJudgement:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			j.Reset()
		}
	}
}

func (j *Jenga) Draw(screen *ebiten.Image, face text.Face) {
	drawText(screen, face, "タンクトップ小隊 ジェンガモード", screenWidth/2, 36, titleColor)
	j.dice.Draw(screen, screenWidth/2, 180)

	switch j.phase {
	case jengaIdle:
		drawText(screen, face, "ジャッジメントダイス降臨．", screenWidth/2, 300, textColor)
		drawText(screen, face, "一が出れば俺のターン．", screenWidth/2, 330, textColor)
		drawText(screen, face, "二から六が出れば．お前のターンとなるぜ．", screenWidth/2, 360, textColor)
		drawText(screen, face, "SPACE でジャッジメントダイスを振れ", screenWidth/2, 400, hintColor)
	case jengaRolling:
		drawText(screen, face, "ジャッジメントダイス！ダイスを振り．．．", screenWidth/2, 330, textColor)
	case jengaJudgement:
		if j.dice.Final() == 1 {
			drawText(screen, face, "くっ．．．俺のターン．．．", screenWidth/2, 330, titleColor)
		} else {
			drawText(screen, face, "よし！！お前のターン！", screenWidth/2, 330, titleColor)
		}
		drawText(screen, face, "SPACE でもう一度", screenWidth/2, 400, hintColor)
	}
}
