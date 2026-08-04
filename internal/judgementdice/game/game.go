package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	bgColor    = color.RGBA{R: 0x1b, G: 0x10, B: 0x3a, A: 0xff}
	titleColor = color.RGBA{R: 0xff, G: 0xd7, B: 0x4a, A: 0xff}
	textColor  = color.RGBA{R: 0xf5, G: 0xf1, B: 0xe3, A: 0xff}
	hintColor  = color.RGBA{R: 0xa0, G: 0x94, B: 0xc8, A: 0xff}
)

type scene int

const (
	sceneTitle scene = iota
	sceneEffect
	sceneDuel
	sceneJenga
)

type Game struct {
	scene scene
	dice  *Dice
	duel  *Duel
	jenga *Jenga
	face  text.Face
}

func New() *Game {
	return &Game{
		scene: sceneTitle,
		dice:  NewDice(),
		duel:  NewDuel(),
		jenga: NewJenga(),
		face:  text.NewGoXFace(bitmapfont.FaceEA),
	}
}

func (g *Game) Update() error {
	g.handle(readInput(g.scene))
	return nil
}

// handle は 1 フレーム分の入力を今のシーンに適用する。入力の読み取り
// (readInput) と分けてあるので、Ebitengine を起動せずにテストできる。
func (g *Game) handle(in input) {
	switch g.scene {
	case sceneTitle:
		switch in.mode {
		case 1:
			g.scene = sceneEffect
		case 2:
			g.scene = sceneDuel
		case 3:
			g.scene = sceneJenga
		}
	case sceneEffect:
		if in.back {
			g.dice.Reset()
			g.scene = sceneTitle
			return
		}
		g.updateEffect(in)
	case sceneDuel:
		if in.back {
			g.duel.Reset()
			g.scene = sceneTitle
			return
		}
		g.duel.Update(in)
	case sceneJenga:
		if in.back {
			g.jenga.Reset()
			g.scene = sceneTitle
			return
		}
		g.jenga.Update(in)
	}
}

func (g *Game) updateEffect(in input) {
	switch g.dice.State() {
	case stateIdle:
		if in.confirm {
			g.dice.Roll()
		}
	case stateResult:
		if in.confirm {
			g.dice.Roll()
			return
		}
		if in.reset {
			g.dice.Reset()
			return
		}
		g.dice.Update()
	default:
		g.dice.Update()
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)

	switch g.scene {
	case sceneTitle:
		g.drawTitleScene(screen)
	case sceneEffect:
		g.drawEffectScene(screen)
	case sceneDuel:
		g.drawDuelScene(screen)
	case sceneJenga:
		g.drawJengaScene(screen)
	}
}

func (g *Game) drawTitleScene(screen *ebiten.Image) {
	drawText(screen, g.face, "ジャッジメントダイス", screenWidth/2, 80, titleColor)
	// 選択肢の y は titleItemY（タップの当たり判定）と同じ順・同じ位置。
	drawText(screen, g.face, "1: 遊戯王ジャッジメントダイス", screenWidth/2, titleItemY[0], textColor)
	drawText(screen, g.face, "2: タンクトップ小隊のジャッジメントダイス", screenWidth/2, titleItemY[1], textColor)
	drawText(screen, g.face, "3: タンクトップ小隊ジェンガモード", screenWidth/2, titleItemY[2], textColor)
	drawText(screen, g.face, "1 / 2 / 3 キー、または選択肢をタップ", screenWidth/2, screenHeight-60, hintColor)
}

// drawBackButton は backRect の位置に「← タイトル」を描く。タップで戻れる
// ことが見えていないと、キーの無い端末ではモードから出られない。
func (g *Game) drawBackButton(screen *ebiten.Image) {
	drawText(screen, g.face, "← タイトル", backRect.x+backRect.w/2, backRect.y+backRect.h/2, hintColor)
}

func (g *Game) drawEffectScene(screen *ebiten.Image) {
	drawText(screen, g.face, "ジャッジメントダイス", screenWidth/2, 36, titleColor)
	g.dice.Draw(screen, screenWidth/2, 170)

	switch g.dice.State() {
	case stateResult:
		drawText(screen, g.face, fmt.Sprintf("出目：%d", g.dice.Final()), screenWidth/2, 290, titleColor)
		drawText(screen, g.face, Effects[g.dice.Final()], screenWidth/2, 330, textColor)
	case stateRolling:
		drawText(screen, g.face, "ダイスを振っています……", screenWidth/2, 310, textColor)
	default:
		drawText(screen, g.face, "SPACE またはタップでダイスを振る", screenWidth/2, 310, textColor)
	}

	drawText(screen, g.face, "SPACE / タップ: Roll    R: Reset    ESC: Title", screenWidth/2, screenHeight-24, hintColor)
	g.drawBackButton(screen)
}

func (g *Game) drawDuelScene(screen *ebiten.Image) {
	g.duel.Draw(screen, g.face)
	drawText(screen, g.face, "SPACE / タップ: Next    ESC: Title", screenWidth/2, screenHeight-24, hintColor)
	g.drawBackButton(screen)
}

func (g *Game) drawJengaScene(screen *ebiten.Image) {
	g.jenga.Draw(screen, g.face)
	drawText(screen, g.face, "SPACE / タップ: Roll    ESC: Title", screenWidth/2, screenHeight-24, hintColor)
	g.drawBackButton(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func drawText(screen *ebiten.Image, face text.Face, s string, x, y float64, clr color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	op.LayoutOptions.PrimaryAlign = text.AlignCenter
	op.LayoutOptions.SecondaryAlign = text.AlignCenter
	text.Draw(screen, s, face, op)
}
