package game

import "testing"

func TestRectContains(t *testing.T) {
	r := rect{x: 10, y: 20, w: 100, h: 40}
	cases := []struct {
		name string
		x, y int
		want bool
	}{
		{"左上の角", 10, 20, true},
		{"真ん中", 60, 40, true},
		{"右下の角は含まない", 110, 60, false},
		{"右端の内側", 109, 59, true},
		{"左に外れる", 9, 40, false},
		{"上に外れる", 60, 19, false},
	}
	for _, c := range cases {
		if got := r.contains(c.x, c.y); got != c.want {
			t.Errorf("%s: contains(%d, %d) = %v, want %v", c.name, c.x, c.y, got, c.want)
		}
	}
}

func TestModeAt(t *testing.T) {
	cases := []struct {
		name string
		x, y int
		want int
	}{
		{"1 行目", 320, int(titleItemY[0]), 1},
		{"2 行目", 10, int(titleItemY[1]), 2},
		{"3 行目", 630, int(titleItemY[2]), 3},
		{"選択肢の上", 320, 100, 0},
		{"選択肢の間", 320, 220, 0},
		{"選択肢の下", 320, 400, 0},
	}
	for _, c := range cases {
		if got := modeAt(c.x, c.y); got != c.want {
			t.Errorf("%s: modeAt(%d, %d) = %d, want %d", c.name, c.x, c.y, got, c.want)
		}
	}
}

// 「← タイトル」は画面内に収まっていて、かつ下端のヒント行と重ならない位置に
// ある。ここがずれると、タップできる場所と見えている文字が食い違う。
func TestBackRectIsOnScreen(t *testing.T) {
	if backRect.x < 0 || backRect.y < 0 {
		t.Errorf("backRect = %+v, 画面の外にはみ出している", backRect)
	}
	if backRect.x+backRect.w > screenWidth || backRect.y+backRect.h > screenHeight {
		t.Errorf("backRect = %+v, 画面の外にはみ出している", backRect)
	}
	for i := range titleItemY {
		if r := titleItemRect(i); r.y+r.h > backRect.y {
			t.Errorf("titleItemRect(%d) = %+v が backRect = %+v と重なる", i, r, backRect)
		}
	}
}
