package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/rin2yh/games/internal/webpointer"
)

// input は 1 フレーム分の操作を、キーとポインタ（マウス / タッチ）を同じ
// 意味にまとめたもの。シーン側はどちらで操作されたかを気にしなくていい。
type input struct {
	confirm bool  // SPACE / 画面のタップ：ダイスを振る、次へ進む
	back    bool  // ESC / 「← タイトル」のタップ
	reset   bool  // R
	mode    scene // タイトルで選ばれたモード。未選択なら sceneTitle
}

// titleItems はタイトルの選択肢。表示・キー・タップの当たり判定がこの 1 つの
// 表から出るので、行を足しても消しても食い違いようがない。
//
// keys に上段の数字とテンキーの両方を並べるのは、「1」を押したつもりの入力が
// キーの位置で別物になるから。Ebitengine にとって上段は Key1、テンキーは
// KeyNumpad1 で、Key1 しか見ていないとテンキーの数字は無視される。
var titleItems = [...]struct {
	label string
	y     float64
	keys  []ebiten.Key
	mode  scene
}{
	{"1: 遊戯王ジャッジメントダイス", 200, []ebiten.Key{ebiten.Key1, ebiten.KeyNumpad1}, sceneEffect},
	{"2: タンクトップ小隊のジャッジメントダイス", 240, []ebiten.Key{ebiten.Key2, ebiten.KeyNumpad2}, sceneDuel},
	{"3: タンクトップ小隊ジェンガモード", 280, []ebiten.Key{ebiten.Key3, ebiten.KeyNumpad3}, sceneJenga},
}

const titleItemHeight = 36

// modeAt はタイトルで y の位置をタップしたときに選ばれるモードを返す。
// どの選択肢でもなければ sceneTitle。行は画面幅いっぱいを受けるので、指でも
// 押しやすいように x は見ない。
func modeAt(y int) scene {
	for _, it := range titleItems {
		if float64(y) >= it.y-titleItemHeight/2 && float64(y) < it.y+titleItemHeight/2 {
			return it.mode
		}
	}
	return sceneTitle
}

// backRect はタイトル以外のシーンで左下に出す「← タイトル」。ESC の代わりで、
// キーを持たない端末から戻る唯一の手段になる。
var backRect = image.Rect(0, screenHeight-48, 160, screenHeight)

// readInput は Ebitengine から 1 フレーム分の入力を読む。タップの意味はシーン
// によって変わるが、シーン側が自分に関係するものだけ見るので、ここでは当たった
// ものを全部立てておけばいい。
func readInput() input {
	in := input{
		confirm: inpututil.IsKeyJustPressed(ebiten.KeySpace),
		back:    inpututil.IsKeyJustPressed(ebiten.KeyEscape),
		reset:   inpututil.IsKeyJustPressed(ebiten.KeyR),
	}
	for _, it := range titleItems {
		if anyJustPressed(it.keys) {
			in.mode = it.mode
			break
		}
	}

	x, y, ok := justTapped()
	if !ok {
		return in
	}
	if m := modeAt(y); m != sceneTitle {
		in.mode = m
	}
	if image.Pt(x, y).In(backRect) {
		in.back = true
	} else {
		in.confirm = true
	}
	return in
}

func anyJustPressed(keys []ebiten.Key) bool {
	for _, k := range keys {
		if inpututil.IsKeyJustPressed(k) {
			return true
		}
	}
	return false
}

// justTapped はこのフレームで新しく押されたマウス左ボタン / タッチの位置を、
// ブラウザでのずれを直したうえで返す。
func justTapped() (x, y int, ok bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y = ebiten.CursorPosition()
	} else {
		ids := inpututil.AppendJustPressedTouchIDs(nil)
		if len(ids) == 0 {
			return 0, 0, false
		}
		x, y = ebiten.TouchPosition(ids[0])
	}
	ox, oy := webpointer.Offset(screenWidth, screenHeight)
	return x - ox, y - oy, true
}
