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

// titleItems はタイトルの選択肢。表示・入力・タップの当たり判定がこの 1 つの
// 表から出るので、行を足しても消しても食い違いようがない。
//
// 選択の判定は ch（入力された文字）を主にする。画面は「1 を押してね」と文字で
// 頼んでいるのだから、押されたキーの位置ではなく打たれた文字を見るのが素直で、
// かつ確実でもある。Ebitengine の js バックエンドはキーイベントを 2 系統で
// 記録していて、keys が当たるのは e.code（キーの物理位置）が期待通りに来た
// ときだけ。e.code は環境によって空だったり別物だったりする（Ebitengine 自身
// が仮想キーボードでの実例を挙げている #2898、ほかにキーリマップ系のツールを
// 挟んだ場合など）。そのとき e.key は "1" のまま届くので、文字で見ていれば
// 選べる。テンキーの「1」も、位置は KeyNumpad1 だが文字は "1" で来る。
//
// keys は逆に文字が来ない環境のための保険で、両方見ておけばどちらか一方が
// 欠けても通る。
var titleItems = [...]struct {
	label string
	y     float64
	ch    rune
	keys  []ebiten.Key
	mode  scene
}{
	{"1: 遊戯王ジャッジメントダイス", 200, '1', []ebiten.Key{ebiten.Key1, ebiten.KeyNumpad1}, sceneEffect},
	{"2: タンクトップ小隊のジャッジメントダイス", 240, '2', []ebiten.Key{ebiten.Key2, ebiten.KeyNumpad2}, sceneDuel},
	{"3: タンクトップ小隊ジェンガモード", 280, '3', []ebiten.Key{ebiten.Key3, ebiten.KeyNumpad3}, sceneJenga},
}

// 全角の数字は ASCII の数字から 0xFEE0 離れている。日本語入力が全角のままでも
// 「１」で選べるように、どちらも同じ文字として扱う。
const toFullWidth = 0xFEE0

func typed(chars []rune, want rune) bool {
	for _, r := range chars {
		if r == want || r == want+toFullWidth {
			return true
		}
	}
	return false
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
	// 文字とキーの両方を見る理由は titleItems のコメントのとおり。ESC だけは
	// 文字が来ない（e.key が "Escape" で、文字ではなくキー名なので）ため
	// e.code 頼みになる。戻る手段が画面内に「← タイトル」としても出ているのは
	// そのため。
	chars := ebiten.AppendInputChars(nil)
	in := input{
		confirm: typed(chars, ' ') || inpututil.IsKeyJustPressed(ebiten.KeySpace),
		back:    inpututil.IsKeyJustPressed(ebiten.KeyEscape),
		reset:   typed(chars, 'r') || typed(chars, 'R') || inpututil.IsKeyJustPressed(ebiten.KeyR),
	}
	for _, it := range titleItems {
		if typed(chars, it.ch) || anyJustPressed(it.keys) {
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
