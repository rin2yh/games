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
	g.handle(readInput())
	return nil
}

// handle は 1 フレーム分の入力を今のシーンに適用する。入力の読み取り
// (readInput) と分けてあるので、Ebitengine を起動せずにテストできる。
func (g *Game) handle(in input) {
	if in.back && g.scene != sceneTitle {
		// 出るときに畳むので、今いるシーン以外はすでに初期状態にある。
		g.dice.Reset()
		g.duel.Reset()
		g.jenga.Reset()
		g.scene = sceneTitle
		return
	}

	switch g.scene {
	case sceneTitle:
		// 未選択なら in.mode は sceneTitle なので、そのまま入れて動かない。
		g.scene = in.mode
	case sceneEffect:
		g.updateEffect(in)
	case sceneDuel:
		// 操作が来たときだけ Confirm。進行の Update より先に呼ぶのは、振り
		// 終わったフレームの操作まで拾って結果を飛ばしてしまわないため。
		if in.confirm {
			g.duel.Confirm()
		}
		g.duel.Update()
	case sceneJenga:
		if in.confirm {
			g.jenga.Confirm()
		}
		g.jenga.Update()
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
	for _, it := range titleItems {
		drawText(screen, g.face, it.label, screenWidth/2, it.y, textColor)
	}
	drawText(screen, g.face, titleHint, screenWidth/2, screenHeight-60, hintColor)
}

// drawFooter はモード共通の下端。操作の案内と、backRect の位置に置く
// 「← タイトル」を必ず一緒に描く。戻り口が見えていないと、キーの無い端末では
// モードから出られないので、片方だけ描くことが起きないようにまとめてある。
func (g *Game) drawFooter(screen *ebiten.Image, hint string) {
	drawText(screen, g.face, hint, screenWidth/2, screenHeight-24, hintColor)

	x := float64(backRect.Min.X+backRect.Max.X) / 2
	y := float64(backRect.Min.Y+backRect.Max.Y) / 2
	drawText(screen, g.face, "← タイトル", x, y, hintColor)
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

	g.drawFooter(screen, "SPACE / タップ: Roll    R: Reset    ESC: Title")
}

func (g *Game) drawDuelScene(screen *ebiten.Image) {
	g.duel.Draw(screen, g.face)
	g.drawFooter(screen, "SPACE / タップ: Next    ESC: Title")
}

func (g *Game) drawJengaScene(screen *ebiten.Image) {
	g.jenga.Draw(screen, g.face)
	g.drawFooter(screen, "SPACE / タップ: Roll    ESC: Title")
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
