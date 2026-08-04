package game

import (
	"fmt"
	"strings"
	"testing"
)

func TestEffectsLength(t *testing.T) {
	if got, want := len(Effects), 7; got != want {
		t.Fatalf("len(Effects) = %d, want %d", got, want)
	}
}

func TestEffectsSlotZeroEmpty(t *testing.T) {
	if Effects[0] != "" {
		t.Fatalf("Effects[0] = %q, want empty string", Effects[0])
	}
}

func TestEffectsFacesNonEmpty(t *testing.T) {
	for i := 1; i <= 6; i++ {
		if Effects[i] == "" {
			t.Errorf("Effects[%d] is empty", i)
		}
		prefix := fmt.Sprintf("%d：", i)
		if !strings.HasPrefix(Effects[i], prefix) {
			t.Errorf("Effects[%d] = %q, want prefix %q", i, Effects[i], prefix)
		}
	}
}
