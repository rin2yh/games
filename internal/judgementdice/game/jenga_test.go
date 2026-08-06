package game

import "testing"

func TestNewJenga(t *testing.T) {
	j := NewJenga()
	if j.phase != jengaIdle {
		t.Errorf("phase = %v, want jengaIdle", j.phase)
	}
	if j.dice == nil {
		t.Fatal("dice is nil")
	}
	if j.dice.state != stateIdle {
		t.Errorf("dice.state = %v, want stateIdle", j.dice.state)
	}
}

func TestJengaReset(t *testing.T) {
	j := NewJenga()
	j.phase = jengaJudgement
	j.dice.state = stateResult
	j.dice.final = 6
	j.dice.face = 6

	j.Reset()

	if j.phase != jengaIdle {
		t.Errorf("phase = %v, want jengaIdle", j.phase)
	}
	if j.dice.state != stateIdle || j.dice.final != 0 || j.dice.face != 1 {
		t.Errorf("dice not reset: state=%v final=%d face=%d", j.dice.state, j.dice.final, j.dice.face)
	}
}

func TestJengaUpdateRollingToJudgement(t *testing.T) {
	j := NewJenga()
	j.phase = jengaRolling
	j.dice.state = stateResult
	j.dice.final = 1
	j.dice.face = 1

	j.Update()

	if j.phase != jengaJudgement {
		t.Errorf("phase = %v, want jengaJudgement", j.phase)
	}
}

func TestJengaConfirmStartsAndRestarts(t *testing.T) {
	j := NewJenga()

	j.Confirm()
	if j.phase != jengaRolling {
		t.Fatalf("phase = %v, want jengaRolling", j.phase)
	}
	if j.dice.state != stateRolling {
		t.Errorf("dice.state = %v, want stateRolling", j.dice.state)
	}

	// 振っている最中の操作は効かない。
	j.Confirm()
	if j.phase != jengaRolling {
		t.Errorf("phase = %v, want jengaRolling", j.phase)
	}

	j.phase = jengaJudgement
	j.Confirm()
	if j.phase != jengaIdle {
		t.Errorf("判定後の Confirm で仕切り直しにならない: phase = %v", j.phase)
	}
}
