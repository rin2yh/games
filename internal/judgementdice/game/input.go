package game

import (
	"image"
	"slices"
	"strings"
	"unicode"

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
// キー判定は逆に文字が来ない環境のための保険で、両方見ておけばどちらか一方が
// 欠けても通る。
var titleItems = [...]struct {
	label string
	y     float64
	ch    rune
	mode  scene
}{
	{"1: 遊戯王ジャッジメントダイス", 200, '1', sceneEffect},
	{"2: タンクトップ小隊のジャッジメントダイス", 240, '2', sceneDuel},
	{"3: タンクトップ小隊ジェンガモード", 280, '3', sceneJenga},
}

// titleHint はタイトル下部の案内。押す文字を titleItems から起こすので、行を
// 足したときに案内だけ古いままにならない。
var titleHint = func() string {
	chars := make([]string, 0, len(titleItems))
	for _, it := range titleItems {
		chars = append(chars, string(it.ch))
	}
	return strings.Join(chars, " / ") + " キー、または選択肢をタップ"
}()

// digitKeys は ch の数字が置かれている 2 か所（上段とテンキー）のキーを返す。
// Ebitengine では別のキーで、定数は 0〜9 が連番に並んでいる。
func digitKeys(ch rune) [2]ebiten.Key {
	d := ebiten.Key(ch - '0')
	return [2]ebiten.Key{ebiten.Key0 + d, ebiten.KeyNumpad0 + d}
}

// 全角の英数字は ASCII から 0xFEE0 離れたところに同じ並びで入っている。
const toFullWidth = 0xFEE0

// normalize は打たれた文字を照合できる形に均す。全角を半角に戻し、大文字を
// 小文字に寄せるので、照合する側は素の文字だけ書けばいい。畳み方をここ 1 か所に
// 押し込めておかないと、キーを足すたびに全角と大小の両方を思い出す羽目になる。
func normalize(chars []rune) []rune {
	for i, r := range chars {
		if '！' <= r && r <= '～' {
			r -= toFullWidth
		}
		chars[i] = unicode.ToLower(r)
	}
	return chars
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
	chars := normalize(ebiten.AppendInputChars(nil))
	in := input{
		confirm: slices.Contains(chars, ' ') || inpututil.IsKeyJustPressed(ebiten.KeySpace),
		back:    inpututil.IsKeyJustPressed(ebiten.KeyEscape),
		reset:   slices.Contains(chars, 'r') || inpututil.IsKeyJustPressed(ebiten.KeyR),
	}
	for _, it := range titleItems {
		keys := digitKeys(it.ch)
		if slices.Contains(chars, it.ch) || slices.ContainsFunc(keys[:], inpututil.IsKeyJustPressed) {
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
