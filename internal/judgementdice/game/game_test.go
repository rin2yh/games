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

func TestHandleTitleSelectsMode(t *testing.T) {
	for _, it := range titleItems {
		g := New()
		g.handle(input{mode: it.mode})
		if g.scene != it.mode {
			t.Errorf("%q: scene = %v, want %v", it.label, g.scene, it.mode)
		}
	}

	g := New()
	g.handle(input{})
	if g.scene != sceneTitle {
		t.Errorf("未選択で scene = %v, want sceneTitle", g.scene)
	}
}

func TestHandleBackReturnsToTitle(t *testing.T) {
	for _, sc := range []scene{sceneEffect, sceneDuel, sceneJenga} {
		g := New()
		g.scene = sc
		g.handle(input{back: true})
		if g.scene != sceneTitle {
			t.Errorf("scene %v: back で戻らず scene = %v", sc, g.scene)
		}
	}
}

// タップは SPACE と同じ意味になる。キーの無い端末でも遊べることの確認。
func TestHandleConfirmRollsDice(t *testing.T) {
	cases := []struct {
		scene   scene
		rolling func(*Game) bool
	}{
		{sceneEffect, func(g *Game) bool { return g.dice.State() == stateRolling }},
		{sceneDuel, func(g *Game) bool { return g.duel.phase == duelRolling }},
		{sceneJenga, func(g *Game) bool { return g.jenga.phase == jengaRolling }},
	}
	for _, c := range cases {
		g := New()
		g.scene = c.scene
		g.handle(input{confirm: true})
		if !c.rolling(g) {
			t.Errorf("scene %v: confirm でダイスが回り出さない", c.scene)
		}
	}
}

func TestHandleResetInEffectScene(t *testing.T) {
	g := New()
	g.scene = sceneEffect
	g.dice.state = stateResult
	g.dice.final = 4
	g.dice.face = 4

	g.handle(input{reset: true})

	if g.dice.state != stateIdle || g.dice.final != 0 || g.dice.face != 1 {
		t.Errorf("dice not reset: state=%v final=%d face=%d", g.dice.state, g.dice.final, g.dice.face)
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
