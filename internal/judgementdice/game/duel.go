package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func (d *Duel) Update() {
	if d.phase != duelIdle {
		d.cpu.Update()
		d.player.Update()
	}
	switch d.phase {
	case duelIdle:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			d.cpu.Roll()
			d.player.Roll()
			d.phase = duelRolling
		}
	case duelRolling:
		if d.cpu.State() == stateResult && d.player.State() == stateResult {
			d.phase = duelJudgement
		}
	case duelJudgement:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			d.Reset()
		}
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
		drawText(screen, face, "SPACE でスタート", screenWidth/2, 390, textColor)
	case duelRolling:
		drawText(screen, face, "サイコロが回っている……", screenWidth/2, 330, textColor)
	case duelJudgement:
		drawText(screen, face, fmt.Sprintf("CPU：%d  /  YOU：%d", d.cpu.Final(), d.player.Final()), screenWidth/2, 330, titleColor)
		drawText(screen, face, d.Result(), screenWidth/2, 370, titleColor)
		drawText(screen, face, "SPACE でもう一勝負", screenWidth/2, 400, textColor)
	}
}

func (d *Duel) Result() string {
	if d.player.Final() < d.cpu.Final() {
		return "あんたの勝ち！"
	}
	return "あんたの負け…"
}
