package game

import (
	"github.com/hajimehoshi/ebiten/v2"
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

// Update は毎フレームの進行。振っている最中のサイコロを進め、止まったら判定に
// 移る。
func (j *Jenga) Update() {
	if j.phase == jengaIdle {
		return
	}
	j.dice.Update()

	if j.phase == jengaRolling && j.dice.State() == stateResult {
		j.phase = jengaJudgement
	}
}

// Confirm は SPACE / タップに当たる 1 回きりの操作。振っている最中は何も
// 起きない。
func (j *Jenga) Confirm() {
	switch j.phase {
	case jengaIdle:
		j.dice.Roll()
		j.phase = jengaRolling
	case jengaJudgement:
		j.Reset()
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
		drawText(screen, face, "SPACE またはタップでジャッジメントダイスを振れ", screenWidth/2, 400, hintColor)
	case jengaRolling:
		drawText(screen, face, "ジャッジメントダイス！ダイスを振り．．．", screenWidth/2, 330, textColor)
	case jengaJudgement:
		if j.dice.Final() == 1 {
			drawText(screen, face, "くっ．．．俺のターン．．．", screenWidth/2, 330, titleColor)
		} else {
			drawText(screen, face, "よし！！お前のターン！", screenWidth/2, 330, titleColor)
		}
		drawText(screen, face, "SPACE またはタップでもう一度", screenWidth/2, 400, hintColor)
	}
}
