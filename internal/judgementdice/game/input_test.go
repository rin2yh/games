package game

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// 打たれた文字での判定。e.code が当てにならない環境ではこちらしか届かないので、
// ここが選択の主役になる。
func TestTyped(t *testing.T) {
	cases := []struct {
		name  string
		chars []rune
		want  rune
		ok    bool
	}{
		{"半角の 1", []rune{'1'}, '1', true},
		{"全角の１", []rune{'１'}, '1', true},
		{"別の数字", []rune{'2'}, '1', false},
		{"何も打たれていない", nil, '1', false},
		{"複数打たれたうちの 1 つ", []rune{'a', '3', 'b'}, '3', true},
		{"半角スペース", []rune{' '}, ' ', true},
	}
	for _, c := range cases {
		if got := typed(c.chars, c.want); got != c.ok {
			t.Errorf("%s: typed(%q, %q) = %v, want %v", c.name, c.chars, c.want, got, c.ok)
		}
	}
}

// 選択肢の文字は、画面に出しているラベルの先頭と一致していないといけない。
// 「1 を押してね」と書いておいて '2' を待つ、のような食い違いを止める。
func TestTitleItemCharMatchesLabel(t *testing.T) {
	for _, it := range titleItems {
		if head := []rune(it.label)[0]; head != it.ch {
			t.Errorf("%q: ラベルの先頭は %q なのに ch = %q", it.label, head, it.ch)
		}
	}
}

// 「1」には上段とテンキーの 2 つがあり、Ebitengine では別のキーになる。
// 片方しか並べていない行があると、そのモードはテンキーから選べない。
func TestTitleItemsAcceptBothNumberRows(t *testing.T) {
	numpad := map[ebiten.Key]bool{
		ebiten.KeyNumpad1: true, ebiten.KeyNumpad2: true, ebiten.KeyNumpad3: true,
	}
	digit := map[ebiten.Key]bool{
		ebiten.Key1: true, ebiten.Key2: true, ebiten.Key3: true,
	}
	for _, it := range titleItems {
		var hasNumpad, hasDigit bool
		for _, k := range it.keys {
			hasNumpad = hasNumpad || numpad[k]
			hasDigit = hasDigit || digit[k]
		}
		if !hasDigit || !hasNumpad {
			t.Errorf("%q: keys = %v（上段=%v テンキー=%v）、両方要る", it.label, it.keys, hasDigit, hasNumpad)
		}
	}
}

func TestModeAt(t *testing.T) {
	cases := []struct {
		name string
		y    int
		want scene
	}{
		{"1 行目", int(titleItems[0].y), sceneEffect},
		{"2 行目", int(titleItems[1].y), sceneDuel},
		{"3 行目", int(titleItems[2].y), sceneJenga},
		{"選択肢の上", 100, sceneTitle},
		{"選択肢の間", 220, sceneTitle},
		{"選択肢の下", 400, sceneTitle},
	}
	for _, c := range cases {
		if got := modeAt(c.y); got != c.want {
			t.Errorf("%s: modeAt(%d) = %v, want %v", c.name, c.y, got, c.want)
		}
	}
}

// 「← タイトル」は画面内に収まっていて、かつ選択肢と重ならない位置にある。
// ここがずれると、タップできる場所と見えている文字が食い違う。
func TestBackRectFits(t *testing.T) {
	if screen := image.Rect(0, 0, screenWidth, screenHeight); !backRect.In(screen) {
		t.Errorf("backRect = %v が画面 %v からはみ出している", backRect, screen)
	}
	for _, it := range titleItems {
		if bottom := it.y + titleItemHeight/2; bottom > float64(backRect.Min.Y) {
			t.Errorf("選択肢 %q の下端 %v が backRect %v と重なる", it.label, bottom, backRect)
		}
	}
}
