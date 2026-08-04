package game

import "testing"

func assertFaceInRange(t *testing.T, face int) {
	t.Helper()
	if face < 1 || face > 6 {
		t.Fatalf("face = %d, want in [1,6]", face)
	}
}

func TestNewDice(t *testing.T) {
	d := NewDice()
	if d.state != stateIdle {
		t.Errorf("state = %v, want stateIdle", d.state)
	}
	if d.face != 1 {
		t.Errorf("face = %d, want 1", d.face)
	}
	if d.final != 0 {
		t.Errorf("final = %d, want 0", d.final)
	}
	if d.tickCount != 0 {
		t.Errorf("tickCount = %d, want 0", d.tickCount)
	}
}

func TestRollTransitions(t *testing.T) {
	d := NewDice()
	d.Roll()
	if d.state != stateRolling {
		t.Errorf("state = %v, want stateRolling", d.state)
	}
	if d.tickCount != 0 {
		t.Errorf("tickCount = %d, want 0", d.tickCount)
	}
	assertFaceInRange(t, d.face)
}

func TestRollIgnoredWhileRolling(t *testing.T) {
	d := NewDice()
	d.Roll()
	d.Update()
	d.Update()
	before := d.tickCount
	if before == 0 {
		t.Fatalf("precondition failed: tickCount should advance, got %d", before)
	}
	d.Roll()
	if d.state != stateRolling {
		t.Errorf("state = %v, want stateRolling", d.state)
	}
	if d.tickCount != before {
		t.Errorf("tickCount = %d, want %d (Roll must be ignored while rolling)", d.tickCount, before)
	}
}

func TestReset(t *testing.T) {
	d := NewDice()
	d.state = stateResult
	d.face = 4
	d.final = 4
	d.tickCount = 10

	d.Reset()

	if d.state != stateIdle {
		t.Errorf("state = %v, want stateIdle", d.state)
	}
	if d.face != 1 {
		t.Errorf("face = %d, want 1", d.face)
	}
	if d.final != 0 {
		t.Errorf("final = %d, want 0", d.final)
	}
	if d.tickCount != 0 {
		t.Errorf("tickCount = %d, want 0", d.tickCount)
	}
}

func TestUpdateRollingReachesResult(t *testing.T) {
	d := NewDice()
	d.Roll()
	for range rollDuration {
		d.Update()
	}
	if d.state != stateResult {
		t.Errorf("state = %v, want stateResult", d.state)
	}
	if d.tickCount != 0 {
		t.Errorf("tickCount = %d, want 0", d.tickCount)
	}
	assertFaceInRange(t, d.final)
	if d.face != d.final {
		t.Errorf("face = %d, want %d (== final)", d.face, d.final)
	}
}

func TestUpdateResultBounceCaps(t *testing.T) {
	d := NewDice()
	d.state = stateResult
	d.face = 3
	d.final = 3
	d.tickCount = 0

	for range bounceFrames * 3 {
		d.Update()
	}
	if d.tickCount != bounceFrames {
		t.Errorf("tickCount = %d, want %d (capped at bounceFrames)", d.tickCount, bounceFrames)
	}
	if d.state != stateResult {
		t.Errorf("state = %v, want stateResult", d.state)
	}
}

func TestUpdateFaceAlwaysValidWhileRolling(t *testing.T) {
	d := NewDice()
	d.Roll()
	assertFaceInRange(t, d.face)
	for range rollDuration {
		d.Update()
		assertFaceInRange(t, d.face)
	}
}
