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
	cases := []struct {
		mode int
		want scene
	}{
		{1, sceneEffect},
		{2, sceneDuel},
		{3, sceneJenga},
		{0, sceneTitle},
	}
	for _, c := range cases {
		g := New()
		g.handle(input{mode: c.mode})
		if g.scene != c.want {
			t.Errorf("mode %d: scene = %v, want %v", c.mode, g.scene, c.want)
		}
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
	g := New()
	g.scene = sceneEffect
	g.handle(input{confirm: true})
	if g.dice.State() != stateRolling {
		t.Errorf("dice.State() = %v, want stateRolling", g.dice.State())
	}

	g = New()
	g.scene = sceneDuel
	g.handle(input{confirm: true})
	if g.duel.phase != duelRolling {
		t.Errorf("duel.phase = %v, want duelRolling", g.duel.phase)
	}

	g = New()
	g.scene = sceneJenga
	g.handle(input{confirm: true})
	if g.jenga.phase != jengaRolling {
		t.Errorf("jenga.phase = %v, want jengaRolling", g.jenga.phase)
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
