package game

import (
	"image"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// 打たれた文字での判定。e.code が当てにならない環境ではこちらしか届かないので、
// ここが選択の主役になる。照合する側が素の文字だけ書けるよう、全角と大小は
// normalize が吸収する。
func TestNormalize(t *testing.T) {
	cases := []struct {
		name  string
		chars []rune
		want  []rune
	}{
		{"半角はそのまま", []rune{'1'}, []rune{'1'}},
		{"全角の数字は半角へ", []rune{'１', '２', '３'}, []rune{'1', '2', '3'}},
		{"大文字は小文字へ", []rune{'R'}, []rune{'r'}},
		{"全角の大文字も", []rune{'Ｒ'}, []rune{'r'}},
		{"スペースはそのまま", []rune{' '}, []rune{' '}},
		{"畳めない文字は触らない", []rune{'あ'}, []rune{'あ'}},
		{"空", nil, nil},
	}
	for _, c := range cases {
		if got := normalize(c.chars); string(got) != string(c.want) {
			t.Errorf("%s: normalize(%q) = %q, want %q", c.name, c.chars, got, c.want)
		}
	}
}

// digitKeys は Ebitengine のキー定数が 0〜9 の連番であることに乗っている。
// 並びが変わったらここで落ちる。
func TestDigitKeys(t *testing.T) {
	cases := []struct {
		ch   rune
		want [2]ebiten.Key
	}{
		{'1', [2]ebiten.Key{ebiten.Key1, ebiten.KeyNumpad1}},
		{'2', [2]ebiten.Key{ebiten.Key2, ebiten.KeyNumpad2}},
		{'3', [2]ebiten.Key{ebiten.Key3, ebiten.KeyNumpad3}},
	}
	for _, c := range cases {
		if got := digitKeys(c.ch); got != c.want {
			t.Errorf("digitKeys(%q) = %v, want %v", c.ch, got, c.want)
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

// 画面下の案内は titleItems から起こす。行を足したのに案内が古いまま、が
// 起きないこと。
func TestTitleHintListsEveryItem(t *testing.T) {
	for _, it := range titleItems {
		if !strings.ContainsRune(titleHint, it.ch) {
			t.Errorf("titleHint = %q に %q が出てこない", titleHint, it.ch)
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
