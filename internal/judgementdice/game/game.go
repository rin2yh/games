package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	switch g.scene {
	case sceneTitle:
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			g.scene = sceneEffect
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			g.scene = sceneDuel
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			g.scene = sceneJenga
		}
	case sceneEffect:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.dice.Reset()
			g.scene = sceneTitle
			return nil
		}
		g.updateEffect()
	case sceneDuel:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.duel.Reset()
			g.scene = sceneTitle
			return nil
		}
		g.duel.Update()
	case sceneJenga:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.jenga.Reset()
			g.scene = sceneTitle
			return nil
		}
		g.jenga.Update()
	}
	return nil
}

func (g *Game) updateEffect() {
	switch g.dice.State() {
	case stateIdle:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.dice.Roll()
		}
	case stateResult:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.dice.Roll()
			return
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
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
	drawText(screen, g.face, "1: 遊戯王ジャッジメントダイス", screenWidth/2, 200, textColor)
	drawText(screen, g.face, "2: タンクトップ小隊のジャッジメントダイス", screenWidth/2, 240, textColor)
	drawText(screen, g.face, "3: タンクトップ小隊ジェンガモード", screenWidth/2, 280, textColor)
	drawText(screen, g.face, "1 / 2 / 3 を押してね", screenWidth/2, screenHeight-60, hintColor)
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
		drawText(screen, g.face, "SPACE でダイスを振る", screenWidth/2, 310, textColor)
	}

	drawText(screen, g.face, "SPACE: Roll    R: Reset    ESC: Title", screenWidth/2, screenHeight-24, hintColor)
}

func (g *Game) drawDuelScene(screen *ebiten.Image) {
	g.duel.Draw(screen, g.face)
	drawText(screen, g.face, "SPACE: Next    ESC: Title", screenWidth/2, screenHeight-24, hintColor)
}

func (g *Game) drawJengaScene(screen *ebiten.Image) {
	g.jenga.Draw(screen, g.face)
	drawText(screen, g.face, "SPACE: Roll    ESC: Title", screenWidth/2, screenHeight-24, hintColor)
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
