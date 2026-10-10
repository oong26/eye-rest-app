package app

import (
	"testing"
	"time"
)

func TestClampRestDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     time.Duration
	}{
		{name: "below minimum", duration: time.Second, want: minRestDuration},
		{name: "inside range", duration: 45 * time.Second, want: 45 * time.Second},
		{name: "above maximum", duration: 3 * time.Minute, want: maxRestDuration},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clampRestDuration(tt.duration); got != tt.want {
				t.Fatalf("clampRestDuration(%v) = %v, want %v", tt.duration, got, tt.want)
			}
		})
	}
}

func TestAdjustRestDurationUsesStepAndBounds(t *testing.T) {
	setRestDuration(defaultRestDuration)
	t.Cleanup(func() {
		setRestDuration(defaultRestDuration)
	})

	if got := adjustRestDuration(restDurationStep); got != 25*time.Second {
		t.Fatalf("adjustRestDuration(+step) = %v, want 25s", got)
	}

	if got := adjustRestDuration(-10 * time.Minute); got != minRestDuration {
		t.Fatalf("adjustRestDuration(below min) = %v, want %v", got, minRestDuration)
	}

	if got := adjustRestDuration(10 * time.Minute); got != maxRestDuration {
		t.Fatalf("adjustRestDuration(above max) = %v, want %v", got, maxRestDuration)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{duration: 20 * time.Second, want: "20s"},
		{duration: time.Minute, want: "1m"},
		{duration: 90 * time.Second, want: "1m30s"},
	}

	for _, tt := range tests {
		if got := formatDuration(tt.duration); got != tt.want {
			t.Fatalf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
		}
	}
}
