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
