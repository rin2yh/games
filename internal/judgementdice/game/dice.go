package game

import (
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type state int

const (
	stateIdle state = iota
	stateRolling
	stateResult
)

const (
	diceSize     float32 = 120
	rollDuration         = 90
	flipInterval         = 5
	bounceFrames         = 30
	pipRadius    float32 = 9
)

var (
	diceBodyColor   = color.RGBA{R: 0xf5, G: 0xf1, B: 0xe3, A: 0xff}
	diceAccentColor = color.RGBA{R: 0x2b, G: 0x1a, B: 0x4a, A: 0xff}

	pipOffsets = [7][][2]float32{
		1: {{diceSize / 2, diceSize / 2}},
		2: {{diceSize / 4, diceSize / 4}, {diceSize * 3 / 4, diceSize * 3 / 4}},
		3: {{diceSize / 4, diceSize / 4}, {diceSize / 2, diceSize / 2}, {diceSize * 3 / 4, diceSize * 3 / 4}},
		4: {{diceSize / 4, diceSize / 4}, {diceSize * 3 / 4, diceSize / 4}, {diceSize / 4, diceSize * 3 / 4}, {diceSize * 3 / 4, diceSize * 3 / 4}},
		5: {{diceSize / 4, diceSize / 4}, {diceSize * 3 / 4, diceSize / 4}, {diceSize / 2, diceSize / 2}, {diceSize / 4, diceSize * 3 / 4}, {diceSize * 3 / 4, diceSize * 3 / 4}},
		6: {{diceSize / 4, diceSize / 4}, {diceSize * 3 / 4, diceSize / 4}, {diceSize / 4, diceSize / 2}, {diceSize * 3 / 4, diceSize / 2}, {diceSize / 4, diceSize * 3 / 4}, {diceSize * 3 / 4, diceSize * 3 / 4}},
	}
)

type Dice struct {
	state     state
	face      int
	final     int
	tickCount int
}

func NewDice() *Dice {
	return &Dice{state: stateIdle, face: 1}
}

func (d *Dice) State() state { return d.state }
func (d *Dice) Final() int   { return d.final }

func (d *Dice) Roll() {
	if d.state == stateRolling {
		return
	}
	d.state = stateRolling
	d.tickCount = 0
	d.face = rand.IntN(6) + 1
}

func (d *Dice) Reset() {
	d.state = stateIdle
	d.tickCount = 0
	d.final = 0
	d.face = 1
}

func (d *Dice) Update() {
	switch d.state {
	case stateRolling:
		d.tickCount++
		if d.tickCount%flipInterval == 0 {
			next := rand.IntN(5) + 1
			if next >= d.face {
				next++
			}
			d.face = next
		}
		if d.tickCount >= rollDuration {
			d.final = rand.IntN(6) + 1
			d.face = d.final
			d.state = stateResult
			d.tickCount = 0
		}
	case stateResult:
		if d.tickCount < bounceFrames {
			d.tickCount++
		}
	}
}

func (d *Dice) Draw(screen *ebiten.Image, cx, cy float32) {
	var offsetY float32
	switch d.state {
	case stateResult:
		if d.tickCount < bounceFrames {
			phase := float64(d.tickCount) / float64(bounceFrames)
			offsetY = float32(-math.Sin(phase*math.Pi) * 12 * (1 - phase))
		}
	case stateRolling:
		offsetY = float32(math.Sin(float64(d.tickCount)*0.6) * 3)
	}

	x := cx - diceSize/2
	y := cy - diceSize/2 + offsetY

	vector.FillRect(screen, x, y, diceSize, diceSize, diceBodyColor, true)
	vector.StrokeRect(screen, x, y, diceSize, diceSize, 3, diceAccentColor, true)

	for _, p := range pipOffsets[d.face] {
		vector.FillCircle(screen, x+p[0], y+p[1], pipRadius, diceAccentColor, true)
	}
}
