//go:build !darwin

package app

import "github.com/gen2brain/beeep"

func showRestPrompt(durationSec int) func() {
	go func() {
		_ = beeep.Alert(
			"👀 WAKTUNYA ISTIRAHAT MATA",
			"Pandanglah objek yang jauh (minimal 6 meter / 20 kaki) selama 20 detik.",
			"",
		)
	}()

	return nil
}
