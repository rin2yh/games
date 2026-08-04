package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// input は 1 フレーム分の操作を、キーとポインタ（マウス / タッチ）を同じ
// 意味にまとめたもの。シーン側はどちらで操作されたかを気にしなくていい。
type input struct {
	confirm bool // SPACE / 画面のタップ：ダイスを振る、次へ進む
	back    bool // ESC / 「← タイトル」のタップ
	reset   bool // R
	mode    int  // タイトルで選ばれたモード。1〜3、未選択なら 0
}

// rect はゲームの論理座標（screenWidth x screenHeight）での当たり判定。
type rect struct{ x, y, w, h float64 }

func (r rect) contains(x, y int) bool {
	fx, fy := float64(x), float64(y)
	return fx >= r.x && fx < r.x+r.w && fy >= r.y && fy < r.y+r.h
}

// タイトルの選択肢の当たり判定。中心の y は drawTitleScene が文字を置く位置と
// そろえる。横は画面幅いっぱいにして、指でも押しやすくしている。
const titleItemHeight = 36

var titleItemY = [3]float64{200, 240, 280} // モード 1 / 2 / 3 の順

func titleItemRect(i int) rect {
	return rect{
		x: 0,
		y: titleItemY[i] - titleItemHeight/2,
		w: screenWidth,
		h: titleItemHeight,
	}
}

// modeAt はタイトルでタップされた位置に対応するモードを返す。どの選択肢でも
// なければ 0。
func modeAt(x, y int) int {
	for i := range titleItemY {
		if titleItemRect(i).contains(x, y) {
			return i + 1
		}
	}
	return 0
}

// backRect はタイトル以外のシーンで左下に出す「← タイトル」。ESC の代わりで、
// キーを持たない端末から戻る唯一の手段になる。
var backRect = rect{x: 0, y: screenHeight - 48, w: 160, h: 48}

// readInput は Ebitengine から 1 フレーム分の入力を読む。タップの意味は
// シーンによって変わる（タイトルではモード選択、それ以外は決定 / 戻る）ので
// sc を取る。
func readInput(sc scene) input {
	in := input{
		confirm: inpututil.IsKeyJustPressed(ebiten.KeySpace),
		back:    inpututil.IsKeyJustPressed(ebiten.KeyEscape),
		reset:   inpututil.IsKeyJustPressed(ebiten.KeyR),
	}
	switch {
	case inpututil.IsKeyJustPressed(ebiten.Key1):
		in.mode = 1
	case inpututil.IsKeyJustPressed(ebiten.Key2):
		in.mode = 2
	case inpututil.IsKeyJustPressed(ebiten.Key3):
		in.mode = 3
	}

	x, y, ok := justTapped()
	if !ok {
		return in
	}
	if sc == sceneTitle {
		if m := modeAt(x, y); m != 0 {
			in.mode = m
		}
		return in
	}
	if backRect.contains(x, y) {
		in.back = true
		return in
	}
	in.confirm = true
	return in
}

// justTapped は、このフレームで新しく押されたマウス左ボタン / タッチの位置を
// 論理座標で返す。
func justTapped() (x, y int, ok bool) {
	switch {
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
		x, y = ebiten.CursorPosition()
	default:
		ids := inpututil.AppendJustPressedTouchIDs(nil)
		if len(ids) == 0 {
			return 0, 0, false
		}
		x, y = ebiten.TouchPosition(ids[0])
	}
	ox, oy := pointerOffset()
	return x - ox, y - oy, true
}
