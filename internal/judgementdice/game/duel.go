package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type duelPhase int

const (
	duelIdle duelPhase = iota
	duelRolling
	duelJudgement
)

type Duel struct {
	phase  duelPhase
	cpu    *Dice
	player *Dice
}

func NewDuel() *Duel {
	return &Duel{
		phase:  duelIdle,
		cpu:    NewDice(),
		player: NewDice(),
	}
}

func (d *Duel) Reset() {
	d.phase = duelIdle
	d.cpu.Reset()
	d.player.Reset()
}

// Update は毎フレームの進行。振っている最中のサイコロを進め、両方止まったら
// 判定に移る。
func (d *Duel) Update() {
	if d.phase == duelIdle {
		return
	}
	d.cpu.Update()
	d.player.Update()

	if d.phase == duelRolling && d.cpu.State() == stateResult && d.player.State() == stateResult {
		d.phase = duelJudgement
	}
}

// Confirm は SPACE / タップに当たる 1 回きりの操作。振っている最中は何も
// 起きない。
func (d *Duel) Confirm() {
	switch d.phase {
	case duelIdle:
		d.cpu.Roll()
		d.player.Roll()
		d.phase = duelRolling
	case duelJudgement:
		d.Reset()
	}
}

func (d *Duel) Draw(screen *ebiten.Image, face text.Face) {
	const cpuX, playerX float32 = 170, 470
	const diceY float32 = 200

	drawText(screen, face, "タンクトップ小隊", screenWidth/2, 36, titleColor)
	drawText(screen, face, "CPU", float64(cpuX), 110, textColor)
	drawText(screen, face, "YOU", float64(playerX), 110, textColor)

	d.cpu.Draw(screen, cpuX, diceY)
	d.player.Draw(screen, playerX, diceY)

	switch d.phase {
	case duelIdle:
		drawText(screen, face, "偽遊戯：ルールは単純！　二人同時にサイコロを振る", screenWidth/2, 330, textColor)
		drawText(screen, face, "オレより小さい目を出せばあんたの勝ち！　同じか大きいなら負けだ", screenWidth/2, 360, textColor)
		drawText(screen, face, "SPACE またはタップでスタート", screenWidth/2, 390, textColor)
	case duelRolling:
		drawText(screen, face, "サイコロが回っている……", screenWidth/2, 330, textColor)
	case duelJudgement:
		drawText(screen, face, fmt.Sprintf("CPU：%d  /  YOU：%d", d.cpu.Final(), d.player.Final()), screenWidth/2, 330, titleColor)
		drawText(screen, face, d.Result(), screenWidth/2, 370, titleColor)
		drawText(screen, face, "SPACE またはタップでもう一勝負", screenWidth/2, 400, textColor)
	}
}

func (d *Duel) Result() string {
	if d.player.Final() < d.cpu.Final() {
		return "あんたの勝ち！"
	}
	return "あんたの負け…"
}
