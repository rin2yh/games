package game

import "testing"

func TestNewDuel(t *testing.T) {
	d := NewDuel()
	if d.phase != duelIdle {
		t.Errorf("phase = %v, want duelIdle", d.phase)
	}
	if d.cpu == nil {
		t.Fatal("cpu is nil")
	}
	if d.player == nil {
		t.Fatal("player is nil")
	}
	if d.cpu.state != stateIdle {
		t.Errorf("cpu.state = %v, want stateIdle", d.cpu.state)
	}
	if d.player.state != stateIdle {
		t.Errorf("player.state = %v, want stateIdle", d.player.state)
	}
}

func TestDuelReset(t *testing.T) {
	d := NewDuel()
	d.phase = duelJudgement
	d.cpu.state = stateResult
	d.cpu.final = 5
	d.cpu.face = 5
	d.player.state = stateResult
	d.player.final = 2
	d.player.face = 2

	d.Reset()

	if d.phase != duelIdle {
		t.Errorf("phase = %v, want duelIdle", d.phase)
	}
	if d.cpu.state != stateIdle || d.cpu.final != 0 || d.cpu.face != 1 {
		t.Errorf("cpu not reset: state=%v final=%d face=%d", d.cpu.state, d.cpu.final, d.cpu.face)
	}
	if d.player.state != stateIdle || d.player.final != 0 || d.player.face != 1 {
		t.Errorf("player not reset: state=%v final=%d face=%d", d.player.state, d.player.final, d.player.face)
	}
}

func TestResultPlayerWins(t *testing.T) {
	d := NewDuel()
	d.player.final = 2
	d.cpu.final = 5
	if got, want := d.Result(), "あんたの勝ち！"; got != want {
		t.Errorf("Result() = %q, want %q", got, want)
	}
}

func TestResultPlayerLosesOnTie(t *testing.T) {
	d := NewDuel()
	d.player.final = 3
	d.cpu.final = 3
	if got, want := d.Result(), "あんたの負け…"; got != want {
		t.Errorf("Result() = %q, want %q", got, want)
	}
}

func TestResultPlayerLosesOnHigher(t *testing.T) {
	d := NewDuel()
	d.player.final = 6
	d.cpu.final = 1
	if got, want := d.Result(), "あんたの負け…"; got != want {
		t.Errorf("Result() = %q, want %q", got, want)
	}
}

func TestDuelUpdateRollingToJudgement(t *testing.T) {
	d := NewDuel()
	d.phase = duelRolling
	d.cpu.state = stateResult
	d.cpu.final = 4
	d.cpu.face = 4
	d.player.state = stateResult
	d.player.final = 2
	d.player.face = 2

	d.Update()

	if d.phase != duelJudgement {
		t.Errorf("phase = %v, want duelJudgement", d.phase)
	}
}

func TestDuelUpdateRollingStaysWhenNotBothDone(t *testing.T) {
	d := NewDuel()
	d.phase = duelRolling
	d.cpu.state = stateResult
	d.cpu.final = 4
	d.cpu.face = 4
	d.player.state = stateRolling
	d.player.tickCount = 10

	d.Update()

	if d.phase != duelRolling {
		t.Errorf("phase = %v, want duelRolling", d.phase)
	}
}
