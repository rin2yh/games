package game

import "testing"

func TestNew(t *testing.T) {
	g := New()
	if g.scene != sceneTitle {
		t.Errorf("scene = %v, want sceneTitle", g.scene)
	}
	if g.dice == nil {
		t.Error("dice is nil")
	}
	if g.duel == nil {
		t.Error("duel is nil")
	}
	if g.jenga == nil {
		t.Error("jenga is nil")
	}
}

func TestLayout(t *testing.T) {
	g := New()
	cases := []struct {
		outW, outH int
	}{
		{0, 0},
		{800, 600},
		{1920, 1080},
	}
	for _, c := range cases {
		w, h := g.Layout(c.outW, c.outH)
		if w != screenWidth || h != screenHeight {
			t.Errorf("Layout(%d, %d) = (%d, %d), want (%d, %d)",
				c.outW, c.outH, w, h, screenWidth, screenHeight)
		}
	}
}
