//go:build !darwin

package app

import (
	"time"

	"github.com/gen2brain/beeep"
)

func showRestPrompt(durationSec int) func() {
	go func() {
		_ = beeep.Alert(
			"👀 WAKTUNYA ISTIRAHAT MATA ("+formatDuration(time.Duration(durationSec)*time.Second)+")",
			"Pandanglah objek yang jauh (minimal 6 meter / 20 kaki) selama "+formatDuration(time.Duration(durationSec)*time.Second)+".",
			"",
		)
	}()

	return nil
}
